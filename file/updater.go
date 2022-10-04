package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
	"unsafe"

	reg "github.com/fumiama/go-registry"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/ttl"

	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/floatbox/web"
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
	ErrEmptyBody   = errors.New("read body len <= 0")
	ErrInvalidPath = errors.New("invalid path")
)

// GetCustomLazyData 获取自定义懒加载数据, 不进行 md5 验证, 忽略 data/Abcde/ 路径.
// 传入的 path 的前缀 data/abcde/
// 在验证完 md5 后将被删去
// 以便进行下载, 但保存时仍位于 data/abcde/xxx
func GetCustomLazyData(dataurl, path string) ([]byte, error) {
	_, p, found := strings.Cut(path[5:], "/")
	if !found {
		return nil, ErrInvalidPath
	}
	u := dataurl + p + "?inline=true"
	if IsExist(path) {
		return os.ReadFile(path)
	}
	// 下载
	data, err := web.RequestDataWith(web.NewTLS12Client(), u, "GET", "gitcode.net", web.RandUA())
	if err != nil {
		return nil, err
	}
	logrus.Printf("[file]从自定义镜像下载数据%d字节...", len(data))
	if len(data) == 0 {
		return nil, ErrEmptyBody
	}
	// 写入数据
	return data, os.WriteFile(path, data, 0644)
}

// GetLazyData 获取公用懒加载数据
// 传入的 path 的前缀 data/
// 在验证完 md5 后将被删去
// 以便进行下载, 但保存时仍位于 data/Abcde/xxx
func GetLazyData(path string, isDataMustEqual bool) ([]byte, error) {
	var data []byte
	var filemd5 *[16]byte
	var ms string
	var err error

	if !unicode.IsUpper([]rune(path[5:])[0]) {
		panic("cannot get private data")
	}

	u := dataurl + path[5:] + "?inline=true"

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
	data, err = web.RequestDataWith(web.NewTLS12Client(), u, "GET", "gitcode.net", web.RandUA())
	if err != nil {
		return nil, err
	}
	logrus.Printf("[file]从镜像下载数据%d字节...", len(data))
	if len(data) == 0 {
		return nil, ErrEmptyBody
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
