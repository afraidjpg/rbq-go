package app

import (
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"github.com/alive1944/qq-robot-go/core/config"
	"sync"
)

var storage *leveldb.DB
var l = &sync.Mutex{}
// 连接到leveldb文件
func connectToLeveldb() {
	l.Lock()
	defer l.Unlock()
	if storage != nil {
		return
	}

	basePath, _ := os.Getwd()
	loginQQ := config.Cfg.GetInt64("account.login_qq")
	dataPath := filepath.Join(basePath, fmt.Sprintf("../../data/storage_%d", loginQQ))
	if _, err:= os.Stat("./leveldb.go"); errors.Is(err, fs.ErrNotExist){
		dataPath = filepath.Join(basePath, fmt.Sprintf("./data/storage_%d", loginQQ))
	}

	db, err := leveldb.OpenFile(dataPath, nil)
	if err != nil {
		log.Fatal("初始化数据库出错，原因: ", err)
	}
	storage = db
}

func GetStorage() *leveldb.DB {
	if storage == nil {
		connectToLeveldb()
	}
	return storage
}
