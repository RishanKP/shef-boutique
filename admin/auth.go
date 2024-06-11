package admin

import("github.com/gin-gonic/gin"
	"shef-boutique/db"
	"shef-boutique/models"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	)

var collection = db.Client.Database("shef-boutique").Collection("admin")

func userIdExist(id string) bool{
	var admin models.Admin
	err := collection.FindOne(context.TODO(),bson.D{{"userid",id}}).Decode(&admin)

	if err!=nil{
		return false
	}else{
		if admin.UserID == ""{
		 return false
		}else{
		 return true
		}
	}
}

func Register(c *gin.Context){
	var admin models.Admin
	c.BindJSON(&admin)

	if userIdExist(admin.UserID){
		c.JSON(400,gin.H{"message":"user id exist"})
	}else{
		admin.Id = db.GenerateId()
		insert,err := collection.InsertOne(context.TODO(), admin)
		if err!=nil{
		  c.JSON(500,gin.H{"message":"failed to create admin user"})
		}else{
		  fmt.Println("inserted ",insert.InsertedID)
		  c.JSON(200,gin.H{"message":"admin created"})
		}

	}
}

func Login(c *gin.Context){
	var admin models.Admin
	var loginData models.Login

	c.BindJSON(&loginData)

	err := collection.FindOne(context.TODO(),bson.D{{"userid",loginData.UserID}}).Decode(&admin)
	if err!=nil{
	  c.JSON(401,gin.H{"message":"unauthorized. userid not found"})
	}else{

	 if(admin.Pass == loginData.Pass){
	  admin.Pass = "xxxxxxxxxxx"
	  c.JSON(200,gin.H{"details":admin})
	  }else{
	   c.JSON(401,gin.H{"message":"unauthorized. wrong password"})
	  }

	}
}
