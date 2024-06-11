package search

import("context"
       "go.mongodb.org/mongo-driver/bson"
       "go.mongodb.org/mongo-driver/bson/primitive"
       "github.com/gin-gonic/gin"
       "shef-boutique/models"
       "shef-boutique/db"
      )

func Products(c *gin.Context){
	var query = c.Query("query")

	var products []*models.Product

	var collection = db.Client.Database("shef-boutique").Collection("product")

	cur,err := collection.Find(context.TODO(),bson.M{
					"$or":[]bson.M{
						bson.M{"name":primitive.Regex{Pattern: query, Options: "i"}},
						bson.M{"id":primitive.Regex{Pattern: query, Options: "i"}},
					},
	})
	if err != nil{
		products = nil
	}
	for cur.Next(context.TODO()){
		var product models.Product
		err := cur.Decode(&product)

		if err != nil{
			continue
		}

		products = append(products,&product)
	}

	c.JSON(200,gin.H{
		"status":"ok",
		"products":products,
	})

}

