package deezer

import (
	"fmt"
	"music-server/utils"
	"os/exec"
	"path/filepath"
)

func Download(trackID int, path string, quality utils.QualityLevel) error {
	filename := filepath.Base(path)
	dir := filepath.Dir(path)

	var qualityArg string
	if quality.Bitrate == 320 {
		qualityArg = "1"
	} else if quality.Bitrate == 16 {
		qualityArg = "2"
	}

	cmd := exec.Command("rip", "-f", dir, "-o", filename, "--quality", qualityArg, "id", "deezer", "track", fmt.Sprintf("%d", trackID))
	err := cmd.Run()
	return err
}
