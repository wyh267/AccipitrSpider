package Plugs

import (
	//"fmt"
	"github.com/outmana/log4jzl"
	//"gopkg.in/xmlpath.v1"
	//"strings"
	u "Utils"
	"errors"
	//"fmt"
	//"github.com/moovweb/gokogiri"
	"regexp"
)


type Avid struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewAvid(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *Avid {
	this := &Avid{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *Avid) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	return nil, errors.New("ERROR")
}

func (this *Avid) GetDetailContent(content, url string) (map[string]interface{}, error) {

	/*
	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search("/html/body/table/tbody/tr[1]/td[1]")
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}
	

	for _,title := range title_xpath{
		this.Logger.Info("[AVID : PRODUCT] %v", title.Content())
	}

	*/
	result:=make(map[string]interface{})
	title_pat, err := regexp.Compile("<a id=\"ctl00_listing_link_.*?>(.*?)<")
	//title_pat := regexp.MustCompile("<a id=\"ctl00_listing_link_.*?>(.*?)<")
	if err != nil {
		return nil, err
	}
	matchs := title_pat.FindAllSubmatch([]byte(content),-1)
	if matchs == nil {
		err = errors.New("Wrong URL")
		return nil, err
	}
	
	for _,m := range matchs{
		if len(m)>0{
			this.Logger.Info("[AVID : PRODUCT]  %v",string(m[1]))
			result[string(m[1])]=this.Name
		}
	}


	return result, nil
}

func (this *Avid) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
