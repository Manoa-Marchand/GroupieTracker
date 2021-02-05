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
	Relations    RelationJSON
}

type LocationsJSON struct {
	Index []struct {
		Id        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type RelationJSON struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type tabLoca struct {
	City    string
	Country string
	Slugh   string
}

type artistlocation struct {
	Id           int
	Image        string
	Name         string
	CreationDate int
}

type Artisteslocation struct {
	Id           int
	Image        string
	Name         string
	CreationDate int
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
	port := "8086"
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
	locationSlugh := r.FormValue("location")
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

	/*on refait une requete à l'api*/
	urlapi2 := "https://groupietrackers.herokuapp.com/api/artists"
	res2, err := http.Get(urlapi2)
	if err != nil {
		fmt.Println(err.Error())
	}
	data2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res2.Body.Close()

	var locations LocationsJSON
	var artists []ArtistesJSON
	var artist []ArtistesJSON
	var ID []int

	json.Unmarshal(data, &locations)
	json.Unmarshal(data2, &artist)
	/*algo pour mettre en commun la ville avec un artiste*/
	for _, loca := range locations.Index {
		for _, localisation := range loca.Locations {
			if localisation == locationSlugh {
				ID = append(ID, loca.Id)
			}
		}
	}
	for _, artistesonly := range artist {
		for _, i := range ID {
			if artistesonly.Id == i {
				artists = append(artists, artistesonly)
			}
		}
	}

	files := []string{"./template/location.html", "./template/base.html"}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		tpl.Execute(w, &artists)
	}
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
	/*Information de l'artiste*/
	urlapi := "https://groupietrackers.herokuapp.com/api/artists/" + idArtiste
	request1, err := http.Get(urlapi)
	if err != nil {
		fmt.Println(err.Error())
	}
	data1, err := ioutil.ReadAll(request1.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer request1.Body.Close()
	var artiste ArtistesJSON
	json.Unmarshal(data1, &artiste)
	/*Relation de l'artiste*/
	urlapi = "https://groupietrackers.herokuapp.com/api/relation/" + idArtiste
	request2, err := http.Get(urlapi)
	if err != nil {
		fmt.Println(err.Error())
	}
	data, err := ioutil.ReadAll(request2.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer request2.Body.Close()
	var relation RelationJSON
	json.Unmarshal(data, &relation)
	artiste.Relations = relation
	/*Affichage de la page*/
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
