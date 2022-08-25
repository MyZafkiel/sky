package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
		if card.Mblog != nil && strings.Contains(card.Mblog.Text, "版本更新") {
			link := "https://m.weibo.cn/detail/" + card.Mblog.Id
			//t, err := time.Parse(card.Mblog.CreatedAt, "Tue Aug 23 20:00:04 +0800 2022")
			//if err != nil {
			//	return err
			//}

			fmt.Println(link)
		}
	}
	return nil
}

func main() {
	//err := Handle()

	//if err != nil {
	//	panic(err)
	//}
}
