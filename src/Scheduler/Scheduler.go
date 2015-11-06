

package Scheduler


import (
	s "Spider"
	"github.com/outmana/log4jzl"
	u "Utils"
	"fmt"
	"time"
	"math/rand"
)



type Scheduler struct{
	SpiderNum  int
	CrawlChan  chan u.CrawlData 
	//Outchan	   chan u.SpiderOut
	Spiders	   map[int]*s.Spider
	Logger *log4jzl.Log4jzl
	kvdb *u.PandionKV
}



func NewScheduler(spider_num int,logger *log4jzl.Log4jzl,db *u.PandionKV)(*Scheduler){
	
	crawlChan :=make(chan u.CrawlData,100)
	this:=&Scheduler{SpiderNum:spider_num,Logger:logger,CrawlChan:crawlChan,kvdb:db}
	
	
	return this	

}





func (this *Scheduler)ConfigScheduler()(chan u.SpiderOut,error){
	
	Outchan:=make(chan u.SpiderOut,100)
	this.Spiders=make(map[int]*s.Spider)
	for i:=0;i<this.SpiderNum;i++{
		name:=fmt.Sprintf("%v",i)
		in:=make(chan u.SpiderIn,100)
		this.Spiders[i]=s.NewSpider(name,in,Outchan,this.Logger)
		this.Spiders[i].StartCrawl()
		
	}
	//this.Logger.Info("[INFO] SPIDER %v",this.Spiders)
	return Outchan,nil
	
	
}




func (this *Scheduler)StartScheduler()error{
	
	u.SpiderSync.Add(1)
	go this.runScheduler()
	this.Logger.Info("[INFO] Start Scheduler now ... ")
	return nil
	
}




func (this *Scheduler)SendUrl(craldata u.CrawlData) error{
	this.CrawlChan<-craldata
	return nil
}


func (this *Scheduler)runScheduler()error{
	defer u.SpiderSync.Done()
	this.Logger.Info("[INFO] Start Running Scheduler ")
	for{
		
		select{
			case in:= <-this.CrawlChan:
				index:=rand.Intn(this.SpiderNum)
				this.Logger.Info("Scheduler Recive URL:[ %v ] Parpare to send to spider[%v]",in.Url,index)
				if in.Type==u.DETAIL_URL {
					if this.CheckUrl(in.Url) {
					this.Spiders[index].PullUrl(u.SpiderIn{Url:in.Url,Timeout:10,HandlerName:in.HandlerName})
				}
				}else{
					this.Spiders[index].PullUrl(u.SpiderIn{Url:in.Url,Timeout:10,HandlerName:in.HandlerName})
				}
				
				//this.Logger.Info("[INFO] index %v",index)
				
			case <-time.After(time.Second*10):
				//this.Logger.Info("[INFO] no CrawlChan to crawl")
			
		}
			
	}
	return nil
	
}



func (this *Scheduler)CheckUrl(url string) bool{
	
	_,err := this.kvdb.Get(url)
	if err != nil {
		this.kvdb.Set(url,"1")
		return true
	}
	
	return false
	
}



