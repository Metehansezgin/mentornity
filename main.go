package main

import (
	"context"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type Data struct {
	Name  string `json:"Name"`
	City string `json:"City"`
	Category   string `json:"Category"`
}

type ScoreData struct {
	Name  string `json:"Name"`
	City string `json:"City"`
	Category   string `json:"Category"`
	Score	int `json:"Score"`
}

func ConnectDB() *mongo.Collection {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		 log.Fatal(err)
	}

	collection := client.Database("odev").Collection("mentor")

	return collection
}

func _getData () []Data{

	var veriler []Data

	cur, err := ConnectDB().Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {		 
		var veri Data
		 
		err := cur.Decode(&veri)
		if err != nil {
			log.Fatal(err)
		}
		veriler = append(veriler, veri)
	}
	return veriler
}

func getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var veriler []Data
	
	veriler = _getData()
	
	json.NewEncoder(w).Encode(veriler)
}

func addData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=ascii")
    var veri Data
	var verilerSonuc []ScoreData
	var verilerDb []Data

	_ = json.NewDecoder(r.Body).Decode(&veri)
	verilerDb = _getData()
	
	for  _,_veri := range verilerDb {
		var Vsonuc ScoreData
		Vsonuc.Name = _veri.Name
		Vsonuc.City = _veri.City
		Vsonuc.Category = _veri.Category

		if _veri.City == veri.City && _veri.Category != veri.Category{
			Vsonuc.Score =100
		}else if _veri.City == veri.City && _veri.Category == veri.Category{
			Vsonuc.Score =50
		}else if _veri.City != veri.City && _veri.Category != veri.Category{
			Vsonuc.Score =50
		}else if _veri.City != veri.City && _veri.Category == veri.Category{
			Vsonuc.Score =0
		}

		verilerSonuc = append(verilerSonuc,Vsonuc)
	}
	result, err := ConnectDB().InsertOne(context.TODO(), veri)
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(result)
	}
	  
	json.NewEncoder(w).Encode(verilerSonuc)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/api/getData", getData).Methods("GET")
	
	r.HandleFunc("/api/addData", addData).Methods("POST")

	http.ListenAndServe(":4000", r)

}