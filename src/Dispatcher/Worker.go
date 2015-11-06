
package Dispatcher

import (
	u "Utils"
	"github.com/outmana/log4jzl"
	//"gopkg.in/xmlpath.v1"
	//"strings"
	//"github.com/moovweb/gokogiri"
	s "Scheduler"
	
)


type Worker struct{
	Plug	map[string]u.PlugInterface
	Name	string
	Logger *log4jzl.Log4jzl
	StatusChan  chan u.WorkStatus
	Sche 	*s.Scheduler
}



func NewWorker(name string,plug map[string]u.PlugInterface,statuschan chan u.WorkStatus,logger *log4jzl.Log4jzl,scheduler *s.Scheduler)(*Worker){
	this := &Worker{Name:name,Plug:plug,Logger:logger,StatusChan:statuschan,Sche:scheduler}
	this.StatusChan <- u.WorkStatus{this.Name,0}
	return this
}




func (this *Worker)GoWork(content,url,handler string)error{
	
	this.StatusChan <- u.WorkStatus{this.Name,1}
	
	this.Logger.Info("Worker[ %v ] Start Working...hander : %v ",this.Name,handler)
	urls,err:=this.Plug[handler].GetNextUrls(content,url)
	//this.Logger.Info("[INFO] urls %v,err:%v",urls,err)
	
	if len(urls)==0 || urls==nil || err !=nil {
		this.Plug[handler].GetDetailContent(content,url)
	}else{
		for _,url := range urls{
			this.Sche.SendUrl(url)
		}
		
	}
	
	
	//this.Logger.Info("[INFO] content %v",content)
	
	
	/////doc, _ := gokogiri.ParseHtml([]byte(content))
	//html:=doc.Root()
  	// perform operations on the parsed page -- consult the tests for examples
	/////res,_:=doc.Search("//*[@class=\"read-more\"]/@href")
  	// important -- don't forget to free the resources when you're done!
	//this.Logger.Info("[INFO]  %v",res)
  	
	/////for i:=range res{
	/////	this.Logger.Info("[INFO]  %v",res[i])
	/////}
	/////doc.Free()
	this.Logger.Info("Worker[ %v ] Work Finish [ %v ]...",this.Name,url)
	this.StatusChan <- u.WorkStatus{this.Name,0}
	return nil
	
}
