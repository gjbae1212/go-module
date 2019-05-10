package async_task

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockTask2 struct {
	Data int
}

func (m *mockTask2) Process(ctx context.Context) error {
	time.Sleep(time.Second * 2)
	return ctx.Err()
}

func TestWorker(t *testing.T) {
	assert := assert.New(t)

	check := 0
	k, err := NewAsyncTask(WithQueueSizeOption(10),
		WithWorkerSizeOption(2),
		WithTimeoutOption(3*time.Second),
		WithErrorHandlerOption(func(err error) {
			check += 1
		}))
	assert.NoError(err)

	ctx := context.Background()
	k.(*keeper).dispatcher.stop()
	err = k.AddTask(ctx, &mockTask2{Data: 10})
	assert.NoError(err)
	time.Sleep(1 * time.Second)
	assert.Equal(1, len(k.(*keeper).dispatcher.taskQueue))

	k.(*keeper).dispatcher.start()
	err = k.AddTask(ctx, &mockTask2{Data: 10})
	assert.NoError(err)
	time.Sleep(1 * time.Second)
	assert.Equal(0, len(k.(*keeper).dispatcher.taskQueue))
	assert.Equal(0, check)

	k, err = NewAsyncTask(WithQueueSizeOption(10),
		WithWorkerSizeOption(2),
		WithTimeoutOption(1*time.Second),
		WithErrorHandlerOption(func(err error) {
			check += 1
		}))
	assert.NoError(err)

	for i := 0; i < 2; i++ {
		err = k.AddTask(ctx, &mockTask2{Data: 10})
		assert.NoError(err)
	}
	time.Sleep(5 * time.Second)
	assert.Equal(2, check)

}