package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prime-cave/mongo-golang/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()
type UserController struct {
	Session *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	u.Id = primitive.NewObjectID()

	results, err := uc.Session.Database("mongo-golang").Collection("users").InsertOne(ctx, &u)
	if err != nil{
		log.Fatal(err)
	}

	uj, err := json.Marshal(&results)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    id := p.ByName("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }
    u := models.User{}

    err = uc.Session.Database("mongo-golang").Collection("users").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&u)
    if err != nil {
		log.Printf("Error finding user with _id %v: %v", oid, err)
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "%s", err)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(u)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	oid, err :=  primitive.ObjectIDFromHex(id)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
        return
	}

	result, err := uc.Session.Database("mongo-golang").Collection("users").DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil{
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
