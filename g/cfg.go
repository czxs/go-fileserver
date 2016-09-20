package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"os"
	"sync"
	"time"
)

type SysFileInfo struct {
	FName  string
	FSize  int64
	FMtime time.Time
	FPerm  os.FileMode
	FMd5   string
	FType  bool
	FPath  string
}

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
}

type RabbitmqConfig struct {
	S_addr      string `json:"s_addr"`
	User        string `json:"user"`
	Pass        string `json:"pass"`
	Enable      bool   `json:"enable"`
	Exchange    string `json:"exchange"`
	Queue       string `json:"queue"`
	RouteingKey string `json:"routing_key"`
}

type RClientConfig struct {
	Agent         []string `json:"agent"`
	Enable        bool     `json:"enable"`
	Agent_Process int      `json:"agent_process"`
}

type RServerConfig struct {
	Enable  bool   `json:"enable"`
	Listen  string `json:"listen"`
	Dirpath string `json:"dirpath"`
	Proxy   bool   `json:"proxy"`
}

type ReciverConfig struct {
	Client *RClientConfig `json:"client"`
	Server *RServerConfig `json:"server"`
}

type GlobalConfig struct {
	Debug    bool            `json:"debug"`
	Rabbitmq *RabbitmqConfig `json:"rabbitmq"`
	Reciver  *ReciverConfig  `json:"reciver"`
	Http     *HttpConfig     `json:"http"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

/*func Hostname() (string, error) {
	hostname := Config().Hostname
	if hostname != "" {
		return hostname, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println("ERROR: os.Hostname() fail", err)
	}
	return hostname, err
}

func IP() string {
	ip := Config().IP
	if ip != "" {
		// use ip in configuration
		return ip
	}

	if len(LocalIps) > 0 {
		ip = LocalIps[0]
	}

	return ip
}
*/
func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)

	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
