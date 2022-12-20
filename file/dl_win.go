//go:build windows
// +build windows

package file

import (
	"syscall"
	"unsafe"
)

var (
	urlmon                   = syscall.MustLoadDLL("Urlmon.dll")
	processURLDownloadToFile = urlmon.MustFindProc("URLDownloadToFileW")
)

func DownloadTo(url, file string) error {
	u := []rune(url)
	u16 := make([]uint16, len(u)+1)
	for i, c := range u {
		u16[i] = uint16(c)
	}
	f := []rune(file)
	f16 := make([]uint16, len(f)+1)
	for i, c := range f {
		f16[i] = uint16(c)
	}
	ret, _, err := processURLDownloadToFile.Call(
		0,
		*(*uintptr)(unsafe.Pointer(&u16)),
		*(*uintptr)(unsafe.Pointer(&f16)),
		0, 0)
	if ret == 0 {
		err = nil
	}
	return err
}
