package global

import (
	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"fmt"
	"github.com/beewit/beekit/utils/convert"
	"encoding/json"
)

var (
	CFG         = conf.New("config.json")
	Log         = log.Logger
	DB          = mysql.DB
	RD          = redis.Cache
	IP          = CFG.Get("server.ip")
	Port        = CFG.Get("server.port")
	Host        = fmt.Sprintf("http://%v:%v", IP, Port)
	FilesPath   = fmt.Sprintf("%v", CFG.Get("files.path"))
	FilesDoMain = fmt.Sprintf("%v", CFG.Get("files.doMain"))
	MaxFileSize = convert.MustFloat64(fmt.Sprintf("%v", CFG.Get("files.maxFileSize")))
	ExtFilter   = fmt.Sprintf("%v", CFG.Get("files.extFilter"))
)

type Account struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Photo    string `json:"photo"`
	Mobile   string `json:"mobile"`
	Status   string `json:"status"`
}

func ToByteAccount(b []byte) *Account {
	var rp = new(Account)
	err := json.Unmarshal(b[:], &rp)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	return rp
}

func ToInterfaceAccount(m interface{}) *Account {
	b := convert.ToInterfaceByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}

func ToMapAccount(m map[string]interface{}) *Account {
	b := convert.ToMapByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}
