package main

import (
	"context"
	//	"os/user"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"sort"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
)

type User struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Ranking struct {
	Rank1 User `json:"rank1"`
	Rank2 User `json:"rank2"`
	Rank3 User `json:"rank3"`
	Rank4 User `json:"rank4"`
	Rank5 User `json:"rank5"`
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

	router.GET("/getRanking/:number", func(c *gin.Context) {
		iter := client.Collection("user").Documents(ctx)
		var resultData []User

		number := c.Params.ByName("number")

		num, _ := strconv.Atoi(number)

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

		var res string

		tmpNum := strconv.Itoa(resultData[num-1].Score)

		res = number + "位" + " " + resultData[num-1].Name + " " + tmpNum

		c.String(200, res)
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
