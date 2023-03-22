package reward

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type View struct {
	Url    string
	Infov  *Infov
	client *http.Client
	Lang   string // 国家
}

// 获取到的分数
type Infov struct {
	AvailablePoints int
	DailyPoints     DailyPoint
	MobiSearch      MobiSearch
	PcSearch        PcSearch
}

type DailyPoint struct {
	PointProgressMax int `json:"pointProgressMax"`
	PointProgress    int `json:"pointProgress"`
	AvailablePoints  int
}

type MobiSearch struct {
	PointMax      int `json:"pointProgressMax"`
	PointProgress int `json:"pointProgress"`
}

type PcSearch struct {
	PointMax      int `json:"pointProgressMax"`
	PointProgress int `json:"pointProgress"`
}

func (v *View) doGet(c *Conn) *http.Response {
	v.client = &http.Client{}
	req, err := http.NewRequest("GET", v.Url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// 向req添加Cookies
	for _, v := range c.Cookie.Cookies {
		req.AddCookie(v)
	}
	//配置真实ip
	req.Header.Set("x-forwarded-for", c.SetIP)
	resp, err := v.client.Do(req)
	if err != nil {
		log.Println(err)
		log.Println("未正确配置代理")
		time.Sleep(time.Second * 5)
		os.Exit(400)
	}
	// 测试使用 临时
	//file, _ := os.OpenFile("./resp.html", os.O_CREATE|os.O_RDWR, 0777)
	//respStr, _ := ioutil.ReadAll(resp.Body)
	//file.Write(respStr)
	return resp
}

func (v *View) Handler(c *Conn) {
	v.Infov = new(Infov)
	resp := v.doGet(c)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("View Get != 200")
		time.Sleep(time.Second)
		v.Handler(c) // 递归执行
	} else {
		respB, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		respStr := string(respB)

		// TODO: 获取到国家（区域）信息
		patternlang, _ := regexp.Compile("lang=\"([\\s\\S]*?)\"")
		langL := patternlang.FindStringSubmatch(respStr)
		if len(langL) == 0 {
			log.Fatalln("长度为0", "json格式化失败(疑似Cookie失效)")
		}
		c.View.Lang = langL[1]

		// 测试使用 临时
		//file, _ := os.OpenFile("./resp.html", os.O_CREATE|os.O_RDWR, 0777)
		//file.Write([]byte(respStr))
		pattern1, _ := regexp.Compile("\"availablePoints\":([\\d]*),")
		targetMark := pattern1.FindStringSubmatch(respStr)
		if len(targetMark) == 0 {
			log.Fatalln("长度为0", "json格式化失败(疑似Cookie失效)")
		}
		point, err := strconv.Atoi(targetMark[1])
		if err != nil {
			log.Fatalln(err, "json格式化失败(疑似Cookie失效)")
		}

		v.Infov.AvailablePoints = point

		pattern2, _ := regexp.Compile("\"dailyPoint\":\\[([\\s\\S]*?)]")
		dailyPointStr := pattern2.FindStringSubmatch(respStr)
		if len(dailyPointStr) == 0 {
			log.Fatalln("长度为0", "json格式化失败(疑似Cookie失效)")
		}
		err = json.Unmarshal([]byte(dailyPointStr[1]), &v.Infov.DailyPoints)
		patternm, _ := regexp.Compile("\"mobileSearch\":\\[([\\s\\S]*?)]")
		mobileSearchstr := patternm.FindStringSubmatch(respStr)
		// #issue1解决新用户不存在手机分数
		if len(mobileSearchstr) != 0 {
			err = json.Unmarshal([]byte(mobileSearchstr[1]), &v.Infov.MobiSearch)
		}
		pattern, _ := regexp.Compile("\"pcSearch\":\\[([\\s\\S]*?}),\\{")
		pcSearchstr := pattern.FindStringSubmatch(respStr)
		if len(pcSearchstr) == 0 {
			log.Fatalln("长度为0", "json格式化失败(疑似Cookie失效)")
		}
		err = json.Unmarshal([]byte(pcSearchstr[1]), &v.Infov.PcSearch)
		if err != nil {
			log.Fatalln(err, "json格式化失败(疑似Cookie失效)")
		}
	}

}
