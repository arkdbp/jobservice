package lib

import (
	"errors"
	"github.com/google/uuid"
	"os/exec"
	"sync"
)

// JobRequest will allow to create start job request
type JobRequest struct {
	Path      string
	Args      []string
	Envs      []string
	Directory string
}

type requestValidator func(job *Job) error

func requestJobValidator(request *JobRequest) requestValidator {
	return func(job *Job) error {
		if request.Path == "" {
			return errors.New("request path is required")
		}

		path, err := exec.LookPath(request.Path)
		if err != nil {
			return err
		}

		job.path = path
		job.args = request.Args
		job.envs = request.Envs
		job.directory = request.Directory

		// setting default values for new job
		job.status = started
		job.jobID = uuid.New().String()
		job.output = NewJobBuffer()
		job.error = NewJobBuffer()
		job.lock = &sync.Mutex{}
		return nil
	}
}
