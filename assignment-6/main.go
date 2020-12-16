package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
)

type Player struct {
	gorm.Model
	Name    string
	Role    string
	JersyNo int
	Age     int
}

var db *gorm.DB

var err error

func main() {
	router := mux.NewRouter()
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=password")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Player{})

	router.HandleFunc("/", GetPlayers).Methods("GET")
	router.HandleFunc("/player/{id}", GetPlayer).Methods("GET")
	router.HandleFunc("/player/{id}", DeletePlayer).Methods("DELETE")
	router.HandleFunc("/addplayer", AddPlayer).Methods("POST")
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	var players []Player
	db.Find(&players)
	json.NewEncoder(w).Encode(&players)
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	db.First(&player, params["id"])
	json.NewEncoder(w).Encode(&player)
}
func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var player Player
	db.First(&player, params["id"])
	db.Delete(&player)
	var players []Player
	db.Find(&players)
	json.NewEncoder(w).Encode(&players)
}

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	json.NewDecoder(r.Body).Decode(&player)
	db.Create(&player)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}
