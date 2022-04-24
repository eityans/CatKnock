package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

type C struct {
  Id int
  Name string
}

func testHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "test.tmpl", gin.H{
		"a": "a",
		"b": []string{"b_todo1","b_todo2"},
		"c": []C{{1,"c_mika"},{2,"c_risa"}},
		"d": C{3,"d_mayu"},
		"e": true,
		"f": false,
		"h": true,
	})
}

func loadEnv() {
	if os.Getenv("ENV") == "production" { return }
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {
	loadEnv()

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.Exec(`INSERT INTO users(name, age) VALUES($1, $2)`, "Bob", 18)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	rows, err := db.Query(`SELECT id, name, age FROM users ORDER BY NAME`)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var name string
		var age int64
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, age)
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")
	router.GET("/ping", pingHandler)
	router.GET("/test", testHandler)
	router.POST("/callback", callbackHandler)
	router.Run()
}
