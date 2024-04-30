package application_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/stretchr/testify/assert"
)

func TestWait(t *testing.T) {
	time.AfterFunc(1*time.Second, func() {
		application.PANIC <- fmt.Errorf("test")
	})

	err := application.Wait()
	assert.Error(t, err)

	time.AfterFunc(1*time.Second, func() {
		application.SIGS <- nil
	})

	err = application.Wait()
	assert.Nil(t, err)
}
