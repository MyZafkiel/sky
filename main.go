package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const Link = "https://m.weibo.cn/api/container/getIndex?type=uid&value=6355968578&containerid=1076036355968578"

type Blog struct {
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
	Id        string `json:"id"`
}

type Card struct {
	Mblog *Blog `json:"mblog,omitempty"`
}

type CardInfo struct {
	Containerid string `json:"containerid"`
	VP          int    `json:"v_p"`
	ShowStyle   int    `json:"show_style"`
	Total       int    `json:"total"`
	SinceId     int64  `json:"since_id"`
}

type Resp struct {
	Ok   int `json:"ok"`
	Data struct {
		CardlistInfo *CardInfo `json:"cardlistInfo"`
		Cards        []*Card   `json:"cards"`
	} `json:"data"`
}

func SendMsg(link string) {
	ico := "https://wx4.sinaimg.cn/orj480/006W8ZMely8gxvtmz0zmwj30e80e8jtf.jpg"
	U := fmt.Sprintf("https://api.day.app/nsKRTSwc2jd39ozxkrnTK6/%s/%s?icon=%s&url=%s",
		url.QueryEscape("光遇"),
		url.QueryEscape("有新版本更新，前往查看"),
		ico,
		link,
	)
	http.Get(U)
}

func Handle() error {
	res, err := http.Get(Link)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var body Resp
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		return err
	}
	if body.Ok != 1 {
		return errors.New("状态码不正确")
	}
	for _, card := range body.Data.Cards {
		if card.Mblog != nil {
			t, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", card.Mblog.CreatedAt)
			if err != nil {
				continue
			}
			//now, _ := time.Parse("Mon Jan 02 15:04:05 -0700 2006", "Tue Aug 23 20:00:04 +0800 2022")
			now := time.Now()
			if now.Year() == t.Year() && now.Month() == t.Month() && now.Day() == t.Day() && strings.Contains(card.Mblog.Text, "版本更新") {
				fmt.Println("检测到版本更新")
				SendMsg("https://m.weibo.cn/detail/" + card.Mblog.Id)
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("服务启动")
	c := time.NewTicker(time.Second * 8)
	for true {
		select {
		case <-c.C:
			fmt.Println("任务开始")
			err := Handle()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("任务完成")
		}
	}
}
