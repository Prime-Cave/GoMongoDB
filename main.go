package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prime-cave/mongo-golang/controllers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main(){
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	fmt.Println("Sever is running on port... 8080")
	http.ListenAndServe(":8080", r)
}

func getSession() *mongo.Client{
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(""))
	 if err != nil {
		panic(err)
	 }
	 return client
}