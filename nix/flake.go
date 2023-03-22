package nix

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ublue-os/fleek/core"
)

type Data struct {
	Config          *core.Config
	UserName        string
	Home            string
	LowPackages     []string
	DefaultPackages []string
	HighPackages    []string
	LowPrograms     []string
	DefaultPrograms []string
	HighPrograms    []string
}

type Flake struct {
	RootDir   string
	Templates *template.Template
	Config    *core.Config
}

func NewFlake(root string, config *core.Config) (*Flake, error) {
	t, err := template.ParseFS(content, "*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %s", err)
	}
	f := &Flake{
		Templates: t,
		Config:    config,
		RootDir:   root,
	}
	return f, nil

}

// Init writes the first flake configuration
func (f *Flake) Init(force bool) error {

	err := f.Config.Validate()
	if err != nil {
		return err
	}

	data := Data{
		Config:          f.Config,
		LowPackages:     core.LowPackages,
		DefaultPackages: core.DefaultPackages,
		HighPackages:    core.HighPackages,
		LowPrograms:     core.LowPrograms,
		DefaultPrograms: core.DefaultPrograms,
		HighPrograms:    core.HighPrograms,
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
func (f *Flake) Write() error {

	data := Data{
		Config:          f.Config,
		LowPackages:     core.LowPackages,
		DefaultPackages: core.DefaultPackages,
		HighPackages:    core.HighPackages,
		LowPrograms:     core.LowPrograms,
		DefaultPrograms: core.DefaultPrograms,
		HighPrograms:    core.HighPrograms,
	}
	err := f.writeFile("flake.nix", data, true)
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
	for _, sys := range data.Config.Systems {
		err = f.writeSystem(sys, true)
		if err != nil {
			return err
		}
	}
	return f.writeFile("shell.nix", data, true)

}

func (f *Flake) Apply() error {

	workdir, err := core.FlakeLocation()
	if err != nil {
		return err
	}
	user, err := core.Username()
	if err != nil {
		return err
	}
	host, err := core.Hostname()
	if err != nil {
		return err
	}
	apply := exec.Command("nix", "run", "--impure", "home-manager/master", "--", "-b", "bak", "switch", "--flake", ".#"+user+"@"+host)
	apply.Stderr = os.Stderr
	apply.Stdin = os.Stdin
	apply.Stdout = os.Stdout
	apply.Dir = workdir
	apply.Env = os.Environ()

	if f.Config.Unfree {
		apply.Env = append(apply.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	err = apply.Run()
	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) Check() error {

	workdir, err := core.FlakeLocation()
	if err != nil {
		return err
	}
	apply := exec.Command("nix", "run", "--impure", "home-manager/master", "build", "--impure", "--", "--flake", ".")
	apply.Stderr = os.Stderr
	apply.Stdin = os.Stdin
	apply.Stdout = os.Stdout
	apply.Dir = workdir
	apply.Env = os.Environ()
	if f.Config.Unfree {
		apply.Env = append(apply.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	err = apply.Run()
	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) Update() error {
	workdir, err := core.FlakeLocation()
	if err != nil {
		return err
	}
	apply := exec.Command("nix", "flake", "update")
	apply.Stderr = os.Stderr
	apply.Stdin = os.Stdin
	apply.Stdout = os.Stdout
	apply.Dir = workdir
	apply.Env = os.Environ()

	err = apply.Run()
	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) writeFile(fname string, d Data, force bool) error {
	fleekPath, err := core.FlakeLocation()
	if err != nil {
		return err
	}
	fpath := filepath.Join(fleekPath, fname)
	_, err = os.Stat(fpath)
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
	fleekPath, err := core.FlakeLocation()
	if err != nil {
		return err
	}
	hostPath := filepath.Join(fleekPath, sys.Hostname)
	err = os.MkdirAll(hostPath, 0755)
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
	if force || os.IsNotExist(err) {

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
