package main

import (
	"fmt"
	"github.com/AMySelf/Microsoft/reward_Remote/reward"
	"log"
	"time"
)

func main() {
	env := reward.Env{}
	env.InitEnv()
	for _, v := range env.SetIPs {
		startOne(v)
		time.Sleep(time.Second * 2)
	}
}

func startOne(setIp string) {
	ViewUrl := "https://rewards.bing.com/"
	conn := reward.New(ViewUrl)
	// 设置刷分地区ip
	conn.SetIP = setIp
	conn.View.Handler(conn)
	fmt.Println("[Info]开始获取积分")
	fmt.Println("当前国家信息: " + conn.View.Lang)
	fmt.Println("当前可用分数: ", conn.View.Infov.AvailablePoints)
	fmt.Println("今日可获取最大分数: ", conn.View.Infov.DailyPoints.PointProgressMax)
	fmt.Println("今日分数: ", conn.View.Infov.DailyPoints.PointProgress)

	// 初始化任务管理器
	manager := conn.NewManager()
	params := reward.Params{
		Conn:   conn,
		UrlGet: "https://www.bing.com/search",
		UaPc:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.51",
		UaMb:   "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
	}
	// init任务管理器处理
	manager.Handler(params)
	// goroutine
	go manager.AddTask(conn.Get.Handler)
	go manager.StartTask()
	func() {
		statusPc, statusMb := 0, 0
		for i := range manager.DoneIndex {
			fmt.Println("Task: ", i)
			conn.View.Handler(conn)
			pcSearch := conn.View.Infov.PcSearch
			mobiSearch := conn.View.Infov.MobiSearch
			if statusPc == 0 && pcSearch.PointProgress == pcSearch.PointMax {
				log.Println("Pc分数刷取完毕")
				statusPc = 1
			}
			if statusMb == 0 && mobiSearch.PointProgress == mobiSearch.PointMax {
				log.Println("手机分数刷取完毕")
				statusMb = 1
			}
		}
		fmt.Println("获取积分完毕！！")
		conn.View.Handler(conn)
		fmt.Println("当前可用分数: ", conn.View.Infov.AvailablePoints)
		fmt.Println("今日可获取最大分数: ", conn.View.Infov.DailyPoints.PointProgressMax)
		fmt.Println("今日分数: ", conn.View.Infov.DailyPoints.PointProgress)
	}()
}
