package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Helper contains shares helper install functions
type Helper struct {
}

// CreateInstallDirectories creates the directories needed for installation
//
// It returns the path where miners will be installed, users need to exclude
// this path from antivirus scanning.
func (helper *Helper) CreateInstallDirectories(
	installDirectory string) (string, error) {

	paths := []string{
		"miner-controller",
		filepath.Join("miner-controller", "miners"),
	}
	avExcludePath := "miners"
	for _, path := range paths {
		path = filepath.Join(installDirectory, path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return avExcludePath, fmt.Errorf(
				"Unable to create installation directory: '%s': %s",
				path,
				err)
		}
		if strings.Contains(path, "miners") {
			avExcludePath = path
		}
	}
	return avExcludePath, nil
}

// GetOSAVGuides returns a list of links and descriptions for antivirus
// directory exclude guides
func (helper *Helper) GetOSAVGuides() string {
	return fmt.Sprintf(`https://www.mininghq.io/help/antivirus`)
}
