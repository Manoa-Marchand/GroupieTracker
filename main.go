package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
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

type LocationsJSON struct {
	Index []struct {
		Id        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates`
	} `json:"index"`
}

type tabLoca struct {
	Id        int
	Locations []string
	Countries []string
	Slugh     []string
}

func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	port := "8080"
	http.HandleFunc("/", index)
	http.HandleFunc("/artists", artists)
	http.HandleFunc("/artist", artist)
	http.HandleFunc("/locations", locations)
	println("Le serveur se lance sur le port " + port)
	http.ListenAndServe(":"+port, nil)

}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	for cpt, str := range list {
		str = strings.Replace(str, "_", " ", -1)
		str = strings.Replace(str, "-", ", ", -1)
		list[cpt] = strings.Title(str)
	}
	return list
}

func locations(w http.ResponseWriter, r *http.Request) {
	urlapi := "https://groupietrackers.herokuapp.com/api/locations"
	res, err := http.Get(urlapi)
	if err != nil {
		fmt.Println(err.Error())
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	var locations LocationsJSON
	json.Unmarshal(data, &locations)

	var tab []string
	var List tabLoca
	for _, loca := range locations.Index {
		for _, uni := range loca.Locations {
			tab = append(tab, uni)
		}
	}
	List.Locations = unique(tab)
	files := []string{"./template/locationList.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, &List)
	}
}

func artist(w http.ResponseWriter, r *http.Request) {
	idArtiste := r.FormValue("artiste")
	urlapi := "https://groupietrackers.herokuapp.com/api/artists/" + idArtiste
	res, err := http.Get(urlapi)
	if err != nil {
		fmt.Println(err.Error())
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	var artiste ArtistesJSON
	json.Unmarshal(data, &artiste)
	files := []string{"./template/artist.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, &artiste)
	}
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
	var artistes []ArtistesJSON
	json.Unmarshal(data, &artistes)
	files := []string{"./template/artistListe.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, &artistes)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	files := []string{"./template/index.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, nil)
	}
}
