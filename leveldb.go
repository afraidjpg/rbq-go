package qq_robot_go

import (
	"errors"
	"fmt"
	"github.com/afraidjpg/qq-robot-go/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

// GetStorage 获取lvdb实例
func GetStorage() *leveldb.DB {
	if storage == nil {
		connectToLeveldb()
	}
	return storage
}

// FetchData 从leveldb中获取数据
func FetchData(db string, ro *opt.ReadOptions) ([]byte, error) {
	return GetStorage().Get([]byte(db), ro)
}

// StoreData 向leveldb中推入数据
func StoreData(dbname string, data []byte, wo *opt.WriteOptions) error {
	return GetStorage().Put([]byte(dbname), data, wo)
}