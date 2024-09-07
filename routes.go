package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/JPTomorrow/obscura/config"
	"github.com/JPTomorrow/obscura/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
This is where all the the routes are defined for the API.
*/

type YoutubeContext struct {
	echo.Context
	yt *YTService
}

func initRoutes(e *echo.Echo) {
	// middleware
	// e.Use(middleware.CORS())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{

		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	yt := NewYoutubeService()
	yt.LoadDatabaseVideos()
	err := yt.PullNewVideos(10)
	if err != nil {
		log.Println("Error pulling new videos -> ", err)
	}
	log.Println("Pulling new videos -> count: ", len(yt.videoPool.videos))

	go func() {
		for range time.Tick(time.Second * 1800) {
			err := yt.PullNewVideos(10)
			if err != nil {
				log.Println("Error pulling new videos -> ", err)
			}
			log.Println("Pulling new videos -> count: ", len(yt.videoPool.videos))
		}
	}()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &YoutubeContext{
				c,
				yt,
			}
			return next(cc)
		}
	})
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// backend routes
	e.GET("/next-vid", FeedNextVideo)

	// debug only routes
	if config.DEBUG {
		e.GET("/update-schema", updateDbSchema)
	}
}

func FeedNextVideo(c echo.Context) error {
	ac := c.(*YoutubeContext)
	vid, err := ac.yt.PitchRandomVideo()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, vid)
}

func updateDbSchema(c echo.Context) error {
	schema, err := db.SchemaString(db.YoutubeVideo{})
	for _, table := range schema {

		if err != nil {
			return c.String(http.StatusBadRequest, table+"\n\n"+err.Error())
		} else {
			_, err := db.Exec(table)
			if err != nil {
				return c.String(http.StatusBadRequest, table+"\n\n"+err.Error())
			}

		}
	}

	msg := strings.Join(schema, "\n\n") + "\nTables created successfully!!!\n\n"
	fmt.Println(msg)
	return c.String(http.StatusOK, msg)
}
