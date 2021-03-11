package lib

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// CmdProcessor will allow to start/stop
type CmdProcessor struct {
	logger *logrus.Logger
	repo   Repo
}

// NewCmdProcessor will allow to create instance of CmdProcessor
func NewCmdProcessor(logger *logrus.Logger, repo Repo) *CmdProcessor {
	return &CmdProcessor{logger: logger, repo: repo}
}

// Start will allow to start the execution of job
func (cp *CmdProcessor) Start(request *JobRequest) (*Job, error) {
	job, err := cp.beforeStart(request)
	if err != nil {
		return nil, err
	}

	jobID := job.jobID
	osCmd := job.getCommand()
	err = osCmd.Start()
	if err != nil {
		cp.logger.Error("failed to run command with error: ", err)

		// ProcessState will not be available here so exit code will be -1 which represent wait hasn't been called yet
		_, _ = cp.repo.updateJobStatus(jobID, jobError, osCmd.ProcessState.ExitCode())
		return nil, err
	}

	cp.updateJobProcess(osCmd.Process, job)
	cp.handleWait(jobID, osCmd)
	cp.logger.Trace("command started with status: ", job.Status())
	return job, nil
}

func (cp *CmdProcessor) beforeStart(request *JobRequest) (*Job, error) {
	var job Job

	validator := requestJobValidator(request)
	err := validator(&job)
	if err != nil {
		cp.logger.Error("failed to validate request with error: ",err)
		return nil, err
	}

	jb, err := cp.repo.saveJob(&job)
	if err != nil {
		cp.logger.Error("failed to save job with error: ",err)
		return nil, err
	}

	return jb, nil
}

func (cp *CmdProcessor) updateJobProcess(process *os.Process, job *Job) {
	job.SetProcess(process)
	_, _ = cp.repo.saveJob(job)
}

func (cp *CmdProcessor) handleWait(jobID string, osCmd *exec.Cmd) {
	go func() {
		status := success
		err := osCmd.Wait()
		if err != nil {
			cp.logger.Errorf("command wait error:%+v status: %d", err, status)
			status = jobError
		}

		_, _ = cp.repo.updateJobStatus(jobID, status, osCmd.ProcessState.ExitCode())
		cp.logger.Trace("command wait completed with status: ", status)
	}()
}

// Stop will allow to stop running job
func (cp *CmdProcessor) Stop(id string) error {
	job, err := cp.repo.getJob(id)
	if err != nil {
		return err
	}

	err = job.Process().Kill()
	if err != nil {
		cp.logger.Error("failed to stop job with error: ", err)
		return err
	}
	_, _ = cp.repo.updateJobStatus(job.JobID(), manualStop, undefinedExitCode)
	return nil
}

// GetJob will allow to retrieve job
func (cp *CmdProcessor) GetJob(id string) (*Job, error) {
	return cp.repo.getJob(id)
}
