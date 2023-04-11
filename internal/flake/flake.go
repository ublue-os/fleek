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
	"text/template"

	"github.com/riywo/loginshell"
	app "github.com/ublue-os/fleek"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/fleek"
)

const nixbin = "nix"

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
		fin.Verbose.Println(app.Trans("flake.initializingTemplates"))
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

func (f *Flake) Update() error {
	spinner, err := fin.Spinner().Start(f.app.Trans("flake.update"))
	if err != nil {
		return err
	}
	updateCmdLine := []string{"flake", "update"}
	out, err := f.runNix(nixbin, updateCmdLine)

	if err != nil {
		return err
	}
	if f.Config.Verbose {
		if len(out) > 0 {
			fin.Verbose.Println(out)
		}
	}
	spinner.Success()
	err = f.mayCommit("fleek: update flake.lock")

	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) Create(force bool, symlink bool) error {
	fin.Info.Println(f.app.Trans("init.writingConfigs"))
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

	fin.Info.Println(f.app.Trans("init.blingLevel", f.Config.Bling))

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

	err = f.Config.WriteInitialConfig(f.Config.Force, symlink)
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

	return nil

}
func (f *Flake) IsJoin() (bool, error) {
	// if the user has a flake.nix, but no fleek.yml, then we assume they want to join
	_, err := os.Stat(filepath.Join(f.Config.UserFlakeDir(), "flake.nix"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	home, _ := os.UserHomeDir()
	_, err = os.Stat(filepath.Join(home, ".fleek.yml"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (f *Flake) Join() error {
	fin.Info.Println(f.app.Trans("init.writingConfigs"))
	err := f.ensureFlakeDir()
	if err != nil {
		return err
	}
	// Symlink the yaml file to home
	cfile, err := f.Config.Location()
	if err != nil {
		fin.Debug.Printfln("location err: %s ", err)
		return err
	}
	fin.Debug.Printfln("init cfile: %s ", cfile)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	csym := filepath.Join(home, ".fleek.yml")
	err = os.Symlink(cfile, csym)
	if err != nil {
		fin.Debug.Println("first symlink attempt failed")
		return err
	}
	err = f.ReadConfig()
	if err != nil {
		return err
	}
	err = f.Config.Validate()
	if err != nil {
		return err
	}
	fin.Debug.Println("new system")

	sys, err := fleek.NewSystem()
	if err != nil {
		return err
	}
	fin.Debug.Println("write system")
	err = f.writeSystem(*sys, true)
	if err != nil {
		return err
	}
	//
	var found bool
	for _, s := range f.Config.Systems {
		if s.Hostname == sys.Hostname && s.Username == sys.Username && s.Arch == sys.Arch {
			fin.Debug.Println("system already exists")
			found = true
		}
	}
	if !found {
		f.Config.Systems = append(f.Config.Systems, sys)

	}
	fin.Debug.Println("write config")

	err = f.Config.Save()
	if err != nil {
		fin.Debug.Println("config save failed")
		return err
	}
	git, err := f.IsGitRepo()
	if err != nil {
		return err
	}
	if git {
		fin.Warning.Println(f.app.Trans("git.warn"))
		err = f.setRebase()
		if err != nil {
			return err
		}
		err = f.mayCommit("fleek: new system")
		if err != nil {
			return err
		}
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
func (f *Flake) Write(includeSystems bool, message string) error {
	spinner, err := fin.Spinner().Start(f.app.Trans("flake.writing"))
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
	fmt.Println(bling.Name)
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
	if f.Config.Ejected {
		err = f.writeFile("README.md", data, true)
		if err != nil {
			return err
		}
	}

	spinner.Success()
	err = f.mayCommit(message)

	if err != nil {
		return err
	}
	return nil

}

func (f *Flake) ensureFlakeDir() error {
	if f.Config.Verbose {
		fin.Verbose.Println(f.app.Trans("flake.ensureDir"))
	}
	err := f.Config.MakeFlakeDir()
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil
		}
	}
	return nil
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
func (f *Flake) Apply() error {
	spinner, err := fin.Spinner().Start(f.app.Trans("flake.apply"))
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

var (
	//go:embed *.tmpl
	content embed.FS
)
