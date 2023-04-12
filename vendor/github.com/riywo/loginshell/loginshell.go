package loginshell

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

func Shell() (string, error) {
	switch runtime.GOOS {
	case "plan9":
		return Plan9Shell()
	case "linux":
		return NixShell()
	case "openbsd":
		return NixShell()
	case "freebsd":
		return NixShell()
	case "android":
		return AndroidShell()
	case "darwin":
		return DarwinShell()
	case "windows":
		return WindowsShell()
	}

	return "", errors.New("Undefined GOOS: " + runtime.GOOS)
}

func Plan9Shell() (string, error) {
	if _, err := os.Stat("/dev/osversion"); err != nil {
		if os.IsNotExist(err) {
			return "", err
		} else {
			return "", errors.New("/dev/osversion check failed")
		}
	}

	return "/bin/rc", nil
}

func NixShell() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	out, err := exec.Command("getent", "passwd", user.Uid).Output()
	if err != nil {
		return "", err
	}

	ent := strings.Split(strings.TrimSuffix(string(out), "\n"), ":")
	return ent[6], nil
}

func AndroidShell() (string, error) {
	shell := os.Getenv("SHELL");
	if shell == "" {
		return "", errors.New("SHELL not defined in android.")
	}
	return shell, nil
}

func DarwinShell() (string, error) {
	dir := "Local/Default/Users/" + os.Getenv("USER")
	out, err := exec.Command("dscl", "localhost", "-read", dir, "UserShell").Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("UserShell: (/[^ ]+)\n")
	matched := re.FindStringSubmatch(string(out))
	shell := matched[1]
	if shell == "" {
		return "", errors.New(fmt.Sprintf("Invalid output: %s", string(out)))
	}

	return shell, nil
}

func WindowsShell() (string, error) {
	consoleApp := os.Getenv("COMSPEC")
	if consoleApp == "" {
		consoleApp = "cmd.exe"
	}

	return consoleApp, nil
}
