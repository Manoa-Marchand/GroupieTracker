package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

/*Création des structure qui vont stocker les éléments*/
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
	} `json:"index"`
}
type dateByIdJSON struct {
	Id   int      `json:"id"`
	Date []string `json:"dates"`
}

type tabLoca struct {
	City    string
	Country string
	Slugh   string
}

type infoLocation struct {
	City    string
	Country string
	date    string
	name    string
}

/*Fonction main qui va lancer le serveur ainsi que gerer les pages */
func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	port := "8080"
	http.HandleFunc("/", index)
	http.HandleFunc("/artists", artists)
	http.HandleFunc("/artist", artist)
	http.HandleFunc("/locations", locations)
	http.HandleFunc("/location", location)
	println("Le serveur se lance sur le port " + port)
	http.ListenAndServe(":"+port, nil)

}

/*fonction qui permet d'éviter les doublons*/
func unique(list []string) []string {
	keys := make(map[string]bool)
	newList := []string{}
	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			newList = append(newList, entry)
		}
	}
	return newList
}

/*fonction pour mettre des espaces pour les locations*/
func space(list []string) []string {
	for index, word := range list {
		word = strings.Replace(word, "_", " ", -1)
		list[index] = strings.Title(word)
	}
	return list
}

/*Fonction qui s'occupe de la page de chaque ville*/
func location(w http.ResponseWriter, r *http.Request) {
	/*recuperation de la ville avec son slugh*/
	//locationSlugh := r.FormValue("location")
	/*on refait une requete à l'api*/
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
	/*algo pour mettre en commun la ville avec un artiste et une date*/
	/*
		for _, loca := range locations.Index {
			for index, localisation := range loca.Locations {
				if localisation == locationSlugh {

				}
			}
		}

	*/
}

/*Fonction qui gere la page de la la liste des villes*/
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
	var List [193]tabLoca
	var Villes []string
	var Pays []string
	var Slughs []string
	for _, loca := range locations.Index {
		for _, uni := range loca.Locations {
			tab = append(tab, uni)
		}
	}
	listUnique := unique(tab)
	for index := range listUnique {
		slugh := listUnique[index]
		Slughs = append(Slughs, slugh)
		for index2 := range slugh {
			if slugh[index2] == '-' {
				ville := slugh[:index2]
				pays := slugh[index2+1:]
				Villes = append(Villes, ville)
				Pays = append(Pays, pays)
			}
		}
	}
	newVilles := space(Villes)
	newPays := space(Pays)

	for index := range listUnique {
		List[index].Slugh = Slughs[index]
		List[index].City = newVilles[index]
		List[index].Country = newPays[index]
	}
	files := []string{"./template/locationList.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, &List)
	}
}

/*fonction qui s'occupe de la page des artistes*/
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

/*fonction qui s'occupe de la page avec la liste des artistes */
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

/*fonction qui s'occupe de la page principale*/
func index(w http.ResponseWriter, r *http.Request) {
	files := []string{"./template/index.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, nil)
	}
}
