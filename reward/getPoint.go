package reward

import (
	"github.com/levigross/grequests"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

type UrlGet string
type UaPc string
type UaMb string
type TypeUa string

type Get struct {
	Url  UrlGet
	Info Infog
	//client *http.Client
	UApc UaPc
	UAmb UaMb
	RO   grequests.RequestOptions
}

// 请求后返回的信息

type Infog struct {
	Pc string
	mb string
}

// 发起请求 _type string: UA头的类型 (pc mb)表示电脑或手机 *grequests.Response *http.Response
func (g Get) do(c *Conn, _type string) *grequests.Response {
	g.RO = grequests.RequestOptions{
		Cookies: c.Cookie.Cookies,
		Headers: map[string]string{
			"x-forwarded-for": c.SetIP,
		},
	}
	// 组合搜索url
	if len(c.Conf.KeyWords) == 0 {
		log.Fatalln("c.Conf.KeyWords == 0", "请配置KeyWords或删除conf.yaml文件(重置配置)")
	}
	rand.Seed(time.Now().Unix()) // 设置随机数种子
	keyword := c.Conf.KeyWords[rand.Intn(len(c.Conf.KeyWords))] + strconv.Itoa(rand.Intn(10000))
	url := string(g.Url) + "?q=" + url.QueryEscape(keyword)

	if _type == "mb" {
		g.RO.UserAgent = string(g.UAmb)
	} else {
		g.RO.UserAgent = string(g.UApc)
	}
	resp, err := grequests.Get(url, &g.RO)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

// _type 为请求的UA类型: pc or mb(mobilePhone)
func (g *Get) Handler(c *Conn, searchUrl UrlGet, UApc UaPc, UAmb UaMb, _type TypeUa) {
	g.Url = searchUrl
	g.UApc = UApc
	g.UAmb = UAmb
	resp := g.do(c, string(_type))
	defer resp.Close()
	// 刷取积分错误相关
	if *c.NF >= 5 {
		// 错误信息打印
		log.Println("[Error]当前区域获取积分失败: " + c.View.Lang)
		// 停止刷分
		c.manager.StopSend <- true
	}
	if c.View.Infov.AvailablePoints == *c.PrePoint {
		*c.NF += 1
	} else {
		*c.NF = 0
	}
	*c.PrePoint = c.View.Infov.AvailablePoints
	// 信息打印
	if resp.StatusCode == 200 {
		log.Println("<"+_type+"> ", "200 OK")
		log.Println("当前分数:", c.View.Infov.AvailablePoints)
	} else {
		log.Println("bad response", "code: ", resp.StatusCode)
	}
}
