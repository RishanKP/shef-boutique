package db

import ("context"
         "go.mongodb.org/mongo-driver/mongo"
         "go.mongodb.org/mongo-driver/mongo/options"
         "fmt"
         "log"
         "os"
	 "time"
       )


var clientOptions = options.Client().ApplyURI("mongodb+srv://admin:"+ os.Getenv("DBPASS") +"@cluster0.rn0nv.mongodb.net/shef-boutique?retryWrites=true&w=majority")
var Client, Err = mongo.Connect(context.TODO(), clientOptions)

func Connect()  {

 if Err != nil {
     log.Fatal(Err)
 }
 Err = Client.Ping(context.TODO(), nil)

 if Err != nil {
     log.Fatal(Err)
 }

 fmt.Println("Connected to MongoDB!")
}

func AddStat(category string){
	collection := Client.Database("stats").Collection(category)

	type dateStruct struct{
		Date time.Time `json:"date"`
	}

	insert,err := collection.InsertOne(context.TODO(), dateStruct{Date:time.Now()})
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(insert)

	return
}
