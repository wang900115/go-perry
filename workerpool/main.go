package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/Jeffail/tunny"
)

func SendEmail(email string, subject string, body string) {
	fmt.Printf("Sending email to %s\n", email)
	fmt.Printf("Subject: %s\n Body: %s\n", subject, body)
	time.Sleep(2 * time.Second) // Simulate time taken to send email
}

func main() {
	numCPUs := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n\n", numCPUs)

	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		m, ok := payload.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unable to extract payload")
		}

		email, ok := m["email"].(string)
		if !ok {
			return fmt.Errorf("Email not found in payload")
		}

		subject, ok := m["subject"].(string)
		if !ok {
			return fmt.Errorf("Subject not found in payload")
		}
		body, ok := m["body"].(string)
		if !ok {
			return fmt.Errorf("Body not found in payload")
		}

		SendEmail(email, subject, body)
		return nil
	})

	defer pool.Close()

	for i := 0; i < 100; i++ {
		var data interface{} = map[string]interface{}{
			"email":   "email" + strconv.Itoa(i+1) + "@example.com",
			"subject": "Welcome",
			"body":    "Thank you for signing up",
		}

		go func() {
			result := pool.Process(data)
			if result == nil {
				fmt.Println("Mail sent!")
			}
		}()
	}

	for {
		qLen := pool.QueueLength()
		fmt.Printf("----------------------------------------------------Queue Length: %d\n", qLen)
		if qLen == 0 {
			break
		}

		time.Sleep(time.Second)

	}

	time.Sleep(3 * time.Second)

}
