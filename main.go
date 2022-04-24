package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"catknock/infrastructure"
	"catknock/model"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ここをキャンプ地とする",
	})
}

// https://github.com/line/line-bot-sdk-go/blob/master/examples/echo_bot/server.go
func callbackHandler(c *gin.Context) {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.Writer.WriteHeader(400)
		} else {
			c.Writer.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}

}

// https://qiita.com/lanevok/items/dbf591a3916070fcba0d
func usersHandler(c *gin.Context) {
	dbMap := infrastructure.GetDb()
	var users []model.User
	_, err := dbMap.Select(&users, `SELECT id, name, age FROM users`)
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "users.tmpl", gin.H{
		"users": users,
	})
}

func main() {
	infrastructure.LoadEnv()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")
	router.GET("/ping", pingHandler)
	router.GET("/users", usersHandler)
	router.POST("/callback", callbackHandler)
	router.Run()
}
