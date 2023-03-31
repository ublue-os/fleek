package nix

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/ublue-os/fleek/internal/core"
)

var ErrPackageConflict = errors.New("package exists in fleek and nix profile")

type Data struct {
	Config   *core.Config
	UserName string
	Home     string
	Bling    *core.Bling
}

type Flake struct {
	RootDir   string
	Templates *template.Template
	Config    *core.Config
}

const nixbin = "nix"

func NewFlake(root string, config *core.Config) (*Flake, error) {

	t, err := template.ParseFS(content, "*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	f := &Flake{
		Templates: t,
		Config:    config,
		RootDir:   root,
	}
	return f, nil

}
func (f *Flake) runNix(cmd string, cmdLine []string) ([]byte, error) {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Dir = f.RootDir
	command.Env = os.Environ()
	if f.Config.Unfree {
		command.Env = append(command.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	return command.Output()

}

func (f *Flake) PackageIndex() ([]byte, error) {
	// nix search nixpkgs --json
	indexCmdLine := []string{"search", "nixpkgs", "--json"}
	out, err := f.runNix(nixbin, indexCmdLine)
	if err != nil {
		return out, fmt.Errorf("nix search: %w", err)
	}

	return out, nil
}

// Exist verifies that the Flake directory exists
func (f *Flake) Exists() (bool, error) {
	_, err := os.Stat(f.RootDir)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Init writes the first flake configuration
func (f *Flake) Init(force bool) error {

	err := f.Config.Validate()
	if err != nil {
		return err
	}
	var bling *core.Bling

	switch f.Config.Bling {
	case "high":
		bling, err = core.HighBling()
	case "default":
		bling, err = core.DefaultBling()
	case "low":
		bling, err = core.LowBling()
	case "none":
		bling, err = core.NoBling()
	default:
		bling, err = core.DefaultBling()
	}
	if err != nil {
		return err
	}

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
	for _, sys := range data.Config.Systems {
		err = f.writeSystem(sys, force)
		if err != nil {
			return err
		}
	}
	return f.writeFile("user.nix", data, force)

}

// Write writes the applied flake configuration
func (f *Flake) Write(includeSystems bool) error {
	var bling *core.Bling
	var err error
	switch f.Config.Bling {
	case "high":
		bling, err = core.HighBling()
	case "default":
		bling, err = core.DefaultBling()
	case "low":
		bling, err = core.LowBling()
	case "none":
		bling, err = core.NoBling()
	default:
		bling, err = core.DefaultBling()
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
			err = f.writeSystem(sys, true)
			if err != nil {
				return err
			}
		}
	}
	return f.writeFile("shell.nix", data, true)

}

func (f *Flake) Apply() ([]byte, error) {

	user, err := core.Username()

	if err != nil {
		return []byte{}, err
	}
	host, err := core.Hostname()
	if err != nil {
		return []byte{}, err
	}
	applyCmdLine := []string{"run", "--impure", "home-manager/master", "--", "-b", "bak", "switch", "--flake", ".#" + user + "@" + host}
	out, err := f.runNix(nixbin, applyCmdLine)
	if err != nil {
		if bytes.Contains(out, []byte("priority")) {
			return out, ErrPackageConflict
		}
		if bytes.Contains(out, []byte("conflict")) {
			return out, ErrPackageConflict
		}
		return out, fmt.Errorf("nix run: %w", err)
	}

	return out, nil
}
func (f *Flake) GC() ([]byte, error) {
	gc := exec.Command("nix-collect-garbage", "-d")

	gc.Dir = f.RootDir
	gc.Env = os.Environ()
	if f.Config.Unfree {
		gc.Env = append(gc.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	out, err := gc.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
func (f *Flake) Check() ([]byte, error) {
	checkCmdLine := []string{"run", "--impure", "home-manager/master", "build", "--impure", "--", "--flake", "."}
	out, err := f.runNix(nixbin, checkCmdLine)

	if err != nil {
		return out, err
	}
	return out, nil
}
func (f *Flake) Update() ([]byte, error) {
	updateCmdLine := []string{"flake", "update"}
	out, err := f.runNix(nixbin, updateCmdLine)

	if err != nil {
		return out, err
	}
	return out, nil
}
func (f *Flake) writeFile(fname string, d Data, force bool) error {

	fpath := filepath.Join(f.RootDir, fname)
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
func (f *Flake) writeSystem(sys core.System, force bool) error {

	hostPath := filepath.Join(f.RootDir, sys.Hostname)
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

var (
	//go:embed *.tmpl
	content embed.FS
)
