package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

type botLINE struct {
	*linebot.Client
}

func initBotLINE() (*botLINE, error) {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &botLINE{
		bot,
	}, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	bot, err := initBotLINE()
	if err != nil {
		panic(err)
	}

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/webhook", func(c echo.Context) error {
		fmt.Println("webhook")
		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			log.Println(err)
			return err
		}

		for _, event := range events {
			log.Printf("Got event %+v", event.Type)
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if err := bot.handleText(message, event.ReplyToken, event.Source); err != nil {
						log.Print(err)
					}
				// case *linebot.ImageMessage:
				// 	if err := app.handleImage(message, event.ReplyToken); err != nil {
				// 		log.Print(err)
				// 	}
				// case *linebot.VideoMessage:
				// 	if err := app.handleVideo(message, event.ReplyToken); err != nil {
				// 		log.Print(err)
				// 	}
				// case *linebot.AudioMessage:
				// 	if err := app.handleAudio(message, event.ReplyToken); err != nil {
				// 		log.Print(err)
				// 	}
				// case *linebot.FileMessage:
				// 	if err := app.handleFile(message, event.ReplyToken); err != nil {
				// 		log.Print(err)
				// 	}
				// case *linebot.LocationMessage:
				// 	if err := app.handleLocation(message, event.ReplyToken); err != nil {
				// 		log.Print(err)
				// 	}
				// case *linebot.StickerMessage:
				// 	if err := app.handleSticker(message, event.ReplyToken); err != nil {
				// 		log.Print(err)
				// 	}
				default:
					log.Printf("Unknown message: %+v", message.Message)
				}
			// case linebot.EventTypeFollow:
			// 	if err := app.replyText(event.ReplyToken, "Got followed event"); err != nil {
			// 		log.Print(err)
			// 	}
			// case linebot.EventTypeUnfollow:
			// 	log.Printf("Unfollowed this bot: %v", event)
			// case linebot.EventTypeJoin:
			// 	if err := app.replyText(event.ReplyToken, "Joined "+string(event.Source.Type)); err != nil {
			// 		log.Print(err)
			// 	}
			// case linebot.EventTypeLeave:
			// 	log.Printf("Left: %v", event)
			// case linebot.EventTypePostback:
			// 	data := event.Postback.Data
			// 	if data == "DATE" || data == "TIME" || data == "DATETIME" {
			// 		data += fmt.Sprintf("(%v)", *event.Postback.Params)
			// 	}
			// 	if err := app.replyText(event.ReplyToken, "Got postback: "+data); err != nil {
			// 		log.Print(err)
			// 	}
			// case linebot.EventTypeBeacon:
			// 	if err := app.replyText(event.ReplyToken, "Got beacon: "+event.Beacon.Hwid); err != nil {
			// 		log.Print(err)
			// 	}
			default:
				log.Printf("Unknown event: %v", event)
			}
		}

		return nil
	})
	e.Logger.Fatal(e.Start(":1323"))

}

func (b botLINE) replyText(replyToken, text string) error {
	if _, err := b.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}

func (b botLINE) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "profile":
		if source.UserID != "" {
			profile, err := b.GetProfile(source.UserID).Do()
			if err != nil {
				return b.replyText(replyToken, err.Error())
			}

			// b.Multicast(
			// 	[]string{"sdsad", "sdas"},
			// 	linebot.NewTextMessage("Display name: "+profile.DisplayName),
			// 	linebot.NewTextMessage("Status message: "+profile.StatusMessage),
			// )

			if _, err := b.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("Display name: "+profile.DisplayName),
				linebot.NewTextMessage("Status message: "+profile.StatusMessage),
			).Do(); err != nil {
				return err
			}

		} else {
			return b.replyText(replyToken, "Bot can't use profile API without user ID")
		}
	case "flex":
		log.Println("flex")
		container, err := linebot.UnmarshalFlexMessageJSON([]byte(carousel()))
		if err != nil {
			log.Print(err)
			return err
		}
		log.Printf("err %+v", err)
		log.Printf("flex %+v", container)

		if _, err := b.ReplyMessage(
			replyToken,
			linebot.NewFlexMessage("AUM", container),
		).Do(); err != nil {
			return err
		}
	default:
		return b.replyText(replyToken, sendFlex())
	}

	return nil
}

