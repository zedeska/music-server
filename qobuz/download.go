package qobuz

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Download(trackID int, path string) error {
	filename := filepath.Base(path)
	dir := filepath.Dir(path)

	cmd := exec.Command("rip", "-f", dir, "-o", filename, "id", "qobuz", "track", fmt.Sprintf("%d", trackID))
	err := cmd.Run()
	return err
}
