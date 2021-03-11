package lib

// Repo provides interface definition for storage adapter implementation
type Repo interface {
	// saveJob will allow to save a job in storage
	saveJob(job *Job) (*Job, error)
	// getJob will allow to get a job from storage
	getJob(ID string) (*Job, error)
	// updateJobStatus will allow to update status and exit code of the job
	updateJobStatus(ID string, status int, exitCode int) (*Job, error)
}
