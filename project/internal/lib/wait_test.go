package lib

import (
	"errors"
	"syscall"
	"testing"
	"time"
)

func TestWithCriticalError(t *testing.T) {
	// test when criticalError is not nil
	err := errors.New("critical error occurred")
	WithCriticalError(err)

	select {
	case e := <-errorsChan:
		if e != err {
			t.Errorf("Expected error %v but got %v", err, e)
		}
	default:
		t.Errorf("Expected error %v but none was received", err)
	}

	// test when criticalError is nil
	WithCriticalError(nil)

	select {
	case e := <-errorsChan:
		t.Errorf("Expected no error but got %v", e)
	default:
		// expected behavior, do nothing
	}
}

func TestWaitStatus(t *testing.T) {
	// test when SIGINT is received
	go func() {
		time.Sleep(100 * time.Millisecond)
		sigsChan <- syscall.SIGINT
	}()

	result := WaitStatus()
	if result != 0 {
		t.Errorf("Expected result 0 but got %v", result)
	}

	// test when an error is received
	err := errors.New("an error occurred")
	go func() {
		time.Sleep(100 * time.Millisecond)
		errorsChan <- err
	}()

	result = WaitStatus()
	if result != 1 {
		t.Errorf("Expected result 1 but got %v", result)
	}
}
