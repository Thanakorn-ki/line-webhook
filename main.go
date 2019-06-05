package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		panic(err)
	}

	e := echo.New()

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Any("/webhook", func(c echo.Context) error {
		fmt.Println("webhook")
		events, err := bot.ParseRequest(c.Request())
		if err != nil {
			log.Println(err)
			return err
		}

		log.Print(events[0].Message)
		return nil
	})
	e.Logger.Fatal(e.Start(":1323"))

}
