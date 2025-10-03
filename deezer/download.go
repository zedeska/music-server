package deezer

import (
	"fmt"
	"os"
	"os/exec"
)

func Download(trackID int, path string) error {
	cmd := exec.Command("python3", "orpheus.py", "-o", path, "download", "deezer", "track", fmt.Sprintf("%d", trackID))
	homedir, _ := os.UserHomeDir()
	cmd.Dir = homedir + "/OrpheusDL"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}
