package gitManager

import (
	"fmt"
	"os"

	files "base/Files"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
)

func GitManager() {

	workPath := files.GetWorkDir()
	config := files.GetConfig()
	sshKey := getSshKey(workPath)

	r, err := git.PlainOpen(workPath + "/repo")
	if err != nil {
		log.Warn().Err(err).Msg("open git repo failed")
		log.Info().Msg("clone repo...")
		r = cloneRepoIfNotExists(workPath, config, sshKey)
	}

	ref, err := r.Head()
	if err != nil {
		log.Error().Err(err).Msg("get git head failed")
		panic(err)
	}
	
	fmt.Println(ref.Hash())
}

func cloneRepoIfNotExists(workPath string, config files.Config, sshKey *ssh.PublicKeys) *git.Repository {
	giti, err := git.PlainClone(workPath+"/repo", false, &git.CloneOptions{
		Auth:     sshKey,
		URL:      config.GitSsh,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Error().Err(err).Msg("Error cloning")
	}

	return giti
}

func getSshKey(workPath string) *ssh.PublicKeys {
	path := workPath + "/ssh-key/ed25519"
	_, err := os.Stat(path)
	if err != nil {
		log.Error().Err(err).Msg("read file privateKeyFile failed")
		panic(err)
	}

	publicKeys, err := ssh.NewPublicKeysFromFile("git", path, "")
	if err != nil {
		log.Error().Err(err).Msg("generate publickeys failed")
		panic(err)
	}

	return publicKeys
}
