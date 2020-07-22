//解析配置参数
package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	WebHook		 string	  `yaml:"webhook"`
	Pid          string   `yaml:"pid"`
	SnapPath     string   `yaml:"snap_path"`
	AlterLimit   float64  `yaml:"alter_limit"`
	Interval     int      `yaml:"interval"`
	FromMail     string   `yaml:"fromMail"`
	FromMailHost string   `yaml:"fromMailHost"`
	FromMailPass string   `yaml:"fromMailPass"`
	FromMailPort string   `yaml:"fromMailPort"`
	ToMail       []string `yaml:"toMail"`
	Files        []string	`yaml:"files"`
	Errfiles	 []string	`yaml:"errfiles"'`
	ProjectName    string	`yaml:"projectname"`
	ExceptionKeywords []string	`yaml:"exceptionkeywords"`
	Processes		[]string	`yaml:"processes"`
	Send 			int 		`yaml:"send"`
	Recv			int 		`yaml:"recv"`
}

var conf Config

func InitConfig(configPath string) error {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return errors.Wrap(err, "Read config file failed")
	}

	if err = yaml.Unmarshal(configFile, &conf); err != nil {
		return errors.Wrap(err, "Unmarshal config file failed.")
	}
	return nil
}

func GetConfig() Config {
	return conf
}
