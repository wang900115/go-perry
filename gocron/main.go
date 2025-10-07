package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func main() {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	j, err := s.NewJob(
		gocron.DurationJob(
			30*time.Second,
		),

		gocron.NewTask(
			func(a string) {
				log.Println(a)
			},
			"Every 30 seconds",
		),

		gocron.WithName("job: Every 30 seconds"),
		gocron.WithEventListeners(
			gocron.BeforeJobRuns(
				func(jobID uuid.UUID, jobName string) {
					log.Printf("Job starting: %s, %s\n", jobID, jobName)
				},
			),
			gocron.AfterJobRuns(
				func(jobID uuid.UUID, jobName string) {
					log.Printf("Job completed: %s, %s", jobID, jobName)
				},
			),
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					log.Printf("Job has an error: %s, %s\n", jobID, jobName)
				},
			),
		),
	)

	log.Println(j.ID())

	s.NewJob(
		gocron.CronJob(
			"*/10 * * * *",
			false,
		),

		gocron.NewTask(
			func(a string) {
				log.Println(a)
			},
			"Cronjob: Every 10 mins",
		),

		gocron.WithName("Cronjob: Every 10 mins"),
	)

	s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(18, 07, 00),
				gocron.NewAtTime(05, 30, 00),
			),
		),

		gocron.NewTask(
			func(a string, b string) {
				log.Println(a, b)
			},
			"Dailyjob", " Runs everyday",
		),

		gocron.WithName("Dailyjob"),
	)

	s.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("\nInterrupt signal received. Exiting...")
		_ = s.Shutdown()
		os.Exit(0)
	}()

	for {

	}

}
