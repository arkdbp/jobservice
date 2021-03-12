package lib

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"syscall"
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

	jobID := job.JobID()
	osCmd := job.getCommand()
	err = osCmd.Start()
	if err != nil {
		cp.logger.Error("failed to run command with error: ", err)

		// ProcessState will not be available here so exit code will be -1 which represent wait hasn't been called yet
		_, err = cp.repo.updateJobStatus(jobID, jobError, osCmd.ProcessState.ExitCode())
		if err != nil {
			cp.logger.Error("failed to update job status with error: ", err)
		}
		return nil, err
	}

	job.SetProcess(osCmd.Process)
	go cp.handleWait(jobID, osCmd)
	cp.logger.Trace("command started with status: ", job.Status())
	return job, nil
}

// Stop will allow to stop running job
func (cp *CmdProcessor) Stop(id string) error {
	job, err := cp.repo.getJob(id)
	if err != nil {
		return err
	}

	processGroupID, err := syscall.Getpgid(job.Process().Pid)
	if err != nil {
		cp.logger.Error("failed to stop job with error: ", err)
		return err
	}

	err = syscall.Kill(-processGroupID, syscall.SIGKILL)
	if err != nil {
		cp.logger.Error("failed to stop job with error: ", err)
		return err
	}

	_, err = cp.repo.updateJobStatus(job.JobID(), manualStop, 130)
	if err != nil {
		cp.logger.Error("failed to update job status with error: ", err)
		return err
	}
	return nil
}

// GetJob will allow to retrieve job
func (cp *CmdProcessor) GetJob(id string) (*Job, error) {
	return cp.repo.getJob(id)
}

func (cp *CmdProcessor) beforeStart(request *JobRequest) (*Job, error) {
	job, err := newJob(request)
	if err != nil {
		cp.logger.Error("failed to validate request with error: ", err)
		return nil, err
	}

	jb, err := cp.repo.saveJob(job)
	if err != nil {
		cp.logger.Error("failed to save job with error: ", err)
		return nil, err
	}

	return jb, nil
}

func (cp *CmdProcessor) handleWait(jobID string, osCmd *exec.Cmd) {
	status := success
	err := osCmd.Wait()
	if err != nil {
		cp.logger.Errorf("command wait error:%+v status: %d", err, status)
		status = jobError
	}

	_, err = cp.repo.updateJobStatus(jobID, status, osCmd.ProcessState.ExitCode())
	if err != nil {
		cp.logger.Error("failed to update job status with error: ", err)
		return
	}

	cp.logger.Trace("command wait completed with status: ", status)
}
