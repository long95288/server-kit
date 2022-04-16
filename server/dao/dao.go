package dao

import (
	"log"
	"os"
	"path"
	"server-kit/server/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type Msg struct {
	Id      int64
	Sender  string
	Body    string
	Created time.Time `xorm:"created notnull"`
}

var engine *xorm.Engine

func InitOrm() error {
	var err error
	err = os.MkdirAll(config.SrvConf.MsgDBPath, os.ModePerm)
	if err != nil {
		return err
	}
	engine, err = xorm.NewEngine("sqlite3", path.Join(config.SrvConf.MsgDBPath, "msg.db"))
	if nil != err {
		return err
	}
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	if ok, err := engine.IsTableExist(Msg{}); !ok {
		err = engine.CreateTables(Msg{})
		if err != nil {
			return err
		}
	}
	log.Println("init xorm engine success")
	return nil
}

func AddMsg(msg Msg) error {
	affect, err := engine.InsertOne(msg)
	if err != nil {
		return err
	}
	log.Printf("AddMsg %d success \n", affect)
	return nil
}

func GetMsgList(pageNum int, pageSize int) ([]Msg, error) {
	msgList := make([]Msg, 0)
	err := engine.Where("1 = 1").Limit(pageSize, (pageNum-1)*pageSize).Desc("created").Find(&msgList)
	return msgList, err
}
