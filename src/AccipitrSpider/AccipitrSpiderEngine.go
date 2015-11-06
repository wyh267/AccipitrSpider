/*****************************************************************************
 *  file name : spider.go
 *  author : Wu Yinghao
 *  email  : wyh817@gmail.com
 *
 *  file description : Spider entry file
 *
******************************************************************************/

package main

import (
	"Scheduler"
	"fmt"
	//"Spider"
	"Dispatcher"
	"Plugs"
	u "Utils"
	"bufio"
	"errors"
	"flag"
	"github.com/outmana/log4jzl"
	"io"
	"os"
	"strings"
	"time"
	//"sync"
	//"time"
)

func main() {

	//var wg sync.WaitGroup

	//读取启动参数
	var configFile string
	flag.StringVar(&configFile, "conf", "config.ini", "configure file full path")
	flag.Parse()

	fmt.Printf("Start AccipitrSpiderEngine...\n")

	//启动日志系统
	logger, err := log4jzl.New("AccipitrSpiderEngine")
	if err != nil {
		fmt.Printf("[ERROR] Create logger Error: %v\n", err)
		//return
	}

	//读取配置文件
	configure, err := u.NewConfigure(configFile)
	if err != nil {
		fmt.Printf("[ERROR] Parse Configure File Error: %v\n", err)
		return
	}

	//初始化数据库适配器
	/*
		dbAdaptor, err := u.NewDBAdaptor(configure, logger)
		if err != nil {
			fmt.Printf("[ERROR] Create DB Adaptor Error: %v\n", err)
			return
		}
		defer dbAdaptor.Release()
	*/
	//初始化KVDB数据库
	DBname, _ := configure.GetKVDB()
	var kvdb *u.PandionKV
	if !Exist(fmt.Sprintf("./DB/%v.idx", DBname)) {
		kvdb = u.NewPandionKV(DBname, logger)
	} else {
		kvdb = u.NewPandionKVWithFile(DBname, logger)
	}

	//初始化自定义的插件
	ps := make(map[string]u.PlugInterface)
	plug := Plugs.NewYXQ("yingxiaoqun", logger, nil) // dbAdaptor)
	plugyo:=Plugs.NewRadialeng("radialeng",logger,nil)
	maudio:=Plugs.NewMAudio("maudio",logger,nil)
	jblpro:=Plugs.NewJblpro("jblpor",logger,nil)
	ikmultimedia:=Plugs.NewIkmultimedia("ikmultimedia",logger,nil)
	avid:=Plugs.NewAvid("avid",logger,nil)
	shure:=Plugs.NewShure("shure",logger,nil)
	ps["yingxiaoqun"] = plug
	ps["radialeng"] = plugyo
	ps["maudio"]= maudio
	ps["jblpro"]= jblpro
	ps["ikmultimedia"]=ikmultimedia
	ps["avid"]=avid
	ps["shure"]=shure
	//启动调度器
	scheduler := Scheduler.NewScheduler(9, logger, kvdb)
	out_chan, _ := scheduler.ConfigScheduler()
	scheduler.StartScheduler()

	//启动分发器
	dispatcher := Dispatcher.NewDispatcher(out_chan, logger)
	dispatcher.ConfigureDispatcher(9, ps, scheduler)
	dispatcher.StartDispatcher()

	//读取种子url
	urls, err := ReadSeedUrls()
	if err != nil {
		fmt.Printf("[ERROR] ReadSeedUrls Error: %v\n", err)
		return
	}

	for k, v := range urls {
		for _, url := range v {
			scheduler.SendUrl(u.CrawlData{Url: url, Type: u.SEED_URL, HandlerName: k})
		}
	}

	interval, _ := configure.GetInterval()
	for {
		select {
		case <-time.After(time.Second * time.Duration(interval)):
			for k, v := range urls {
				for _, url := range v {
					scheduler.SendUrl(u.CrawlData{Url: url, Type: u.SEED_URL, HandlerName: k})
				}
			}

		}
	}

	//out<-u.SpiderOut{Url:"http://www.yingxiaoqun.com111",Content:"",Header:""}
	u.SpiderSync.Wait()
	//time.Sleep(100*time.Second)

}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func ReadSeedUrls() (map[string][]string, error) {
	urls := make(map[string][]string, 0)
	f, err := os.Open("./seed_urls.txt")
	defer f.Close()
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return urls, nil
			}
			return nil, err
		}
		splits := strings.Split(strings.TrimSpace(line), "\t")
		//logger.Info(" splits : %v ",splits)
		if len(splits) != 2 {
			continue
		}

		v, ok := urls[splits[1]]
		if !ok {
			urls[splits[1]] = make([]string, 0)
			urls[splits[1]] = append(urls[splits[1]], splits[0])
		} else {
			v = append(v, splits[0])
			urls[splits[1]] = v
		}

	}

	return nil, errors.New("no urls")
}
