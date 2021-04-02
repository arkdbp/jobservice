package main

import (
	"context"
	"github.com/arkdbp/jobservice/api"
	c "github.com/arkdbp/jobservice/client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	config := c.ProvideConfig()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(config.TLSConfig())),
	}

	conn, err := grpc.Dial("localhost:12000", opts...)
	if err != nil {
		grpclog.Fatal("failed to dial OAuth: ", err)
	}
	defer conn.Close()
	client := api.NewJobServiceClient(conn)
	job, err := client.StartJob(ctx, &api.StartJobRequest{
		Path:      "ls",
		Args:      []string{"-l", "-a"},
		Directory: "/home/dpanchal/Documents",
	})

	if err != nil {
		log.Println(err)
	}

	stream, err := client.StreamJobOutput(ctx, &api.IDRequest{JobId: job.Job.GetJobId()})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// reading output
		for {
			jobOutput, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("getting error ", err)
				break
			}

			log.Println("read perform data: ", jobOutput)
		}
	}()
	wg.Wait()
}
