package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type ArtistesJSON struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    string   `json:"relations"`
}

var artistes []ArtistesJSON

type LocationsJSON struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

var localisation []LocationsJSON

func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	port := "8080"
	http.HandleFunc("/", index)
	http.HandleFunc("/artists", artists)
	println("Le serveur se lance sur le port " + port)
	http.ListenAndServe(":"+port, nil)

}

func artists(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	json.Unmarshal(data, &artistes)
	tpl, err := template.ParseFiles("template/artist.html")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, &artistes)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, nil)
	}
}
