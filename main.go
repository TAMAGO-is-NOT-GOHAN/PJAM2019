package main

import (
	"context"
	"log"
  "fmt"
	"net/http"

	firebase "firebase.google.com/go"
  "google.golang.org/api/iterator"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
  "sort"
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

	router.GET("/getRanking", func(c *gin.Context) {
    iter := client.Collection("user").Documents(ctx)
    var resultData []User

    for {
      doc, err := iter.Next()
      var tmp User
      if err == iterator.Done {
        break
      }
      if err != nil {
        return
      }
      fmt.Println(doc.Data())
      doc.DataTo(&tmp)
      resultData = append(resultData, tmp)
    }

    sort.Slice(resultData, func(i, j int) bool {
      return resultData[i].Score > resultData[j].Score
    })

    c.JSON(200, resultData)
	})

	router.POST("/postResult", func(c *gin.Context) {
    var user User
    c.BindJSON(&user)

		_, _, err := client.Collection("user").Add(ctx, map[string]interface{}{
			"name":  user.Name,
			"score": user.Score,
		})
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		c.String(http.StatusOK, "Hello!")
	})

	http.Handle("/", router)
	appengine.Main()
}
