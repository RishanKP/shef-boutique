package product

import("github.com/gin-gonic/gin"
	"os"
	"fmt"
	"context"
	"shef-boutique/models"
	"shef-boutique/db"
	"time"
	"strconv"
	"strings"
	myaws "shef-boutique/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/aws"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	)

var collection = db.Client.Database("shef-boutique").Collection("product")

func Add(c *gin.Context){
	var product models.Product

	sess := myaws.ConnectAws()
	uploader := s3manager.NewUploader(sess)
	bucket := os.Getenv("AWS_BUCKET")

	for i:=1;i<=5;i++{
		key := "boutique/products/"
		fieldname := "image"+strconv.Itoa(i)
		file,header,err := c.Request.FormFile(fieldname)
		if err!=nil{
			break
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
		product.Image = append(product.Image,up.Location)
	}

	allsizes := c.Request.FormValue("sizes")
	allsizes = strings.Replace(allsizes,"\"","",-1)
	allsizes = strings.Replace(allsizes,"[","",-1)
	allsizes = strings.Replace(allsizes,"]","",-1)

	sizes := strings.Split(allsizes,",")
	for _,size := range sizes{
		product.Sizes = append(product.Sizes,size)
	}

	product.Name = c.Request.FormValue("name")
	product.Price,_ = strconv.ParseFloat(c.Request.FormValue("price"),64)
	product.Desc = c.Request.FormValue("description")
	product.Category = c.Request.FormValue("category")
	product.CategoryId = c.Request.FormValue("category_id")
	product.Instock = true
	product.IsNewArrival,_ = strconv.ParseBool(c.Request.FormValue("newarrival"))
	product.Id = db.GenerateId()
	product.DateAdded = time.Now()

	insert,err := collection.InsertOne(context.TODO(),product)
	if err!=nil{
	 c.JSON(500,gin.H{
	    "status":"failed",
	    "error":"failed to insert document",
	 })

	 return
	}
	categoryCollection := db.Client.Database("shef-boutique").Collection("category")
	subcategoryCollection := db.Client.Database("shef-boutique").Collection("subcategory")

	if product.IsNewArrival == true{
		categoryUpdate := bson.D{
			{"$inc",bson.D{
				{"newarrivals",1},
			}},
		}

		var subcategory models.SubCategory
		err = subcategoryCollection.FindOne(context.TODO(),bson.D{{"id",product.CategoryId}}).Decode(&subcategory)
		if err!=nil{
			fmt.Println(err)
		}

		categoryUpdateResult,_ := categoryCollection.UpdateOne(context.TODO(),bson.D{{"id",subcategory.ParentId}},categoryUpdate)
		subcategoryUpdateResult,_ := subcategoryCollection.UpdateOne(context.TODO(),bson.D{{"id",product.CategoryId}},categoryUpdate)

		fmt.Println(categoryUpdateResult)
		fmt.Println(subcategoryUpdateResult)
	}

	c.JSON(200,gin.H{"status":"ok","insertionId":insert.InsertedID})
}


func AllByCategory(c *gin.Context){
	var products []*models.Product

	var size string
	var lrange,rrange float64

	if c.Query("size") != ""{
		size = c.Query("size")
	}
	if c.Query("price") != ""{
		ranges := strings.Split(c.Query("price"),"-")
		lrange,_ = strconv.ParseFloat(ranges[0],64)
		rrange,_ = strconv.ParseFloat(ranges[1],64)
	}
	cur,err := collection.Find(context.TODO(),bson.D{{}})

	if size == "" && rrange == 0.0{
		cur,err = collection.Find(context.TODO(),bson.D{{"categoryid",c.Param("id")}})
	}else if size != "" && rrange == 0.0{
		cur,err = collection.Find(context.TODO(),bson.D{{"categoryid",c.Param("id")},{"sizes",size},})
	}else if size == "" && rrange!= 0.0{
		cur,err = collection.Find(context.TODO(),bson.D{{"categoryid",c.Param("id")},
								{
									"price",bson.M{
										"$gte":lrange,
										"$lte":rrange,
									},
								},
							 })
	}else{
		cur,err = collection.Find(context.TODO(),bson.D{{"categoryid",c.Param("id")},{"sizes",size},
								{
									"price",bson.M{
										"$gte":lrange,
										"$lte":rrange,
									},
								},
							 })

	}

	if err!=nil{
		c.JSON(500,gin.H{
			"status":"failed",
			"error":"internal server error. failed to fetch data",
		})

	  return
	}

	for cur.Next(context.TODO()){
		var product models.Product
		err := cur.Decode(&product)

		if err!=nil{
			c.JSON(500,gin.H{
				"status":"failed",
				"error":"internal server error. failed to fetch data",
			})

		   return
		}

		products = append(products,&product)
	}

	var contact models.Contact
	err = db.Client.Database("shef-boutique").Collection("contact").FindOne(context.TODO(),bson.D{{"id","621e45bead41b28af1de18e1"}}).Decode(&contact)
	if err != nil{
		c.JSON(500,gin.H{
                        "status":"failed",
                        "message":"failed to fetch contact details",
                })

                return

	}

	var hideprice models.HidePrice
	err = db.Client.Database("shef-boutique").Collection("hideprice").FindOne(context.TODO(),bson.D{{"id","6252873d36dec99be69166ad"}}).Decode(&hideprice)
	if err != nil{
		c.JSON(500,gin.H{
                        "status":"failed",
                        "message":"failed to fetch price details",
                })
                return

	}

	c.JSON(200,gin.H{
		"status":"ok",
		"products":products,
		"hide_price":hideprice.HidePrice,
		"contact":contact,
	})
}

func GetOne(c *gin.Context){
	id := c.Param("id")
	var product models.Product
	err := collection.FindOne(context.TODO(),bson.D{{"id",id}}).Decode(&product)

	if err!=nil{
		c.JSON(500,gin.H{
			"status":"failed",
			"message":"failed to fetch document",
		})

		return
	}

	var contact models.Contact
	err = db.Client.Database("shef-boutique").Collection("contact").FindOne(context.TODO(),bson.D{{"id","621e45bead41b28af1de18e1"}}).Decode(&contact)
	if err != nil{
		c.JSON(500,gin.H{
                        "status":"failed",
                        "message":"failed to fetch contact details",
                })

                return

	}


	var hideprice models.HidePrice
	err = db.Client.Database("shef-boutique").Collection("hideprice").FindOne(context.TODO(),bson.D{{"id","6252873d36dec99be69166ad"}}).Decode(&hideprice)
	if err != nil{
		c.JSON(500,gin.H{
                        "status":"failed",
                        "message":"failed to fetch price details",
                })

                return

	}


	c.JSON(200,gin.H{
		"status":"ok",
		"product":product,
		"hide_price":hideprice.HidePrice,
		"contact":contact,
	})
}

func Update(c *gin.Context){
	var product models.Product
	c.BindJSON(&product)
	product.DateAdded = time.Now()

	filter := bson.D{{"id",c.Param("id")}}
	update := bson.D{{"$set",bson.D{
					{"name",product.Name},
					{"price",product.Price},
					{"desc",product.Desc},
					{"category",product.Category},
					{"categoryid",product.CategoryId},
					{"sizes",product.Sizes},
					{"dateadded",product.DateAdded},
				},
			}}

	updateResult,err:=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to update data"})
		return
	}

	c.JSON(200,gin.H{"msg":"updated",
			 "result":updateResult,
			})

}

