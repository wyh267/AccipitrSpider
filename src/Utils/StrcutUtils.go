

package Utils


import (
	"sync"
)

var SpiderSync sync.WaitGroup

const SEED_URL int64 = 1
const DETAIL_URL int64 = 2
const LIST_URL int64 = 3


type CrawlData struct{
	Url string
	Type int64
	HandlerName string
}





type SpiderIn struct{
	Url 	string
	Timeout int64
	HandlerName string
}




type SpiderOut struct{
	Content  string
	Header	 string
	Url		 string
	HandlerName string
}





type SpiderContrl struct{
	ContrlCMD	string
}


type DispatcherContrl struct{
	ContrlCMD	string
}


type PlugInterface interface{
	
	 GetNextUrls(content,url string)([]CrawlData,error)
	 GetDetailContent(content,url string)(map[string]interface{},error)
	 SaveDetailContent(details map[string]interface{}) error
	
}







type WorkStatus struct{
	Name	string
	Status  int64 // 0:OK 1:BUSY 2:ERROR
}