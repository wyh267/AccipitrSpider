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


type Radialeng struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewRadialeng(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *Radialeng {
	this := &Radialeng{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *Radialeng) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	return nil, errors.New("ERROR")
}

func (this *Radialeng) GetDetailContent(content, url string) (map[string]interface{}, error) {

	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search("//*[@id=\"reviewBlock\"]/div[*]/span")
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}

	for _,title := range title_xpath{
		this.Logger.Info("[RADIALENG : PRODUCT] %v", title.Content())
	}

	return nil, nil
}

func (this *Radialeng) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
