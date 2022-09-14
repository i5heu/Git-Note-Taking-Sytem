package gitManager

import (
	"os"

	configHelper "Tyche/Config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/rs/zerolog/log"
)

var gitInstance *git.Repository

func GitManager(config configHelper.Config) {

	workPath := configHelper.GetWorkDir()
	sshKey := getSSHKey(workPath)

	var err error
	gitInstance, err = git.PlainOpen(workPath + "/repo")
	if err != nil {
		log.Warn().Err(err).Msg("open git repo failed")
		log.Info().Msg("clone repo...")
		gitInstance = cloneRepoIfNotExists(workPath, config, sshKey)
	}

	Pull()
}

func Pull() {
	log.Info().Msg("Pulling...")
	wt, err2 := gitInstance.Worktree()
	if err2 != nil {
		log.Error().Err(err2).Msg("get worktree failed")
		panic(err2)
	}
	wt.Pull(&git.PullOptions{RemoteName: "origin"})
	log.Info().Msg("Pulled")
}

func cloneRepoIfNotExists(workPath string, config configHelper.Config, sshKey *ssh.PublicKeys) *git.Repository {
	giti, err := git.PlainClone(workPath+"/repo", false, &git.CloneOptions{
		Auth:     sshKey,
		URL:      config.GitSSH,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Error().Err(err).Msg("Error cloning")
	}

	return giti
}

func getSSHKey(workPath string) *ssh.PublicKeys {
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
