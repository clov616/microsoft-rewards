package reward

import (
	"log"
	"os"
	"regexp"
	"strings"
)

type Conf struct {
	KeyWords []string `yaml:"key_words"`
}

type Env struct {
	SetIPs []string `yaml:"set_ips"`
}

var (
	// 默认配置
	keyWords  = []string{"a", "AB", "bingNew", "iaimi", "爱美", "孤独摇滚", "b站", "bilibili", "大佐de苦de手", "柚子社"}
	defaultIP = "[14.102.128.0,2.59.154.0]"
)

// 处理conf
func (c *Conf) Handler() {
	c.KeyWords = keyWords
}

func (e *Env) InitEnv() {
	// 获取setURL
	SetIpStr := os.Getenv("IPS")
	if SetIpStr == "" {
		// 默认刷国区和日区
		log.Println("[Warn]未配置IPS，使用默认IPS")
		SetIpStr = defaultIP
	}
	// 匹配ip
	e.matchIp(SetIpStr)
}

func (e *Env) matchIp(SetIpStr string) {
	pattern, _ := regexp.Compile("\\[([\\s\\S]*)]")
	target := pattern.FindStringSubmatch(SetIpStr)
	tempStr := strings.TrimSpace(target[1]) // 去头尾空格
	e.SetIPs = strings.Split(tempStr, ",")
}
