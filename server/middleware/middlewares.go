package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jatin-02/todo/model"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	loadTheEnv()
	createDBInstance()
}

func loadTheEnv() {
	err := godotenv.Load(".env")
	checkNilError(err)
}

func createDBInstance() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	colName := os.Getenv("DB_COLLECTION")

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	checkNilError(err)

	err = client.Ping(context.TODO(), nil)
	checkNilError(err)
	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance created")
}

func getAllTasks() []primitive.M {
	curr, err := collection.Find(context.Background(), bson.D{{}})
	checkNilError(err)

	var results []primitive.M

	for curr.Next(context.Background()) {
		var result bson.M
		err := curr.Decode(&result)
		checkNilError(err)
		results = append(results, result)
	}

	defer curr.Close(context.Background())
	return results
}

func taskComplete(task string) {
	id, _ := primitive.ObjectIDFromHex(task)

	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	checkNilError(err)

	fmt.Println("modified count: ", result.ModifiedCount)
}

func insertOneTask(task model.TodoList) {
	insertResult, err := collection.InsertOne(context.Background(), task)
	checkNilError(err)
	fmt.Println("Inserted one record", insertResult.InsertedID)
}

func undoTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)

	filter := bson.M{"_id":id}
	update := bson.M{"$set":bson.M{"status":false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	checkNilError(err)

	fmt.Println("Modified count: ", result.ModifiedCount)
}

func deleteOneTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)

	filter := bson.M{"_id":id}
	result, err := collection.DeleteOne(context.Background(), filter)
	checkNilError(err)

	fmt.Println("Deleted count: ", result.DeletedCount)
}

func deleteAllTasks() int64 {
	result, err := collection.DeleteMany(context.Background(), bson.D{{}})
	checkNilError(err)
	fmt.Println("Deleted count: ", result.DeletedCount)

	return result.DeletedCount
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")

	payload := getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	taskComplete(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func InsertOneTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var task model.TodoList
	json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)
}

func UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	deleteOneTask(params["id"])
}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin","*")
	count := deleteAllTasks()
	json.NewEncoder(w).Encode(count)
}

func checkNilError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}