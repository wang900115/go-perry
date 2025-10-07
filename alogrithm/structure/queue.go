package queue

type Queue struct {
	data []int
}

func (q *Queue) Enqueue(val int) {
	q.data = append(q.data, val)
}

func (q *Queue) Dequeue() (int, bool) {
	if len(q.data) == 0 {
		return 0, false
	}
	val := q.data[0]
	q.data = q.data[1:]
	return val, true
}

func (q *Queue) Peek() (int, bool) {
	if len(q.data) == 0 {
		return 0, false
	}
	return q.data[0], true
}

func (q *Queue) Isempty() bool {
	return len(q.data) == 0
}

func (q *Queue) Size() int {
	return len(q.data)
}
