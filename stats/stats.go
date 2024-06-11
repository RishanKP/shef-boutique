package stats

import("github.com/gin-gonic/gin"
       "time"
       "shef-boutique/db"
       "context"
       "fmt"
       "go.mongodb.org/mongo-driver/bson"
       "go.mongodb.org/mongo-driver/mongo"
       "go.mongodb.org/mongo-driver/bson/primitive"
      )

type Dates struct{
	SDate time.Time `json:"sdate"`
	EDate time.Time `json:"edate"`
}
type Stat struct{
	Hour int32 `json:"hour"`
	Total int32 `json:"total"`
}

func StatsByDay(c *gin.Context){
	var date Dates
	c.BindJSON(&date)
	category := c.Param("id")

	collection := db.Client.Database("stats").Collection(category)
	count,err := collection.Aggregate(context.TODO(),mongo.Pipeline{
		 bson.D{{
			"$match",bson.M{
			"date":bson.M{
				"$gte": primitive.NewDateTimeFromTime(date.SDate.AddDate(0,0,1)),
			        "$lt": primitive.NewDateTimeFromTime(date.EDate.AddDate(0,0,1)),
			},
	            }}},
		 bson.D{{
			"$group", bson.M{
			"_id": bson.M{"$hour":"$date"},
			"total": bson.M{"$sum":1},
		     }}},
		bson.D{{
			"$sort",bson.D{
				{"_id",1},
			},
		 }},
		})

	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"error fetching data"})
		return
	}

	var response []bson.M

	if err = count.All(context.TODO(), &response); err != nil {
		panic(err)
	}

	var stats []*Stat
	for k := range response{
		var stat Stat
		stat.Hour = response[k]["_id"].(int32)
		stat.Total = response[k]["total"].(int32)

		stats = append(stats,&stat)
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"date":date,
		"stats":stats,
	})

}
