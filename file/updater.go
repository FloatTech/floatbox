package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
	"unsafe"

	reg "github.com/fumiama/go-registry"
	"github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/extension/ttl"

	"github.com/FloatTech/floatbox/process"
)

const (
	dataurl = "https://gitcode.net/u011570312/zbpdata/-/raw/main/"
)

var (
	mu      sync.Mutex
	connmap = ttl.NewCacheOn(time.Second*3, [4]func(struct{}, *reg.Regedit){
		func(s struct{}, r *reg.Regedit) {
			logrus.Infoln("[file]已连接md5验证服务器")
		},
		func(s struct{}, r *reg.Regedit) {
			logrus.Infoln("[file]重用md5验证服务器")
		},
		func(s struct{}, r *reg.Regedit) {
			process.GlobalInitMutex.Lock()
			_ = r.Close()
			logrus.Infoln("[file]关闭到md5验证服务器的连接")
			process.GlobalInitMutex.Unlock()
		}, nil,
	})
)

// GetLazyData 获取懒加载数据
// 传入的 path 的前缀 data/
// 在验证完 md5 后将被删去
// 以便进行下载
func GetLazyData(path string, isDataMustEqual bool) ([]byte, error) {
	var data []byte
	var resp *http.Response
	var filemd5 *[16]byte
	var ms string
	var err error

	u := dataurl + path[5:]

	r := connmap.Get(struct{}{})
	if r == nil {
		mu.Lock()
		if r == nil {
			r = reg.NewRegReader("reilia.fumiama.top:32664", "fumiama")
			connerr := r.ConnectIn(time.Second * 4)
			if connerr != nil {
				logrus.Warnln("[file]连接md5验证服务器失败:", connerr)
				mu.Unlock()
				return nil, connerr
			}
			connmap.Set(struct{}{}, r)
		}
		mu.Unlock()
	}

	if r == nil {
		logrus.Warnln("[file]无法连接到md5验证服务器, 请自行确保下载文件", path, "的正确性")
	} else {
		ms, err = r.Get(path)
		if err != nil || len(ms) != 16 {
			logrus.Warnln("[file]获取md5失败, 请自行确保下载文件", path, "的正确性:", err)
		} else {
			filemd5 = (*[16]byte)(*(*unsafe.Pointer)(unsafe.Pointer(&ms)))
			logrus.Debugln("[file]从验证服务器获得文件", path, "md5:", hex.EncodeToString(filemd5[:]))
		}
	}

	if IsExist(path) {
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if filemd5 != nil {
			if md5.Sum(data) == *filemd5 {
				logrus.Debugln("[file]文件", path, "md5匹配, 文件已存在且为最新")
				goto ret
			} else if !isDataMustEqual {
				logrus.Warnln("[file]文件", path, "md5不匹配, 但不主动更新")
				goto ret
			}
			logrus.Debugln("[file]文件", path, "md5不匹配, 开始更新文件")
		} else {
			logrus.Warnln("[file]文件", path, "存在, 已跳过md5检查")
			goto ret
		}
	}

	// 下载
	resp, err = http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.ContentLength <= 0 {
		return nil, errors.New("resp body len <= 0")
	}
	logrus.Printf("[file]从镜像下载数据%d字节...", resp.ContentLength)
	// 读取数据
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("read body len <= 0")
	}
	if filemd5 != nil {
		if md5.Sum(data) == *filemd5 {
			logrus.Debugln("[file]文件", path, "下载完成, md5匹配, 开始保存")
		} else {
			logrus.Errorln("[file]文件", path, "md5不匹配, 下载失败")
			return nil, errors.New("file md5 mismatch")
		}
	} else {
		logrus.Warnln("[file]文件", path, "下载完成, 已跳过md5检查, 开始保存")
	}
	// 写入数据
	err = os.WriteFile(path, data, 0644)
ret:
	return data, err
}