func ToggleNewArrival(c *gin.Context){

	var product models.Product

	err := collection.FindOne(context.TODO(),bson.D{{"id",c.Param("id")}}).Decode(&product)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to fetch data"})
		return
	}

	categoryCollection := db.Client.Database("shef-boutique").Collection("category")
	subcategoryCollection := db.Client.Database("shef-boutique").Collection("subcategory")

	var subcategory models.SubCategory
	err = subcategoryCollection.FindOne(context.TODO(),bson.D{{"id",product.CategoryId}}).Decode(&subcategory)

	if err!=nil{
		c.JSON(400,gin.H{"msg":"failed to fetch category details"})
		return

	}

	update := bson.D{{}}
	categoryUpdate := bson.D{{}}

	if product.IsNewArrival == true{
		update = bson.D{{"$set",bson.D{{"isnewarrival",false}}}}
		categoryUpdate = bson.D{
			{"$inc",bson.D{
				{"newarrivals",-1},
			}},
		}
	}else{
		update = bson.D{{"$set",bson.D{{"isnewarrival",true}}}}
		categoryUpdate = bson.D{
			{"$inc",bson.D{
				{"newarrivals",1},
			}},
		}
	}

	updateResult,err := collection.UpdateOne(context.TODO(),bson.D{{"id",c.Param("id")}},update)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to update data"})
		return
	}

	categoryUpdateResult,_ := categoryCollection.UpdateOne(context.TODO(),bson.D{{"id",subcategory.ParentId}},categoryUpdate)
	subCategoryUpdateResult,_ := subcategoryCollection.UpdateOne(context.TODO(),bson.D{{"id",subcategory.Id}},categoryUpdate)

	c.JSON(200,gin.H{"msg":"updated",
			 "updateResult":updateResult,
			 "categoryUpdateResult":categoryUpdateResult,
			 "subCategoryUpdateResult":subCategoryUpdateResult,
	})

}

