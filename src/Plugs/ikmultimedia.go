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


type Ikmultimedia struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewIkmultimedia(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *Ikmultimedia {
	this := &Ikmultimedia{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *Ikmultimedia) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	return nil, errors.New("ERROR")
}

func (this *Ikmultimedia) GetDetailContent(content, url string) (map[string]interface{}, error) {

	result:=make(map[string]interface{})

	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search("id(\"products-maindiv\")/div[*]/div[*]/div[*]/div/p/a/span")
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}

	for _,title := range title_xpath{
		this.Logger.Info("[Ikmultimedia : PRODUCT] %v", title.Content())
		result[title.Content()]=this.Name
	}

	


	return result, nil
}

func (this *Ikmultimedia) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
