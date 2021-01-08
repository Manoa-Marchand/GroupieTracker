package main

import "net/http"

func main() {
	fs := http.FileServer(http.Dir("./template/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	port := "8080"
	http.HandleFunc("/", index)
	println("Le serveur se lance sur le port " + port)
	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}
