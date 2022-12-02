package ctxext

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// State store the context of a matcher.
type State map[string]interface{}

// Ctx represents the Context which hold the event.
// 代表上下文
type Ctx struct {
	State State

	message string
}

func TestCtxState(t *testing.T) {
	ctx := &Ctx{State: State{}, message: "23333"}
	p := reflect.ValueOf(ctx).Elem().FieldByName("State").UnsafePointer()
	(*(*map[string]interface{})(unsafe.Pointer(&p)))["matched"] = "test"
	assert.Equal(t, "test", ctx.State["matched"])
}
