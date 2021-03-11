package lib

import (
	"github.com/go-playground/assert/v2"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"testing"
	"time"
)

func Test_cmdProcessor_StartGetJob(t *testing.T) {
	logger := logrus.New()
	type fields struct {
		logger *logrus.Logger
		repo   Repo
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				logger: logger,
				repo:   NewMemRepo(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &CmdProcessor{
				logger: tt.fields.logger,
				repo:   tt.fields.repo,
			}
			var wg sync.WaitGroup
			jb, err := cp.Start(&JobRequest{
				Path:      "ls",
				Args:      []string{"-l", "-a"},
				Directory: "/home"})
			if err != nil {
				wg.Done()
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				// getting job
				for {
					select {
					case <-time.After(time.Second * 120):
						break
					default:
					}
					jb, err = cp.GetJob(jb.JobID())
					if err != nil {
						t.Error("failed to get logs,", err)
					}

					if jb != nil && jb.Status() > 1 {
						break
					}
				}

				// reading output
				for {
					data := make([]byte, 1000)
					n, err := jb.Output().Read(data)
					if err != nil && err == io.EOF && jb.Status() > 1 {
						t.Log("failed to read output with error: ", err)
						break
					}
					t.Logf("read perform: %d size of data data: %s", n, data)
				}

				//reading error
				for {
					data := make([]byte, 10)
					n, err := jb.Error().Read(data)
					if err != nil && err == io.EOF {
						break
					}
					t.Logf("read perform: %d size of data data: %s", n, data)
				}

			}()
			wg.Wait()
		})
	}
}

func Test_cmdProcessor_StartStopJob(t *testing.T) {
	logger := logrus.New()
	type fields struct {
		logger *logrus.Logger
		repo   Repo
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				logger: logger,
				repo:   NewMemRepo(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := &CmdProcessor{
				logger: tt.fields.logger,
				repo:   tt.fields.repo,
			}

			jb, err := cp.Start(&JobRequest{
				Path:      "ls",
				Args:      []string{"-l", "-a"},
				Directory: "/home"})
			if err != nil {
				t.Error(err)
			}

			err = cp.Stop(jb.JobID())
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, jb.Status(), manualStop)
		})
	}
}
