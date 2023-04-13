package domains

import "github.com/sashabaranov/go-openai"

type ChatCompletionMessageQueue struct {
	queue   []openai.ChatCompletionMessage
	maxSize int
}

func NewFixedQueue(maxSize int) *ChatCompletionMessageQueue {
	return &ChatCompletionMessageQueue{
		queue:   make([]openai.ChatCompletionMessage, 0, maxSize),
		maxSize: maxSize,
	}
}

func (q *ChatCompletionMessageQueue) Push(v openai.ChatCompletionMessage) {
	if len(q.queue) == q.maxSize {
		copy(q.queue[0:], q.queue[1:])
	}
	q.queue = append(q.queue, v)
}

func (q *ChatCompletionMessageQueue) Len() int {
	return len(q.queue)
}

func (q *ChatCompletionMessageQueue) Cap() int {
	return q.maxSize
}

func (q *ChatCompletionMessageQueue) Get(i int) openai.ChatCompletionMessage {
	return q.queue[i]
}
