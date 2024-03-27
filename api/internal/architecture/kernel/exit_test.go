package kernel_test

import (
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/kodmain/thetiptop/api/internal/architecture/kernel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var err = fmt.Errorf("test error")

type MockExiter struct {
	mock.Mock
	LastExitCode int
}

func (m *MockExiter) Exit(code int) {
	m.LastExitCode = code
	m.Called(code)
}

func TestWaitSigKill(t *testing.T) {
	mockExiter := new(MockExiter)
	mockExiter.On("Exit", mock.Anything)

	go kernel.Wait(mockExiter)

	kernel.SIGS <- syscall.SIGKILL
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 0, mockExiter.LastExitCode)
}

func TestWaitPanic(t *testing.T) {
	mockExiter := new(MockExiter)
	mockExiter.On("Exit", mock.Anything)

	go kernel.Wait(mockExiter)
	time.Sleep(100 * time.Millisecond)
	kernel.PANIC <- err
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, mockExiter.LastExitCode)
}
