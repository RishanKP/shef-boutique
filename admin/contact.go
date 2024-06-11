package admin

import("github.com/gin-gonic/gin"
	"shef-boutique/db"
	"shef-boutique/models"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	)

var contactCollection = db.Client.Database("shef-boutique").Collection("contact")

func GetContact(c *gin.Context){
	var contact models.Contact
	err := contactCollection.FindOne(context.TODO(),bson.D{{"id","621e45bead41b28af1de18e1"}}).Decode(&contact)
	if err!=nil {
		c.JSON(400,gin.H{"msg":"failed to fetch"})
		return
	}

	c.JSON(200,gin.H{"status":"ok","contact":contact})

}
func UpdateContact(c *gin.Context){
	var contact models.Contact
	c.BindJSON(&contact)

	newData := bson.D{{"$set",bson.D{
					{"contact",contact.Contact},
					{"whatsapp",contact.Whatsapp},
				},
			}}

	update,err := contactCollection.UpdateOne(context.TODO(),bson.D{{"id","621e45bead41b28af1de18e1"}},newData)
	if err!=nil {
		c.JSON(400,gin.H{"msg":"failed to update"})
		return
	}

	c.JSON(200,gin.H{"msg":"updated","result":update})
}

