package ctxext

import (
	"math/rand"
	"reflect"
	"strings"
	"unsafe"

	"github.com/fumiama/jieba"
)

// ListGetter 获得实时刷新的 list
type ListGetter interface {
	List() []string
}

// ValueInList 判断参数是否在列表中
func ValueInList[Ctx any](getval func(Ctx) string, list ListGetter) func(Ctx) bool {
	return func(ctx Ctx) bool {
		val := getval(ctx)
		for _, v := range list.List() {
			if val == v {
				return true
			}
		}
		return false
	}
}

func JiebaFullMatch[Ctx any](seg *jieba.Segmenter, getmsg func(Ctx) string, src ...string) func(Ctx) bool {
	return func(ctx Ctx) bool {
		msgs := seg.Cut(getmsg(ctx), true)
		rand.Shuffle(len(msgs), func(i, j int) { msgs[i], msgs[j] = msgs[j], msgs[i] })
		for _, str := range src {
			for _, msg := range msgs {
				if str == msg {
					p := reflect.ValueOf(ctx).Elem().FieldByName("State").UnsafePointer()
					(*(*map[string]interface{})(unsafe.Pointer(&p)))["matched"] = msg
					return true
				}
			}
		}
		return false
	}
}

func JiebaKeyword[Ctx any](seg *jieba.Segmenter, getmsg func(Ctx) string, src ...string) func(Ctx) bool {
	return func(ctx Ctx) bool {
		msgs := seg.Cut(getmsg(ctx), true)
		rand.Shuffle(len(msgs), func(i, j int) { msgs[i], msgs[j] = msgs[j], msgs[i] })
		for _, str := range src {
			for _, msg := range msgs {
				if strings.Contains(msg, str) {
					p := reflect.ValueOf(ctx).Elem().FieldByName("State").UnsafePointer()
					(*(*map[string]interface{})(unsafe.Pointer(&p)))["keyword"] = str
					return true
				}
			}
		}
		return false
	}
}
