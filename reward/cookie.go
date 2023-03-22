package reward

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Cookie struct {
	Cookies   []*http.Cookie
	cookieStr string
	cookieL   []string
}

func (c *Cookie) txt2Cookies() {
	cookieL := make([]string, 10)
	cookieL = strings.Split(c.cookieStr, "; ")
	c.cookieL = cookieL
	for _, v := range c.cookieL {
		tmpL := ownSplit(v, "=")
		c.Cookies = append(c.Cookies, &(http.Cookie{Name: url.QueryEscape(tmpL[0]), Value: url.QueryEscape(tmpL[1])}))
	}
}

// cookie切片规则
func ownSplit(preStr string, pattern string) (preL []string) {
	firstend := -1
	for i, v := range preStr {
		if string(v) == pattern {
			firstend = i
			break
		}
	}
	preL = make([]string, 10)
	preL[0] = preStr[0:firstend]
	preL[1] = preStr[firstend+1:]
	return
}

// 获取环境变量中的cookie 如果获取为空返回false
func (c *Cookie) getCookie() bool {
	secret_cookie := os.Getenv("MY_COOKIE")
	if secret_cookie == "" {
		log.Println("[Error]环境变量(secret)未正确设置")
		log.Println("[Error]MY_COOKIE is nil")
		return false
	}
	c.cookieStr = strings.TrimSpace(secret_cookie)
	return true
}

func (c *Cookie) Handler() {
	if !c.getCookie() {
		log.Println("[Error]cookie为空,请配置cookie")
		time.Sleep(time.Second * 5)
		os.Exit(400)
	}
	c.txt2Cookies()
}
