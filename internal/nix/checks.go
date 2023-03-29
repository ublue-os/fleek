package nix

import (
	"os/exec"
)

// CheckNix verifies that the nix
// command is available in user's PATH
func CheckNix() bool {
	_, err := exec.LookPath("nix")
	return err == nil
}
