//go:build windows
// +build windows

package file

import (
	"strings"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

var (
	urlmon                   = syscall.MustLoadDLL("Urlmon.dll")
	processURLDownloadToFile = urlmon.MustFindProc("URLDownloadToFileW")
	processIsValidURL        = urlmon.MustFindProc("IsValidURL")
)

const (
	CP_UTF8 = 65001
)

var (
	k32                        = syscall.MustLoadDLL("Kernel32.dll")
	processMultiByteToWideChar = k32.MustFindProc("MultiByteToWideChar")
)

var (
	ErrInvalidURL = errors.New("invalid url")
)

func DownloadTo(url, file string) error {
	file = strings.ReplaceAll(file, "/", "\\")
	u16 := make([]uint16, len(url)*4)
	ret, _, err := processMultiByteToWideChar.Call(CP_UTF8, 0, *(*uintptr)(unsafe.Pointer(&url)), uintptr(len(url)), *(*uintptr)(unsafe.Pointer(&u16)), uintptr(len(u16)))
	if ret == 0 {
		return errors.Wrap(err, "url")
	}
	u16 = u16[:ret]
	ret, _, _ = processIsValidURL.Call(0, *(*uintptr)(unsafe.Pointer(&u16)), 0)
	if ret != 0 {
		return ErrInvalidURL
	}
	f16 := make([]uint16, len(file)*4)
	ret, _, err = processMultiByteToWideChar.Call(CP_UTF8, 0, *(*uintptr)(unsafe.Pointer(&file)), uintptr(len(file)), *(*uintptr)(unsafe.Pointer(&f16)), uintptr(len(f16)))
	if ret == 0 {
		return errors.Wrap(err, "file")
	}
	f16 = f16[:ret]
	ret, _, _ = processURLDownloadToFile.Call(
		0,
		*(*uintptr)(unsafe.Pointer(&u16)),
		*(*uintptr)(unsafe.Pointer(&f16)),
		0, 0)
	if ret == 0 {
		return nil
	}
	return syscall.Errno(ret)
}
