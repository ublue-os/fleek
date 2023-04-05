package flake

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/pterm/pterm"

	"github.com/riywo/loginshell"
	app "github.com/ublue-os/fleek"
	"github.com/ublue-os/fleek/internal/fleek"
	"github.com/ublue-os/fleek/internal/ux"
)

const nixbin = "nix"
const gitbin = "git"

var ErrPackageConflict = errors.New("package exists in fleek and nix profile")

type Flake struct {
	Templates *template.Template
	Config    *fleek.Config
	app       *app.App
}
type Data struct {
	Config   *fleek.Config
	UserName string
	Home     string
	Bling    *fleek.Bling
}

func Load(cfg *fleek.Config, app *app.App) (*Flake, error) {
	if cfg.Verbose {
		ux.Verbose.Println(app.Trans("flake.initializingTemplates"))
	}
	t, err := template.ParseFS(content, "*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return &Flake{
		Templates: t,
		Config:    cfg,
		app:       app,
	}, nil
}
func (f *Flake) runGit(cmd string, cmdLine []string) ([]byte, error) {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Dir = f.Config.UserFlakeDir()
	command.Env = os.Environ()
	return command.Output()

}

// runGitI runs the command without capturing stdin/out so
// user can accept host keys interactively
func (f *Flake) runGitI(cmd string, cmdLine []string) error {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	command.Dir = f.Config.UserFlakeDir()
	command.Env = os.Environ()
	return command.Run()

}
func (f *Flake) Clone(repo, branch, privateKey string, promptPass bool) error {
	ux.Info.Println(f.app.Trans("flake.cloning", repo))

	var gco *git.CloneOptions
	refName := plumbing.NewBranchReferenceName(branch)
	if strings.Contains(repo, "git@") {
		var password string
		if promptPass {
			// ask for it
			// TODO: don't echo to terminal
			password, _ = pterm.DefaultInteractiveTextInput.Show("Private key password:")

		} else {
			password = ""
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		keyPath := filepath.Join(home, privateKey)
		ux.Info.Println(f.app.Trans("flake.gettingKeys"))
		publicKeys, err := ssh.NewPublicKeysFromFile("git", keyPath, password)
		if err != nil {
			ux.Error.Println(f.app.Trans("flake.gettingKeys"))
			return err
		}
		gco = &git.CloneOptions{
			Auth:          publicKeys,
			URL:           repo,
			ReferenceName: refName,
		}
	} else {
		gco = &git.CloneOptions{
			URL:               repo,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			ReferenceName:     refName,
		}
	}

	err := f.ensureFlakeDir()
	if err != nil {
		return err
	}

	r, err := git.PlainClone(f.Config.UserFlakeDir(), false, gco)
	if err != nil {
		return err
	}
	// ... retrieving the branch being pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		return err
	}
	if f.Config.Verbose {
		ux.Info.Println(ref.Name())
	}
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}
	if f.Config.Verbose {
		ux.Info.Println(commit)
	}

	ux.Info.Println(f.app.Trans("flake.gitConfigs"))
	if err != nil {
		return err
	}

	err = f.setRebase()
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	yamlPath := filepath.Join(f.Config.FlakeDir, ".fleek.yml")
	csym := filepath.Join(home, ".fleek.yml")

	err = os.Symlink(yamlPath, csym)
	if err != nil {
		return err
	}

	// reread config file
	err = f.ReadConfig()
	if err != nil {
		return err
	}
	_, err = f.Config.CurrentSystem()
	if err != nil {
		if strings.Contains(err.Error(), "not") {
			ux.Info.Println(f.app.Trans("apply.newSystem"))

			// make a new system
			// prompt for git configuration
			email, err := pterm.DefaultInteractiveTextInput.Show(f.app.Trans("init.gitEmail"))
			if err != nil {
				return err
			}

			name, err := pterm.DefaultInteractiveTextInput.Show(f.app.Trans("init.gitName"))
			if err != nil {
				return err
			}
			// create new system struct
			sys, err := fleek.NewSystem(email, name, privateKey)
			if err != nil {
				return err
			}
			ux.Info.Printfln(f.app.Trans("init.newSystem", sys.Username, sys.Hostname))

			// append new(current) system
			f.Config.Systems = append(f.Config.Systems, sys)
			// save it
			err = f.Config.Save()
			if err != nil {
				return err
			}
			f.writeSystem(*sys, f.Config.Force)
			if err != nil {
				return err
			}
			err = f.Commit("add system: " + sys.Hostname)
			if err != nil {
				return err
			}
		}
	}
	if f.Config.Branch != branch {
		f.Config.Branch = branch
		err = f.Config.Save()
		if err != nil {
			return err
		}
	}
	err = f.Commit("fleek clone")
	if err != nil {
		return err
	}

	return err
}
func (f *Flake) Update() error {
	spinner, err := ux.Spinner().Start(f.app.Trans("flake.update"))
	if err != nil {
		return err
	}
	updateCmdLine := []string{"flake", "update"}
	out, err := f.runNix(nixbin, updateCmdLine)

	if err != nil {
		return err
	}
	if f.Config.Verbose {
		ux.Verbose.Println(out)
	}
	spinner.Success()
	return nil
}
func (f *Flake) Create(force bool) error {
	ux.Info.Println(f.app.Trans("init.writeConfigs"))
	err := f.ensureFlakeDir()
	if err != nil {
		return err
	}
	fullShell, err := loginshell.Shell()
	if err != nil {
		return err
	}
	shell := filepath.Base(fullShell)
	f.Config.Shell = shell
	err = f.Config.Validate()
	if err != nil {
		return err
	}
	var bling *fleek.Bling

	switch f.Config.Bling {
	case "high":
		bling, err = fleek.HighBling()
	case "default":
		bling, err = fleek.DefaultBling()
	case "low":
		bling, err = fleek.LowBling()
	case "none":
		bling, err = fleek.NoBling()
	default:
		bling, err = fleek.DefaultBling()
	}
	if err != nil {
		return err
	}

	ux.Info.Println(f.app.Trans("init.blingLevel", f.Config.Bling))

	data := Data{
		Config: f.Config,
		Bling:  bling,
	}

	err = f.writeFile("flake.nix", data, force)
	if err != nil {
		return err
	}

	err = f.writeFile("home.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("aliases.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("path.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("programs.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("shell.nix", data, force)
	if err != nil {
		return err
	}
	// make a new system
	// prompt for git configuration
	email, err := pterm.DefaultInteractiveTextInput.Show(f.app.Trans("init.gitEmail"))
	if err != nil {
		return err
	}

	name, err := pterm.DefaultInteractiveTextInput.Show(f.app.Trans("init.gitName"))
	if err != nil {
		return err
	}
	err = f.Config.WriteInitialConfig(email, name, f.Config.Force)
	if err != nil {
		return err
	}
	sys, err := f.Config.CurrentSystem()
	if err != nil {
		return err
	}
	err = f.writeSystem(*sys, force)
	if err != nil {
		return err
	}

	err = f.writeFile("user.nix", data, force)
	if err != nil {
		return err
	}

	err = f.initGit("initial commit")
	if err != nil {
		return err
	}
	return nil

}

func (f *Flake) Check() ([]byte, error) {
	checkCmdLine := []string{"run", "--impure", "home-manager/master", "build", "--impure", "--", "--flake", "."}
	out, err := f.runNix(nixbin, checkCmdLine)

	if err != nil {
		return out, err
	}
	return out, nil
}

// Write writes the applied flake configuration
func (f *Flake) Write(includeSystems bool) error {
	spinner, err := ux.Spinner().Start(f.app.Trans("flake.writing"))
	if err != nil {
		return err
	}

	var bling *fleek.Bling
	switch f.Config.Bling {
	case "high":
		bling, err = fleek.HighBling()
	case "default":
		bling, err = fleek.DefaultBling()
	case "low":
		bling, err = fleek.LowBling()
	case "none":
		bling, err = fleek.NoBling()
	default:
		bling, err = fleek.DefaultBling()
	}
	if err != nil {
		return err
	}
	data := Data{
		Config: f.Config,
		Bling:  bling,
	}
	err = f.writeFile("flake.nix", data, true)
	if err != nil {
		return err
	}
	err = f.writeFile("home.nix", data, true)
	if err != nil {
		return err
	}
	err = f.writeFile("aliases.nix", data, true)
	if err != nil {
		return err
	}
	err = f.writeFile("path.nix", data, true)
	if err != nil {
		return err
	}
	err = f.writeFile("programs.nix", data, true)
	if err != nil {
		return err
	}
	if includeSystems {
		for _, sys := range data.Config.Systems {
			err = f.writeSystem(*sys, true)
			if err != nil {
				return err
			}
		}
	}

	err = f.writeFile("shell.nix", data, true)
	if err != nil {
		return err
	}

	spinner.Success()
	return nil

}
func (f *Flake) initGit(msg string) error {
	floc := f.Config.UserFlakeDir()

	dotGit := filepath.Join(floc, ".git")
	store := osfs.New(dotGit)
	r, err := git.Init(filesystem.NewStorage(store, cache.NewObjectLRUDefault()), store)
	if err != nil {
		return err
	}

	br := plumbing.NewBranchReferenceName("main")

	h := plumbing.NewSymbolicReference(plumbing.HEAD, br)

	// The created reference is saved in the storage.
	err = r.Storer.SetReference(h)
	if err != nil {
		return err
	}

	gitIgnore, err := os.Create(filepath.Join(floc, ".gitignore"))
	if err != nil {
		return err
	}
	defer gitIgnore.Close()
	_, err = gitIgnore.WriteString("result")
	if err != nil {
		return err
	}

	return f.Commit(msg)

}

func (f *Flake) ensureFlakeDir() error {
	if f.Config.Verbose {
		ux.Verbose.Println(f.app.Trans("flake.ensureDir"))
	}
	err := f.Config.MakeFlakeDir()
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil
		}
	}
	return nil
}
func (f *Flake) setRebase() error {
	var output []byte
	configCmdLine := []string{"config", "pull.rebase", "true"}
	out, err := f.runGit(gitbin, configCmdLine)
	output = append(output, out...)
	if err != nil {
		if f.Config.Verbose {
			ux.Verbose.Println(string(output))
		}
		return fmt.Errorf("git config: %w", err)
	}
	return err
}
func (f *Flake) Commit(msg string) error {
	r, err := git.PlainOpen(f.Config.UserFlakeDir())
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	status, err := w.Status()
	if err != nil {
		return err
	}
	if len(status) > 0 {
		ux.Info.Printfln("found %d local changes to commit", len(status))
		if f.Config.Verbose {
			ux.Verbose.Println(status)
		}
		for file := range status {
			_, err = w.Add(file)
			if err != nil {
				return err
			}
		}

		_, err = w.Commit(msg, &git.CommitOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
func (f *Flake) Dirty() (bool, error) {
	r, err := git.PlainOpen(f.Config.UserFlakeDir())
	if err != nil {
		return false, err
	}
	w, err := r.Worktree()
	if err != nil {
		return false, err
	}
	status, err := w.Status()
	if err != nil {
		return false, err
	}
	if len(status) > 0 {
		ux.Info.Printfln("found %d local changes to commit", len(status))
		if f.Config.Verbose {
			ux.Verbose.Println(status)
		}

		return true, nil
	}
	return false, nil
}

func (f *Flake) ReadConfig() error {
	// load the new config
	config, err := fleek.ReadConfig()
	if err != nil {
		return err
	}
	f.Config = config
	return nil
}
func (f *Flake) writeFile(fname string, d Data, force bool) error {

	fpath := filepath.Join(f.Config.UserFlakeDir(), fname)
	_, err := os.Stat(fpath)
	if force || os.IsNotExist(err) {

		file, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer file.Close()
		tmplName := fname + ".tmpl"
		if err = f.Templates.ExecuteTemplate(file, tmplName, d); err != nil {
			return err
		}
	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}
	return nil
}
func (f *Flake) writeSystem(sys fleek.System, force bool) error {

	hostPath := filepath.Join(f.Config.UserFlakeDir(), sys.Hostname)
	err := os.MkdirAll(hostPath, 0755)
	if err != nil {
		return err
	}
	fpath := filepath.Join(hostPath, sys.Hostname+".nix")
	_, err = os.Stat(fpath)
	if force || os.IsNotExist(err) {

		file, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer file.Close()
		tmplName := "host.nix.tmpl"
		if err = f.Templates.ExecuteTemplate(file, tmplName, sys); err != nil {
			return err
		}
	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}
	upath := filepath.Join(hostPath, "user.nix")
	_, err = os.Stat(upath)
	if os.IsNotExist(err) {

		file, err := os.Create(upath)
		if err != nil {
			return err
		}
		defer file.Close()
		tmplName := "user.nix.tmpl"
		if err = f.Templates.ExecuteTemplate(file, tmplName, sys); err != nil {
			return err
		}
	}
	return nil
}
func (f *Flake) Apply(msg string) error {
	spinner, err := ux.Spinner().Start(f.app.Trans("flake.apply"))
	if err != nil {
		return err
	}
	if msg == "" {
		msg = "pre-apply"
	}

	err = f.Commit("pre-apply")
	if err != nil {
		return err
	}
	user, err := fleek.Username()

	if err != nil {
		return err
	}
	host, err := fleek.Hostname()
	if err != nil {
		return err
	}
	applyCmdLine := []string{"run", "--impure", "home-manager/master", "--", "-b", "bak", "switch", "--flake", ".#" + user + "@" + host}
	out, err := f.runNix(nixbin, applyCmdLine)
	if err != nil {
		if bytes.Contains(out, []byte("priority")) {
			return ErrPackageConflict
		}
		if bytes.Contains(out, []byte("conflict")) {
			return ErrPackageConflict
		}
		return err
	}
	spinner.Success()
	return nil
}
func (f *Flake) runNix(cmd string, cmdLine []string) ([]byte, error) {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Dir = f.Config.UserFlakeDir()
	command.Env = os.Environ()
	if f.Config.Unfree {
		command.Env = append(command.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	return command.Output()

}
func (f *Flake) Pull() error {
	spinner, err := ux.Spinner().Start(f.app.Trans("flake.pulling"))
	if err != nil {
		return err
	}

	remote, err := f.Remote()
	if err != nil {
		return err
	}
	// if no remote, no need to push
	if remote == "" {
		return err
	}
	pullCmdline := []string{"pull", "origin", "main"}
	out, err := f.runGit(gitbin, pullCmdline)
	if err != nil {
		if f.Config.Verbose {
			ux.Verbose.Println(string(out))
		}
		return fmt.Errorf("git pull: %w", err)
	}
	if f.Config.Verbose {
		ux.Verbose.Println(string(out))
	}
	spinner.Success()

	return nil
}
func (f *Flake) Push() error {
	spinner, err := ux.Spinner().Start(f.app.Trans("flake.pushing"))
	if err != nil {
		return err
	}

	remote, err := f.Remote()
	if err != nil {
		return err
	}
	// if no remote, no need to push
	if remote == "" {
		return nil
	}
	pushCmdline := []string{"push", "-u", "origin", "main"}
	out, err := f.runGit(gitbin, pushCmdline)
	if err != nil {
		if f.Config.Verbose {
			ux.Verbose.Println(string(out))
		}
		return fmt.Errorf("git push: %w", err)
	}
	if f.Config.Verbose {
		ux.Verbose.Println(string(out))
	}
	spinner.Success()
	return nil

}
func (f *Flake) Sync(msg string) error {
	err := f.setRebase()
	if err != nil {
		return err
	}

	r, err := git.PlainOpen(f.Config.UserFlakeDir())
	if err != nil {
		return err
	}
	// ensure we have a remote
	list, err := r.Remotes()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return errors.New(f.app.Trans("flake.errNoRemotes"))
	}

	err = f.Commit(msg)
	if err != nil {
		return err
	}
	spinner, err := ux.Spinner().Start(f.app.Trans("flake.fetching"))
	if err != nil {
		return err
	}

	// fetch and merge
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			if f.Config.Verbose {
				ux.Verbose.Println(f.app.Trans("flake.uptodate"))
			}
		} else {
			ux.Error.Println("fetch:", err)
			return err
		}
	}
	spinner.Success()
	// pull
	err = f.Pull()
	if err != nil {
		return err
	}

	// push
	err = f.Push()
	if err != nil {
		return err
	}

	return nil
}
func (f *Flake) Remote() (string, error) {
	var err error
	r, err := git.PlainOpen(f.Config.UserFlakeDir())
	if err != nil {
		return "", err
	}
	list, err := r.Remotes()
	if err != nil {
		return "", fmt.Errorf("getting remotes	: %w", err)
	}
	var urls string
	for _, r := range list {
		for _, upstream := range r.Config().URLs {
			urls = urls + upstream + "\n"
		}
	}
	return urls, nil
}
func (f *Flake) RemoteAdd(remote string, name string) error {
	var err error
	r, err := git.PlainOpen(f.Config.UserFlakeDir())
	if err != nil {
		return err
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{remote},
	})
	return err
}

var (
	//go:embed *.tmpl
	content embed.FS
)
