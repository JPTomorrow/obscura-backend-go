package main

import (
	"encoding/json"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{

		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// init youtube service
	// this includes a running goroutine that pulls new videos every 8 hours
	yt := NewYoutubeService()
	videosPerInteval := 5
	yt.LoadDatabaseVideos()
	err := yt.PullNewVideos(int64(videosPerInteval))
	if err != nil {
		log.Println("Error pulling new videos -> ", err)
	}
	log.Println("Pulling new videos -> count: ", len(yt.videoPool.videos))

	go func() {
		for range time.Tick(time.Hour * 8) { // run every 8 hours
			err := yt.PullNewVideos(int64(videosPerInteval))
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
	e.POST("/upvote", UpvoteVideo)
	e.POST("/downvote", DownvoteVideo)

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

func UpvoteVideo(c echo.Context) error {
	ac := c.(*YoutubeContext)
	bd := bodyAsJSON(c)
	videoId, idOk := bd["videoId"].(string)
	if !idOk {
		return c.String(http.StatusBadRequest, "videoId is not valid. must be a string.")
	}

	err := ac.yt.VoteVideo(videoId, true)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "successfully upvoted video")
}
func DownvoteVideo(c echo.Context) error {
	ac := c.(*YoutubeContext)
	bd := bodyAsJSON(c)
	videoId, idOk := bd["videoId"].(string)
	if !idOk {
		return c.String(http.StatusBadRequest, "videoId is not valid. must be a string.")
	}

	err := ac.yt.VoteVideo(videoId, false)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "successfully upvoted video")
}

/*
	UTILITY FUNCTIONS
*/

// update the database with the tables defined by the structs provided to db.SchemaString
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

func bodyAsJSON(c echo.Context) map[string]interface{} {

	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {

		log.Fatalln("empty json body")
		return nil
	}

	return jsonBody
}
