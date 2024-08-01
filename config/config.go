package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var FileName string = "map_lan_ip.yml"

type Config struct {
	Interval int      `yaml:"interval"` // 查询间隔
	CIDRs    []string `yaml:"cidrs"`    // 查询网段
	Ports    []Ports  `yaml:"ports"`    // 端口
	PingAddr []string `yaml:"pingaddr"` // fping的目标地址
	Target   string   `yaml:"target"`   // Prometheus 动态加载fd文件路径
	LogFile  string   `yaml:"logfile"`  // 日志文件保存地址
}

type Ports struct {
	Prot int    `yaml:"port"` // 查询端口
	Type string `yaml:"type"` // 端口类型，如：node_exporter fping
}

func init() {
	_, err := os.Stat(FileName)
	if os.IsNotExist(err) {
		config := Config{
			Interval: 60,
			CIDRs:    []string{"192.168.0.0/16"},
			Ports:    []Ports{{Prot: 9800, Type: "node_exporter"}, {Prot: 9605, Type: "fping"}},
			PingAddr: []string{"www.baidu.com"},
			Target:   "targets.yml",
			LogFile:  "map_lan_ip.log",
		}
		err = WriteConfig(config)
		if err != nil {
			log.Panic("初始化配置文件失败:", err.Error())
		}
	}
}

func ReadConfig() (Config, error) {
	data, err := os.ReadFile(FileName)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func WriteConfig(c Config) error {
	data, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	err = os.WriteFile(FileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
