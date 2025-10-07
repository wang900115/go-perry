package main

func main() {

	tasks := []Task{
		&EmailTask{Email: "email1@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "1/images/example.jpg"},
		&EmailTask{Email: "email2@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "2/images/example.jpg"},
		&EmailTask{Email: "email3@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "3/images/example.jpg"},
		&EmailTask{Email: "email4@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "4/images/example.jpg"},
		&EmailTask{Email: "email5@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "5/images/example.jpg"},
		&EmailTask{Email: "email6@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "6/images/example.jpg"},
		&EmailTask{Email: "email7@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "7/images/example.jpg"},
		&EmailTask{Email: "email8@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "8/images/example.jpg"},
		&EmailTask{Email: "email9@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "9/images/example.jpg"},
		&EmailTask{Email: "email10@example.com", Subject: "test", MessageBody: "test"},
		&ImageProcessingTask{ImageUrl: "10/images/example.jpg"},
	}

	wp := WorkerPool{
		Tasks:       tasks,
		concurrency: 5,
	}

	wp.Run()
}
