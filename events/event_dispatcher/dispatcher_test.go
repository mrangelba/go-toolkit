package dispatcher_test

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	dispatcher "github.com/mrangelba/go-toolkit/events/event_dispatcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(ctx context.Context, payload interface{}, wg *sync.WaitGroup) {
	m.Called(payload)
	wg.Done()
}

func TestEventDispatcher_Dispatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	ctx := context.Background()
	event := "test_event"
	payload := "test_payload"

	mockHandler.On("Handle", payload).Return()

	dispatcher := dispatcher.NewEventDispatcher()
	dispatcher.Register(event, mockHandler)

	dispatcher.Dispatch(ctx, event, payload)
}

func TestEventDispatcher_Dispatch_NoHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	event := "test_event"
	payload := "test_payload"

	dispatcher := dispatcher.NewEventDispatcher()

	dispatcher.Dispatch(ctx, event, payload)
}

func TestEventDispatcher_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	event := "test_event"

	dispatcher := dispatcher.NewEventDispatcher()
	err := dispatcher.Register(event, mockHandler)

	assert.NoError(t, err)
}

func TestEventDispatcher_Register_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	event := "test_event"

	dispatcher := dispatcher.NewEventDispatcher()
	dispatcher.Register(event, mockHandler)
	err := dispatcher.Register(event, mockHandler)

	assert.Error(t, err)
}

func TestEventDispatcher_Has(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	event := "test_event"

	dispatcher := dispatcher.NewEventDispatcher()
	dispatcher.Register(event, mockHandler)

	assert.True(t, dispatcher.Has(event, mockHandler))
}

func TestEventDispatcher_Has_False(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	event := "test_event"

	dispatcher := dispatcher.NewEventDispatcher()

	assert.False(t, dispatcher.Has(event, mockHandler))
}

func TestEventDispatcher_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	event := "test_event"

	dispatcher := dispatcher.NewEventDispatcher()
	dispatcher.Register(event, mockHandler)
	dispatcher.Remove(event, mockHandler)

	assert.False(t, dispatcher.Has(event, mockHandler))
}

func TestEventDispatcher_Clear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHandler := &MockHandler{}

	event := "test_event"

	dispatcher := dispatcher.NewEventDispatcher()
	dispatcher.Register(event, mockHandler)
	dispatcher.Clear()

	assert.False(t, dispatcher.Has(event, mockHandler))
}
