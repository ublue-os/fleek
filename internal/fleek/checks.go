package fleek

import (
	"os"
	"os/exec"
)

// CheckNix verifies that the nix
// command is available in user's PATH
func CheckNix() bool {
	_, err := exec.LookPath("nix")
	return err == nil
}

func SSHAuthSock() bool {
	sock := os.Getenv("SSH_AUTH_SOCK")
	return sock != ""

}
