//go:build !windows
// +build !windows

package file

import (
	"io"
	"net/http"
	"os"
)

// DownloadTo 下载到路径
func DownloadTo(url, file string) error {
	resp, err := http.Get(url)
	if err == nil {
		var f *os.File
		f, err = os.Create(file)
		if err == nil {
			_, err = io.Copy(f, resp.Body)
			f.Close()
		}
		resp.Body.Close()
	}
	return err
}
