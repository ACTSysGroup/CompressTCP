package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/rand"
)

var db *gorm.DB

type LibgenMeta struct {
	Id    int    `gorm:"column:id"`
	Md5   string `gorm:"column:md5"`
	Descr string `gorm:"column:descr"`
	Toc   string `gorm:"column:toc"`
}

type Credential struct {
	Username  string
	Password  string
	Host      string
	Database  string
	Port      string
	Batch     int
	Iteration int
	Routine   int
}

func (c Credential) ToDSN() string {
	res := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database)
	fmt.Printf(res)
	return res
}

var cfg Credential

func init() {
	var cfgFile string
	var err error

	flag.StringVar(&cfgFile, "conf", "./config.yaml", "config path")
	flag.Parse()

	if bytes, err := ioutil.ReadFile(cfgFile); err != nil {
		panic(err.Error())
	} else {
		if err2 := yaml.Unmarshal(bytes, &cfg); err2 != nil {
			panic(err.Error())
		}
	}
	if db, err = gorm.Open(mysql.Open(cfg.ToDSN()), &gorm.Config{}); err != nil {
		panic("can not connect to db")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
}
func main() {
	var count int64
	var wg sync.WaitGroup

	db.Table("description").Count(&count)
	sql := fmt.Sprintf("SELECT * FROM description WHERE id > %d LIMIT %d;", rand.Int63nRange(0, count), cfg.Batch)

	fmt.Printf("%v", cfg)
	PerformFetch := func(idx int) {
		var records []LibgenMeta
		for j := 0; j < cfg.Iteration; j++ {
			if db.Raw(sql).Scan(&records); db.Error != nil {
				fmt.Errorf(db.Error.Error())
			}
			fmt.Printf("%d\n", idx)
		}
		wg.Done()
	}
	wg.Add(cfg.Routine)
	for i := 0; i < cfg.Routine; i++ {
		go PerformFetch(i)
	}
	wg.Wait()
}
