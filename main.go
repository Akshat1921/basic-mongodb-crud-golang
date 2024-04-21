package main

import (
	"context"
	"log"
	"net/http"
	"projectx/usecase"
	"time"
	"github.com/gorilla/mux"
	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init(){
	// err:= godotenv.Load()
	// if err!=nil{
	// 	log.Fatal("env load err", err)
	// }
	// log.Println("env file loaded")

	// mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	// if err!=nil{
	// 	log.Fatal("connection error", err)
	// }
	// err = mongoClient.Ping(context.Background(), readpref.Primary())
	// if err!=nil{
	// 	log.Fatal("ping failed", err)
	// }
	// log.Println("mongo connected")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

    err = mongoClient.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
}

func main(){
	defer mongoClient.Disconnect(context.Background())

	r := mux.NewRouter()
	// create employee service
	coll := mongoClient.Database("companydb").Collection("employee")
	empService := usecase.EmployeeService{MongoCollection: coll}

	r.HandleFunc("/health",healthHandler).Methods(http.MethodGet)
	
	r.HandleFunc("/employee", empService.CreateEmployee).Methods(http.MethodPost)
	r.HandleFunc("/employee/{id}", empService.GetEmployeeByID).Methods(http.MethodGet)
	r.HandleFunc("/employee", empService.GetAllEmployee).Methods(http.MethodGet)
	r.HandleFunc("/employee/{id}", empService.UpdateEmployeeByID).Methods(http.MethodPut)
	r.HandleFunc("/employee/{id}", empService.DeleteEmployeeByID).Methods(http.MethodDelete)
	r.HandleFunc("/employee", empService.DeleteAllEmployee).Methods(http.MethodDelete)

	log.Println("server is running on port 4000")

	http.ListenAndServe(":4000", r)

}

func healthHandler(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running"))
}