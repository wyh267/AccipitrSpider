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


type MAudio struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewMAudio(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *MAudio {
	this := &MAudio{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *MAudio) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	return nil, errors.New("ERROR")
}

func (this *MAudio) GetDetailContent(content, url string) (map[string]interface{}, error) {

	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search("id(\"content\")/div/div/ul/li[*]/div/div[3]/h6/a")
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}

	for _,title := range title_xpath{
		this.Logger.Info("[MADUDIO : PRODUCT] %v", title.Content())
	}

	return nil, nil
}

func (this *MAudio) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
