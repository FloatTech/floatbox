// Package file 文件实用工具
package file

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/RomiChan/syncx"
)

type dlcache syncx.Map[string, error]

func (dlc *dlcache) wait(url string) error {
	if err, loaded := (*syncx.Map[string, error])(dlc).LoadOrStore(url, errDlStatusDoing); loaded {
		if err != errDlStatusDoing {
			return err
		}
		t := time.NewTicker(time.Second)
		defer t.Stop()
		i := 0
		for range t.C {
			if err, ok := (*syncx.Map[string, error])(dlc).Load(url); ok && err != errDlStatusDoing {
				return err
			}
			i++
			if i > 60 {
				break
			}
		}
		return errDlStatusTimeout
	}
	_ = time.AfterFunc(time.Minute, func() {
		(*syncx.Map[string, error])(dlc).Delete(url)
	})
	return errDlContinue
}

func (dlc *dlcache) set(url string, err error) {
	(*syncx.Map[string, error])(dlc).Store(url, err)
}

var dlmap = dlcache{}

var (
	errDlContinue      = errors.New("continue")
	errDlStatusDoing   = errors.New("downloading")
	errDlStatusTimeout = errors.New("download timeout")
)

// DownloadTo 下载到路径
func DownloadTo(url, file string) error {
	err := dlmap.wait(url)
	if err != errDlContinue {
		return err
	}
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
	dlmap.set(url, err)
	return err
}
