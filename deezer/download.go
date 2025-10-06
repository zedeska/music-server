package deezer

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var downloadMutex sync.Mutex

func Download(trackID int, path string) error {
	downloadMutex.Lock()
	defer downloadMutex.Unlock()

	cmd := exec.Command("python3", "orpheus.py", "-o", path, "download", "deezer", "track", fmt.Sprintf("%d", trackID))
	homedir, _ := os.UserHomeDir()
	cmd.Dir = homedir + "/OrpheusDL"
	err := cmd.Run()
	return err
}
