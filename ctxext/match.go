package ctxext

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/go-ego/gse"

	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/ttl"
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

var strcutcache = ttl.NewCache[string, []string](time.Hour * 24)

// JiebaSimilarity sameper from 0.0 to 1.0
func JiebaSimilarity[Ctx any](sameper float64, seg *gse.Segmenter, getmsg func(Ctx) string, src ...string) func(Ctx) bool {
	return func(ctx Ctx) bool {
		msgs := seg.Cut(getmsg(ctx), true)
		msgv := make(map[string]uint8, len(msgs)*2)
		for _, msg := range msgs {
			msgv[msg]++
		}
		for _, str := range src {
			words := strcutcache.Get(str)
			if len(words) == 0 {
				words = seg.Cut(str, true)
				strcutcache.Set(str, words)
			}
			testv := make(map[string]uint8, len(words)*2)
			for _, word := range words {
				testv[word]++
			}
			msgspace := make([]uint8, 0, len(msgv)+len(testv))
			strspace := make([]uint8, 0, len(msgv)+len(testv))
			for k, v := range msgv {
				msgspace = append(msgspace, v)
				if tv, ok := testv[k]; ok {
					strspace = append(strspace, tv)
					delete(testv, k)
				} else {
					strspace = append(strspace, 0)
				}
			}
			for _, v := range testv {
				msgspace = append(msgspace, 0)
				strspace = append(strspace, v)
			}
			if math.Similarity(msgspace, strspace) > sameper {
				p := reflect.ValueOf(ctx).Elem().FieldByName("State").UnsafePointer()
				(*(*map[string]interface{})(unsafe.Pointer(&p)))["matched"] = str
				return true
			}
		}
		return false
	}
}
