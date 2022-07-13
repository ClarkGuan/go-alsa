package alsa

import (
	"testing"
)

func TestConfigTopDir(t *testing.T) {
	topDir := ConfigTopDir()
	if topDir != "/usr/share/alsa" {
		t.Fatal("ConfigTopDir() return", topDir)
	}
}
