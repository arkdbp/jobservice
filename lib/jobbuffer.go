package lib

import (
	"bytes"
	"sync"
)

// JobBuffer will allow to use bytes.Buffer for concurrent use
type JobBuffer struct {
	bytesBuffer bytes.Buffer
	lock        sync.RWMutex
}

// NewJobBuffer will allow to create instance of JobBuffer
func NewJobBuffer() *JobBuffer {
	return &JobBuffer{bytesBuffer: bytes.Buffer{}, lock: sync.RWMutex{}}
}

// Write will allow to concurrent write
func (jobBuffer *JobBuffer) Write(p []byte) (n int, err error) {
	jobBuffer.lock.Lock()
	defer jobBuffer.lock.Unlock()
	return jobBuffer.bytesBuffer.Write(p)
}

// Read will allow to concurrent read
func (jobBuffer *JobBuffer) Read(p []byte) (n int, err error) {
	jobBuffer.lock.RLock()
	defer jobBuffer.lock.RUnlock()
	return jobBuffer.bytesBuffer.Read(p)
}

// Len will allow to check len
func (jobBuffer *JobBuffer) Len() int {
	jobBuffer.lock.RLock()
	defer jobBuffer.lock.RUnlock()
	return jobBuffer.bytesBuffer.Len()
}

// String will allow to read buffer as string
func (jobBuffer *JobBuffer) String() string {
	jobBuffer.lock.RLock()
	defer jobBuffer.lock.RUnlock()
	return jobBuffer.bytesBuffer.String()
}
