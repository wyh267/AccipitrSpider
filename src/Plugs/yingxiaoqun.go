package Plugs

import (
	//"fmt"
	"github.com/outmana/log4jzl"
	//"gopkg.in/xmlpath.v1"
	//"strings"
	u "Utils"
	"errors"
	"fmt"
	"github.com/moovweb/gokogiri"
	"regexp"
)

const SQL_YXQ string = "replace into reco_contents (cid,content_id,content_type,keywords,title,abstract,content,author,edit_date,url,create_time,last_modify_time,is_delete) values(?,?,?,?,?,?,?,?,?,?,NOW(),NOW(),?)"

type Yingxiaoqun struct {
	Name      string
	Logger    *log4jzl.Log4jzl
	DbAdaptor *u.DBAdaptor
}

func NewYXQ(name string, logger *log4jzl.Log4jzl, dba *u.DBAdaptor) *Yingxiaoqun {
	this := &Yingxiaoqun{Name: name, Logger: logger, DbAdaptor: dba}
	return this
}

func (this *Yingxiaoqun) GetNextUrls(content, url string) ([]u.CrawlData, error) {

	crawls := make([]u.CrawlData, 0)
	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	res, err := doc.Search("//*[@class=\"read-more\"]/@href")
	if err != nil {
		return nil, err
	}
	for i := range res {
		//this.Logger.Info("[INFO]  %v",res[i])
		crawls = append(crawls, u.CrawlData{Url: res[i].String(), Type: u.DETAIL_URL, HandlerName: this.Name})
	}

	return crawls, nil
}

func (this *Yingxiaoqun) GetDetailContent(content, url string) (map[string]interface{}, error) {

	//result:=make(map[string]interface{})
	num_pat, err := regexp.Compile("(\\d+)$")
	if err != nil {
		return nil, err
	}
	matchs := num_pat.FindStringSubmatch(url)
	if matchs == nil {
		err = errors.New("Wrong URL")
		return nil, err
	}

	//this.Logger.Info("[INFO] matchs[1] %v",matchs[1])

	doc, _ := gokogiri.ParseHtml([]byte(content))
	defer doc.Free()
	title_xpath, err := doc.Search(fmt.Sprintf("//*[@id=\"post-%v\"]/div/div[2]/header/h1", matchs[1]))
	if len(title_xpath) == 0 || err != nil {
		return nil, err
	}
	title := title_xpath[0].Content()

	editdate_xpath, err := doc.Search(fmt.Sprintf("//*[@id=\"post-%s\"]/div/div[1]/div/time/@datetime", matchs[1]))
	if len(editdate_xpath) == 0 || err != nil {
		return nil, err
	}
	editdate := editdate_xpath[0].String()

	content_xpath, err := doc.Search(fmt.Sprintf("//*[@id=\"post-%s\"]/div/div[2]/div", matchs[1]))
	if len(content_xpath) == 0 || err != nil {

		return nil, err
	}
	contents := content_xpath[0].Content()

	editor_xpath, err := doc.Search(fmt.Sprintf("//*[@id=\"post-%s\"]/div/div[2]/header/span/span/a", matchs[1]))
	if len(editor_xpath) == 0 || err != nil {
		return nil, err
	}
	editor := editor_xpath[0].Content()

	this.Logger.Info("[INFO] Title %v", title)
	this.Logger.Info("[INFO] EditDate %v", editdate)
	this.Logger.Info("[INFO] Content %v", contents)
	this.Logger.Info("[INFO] Editor %v", editor)
	/*
		err = this.DbAdaptor.ExecFormat(SQL_YXQ,"111",matchs[1],"2","",title,"",contents,editor,editdate[:10],url,"0")
		if err != nil {
			this.Logger.Error(" MYSQL Error : %v ",err)
			return nil,err
		}
	*/
	return nil, nil
}

func (this *Yingxiaoqun) SaveDetailContent(details map[string]interface{}) error {

	return nil
}
