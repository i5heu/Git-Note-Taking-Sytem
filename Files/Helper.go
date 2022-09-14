package files

import (
	config "Tyche/Config"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func getCleanRepoPath(path string) string {
	cleanPath := filepath.Join(config.GetRepoDir(), path)

	// check if the path is inside the repo
	if !strings.HasPrefix(cleanPath, config.GetRepoDir()) {
		log.Info().Msg("Path is not inside the repo, returning repo dir")
		return config.GetRepoDir()
	}

	return cleanPath
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Error().Err(err).Msg("Error getting file info")
		return false
	}

	return fileInfo.IsDir()
}
