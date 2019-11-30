package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

type User struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func main() {
	router := gin.Default()

	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "tng-pjam2019"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	router.GET("/", func(c *gin.Context) {
		var user User
		user.Name = c.Query("name")
		user.Score, _ = strconv.Atoi(c.Query("score"))

		_, _, err := client.Collection("user").Add(ctx, map[string]interface{}{
			"name":  user.Name,
			"score": user.Score,
		})
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		c.String(http.StatusOK, "Hello World!")
	})

	http.Handle("/", router)
	appengine.Main()
}
