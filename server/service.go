package server

import (
	"context"
	"errors"
	"github.com/arkdbp/jobservice/api"
	"github.com/arkdbp/jobservice/lib"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"sync"
)

var (
	ServiceSet = wire.NewSet(NewServer, NewServiceServer, NewLogger, wire.Bind(new(api.JobServiceServer), new(*ServiceServer)))
)

type ServiceServer struct {
	lib    *lib.CmdProcessor
	logger *logrus.Logger
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return logger
}

func NewServiceServer(lib *lib.CmdProcessor, logger *logrus.Logger) *ServiceServer {
	return &ServiceServer{lib: lib, logger: logger}
}

func (s *ServiceServer) StartJob(ctx context.Context, request *api.StartJobRequest) (*api.StartJobResponse, error) {
	s.logger.Infof("request: %+v", request)
	job, err := s.lib.Start(&lib.JobRequest{
		Path:      request.GetPath(),
		Args:      request.GetArgs(),
		Envs:      request.GetArgs(),
		Directory: request.GetDirectory(),
	})
	//s.logger.Infof("job: %+v  err: %+v", job, err)
	if err != nil {
		return nil, err
	}
	return &api.StartJobResponse{
		Job: &api.Job{
			Path:      job.Path(),
			Args:      job.Args(),
			Envs:      job.Envs(),
			Directory: job.Directory(),
			Status:    api.JobStatus(job.Status()),
			JobId:     job.JobID(),
			ExitCode:  int32(job.ExitCode()),
		},
	}, nil
}

func (s *ServiceServer) StopJob(ctx context.Context, request *api.IDRequest) (*api.StopResponse, error) {
	err := s.lib.Stop(request.GetJobId())
	if err != nil {
		return nil, err
	}
	return &api.StopResponse{}, nil
}

func (s *ServiceServer) GetJob(ctx context.Context, request *api.IDRequest) (*api.Job, error) {
	job, err := s.lib.GetJob(request.GetJobId())
	if err != nil {
		return nil, err
	}
	return &api.Job{
		Path:      job.Path(),
		Args:      job.Args(),
		Envs:      job.Envs(),
		Directory: job.Directory(),
		Status:    api.JobStatus(job.Status()),
		JobId:     job.JobID(),
		ExitCode:  int32(job.ExitCode()),
	}, nil
}

func (s *ServiceServer) StreamJobOutput(request *api.IDRequest, server api.JobService_StreamJobOutputServer) error {
	job, err := s.lib.GetJob(request.GetJobId())
	if err != nil {
		return err
	}
	if job.Status() < 1 {
		return errors.New("job is not started yet")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// reading output
		for {
			data := make([]byte, 100)
			n, err := job.Output().Read(data)
			if err != nil && job.Status() > 1 {
				if err != io.EOF {
					s.logger.Warnf("failed to read output with error: %v", err)
				}
				break
			}
			s.logger.Infof("read data: %s", data)
			if len(data) < n {
				data = data[:n]
			}

			resp := api.JobOutput{
				Output: string(data),
			}
			if err := server.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
		}
	}()
	wg.Wait()
	return nil
}
