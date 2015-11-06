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


type Jblpro struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewJblpro(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *Jblpro {
	this := &Jblpro{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *Jblpro) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	return nil, errors.New("ERROR")
}

func (this *Jblpro) GetDetailContent(content, url string) (map[string]interface{}, error) {

	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search("//*[@id=\"contentPlaceholder_C002_Col00\"]/div[*]/div/ul/li/a")
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}

	for _,title := range title_xpath{
		this.Logger.Info("[JBLPRO : PRODUCT] %v", title.Content())
	}

	products_xpath2, err := doc.Search("//*[@id=\"contentPlaceholder_C002_Col01\"]/div[*]/div/ul/li/a")
	if len(products_xpath2) == 0 || err != nil {
		return nil, err
	}

	for _,title := range products_xpath2{
		this.Logger.Info("[JBLPRO : PRODUCT] %v", title.Content())
	}
	
	products_xpath3, err := doc.Search("//*[@id=\"contentPlaceholder_C002_Col02\"]/div[*]/div/ul/li/a")
	if len(products_xpath3) == 0 || err != nil {
		return nil, err
	}

	for _,title := range products_xpath3{
		this.Logger.Info("[JBLPRO : PRODUCT] %v", title.Content())
	}


	return nil, nil
}

func (this *Jblpro) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
