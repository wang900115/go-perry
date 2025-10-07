package main

import (
	"context"
	"fmt"
	"log"
	"time"

	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
)

func sendEmail(ctx context.Context, args ...interface{}) error {
	help := worker.HelperFor(ctx)
	log.Printf("Working on job with ID: %s\n", help.Jid())
	addr := args[0].(string)
	subject := args[1].(string)

	fmt.Println("Sending mail to " + addr + " with subject " + subject)
	time.Sleep(time.Second * 5)
	return nil
}

func prepareReport(ctx context.Context, args ...interface{}) error {
	help := worker.HelperFor(ctx)
	log.Printf("Working on job with ID: %s\n", help.Jid())
	addr := args[0].(string)

	fmt.Println("Preparing report for the user: " + addr)
	time.Sleep(time.Second * 10)

	return help.With(func(cl *faktory.Client) error {
		job := faktory.NewJob("email", addr, "Report is ready!")
		return cl.Push(job)
	})
}

func main() {
	mgr := worker.NewManager()

	mgr.Register("email", sendEmail)
	mgr.Register("report", prepareReport)

	mgr.Concurrency = 5
	mgr.ShutdownTimeout = 25 * time.Second

	mgr.ProcessStrictPriorityQueues("critical", "default", "low_priority")

	mgr.Run()
}
