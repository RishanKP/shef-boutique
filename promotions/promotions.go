package promotions

import("github.com/gin-gonic/gin"
	"fmt"
	"os"
	"context"
	"shef-boutique/models"
	"shef-boutique/db"
	"time"
	myaws "shef-boutique/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	)

var collection = db.Client.Database("shef-boutique").Collection("promotions")

func Add(c *gin.Context){
	var promotion models.Promotion

	sess := myaws.ConnectAws()
	uploader := s3manager.NewUploader(sess)
	bucket := os.Getenv("AWS_BUCKET")

	key := "promotions/"
	file,header,_ := c.Request.FormFile("image")
	key += header.Filename

	filetype := header.Header["Content-Type"][0]

	up,err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL: aws.String("public-read"),
		Key: aws.String(key),
		ContentType: aws.String(filetype),
		ContentDisposition: aws.String("inline"),
		Body: file,
	})

	if err != nil {
		fmt.Println(err)
		 c.JSON(500, gin.H{
			  "error":"Failed to upload file",
			  "uploader": up,
	 })
	  return
	}

	promotion.Image = up.Location
	promotion.FromDate,_ = time.Parse(time.RFC3339,c.Request.FormValue("from"))
	promotion.ToDate,_  = time.Parse(time.RFC3339,c.Request.FormValue("to"))

	promotion.Id = db.GenerateId()
	promotion.DateAdded = time.Now()

	insert,err := collection.InsertOne(context.TODO(),promotion)
	if err!=nil{
	 c.JSON(500,gin.H{
	    "status":"failed",
	    "error":"failed to insert document",
	 })

	 return
	}

	c.JSON(200,gin.H{"status":"ok","details":promotion,"insertionId":insert.InsertedID})
}

func All(c *gin.Context){
	var promotions []*models.Promotion

	cur,err := collection.Find(context.TODO(),bson.M{"todate": bson.M{
		"$gte": primitive.NewDateTimeFromTime(time.Now()),
		}})

	if err!=nil{
		c.JSON(500,gin.H{
			"status":"failed",
			"error":"internal server error. failed to fetch data",
		})

	  return
	}

	for cur.Next(context.TODO()){
		var promotion models.Promotion
		err := cur.Decode(&promotion)

		if err!=nil{
			c.JSON(500,gin.H{
				"status":"failed",
				"error":"internal server error. failed to fetch data",
			})

		   return
		}

		promotions = append(promotions,&promotion)
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"promotions":promotions,
	})
}

func GetOne(c *gin.Context){
	id := c.Param("id")
	var promotion models.Promotion
	err := collection.FindOne(context.TODO(),bson.D{{"id",id}}).Decode(&promotion)

	if err!=nil{
		c.JSON(500,gin.H{
			"status":"failed",
			"message":"failed to fetch document",
		})

		return
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"promotion":promotion,
	})
}

func Delete(c *gin.Context){
	id := c.Param("id")

	Delete,err := collection.DeleteOne(context.TODO(),bson.D{{"id",id}})

	if err!=nil{
		c.JSON(500,gin.H{
			"status":"failed",
			"message":"failed to delete document",
		})

		return
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"result":Delete,
	})

}

