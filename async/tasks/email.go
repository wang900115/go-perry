package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

type EmailPayload struct {
	EmailAddr string
	Subject   string
	Body      string
}

func NewEmailTask(emailAddr string, subject string, body string) (*asynq.Task, error) {

	payload, err := json.Marshal(EmailPayload{
		EmailAddr: emailAddr,
		Subject:   subject,
		Body:      body,
	})

	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeEmail, payload)
	return task, nil
}

func EmailTaskHandler(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Umarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Sending Email to User: email=%s, subject=%s", p.EmailAddr, p.Subject)
	return nil
}
