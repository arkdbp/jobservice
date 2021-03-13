package lib

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

// MemRepo memory storage implementation
type MemRepo struct {
	jobStorage map[string]*Job
	lock       sync.RWMutex
	logger     *logrus.Logger
}

// NewMemRepo will allow to create an instance of MemRepo
func NewMemRepo() *MemRepo {
	return &MemRepo{jobStorage: make(map[string]*Job), lock: sync.RWMutex{}, logger: logrus.New()}
}

func (m *MemRepo) saveJob(job *Job) (*Job, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.jobStorage[job.jobID] = job
	return job, nil
}

func (m *MemRepo) getJob(ID string) (*Job, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	job, ok := m.jobStorage[ID]
	if !ok {
		return nil, errors.New("job not available")
	}
	return job, nil
}

func (m *MemRepo) updateJobStatus(ID string, status int, exitCode int) (*Job, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	job, ok := m.jobStorage[ID]
	if !ok {
		return nil, errors.New("job not available")
	}

	job.SetStatus(status)
	job.SetExitCode(exitCode)
	m.jobStorage[job.jobID] = job
	return job, nil
}
