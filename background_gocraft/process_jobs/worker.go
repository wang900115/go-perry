package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

type User struct {
	ID    int64
	Email string
	Name  string
}

type Context struct {
	currentUser *User
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting a new Job:", job.Name, "with ID:", job.ID)
	return next()
}

func (c *Context) FindCurrentUser(job *work.Job, next work.NextMiddlewareFunc) error {
	if _, ok := job.Args["userID"]; ok {
		userID := job.ArgInt64("userID")
		c.currentUser = &User{ID: userID, Email: "test" + strconv.Itoa(int(userID)) + "@test.com", Name: "Test User"}
		if err := job.ArgError(); err != nil {
			return err
		}
	}

	return next()
}

var enqueuer = work.NewEnqueuer("demo_app", redisPool)

func main() {
	pool := work.NewWorkerPool(Context{}, 10, "demo_app", redisPool)

	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCurrentUser)

	pool.JobWithOptions("email", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).sendEmail)
	pool.JobWithOptions("report", work.JobOptions{Priority: 10, MaxFails: 1}, (*Context).report)
	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	pool.Stop()
}

func (c *Context) sendEmail(job *work.Job) error {
	addr := c.currentUser.Email

	subject := job.ArgString("subject")

	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Sending email to " + addr + " with subject " + subject)
	time.Sleep(time.Second * 2)
	return nil
}

func (c *Context) report(job *work.Job) error {
	fmt.Println("Preparing report ...")
	for i := range 360 {
		time.Sleep(time.Second * 10)
		job.Checkin("i = " + fmt.Sprint(i))
	}
	enqueuer.Enqueue("email", work.Q{"userID": c.currentUser.ID, "subject": "Report is ready"})
	return nil
}
