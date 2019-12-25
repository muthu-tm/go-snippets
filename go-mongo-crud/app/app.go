package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-snippets/go-mongo-crud/db"
	"github.com/go-snippets/go-mongo-crud/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conn db.Connection

func SetupConn() {
	conn = db.EstablishConn(os.Getenv("DBHost"), os.Getenv("DBPort"))
}

func AddUser(c *gin.Context) {
	user := models.NewUser(c.PostForm("name"), c.PostForm("id"), c.PostForm("email"))
	conn := conn.Use(os.Getenv("DBName"), os.Getenv("DBCollection"))

	insertResult, err := conn.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": insertResult.InsertedID})
}

func GetUser(c *gin.Context) {
	var user models.User
	filter := bson.D{{"id", c.Param("id")}}

	conn := conn.Use(os.Getenv("DBName"), os.Getenv("DBCollection"))
	err := conn.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", user)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": user})
}

func GetAllUsers(c *gin.Context) {
	var users []*models.User

	conn := conn.Use(os.Getenv("DBName"), os.Getenv("DBCollection"))
	findOptions := options.Find()
	cur, err := conn.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, &user)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found documents: %+v\n", users)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}

func RemoveUser(c *gin.Context) {
	conn := conn.Use(os.Getenv("DBName"), os.Getenv("DBCollection"))
	filter := bson.D{{"id", c.Param("id")}}
	deleteResult, err := conn.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the users collection\n", deleteResult.DeletedCount)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": deleteResult.DeletedCount})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	userEmail := c.PostForm("email")
	filter := bson.D{{"id", userID}}

	update := bson.D{
		{"$set", bson.D{
			{"email", userEmail},
		}},
	}
	conn := conn.Use(os.Getenv("DBName"), os.Getenv("DBCollection"))
	updateResult, err := conn.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": updateResult.ModifiedCount})
}
