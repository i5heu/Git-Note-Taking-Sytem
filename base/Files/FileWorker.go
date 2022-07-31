package files

import (
	registry "base/Registry"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type FileJob struct {
	ID                uuid.UUID
	Service           registry.Service
	GetAndLockFile    GetFile
	CreateFile        CreateFile
	TerminationChan   chan bool
	ttl               uint
	heartBeat         bool
	Description       string
	RequestingService registry.Service
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
	BackChan          chan File
	RequestingService registry.Service
	Error             string
}

func (f File) GetContent() string {
	return "LOREM IPSUM"
}

func FileWorker(jobs chan FileJob) {
	//TODO remove expired locks
	//TODO FileAuthorityToken
	var lockedFiles []LockedFile

	go heartBeat(jobs)

	for job := range jobs {
		if !job.heartBeat {
			log.Debug().Timestamp().Uint("ttl", job.ttl).Str("ID", job.ID.String()).Str("Description", job.Description).Msg("New FileWorker Job")
		}

		if job.ID == uuid.Nil {
			job.ID = uuid.New()
		}

		if job.ttl >= 100 {
			job.TerminationChan <- true
			continue
		}
		job.ttl++

		if job.GetAndLockFile.BackChan != nil {
			handleGetAndLockFile(jobs, job, &lockedFiles)
		}

		if job.heartBeat {
			removeExpiredLocks(&lockedFiles)
		}
	}
}

// delete expired locks
func removeExpiredLocks(lockedFiles *[]LockedFile) {
	for i, lockedFile := range *lockedFiles {
		if lockedFile.LockExpire.Before(time.Now()) {
			log.Debug().Str("FilePath", lockedFile.FilePath).Str("RequestingService", lockedFile.RequestingService.Name).Msg("File expired")
			*lockedFiles = append((*lockedFiles)[:i], (*lockedFiles)[i+1:]...)
		}
	}
}

func heartBeat(jobs chan FileJob) {
	tick := 0

	for range time.Tick(time.Second * 1) {
		jobs <- FileJob{
			Description: "HeartBeat FileWorker",
			heartBeat:   true,
		}
		tick++
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

		*lockedFiles = append(*lockedFiles, lockFileCreator(job, file, job.GetAndLockFile.BackChan))

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

func lockFileCreator(job FileJob, file File, backChan chan File) LockedFile {
	fmt.Println("job.RequestingService.Name", job.RequestingService.Name)

	// check if FileAuthority is granted for this service
	if !job.Service.HasFileAuthority {
		log.Debug().Str("FilePath", file.FilePath).Str("RequestingService", job.RequestingService.Name).Msg("FileAuthority not granted")
		return LockedFile{
			Error: "FileAuthority not granted",
		}
	}

	return LockedFile{
		File:              file,
		BackChan:          backChan,
		RequestingService: job.RequestingService,
	}
}
