package main

import (
	"context"
	"net/http"
	"encoding/json"
	"log"
	"os"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
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

	dbString := os.Getenv("DB")
	log.Println("ENV-DB :"+dbString)
	clientOptions := options.Client().ApplyURI(dbString)
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

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println("handl call homePage")
	path := r.URL.Path

	if path == "/"{

		path = "./index.html"
	}else{

		path = "."+path
	}

	http.ServeFile(w,r,path)
 
}

func getData(w http.ResponseWriter, r *http.Request) {
	log.Println("handl call getData")
	w.Header().Add("Content-Type", "application/json")

	var veriler []Data
	
	veriler = _getData()
	
	json.NewEncoder(w).Encode(veriler)
}

func addData(w http.ResponseWriter, r *http.Request) {
	log.Println("handl call addData")
	w.Header().Set("Content-Type", "application/json")
    var veri Data
	var verilerSonuc []ScoreData
	var verilerDb []Data

	_ = json.NewDecoder(r.Body).Decode(&veri)

	log.Println(veri)
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
	log.Println("App start")
	err := godotenv.Load(".env")
	
	log.Println("App start")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}

	r := mux.NewRouter()
	port := os.Getenv("PORT")
	log.Println("ENV-PORT :"+port)
	r.HandleFunc("/api/getData", getData).Methods("GET")
	r.HandleFunc("/", homePage).Methods("GET")
	
	r.HandleFunc("/api/addData", addData).Methods("POST")

	http.ListenAndServe(":"+port, r)

}