func carousel() string {
	return `{
		"type": "carousel",
		"contents": [
		  {
			"type": "bubble",
			"hero": {
			  "type": "image",
			  "size": "full",
			  "aspectRatio": "20:13",
			  "aspectMode": "cover",
			  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_5_carousel.png"
			},
			"body": {
			  "type": "box",
			  "layout": "vertical",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "text",
				  "text": "Arm Chair, White",
				  "wrap": true,
				  "weight": "bold",
				  "size": "xl"
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "contents": [
					{
					  "type": "text",
					  "text": "$49",
					  "wrap": true,
					  "weight": "bold",
					  "size": "xl",
					  "flex": 0
					},
					{
					  "type": "text",
					  "text": ".99",
					  "wrap": true,
					  "weight": "bold",
					  "size": "sm",
					  "flex": 0
					}
				  ]
				}
			  ]
			},
			"footer": {
			  "type": "box",
			  "layout": "vertical",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "button",
				  "style": "primary",
				  "action": {
					"type": "uri",
					"label": "Add to Cart",
					"uri": "https://linecorp.com"
				  }
				},
				{
				  "type": "button",
				  "action": {
					"type": "uri",
					"label": "Add to wishlist",
					"uri": "https://linecorp.com"
				  }
				}
			  ]
			}
		  },
		  {
			"type": "bubble",
			"hero": {
			  "type": "image",
			  "size": "full",
			  "aspectRatio": "20:13",
			  "aspectMode": "cover",
			  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_6_carousel.png"
			},
			"body": {
			  "type": "box",
			  "layout": "vertical",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "text",
				  "text": "Metal Desk Lamp",
				  "wrap": true,
				  "weight": "bold",
				  "size": "xl"
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "flex": 1,
				  "contents": [
					{
					  "type": "text",
					  "text": "$11",
					  "wrap": true,
					  "weight": "bold",
					  "size": "xl",
					  "flex": 0
					},
					{
					  "type": "text",
					  "text": ".99",
					  "wrap": true,
					  "weight": "bold",
					  "size": "sm",
					  "flex": 0
					}
				  ]
				},
				{
				  "type": "text",
				  "text": "Temporarily out of stock",
				  "wrap": true,
				  "size": "xxs",
				  "margin": "md",
				  "color": "#ff5551",
				  "flex": 0
				}
			  ]
			},
			"footer": {
			  "type": "box",
			  "layout": "vertical",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "button",
				  "flex": 2,
				  "style": "primary",
				  "color": "#aaaaaa",
				  "action": {
					"type": "uri",
					"label": "Add to Cart",
					"uri": "https://linecorp.com"
				  }
				},
				{
				  "type": "button",
				  "action": {
					"type": "uri",
					"label": "Add to wish list",
					"uri": "https://linecorp.com"
				  }
				}
			  ]
			}
		  },
		  {
			"type": "bubble",
			"body": {
			  "type": "box",
			  "layout": "vertical",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "button",
				  "flex": 1,
				  "gravity": "center",
				  "action": {
					"type": "uri",
					"label": "See more",
					"uri": "https://linecorp.com"
				  }
				}
			  ]
			}
		  }
		]
	  }`
}
func sendFlex() string {
	return `
	{
		"type": "bubble",
		"hero": {
		  "type": "image",
		  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/01_3_movie.png",
		  "size": "full",
		  "aspectRatio": "20:13",
		  "aspectMode": "cover",
		  "action": {
			"type": "uri",
			"uri": "http://linecorp.com/"
		  }
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "spacing": "md",
		  "contents": [
			{
			  "type": "text",
			  "text": "BROWN'S ADVENTURE\nIN MOVIE",
			  "wrap": true,
			  "weight": "bold",
			  "gravity": "center",
			  "size": "xl"
			},
			{
			  "type": "box",
			  "layout": "baseline",
			  "margin": "md",
			  "contents": [
				{
				  "type": "icon",
				  "size": "sm",
				  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
				},
				{
				  "type": "icon",
				  "size": "sm",
				  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
				},
				{
				  "type": "icon",
				  "size": "sm",
				  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
				},
				{
				  "type": "icon",
				  "size": "sm",
				  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
				},
				{
				  "type": "icon",
				  "size": "sm",
				  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gray_star_28.png"
				},
				{
				  "type": "text",
				  "text": "4.0",
				  "size": "sm",
				  "color": "#999999",
				  "margin": "md",
				  "flex": 0
				}
			  ]
			},
			{
			  "type": "box",
			  "layout": "vertical",
			  "margin": "lg",
			  "spacing": "sm",
			  "contents": [
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "Date",
					  "color": "#aaaaaa",
					  "size": "sm",
					  "flex": 1
					},
					{
					  "type": "text",
					  "text": "Monday 25, 9:00PM",
					  "wrap": true,
					  "size": "sm",
					  "color": "#666666",
					  "flex": 4
					}
				  ]
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "Place",
					  "color": "#aaaaaa",
					  "size": "sm",
					  "flex": 1
					},
					{
					  "type": "text",
					  "text": "7 Floor, No.3",
					  "wrap": true,
					  "color": "#666666",
					  "size": "sm",
					  "flex": 4
					}
				  ]
				},
				{
				  "type": "box",
				  "layout": "baseline",
				  "spacing": "sm",
				  "contents": [
					{
					  "type": "text",
					  "text": "Seats",
					  "color": "#aaaaaa",
					  "size": "sm",
					  "flex": 1
					},
					{
					  "type": "text",
					  "text": "C Row, 18 Seat",
					  "wrap": true,
					  "color": "#666666",
					  "size": "sm",
					  "flex": 4
					}
				  ]
				}
			  ]
			},
			{
			  "type": "box",
			  "layout": "vertical",
			  "margin": "xxl",
			  "contents": [
				{
				  "type": "spacer"
				},
				{
				  "type": "image",
				  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/linecorp_code_withborder.png",
				  "aspectMode": "cover",
				  "size": "xl"
				},
				{
				  "type": "text",
				  "text": "You can enter the theater by using this code instead of a ticket",
				  "color": "#aaaaaa",
				  "wrap": true,
				  "margin": "xxl",
				  "size": "xs"
				}
			  ]
			}
		  ]
		}
	  }`
}
