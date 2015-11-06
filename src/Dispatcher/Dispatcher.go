


package Dispatcher


import (
	u "Utils"
	"github.com/outmana/log4jzl"
	"fmt"
	//"Plugs"
	s "Scheduler"
)


type Dispatcher struct{
	Contrl    chan u.DispatcherContrl
	UrlDetail chan u.SpiderOut
	Logger *log4jzl.Log4jzl
	Workers	  map[string]*Worker
	WorkerStatus chan u.WorkStatus
}


func NewDispatcher(urlDetail chan u.SpiderOut,logger *log4jzl.Log4jzl)(*Dispatcher){
	contrl:=make(chan u.DispatcherContrl,100)
	workerStatus:=make(chan u.WorkStatus,100)
	this := &Dispatcher{UrlDetail:urlDetail,Logger:logger,Contrl:contrl,WorkerStatus:workerStatus}
	return this
}



func (this *Dispatcher)ConfigureDispatcher(workernum int,plugs map[string]u.PlugInterface,scheduler *s.Scheduler)error{
	this.Workers=make(map[string]*Worker)
	
	for i:=0;i<workernum;i++{
		name:=fmt.Sprintf("%v",i)
		this.Workers[name]=NewWorker(name,plugs,this.WorkerStatus,this.Logger,scheduler)
	}
	return nil
	
}



func (this *Dispatcher)StartDispatcher()error{
	
	this.Logger.Info("[INFO] Start Dispatcher... ")
	u.SpiderSync.Add(1)
	go this.runDispatcher()
	return nil
}




func (this *Dispatcher)runDispatcher() error {
	this.Logger.Info("[INFO] Run Dispatcher ... ")
	defer u.SpiderSync.Done()
	for {
		select {
			case in:= <-this.UrlDetail:
			this.Logger.Info("[INFO] dispatcher recive Url : %v",in.Url)
			ws:= <-this.WorkerStatus
			for ws.Status == 1 || ws.Status == 2{
				ws = <-this.WorkerStatus
			}
			go this.Workers[ws.Name].GoWork(in.Content,in.Url,in.HandlerName)
			case contrl:= <-this.Contrl:
			this.Logger.Info("[INFO] dispatcher recive contrl : %v",contrl)
			return nil
			
		}
	}
	
	return nil
}




