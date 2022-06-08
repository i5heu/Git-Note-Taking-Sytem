package files

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type Config struct {
	GitSsh string `json:"gitSsh"`
}

func GetConfig() Config {
	jsonFile, err := os.Open(GetWorkDir() + "/config.json")
	if err != nil {
		log.Error().Err(err).Msg("open config.json failed")
		panic(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	defer jsonFile.Close()

	return config
}

func GetWorkDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("Error getting user home dir")
		panic(err)
	}
	return dirname + "/.Tyche"
}

type FileJob struct {
	GetAndLockFile  GetFile
	CreateFile      CreateFile
	TerminationChan chan bool
	ttl             uint
}

type File struct {
	FilePath   string
	LockStatus string // locked or expired
	LockExpire time.Time
}

// locking a file and return the content
type GetFile struct {
	FilePath string
	Lock     bool
	BackChan chan File
}

// not locking the file
type CreateFile struct {
	FilePath string
	Content  string
	BackChan chan File
}

type LockedFile struct {
	File
	BackChan chan File
}

func (f File) GetContent() string {
	return "f.Content"
}

func FileWorker(jobs chan FileJob) {
	//TODO remove expired locks
	var lockedFiles []LockedFile

	for job := range jobs {

		if job.ttl > 15 {
			job.TerminationChan <- true
			continue
		}
		job.ttl++

		if job.GetAndLockFile.BackChan != nil {
			handleGetAndLockFile(jobs, job, &lockedFiles)
		}
	}
}

func handleGetAndLockFile(jobs chan FileJob, job FileJob, lockedFiles *[]LockedFile) {

	if checkIfFileIsLocked(job.GetAndLockFile.FilePath, *lockedFiles) {
		jobs <- job
		return
	}

	var file File
	if job.GetAndLockFile.Lock {
		file = File{
			FilePath:   job.GetAndLockFile.FilePath,
			LockStatus: "locked",
			LockExpire: time.Now().Add(time.Second * 15),
		}

		*lockedFiles = append(*lockedFiles, lockFileCreator(file, job.GetAndLockFile.BackChan))

	} else {
		file = File{
			FilePath:   job.GetAndLockFile.FilePath,
			LockStatus: "expired",
			LockExpire: time.Now(),
		}
	}

	job.GetAndLockFile.BackChan <- file
}

func checkIfFileIsLocked(filePath string, lockedFiles []LockedFile) bool {
	for _, lockedFile := range lockedFiles {
		if lockedFile.FilePath == filePath {
			return true
		}
	}
	return false
}

func lockFileCreator(file File, backChan chan File) LockedFile {
	return LockedFile{
		File:     file,
		BackChan: backChan,
	}
}
