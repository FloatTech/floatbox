package binary

import (
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// BytesToString 没有内存开销的转换
//
// github.com/wdvxdr1123/ZeroBot/utils/helper.BytesToString
func BytesToString(b []byte) string {
	return helper.BytesToString(b)
}

// StringToBytes 没有内存开销的转换
//
// github.com/wdvxdr1123/ZeroBot/utils/helper.StringToBytes
func StringToBytes(s string) (b []byte) {
	return helper.StringToBytes(s)
}
