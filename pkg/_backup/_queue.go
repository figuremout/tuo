package command

type IQueue interface {
	Enqueue(item Item)
	Dequeue() Item
	IsEmpty() bool
	Size() int
	Flush()
}
type Item interface {
}
type Queue struct {
	items []Item
}

var _ IQueue = &Queue{} // allow IDE to check if interface is realized

func NewQueue() *Queue {
	return &Queue{
		items: []Item{},
	}
}

func (q *Queue) Enqueue(item Item) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() Item {
	length := len(q.items)
	if length == 0 {
		return nil
	}
	item := q.items[0]
	q.items = q.items[1:length]

	return item
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue) Size() int {
	return len(q.items)
}

func (q *Queue) Flush() {
	q.items = []Item{}
}
