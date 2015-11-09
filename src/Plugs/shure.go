package Plugs

import (
	//"fmt"
	"github.com/outmana/log4jzl"
	//"gopkg.in/xmlpath.v1"
	//"strings"
	u "Utils"
	"errors"
	//"fmt"
	"github.com/moovweb/gokogiri"
	//"regexp"
)


type Shure struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewShure(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *Shure {
	this := &Shure{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *Shure) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	return nil, errors.New("ERROR")
}

func (this *Shure) GetDetailContent(content, url string) (map[string]interface{}, error) {

	result:=make(map[string]interface{})
	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search("id(\"main\")/div/div[2]/div/div[1]/div[2]/div[*]/div[*]/article/div[2]/h1/a")
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}

	for _,title := range title_xpath{
		this.Logger.Info("[Shure : PRODUCT] %v", title.Content())
		result[title.Content()]=this.Name
	}

	


	return result, nil
}

func (this *Shure) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