func ToggleStock(c *gin.Context){

	type stock struct{
		Instock bool `bson:"instock"`
	}
	var s stock

	opts := options.FindOne().SetProjection(bson.D{{"instock",1}})
	err := collection.FindOne(context.TODO(),bson.D{{"id",c.Param("id")}},opts).Decode(&s)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to fetch data"})
		return
	}

	update := bson.D{{}}

	if s.Instock == true{
		update = bson.D{{"$set",bson.D{{"instock",false}}}}
	}else{
		update = bson.D{{"$set",bson.D{{"instock",true}}}}
	}

	updateResult,err := collection.UpdateOne(context.TODO(),bson.D{{"id",c.Param("id")}},update)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to update data"})
		return
	}

	c.JSON(200,gin.H{"msg":"updated",
			 "result":updateResult,
	})

}

func HidePrice(c *gin.Context){
	var h models.HidePrice
	err := db.Client.Database("shef-boutique").Collection("hideprice").FindOne(context.TODO(),bson.D{{"id","6252873d36dec99be69166ad"}}).Decode(&h)
	if err!=nil {
		c.JSON(400,gin.H{"msg":"failed to fetch"})
		return
	}

	c.JSON(200,gin.H{"status":"ok","hideprice":h.HidePrice})

}

func ToggleHidePrice(c *gin.Context){

	var s models.HidePrice

	err := db.Client.Database("shef-boutique").Collection("hideprice").FindOne(context.TODO(),bson.D{{"id","6252873d36dec99be69166ad"}}).Decode(&s)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to fetch data"})
		return
	}

	update := bson.D{{}}

	if s.HidePrice == true{
		update = bson.D{{"$set",bson.D{{"hideprice",false}}}}
	}else{
		update = bson.D{{"$set",bson.D{{"hideprice",true}}}}
	}

	updateResult,err := db.Client.Database("shef-boutique").Collection("hideprice").UpdateOne(context.TODO(),bson.D{{"id","6252873d36dec99be69166ad"}},update)
	if err!=nil{
		fmt.Println(err)
		c.JSON(400,gin.H{"msg":"failed to update data"})
		return
	}

	c.JSON(200,gin.H{"msg":"updated",
			 "result":updateResult,
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
