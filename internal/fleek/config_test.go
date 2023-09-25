package fleek

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ublue-os/fleek/internal/xdg"
)

func TestHostname(t *testing.T) {
	h, err := Hostname()
	if err != nil {
		t.Error(err)
	}
	overrideHost := "fleekhost"
	os.Setenv("FLEEK_HOST_OVERRIDE", overrideHost)
	hOverride, err := Hostname()
	if err != nil {
		t.Error(err)
	}
	if h == hOverride {
		t.Fatalf("host override: expected %s got %s", overrideHost, hOverride)
	}
}

func TestUniqueSystems(t *testing.T) {
	c := &Config{
		Systems: []*System{
			{
				Arch: "x86-64",
				OS:   "darwin",
			},
		},
	}
	want := "x86-64-darwin"
	got := c.UniqueSystems()
	if got[0] != want {
		t.Fatalf("unique systems: expected %s got %s", want, got)
	}
	c = &Config{
		Systems: []*System{
			{
				Arch: "x86-64",
				OS:   "darwin",
			},
			{
				Arch: "aarch64",
				OS:   "linux",
			},
		},
	}
	want = "aarch64-linux"
	got = c.UniqueSystems()
	if got[0] != want {
		t.Fatalf("unique systems: expected %s got %s", want, got)
	}
	c = &Config{
		Systems: []*System{
			{
				Hostname: "a",
				Arch:     "x86-64",
				OS:       "darwin",
			},
			{
				Hostname: "b",
				Arch:     "aarch64",
				OS:       "linux",
			},
			{
				Hostname: "c",
				Arch:     "aarch64",
				OS:       "darwin",
			},
			{
				Hostname: "d",
				Arch:     "aarch64",
				OS:       "darwin",
			},
		},
	}
	wantCount := 3
	got = c.UniqueSystems()
	if len(got) != wantCount {
		t.Fatalf("unique systems count: expected %d got %d", wantCount, len(got))
	}
}
func TestUserFlakeDir(t *testing.T) {
	c := &Config{}
	home, _ := os.UserHomeDir()

	blankFlakeDir := c.UserFlakeDir()
	want := filepath.Join(home, xdg.DataSubpathRel("fleek"))

	if blankFlakeDir != want {
		t.Fatalf("user flake dir: expected %s, got %s", want, blankFlakeDir)
	}
	want = filepath.Join(home, ".config", "fleek", "flake")

	c = &Config{
		FlakeDir: ".config/fleek/flake",
	}
	manualFlakeDir := c.UserFlakeDir()
	if manualFlakeDir != want {
		t.Fatalf("manual user flake dir: expected %s, got %s", want, manualFlakeDir)
	}

}
