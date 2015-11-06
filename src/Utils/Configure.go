package Utils

import (
	"errors"
	"github.com/ewangplay/config"
	"strconv"
)

type Configure struct {
	ConfigureMap map[string]string
}

func NewConfigure(filename string) (*Configure, error) {
	config := &Configure{}

	config.ConfigureMap = make(map[string]string)
	err := config.ParseConfigure(filename)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (this *Configure) loopConfigure(sectionName string, cfg *config.Config) error {

	if cfg.HasSection(sectionName) {
		section, err := cfg.SectionOptions(sectionName)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(sectionName, v)
				if err == nil {
					this.ConfigureMap[v] = options
				}
			}

			return nil
		}
		return errors.New("Parse Error")
	}

	return errors.New("No Section")
}

func (this *Configure) ParseConfigure(filename string) error {
	cfg, err := config.ReadDefault(filename)
	if err != nil {
		return err
	}

	this.loopConfigure("global", cfg)
	

	return nil
}



//数据库连接配置信息
func (this *Configure) GetMysqlUserName() (string, error) {

	mysqlusername, ok := this.ConfigureMap["mysqlusername"]

	if ok == false {
		return "root", errors.New("No mysqlusername,use defualt")
	}

	return mysqlusername, nil
}

func (this *Configure) GetMysqlPassword() (string, error) {

	mysqlpassword, ok := this.ConfigureMap["mysqlpassword"]

	if ok == false {
		return "12345", errors.New("No mysqlpassword,use defualt")
	}

	return mysqlpassword, nil
}

func (this *Configure) GetMysqlHost() (string, error) {

	mysqlhost, ok := this.ConfigureMap["mysqlhost"]

	if ok == false {
		return "127.0.0.1", errors.New("No mysqlhost,use defualt")
	}

	return mysqlhost, nil
}

func (this *Configure) GetMysqlPort() (string, error) {

	mysqlport, ok := this.ConfigureMap["mysqlport"]

	if ok == false {
		return "3306", errors.New("No mysqlport,use defualt")
	}

	return mysqlport, nil
}

func (this *Configure) GetMysqlDBname() (string, error) {

	mysqlDBname, ok := this.ConfigureMap["mysqlDBname"]

	if ok == false {
		return "test", errors.New("No mysqlDBname,use defualt")
	}

	return mysqlDBname, nil
}

func (this *Configure) GetMysqlCharset() (string, error) {

	mysqlcharset, ok := this.ConfigureMap["mysqlcharset"]

	if ok == false {
		return "utf8", errors.New("No mysqlcharset,use defualt")
	}

	return mysqlcharset, nil
}

func (this *Configure) GetMysqlMaxConns() (int, error) {

	mysqlmaxconnsstr, ok := this.ConfigureMap["mysqlmaxconns"]
	if ok == false {
		return 9090, errors.New("No mysqlmaxconns set, use default")
	}

	mysqlmaxconns, err := strconv.Atoi(mysqlmaxconnsstr)
	if err != nil {
		return 2000, err
	}

	return mysqlmaxconns, nil
}

func (this *Configure) GetMysqlMaxIdleConns() (int, error) {

	mysqlmaxidleconnsstr, ok := this.ConfigureMap["mysqlmaxidleconns"]
	if ok == false {
		return 9090, errors.New("No mysqlmaxidleconns set, use default")
	}

	mysqlmaxidleconns, err := strconv.Atoi(mysqlmaxidleconnsstr)
	if err != nil {
		return 1000, err
	}

	return mysqlmaxidleconns, nil
}



func (this *Configure) GetInterval() (int,error){
	Intervalstr, ok := this.ConfigureMap["Interval"]
	if ok == false {
		return 600, errors.New("No Interval set, use default")
	}

	Interval, err := strconv.Atoi(Intervalstr)
	if err != nil {
		return 600, err
	}

	return Interval, nil
}


func (this *Configure) GetKVDB()(string,error) {
	
	KVDBname, ok := this.ConfigureMap["KVDBname"]

	if ok == false {
		return "KVDB", errors.New("No KVDBname,use defualt")
	}

	return KVDBname, nil
	
}