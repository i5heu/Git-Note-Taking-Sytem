package files

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type FileJob struct {
	ID              uuid.UUID
	GetAndLockFile  GetFile
	CreateFile      CreateFile
	TerminationChan chan bool
	ttl             uint
}

type File struct {
	FilePath   string
	LockStatus string // locked or expired
	LockExpire time.Time
	Content    string
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
	return "LOREM IPSUM"
}

func FileWorker(jobs chan FileJob) {
	//TODO remove expired locks
	var lockedFiles []LockedFile

	for job := range jobs {
		log.Debug().Timestamp().Uint("ttl", job.ttl).Str("ID", job.ID.String()).Msg("New FileWorker Job")

		if job.ID == uuid.Nil {
			job.ID = uuid.New()
		}

		if job.ttl > 100 {
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
		go func() {
			time.Sleep(time.Millisecond * 15)
			jobs <- job
		}()
		return
	}

	var file File
	if job.GetAndLockFile.Lock {
		file = File{
			FilePath:   job.GetAndLockFile.FilePath,
			LockStatus: "locked",
			LockExpire: time.Now().Add(time.Second * 15),
		}

		fmt.Println("lockedFiles", lockedFiles)
		*lockedFiles = append(*lockedFiles, lockFileCreator(file, job.GetAndLockFile.BackChan))

	} else {
		file = File{
			FilePath:   job.GetAndLockFile.FilePath,
			LockStatus: "expired",
			LockExpire: time.Now(),
		}
	}

	job.TerminationChan <- false
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
