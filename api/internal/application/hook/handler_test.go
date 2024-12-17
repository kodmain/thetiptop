package hook_test

import (
	"testing"
	"time"

	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/hook"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	env.ForceTest()
	var ENAME hook.Event = "EVENT"

	t.Run("Test normal and once handlers", func(t *testing.T) {
		var countHandler int = 0
		var countOnceHandler int = 0

		var handler hook.Handler = func(tags ...string) {
			countHandler++
		}

		var onceHandler hook.OnceHandler = func(tags ...string) {
			countOnceHandler++
		}

		hook.Register(ENAME, handler)
		hook.Register(ENAME, onceHandler)

		hook.Call(ENAME)
		hook.Call(ENAME)

		time.Sleep(1000 * time.Millisecond)
		assert.Equal(t, 2, countHandler, "Le handler normal doit être appelé deux fois")
		assert.Equal(t, 1, countOnceHandler, "Le once handler doit être appelé une seule fois")
	})

	t.Run("Test sync and once sync handlers", func(t *testing.T) {
		var countHandler int = 0
		var countOnceHandler int = 0
		var handlerSync hook.HandlerSync = func(tags ...string) {
			countHandler++
		}

		var onceHandlerSync hook.OnceHandlerSync = func(tags ...string) {
			countOnceHandler++
		}

		hook.Register(ENAME, handlerSync)
		hook.Register(ENAME, onceHandlerSync)

		hook.Call(ENAME)
		hook.Call(ENAME)
		time.Sleep(1000 * time.Millisecond)
		assert.Equal(t, 4, countHandler, "Le handler sync doit être appelé deux fois")
		assert.Equal(t, 1, countOnceHandler, "Le once handler sync doit être appelé une seule fois")
	})

	t.Run("Test mixing previous scenario", func(t *testing.T) {
		var countHandler int = 0
		var countOnceHandler int = 0

		var handler hook.Handler = func(tags ...string) {
			countHandler++
		}

		var onceHandler hook.OnceHandler = func(tags ...string) {
			countOnceHandler++
		}

		var handlerSync hook.HandlerSync = func(tags ...string) {
			countHandler++
		}

		var onceHandlerSync hook.OnceHandlerSync = func(tags ...string) {
			countOnceHandler++
		}

		hook.Register(ENAME, handler)
		hook.Register(ENAME, onceHandler)

		hook.Call(ENAME)
		hook.Call(ENAME)
		time.Sleep(1000 * time.Millisecond)

		assert.Equal(t, 4, countHandler, "Le handler normal doit être appelé deux fois")
		assert.Equal(t, 1, countOnceHandler, "Le once handler doit être appelé une seule fois")

		hook.Register(ENAME, handlerSync)
		hook.Register(ENAME, onceHandlerSync)

		hook.Call(ENAME)
		hook.Call(ENAME)
		time.Sleep(1000 * time.Millisecond)

		assert.Equal(t, 12, countHandler, "Handler total (normal + sync) doit être 12")
		assert.Equal(t, 2, countOnceHandler, "OnceHandler total (once + once sync) doit être 2")
	})
}
