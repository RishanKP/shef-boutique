package subcategory

import("github.com/gin-gonic/gin"
	"os"
	"fmt"
	"context"
	"shef-boutique/models"
	"shef-boutique/db"
	"time"
	myaws "shef-boutique/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	)

var collection = db.Client.Database("shef-boutique").Collection("subcategory")

func Add(c *gin.Context){
	var category models.SubCategory

	sess := myaws.ConnectAws()
	uploader := s3manager.NewUploader(sess)
	bucket := os.Getenv("AWS_BUCKET")

	key := "boutique/subcategory"
	fieldname := "image"
	file,header,err := c.Request.FormFile(fieldname)

	if err!=nil{
		c.JSON(500, gin.H{
			  "error":"error reading file",
		})
		return

	}
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
	category.Image = up.Location

	category.Name = c.Request.FormValue("name")
	category.Description = c.Request.FormValue("description")
	category.ParentId = c.Request.FormValue("parent_id")
	category.NewArrivals = 0
	category.Id = db.GenerateId()
	category.DateAdded = time.Now()

	insert,err := collection.InsertOne(context.TODO(),category)
	if err!=nil{
	 c.JSON(500,gin.H{
	    "status":"failed",
	    "error":"failed to insert document",
	 })

	 return
	}

	c.JSON(200,gin.H{"status":"ok","insertionId":insert.InsertedID})
}


func All(c *gin.Context){
	var categories []*models.SubCategory

	cur,err := collection.Find(context.TODO(),bson.D{{"parentid",c.Param("id")}})
	db.AddStat(c.Param("id"))

	if err!=nil{
		c.JSON(500,gin.H{
			"status":"failed",
			"error":"internal server error. failed to fetch data",
		})

	  return
	}

	for cur.Next(context.TODO()){
		var category models.SubCategory
		err := cur.Decode(&category)

		if err!=nil{
			c.JSON(500,gin.H{
				"status":"failed",
				"error":"internal server error. failed to fetch data",
			})

		   return
		}

		categories = append(categories,&category)
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"categories":categories,
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
