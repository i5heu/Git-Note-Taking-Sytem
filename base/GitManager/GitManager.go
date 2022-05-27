package gitManager

import (
	"fmt"
	"os"

	channels "base/Channels"
	files "base/Files"
	queue "base/Queue"
	registry "base/Registry"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var gitInstance *git.Repository

func GitManager() {

	workPath := files.GetWorkDir()
	config := files.GetConfig()
	sshKey := getSshKey(workPath)

	var err error
	gitInstance, err = git.PlainOpen(workPath + "/repo")
	if err != nil {
		log.Warn().Err(err).Msg("open git repo failed")
		log.Info().Msg("clone repo...")
		gitInstance = cloneRepoIfNotExists(workPath, config, sshKey)
	}

	ref, err := gitInstance.Head()
	if err != nil {
		log.Error().Err(err).Msg("get git head failed")
		panic(err)
	}

	fmt.Println(ref.Hash())
}

func Pull() {

	job := queue.QueueJob{
		QueueItem: queue.QueueItem{
			Service: registry.Service{
				Name: "PULL",
				Id:   uuid.New(),
			},
			Priority:      100,
			RunAfterwards: []registry.Service{},
			SlotOpen:      make(chan bool),
		},
	}

	channels.QueueChan <- job
	<-job.QueueItem.SlotOpen
	// okay lets Pull!

	wt, err2 := gitInstance.Worktree()
	if err2 != nil {
		log.Error().Err(err2).Msg("get worktree failed")
		panic(err2)
	}
	wt.Pull(&git.PullOptions{RemoteName: "origin"})

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
