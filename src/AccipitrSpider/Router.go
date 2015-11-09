//	路由分发器
//
package main

import (
	//"encoding/json"
	//"errors"
	"fmt"
	"io"
	"net/http"
	//"regexp"
	//"strconv"
	u "Utils"
)

type Router struct {
	Kvdb *u.PandionKV
}

//路由设置
//数据分发
func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//路由分发设置，用来判断url是否合法，通过配置文件的正则表达式配置
	
	res,_:=this.Kvdb.GetAllKeys()
		
	var result string = "<html><body><table border=\"1\">"
	for k,v := range res{
		kn:=fmt.Sprintf("<tr><th> %v </th></tr>",k)
		result = result + kn
		for _,vv:=range v{
			n:=fmt.Sprintf("<tr><td> %v </td></tr>",vv)
			result = result + n
		}
		
		
		
	}
	result=result+"</table></body></html>"
	//fmt.Printf("%v\n",result)
	io.WriteString(w, result)
	
	
/*
	header := w.Header()
	header.Add("Content-Type", "application/json")
	header.Add("charset", "UTF-8")
	

	resource, version, err := this.ParseURL(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, MakeErrorResult(-1, err.Error()))
	} else {
		processor, err := this.GetProcessor(resource)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, MakeErrorResult(-1, err.Error()))
		} else {
			//获取缓存

			//处理业务逻辑
			result, err := this.Service.Process(r, version, resource, processor)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, result) //MakeErrorResult(-1, err.Error()))
			} else {
				header.Add("Content-Length", fmt.Sprintf("%v",len(result)))
				io.WriteString(w, result)
			}

			//写缓存

		}
	}
*/
	//fmt.Printf("Hello web\n")
	return 
}

