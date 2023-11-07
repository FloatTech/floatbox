package file

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"
	"unicode"
	"unsafe"

	reg "github.com/fumiama/go-registry"
	"github.com/pbnjay/memory"
	"github.com/sirupsen/logrus"

	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/ttl"
)

const (
	dataurl = "https://gitea.seku.su/fumiama/zbpdata/raw/branch/main/"
	wifeurl = "https://gitea.seku.su/fumiama/zbpwife/raw/branch/main/"
)

var (
	o              = process.NewOnce()
	s              *reg.Storage
	ErrEmptyBody   = errors.New("read body len <= 0")
	ErrInvalidPath = errors.New("invalid path")
	lazycache      = ttl.NewCache[string, []byte](func() time.Duration {
		d := time.Duration(memory.TotalMemory()-memory.FreeMemory()) * 60 // 1G: 1073741824 * 60 ns
		if d <= time.Millisecond {
			return time.Millisecond // min cache time is 1 ms
		}
		return d // 1G is about 1 min
	}())
)

// GetCustomLazyData 获取自定义懒加载数据, 不进行 md5 验证, 忽略 data/Abcde/ 路径.
// 传入的 path 的前缀 data/abcde/ 将被删去
// 以便进行下载, 但保存时仍位于 data/abcde/xxx
func GetCustomLazyData(dataurl, path string) ([]byte, error) {
	data := lazycache.Get(path)
	if data != nil {
		logrus.Debugln("[file]获取缓存的文件数据", path)
		return data, nil
	}
	logrus.Debugln("[file]从自定义镜像下载", path)
	_, p, found := strings.Cut(path[5:], "/")
	if !found {
		return nil, ErrInvalidPath
	}
	u := dataurl + p + "?inline=true"
	if IsExist(path) {
		return os.ReadFile(path)
	}
	// 下载
	data, err := web.GetData(u)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, ErrEmptyBody
	}
	lazycache.Set(path, data)
	// 写入数据
	return data, os.WriteFile(path, data, 0644)
}

// GetLazyData 获取公用懒加载数据
// 传入的 path 的前缀 data/
// 在验证完 md5 后将被删去
// 以便进行下载, 但保存时仍位于 data/Abcde/xxx
func GetLazyData(path, stor string, isDataMustEqual bool) ([]byte, error) {
	var data []byte
	var filemd5 *[16]byte
	var ms string
	var err error

	if !unicode.IsUpper([]rune(path[5:])[0]) {
		panic("cannot get private data")
	}

	data = lazycache.Get(path)
	if data != nil {
		logrus.Debugln("[file]获取缓存的文件数据", path)
		return data, nil
	}

	u := path[5:] + "?inline=true"
	if strings.HasPrefix(path, "data/Wife/") {
		u = wifeurl + u[5:]
	} else {
		u = dataurl + u
	}

	o.Do(func() {
		r := reg.NewRegReader("reilia.fumiama.top:32664", stor, "fumiama")
		s, err = r.Load()
		if err != nil {
			err = r.ConnectIn(time.Second * 4)
			if err != nil {
				logrus.Warnln("[file]连接md5验证服务器失败:", err)
				return
			}
			s, err = r.Cat()
			if err != nil {
				logrus.Warnln("[file]获取md5数据库失败:", err)
				return
			}
			logrus.Infoln("[file]获取md5数据库")
		} else {
			logrus.Infoln("[file]加载md5数据库...")
			if err = r.ConnectIn(time.Second * 4); err == nil {
				if ok, _ := r.IsMd5Equal(s.Md5); !ok {
					logrus.Infoln("[file]md5数据库不是最新, 更新中...")
					if ns, err := r.Cat(); err == nil {
						s = ns
						logrus.Infoln("[file]md5数据库已更新")
					} else {
						logrus.Warnln("[file]md5数据库更新失败, 数据库不是最新:", err)
					}
				} else {
					logrus.Infoln("[file]md5数据库已是最新")
				}
			} else {
				logrus.Warnln("[file]连接md5验证服务器失败, 数据库可能不是最新:", err)
			}
		}
		go func() {
			for range time.NewTicker(time.Hour).C {
				r := reg.NewRegReader("reilia.fumiama.top:32664", stor, "fumiama")
				err := r.ConnectIn(time.Second * 4)
				if err != nil {
					logrus.Warnln("[file]连接md5验证服务器失败:", err)
					continue
				}
				ok, _ := r.IsMd5Equal(s.Md5)
				if ok {
					logrus.Infoln("[file]md5无变化")
					continue
				}
				ns, err := r.Cat()
				if err != nil {
					logrus.Warnln("[file]获取md5数据库失败:", err)
					return
				}
				s = ns
			}
		}()
	})

	if s == nil {
		logrus.Warnln("[file]无法连接到md5验证服务器, 请自行确保下载文件", path, "的正确性")
	} else {
		ms, err = s.Get(path)
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
	data, err = web.GetData(u)
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
	lazycache.Set(path, data)
	return data, err
}
