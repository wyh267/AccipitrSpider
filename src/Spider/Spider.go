/*****************************************************************************
 *  file name : Spider.go
 *  author : Wu Yinghao
 *  email  : wyh817@gmail.com
 *
 *  file description : spider detail file	
 *
******************************************************************************/



package Spider

import (
	u "Utils"
	"github.com/outmana/log4jzl"
	"time"
)


type Spider struct {
	In chan u.SpiderIn
	Out chan u.SpiderOut
	Contrl chan u.SpiderContrl
	Logger *log4jzl.Log4jzl
	Name string
}



func NewSpider(name string,in chan u.SpiderIn,out chan u.SpiderOut,logger *log4jzl.Log4jzl)(*Spider){
	contrl:=make(chan u.SpiderContrl,100)
	this := &Spider{Name:name,In:in,Out:out,Logger:logger,Contrl:contrl}
	return this
	
}



func (this *Spider)CrawlSync(url string)(string,error){
	
	
	content,err := u.RequestUrl(url)
	//this.Logger.Info("[INFO] Crawl :%v , err :%v",url,err)
	return content,err
}




func (this *Spider)StopCrawl() error {
	
	
	
	return nil
}




func (this *Spider)StartCrawl()error {
	
	this.Logger.Info(" Start Spider[ %v ]....",this.Name)
	u.SpiderSync.Add(1)
	go this.runCrawl()
	return nil
}






func (this *Spider)PullUrl(url u.SpiderIn)error{
	this.In<-url
	return nil
}





func (this *Spider)ChanLen()int{
	return len(this.In)
}



func (this *Spider)runCrawl() error {
	defer u.SpiderSync.Done()
	this.Logger.Info("Spider[ %v ] Start Running Crawl ",this.Name)
	for{
		
		select{
			case in:= <-this.In:
				//this.Logger.Info("Spider[ %v ] Recive In : [ %v ] ",this.Name,in)

				content,err:=this.CrawlSync(in.Url)
				if err != nil {
					this.Logger.Info("Spider[ %v ] No Data ",this.Name)
					continue
				}
				this.Logger.Info("Spider[ %v ] Get URL Data [ %v ] Send To Dispatcher ...",this.Name,in.Url)
				this.Out <-u.SpiderOut{Content:content,Header:"",Url:in.Url,HandlerName:in.HandlerName}
				

			case contrl:= <-this.Contrl:
				this.Logger.Info("Spider[ %v ] Recive Contrl : %v ",this.Name,contrl)
				return nil
				
			case <-time.After(time.Second*10):
				//this.Logger.Info("[INFO] no url to crawl")
			
		}
			
	}
	return nil
	
}





