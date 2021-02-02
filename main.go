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

type LocationsJSON struct {
	Index []struct {
		Id        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type tabLoca struct {
	City    string
	Country string
	Slugh   string
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
	for index := 0; index < len(listUnique); index++ {
		slugh := listUnique[index]
		Slughs = append(Slughs, slugh)
		for index2 := 0; index2 < len(slugh); index2++ {
			if slugh[index2] == '-' {
				ville := slugh[:index2]
				pays := slugh[index2+1:]
				Villes = append(Villes, ville)
				Pays = append(Pays, pays)
			}
		}
	}
	for index := range listUnique {
		List[index].Slugh = Slughs[index]
		List[index].City = Villes[index]
		List[index].Country = Pays[index]
	}
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
