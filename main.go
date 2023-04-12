package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type artiste struct {
	Id             int
	Image          string
	Name           string
	Members        []string
	Creationdate   int
	Creadatefilter int
	Firstalbum     string
	Locations      string
	Concertdates   string
	Relations      string
	CIndex         map[string][]string
}

type relation struct {
	Id             int
	DatesLocations map[string][]string
}

func Index(id string) map[string][]string {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var indextest relation
	json.Unmarshal(body, &indextest)
	var newmap = make(map[string][]string)

	for index := range indextest.DatesLocations {
		newCity := FormatLocation(index)

		for date := range indextest.DatesLocations[index] {
			indextest.DatesLocations[index][date] = FormatDate(indextest.DatesLocations[index][date])
		}
		newmap[newCity] = indextest.DatesLocations[index]

	}

	return newmap

}

func GetArtists() []artiste {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var artistes []artiste
	json.Unmarshal(body, &artistes)

	return artistes
}

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/Artistes", func(w http.ResponseWriter, r *http.Request) {
		var artistes []artiste
		if r.Method == http.MethodPost && r.FormValue("searchbar") != "" {
			artistes = SearchArtistsByName(r.FormValue("searchbar"), GetArtists())
		} else {
			artistes = GetArtists()
		}
		artiste := artistes
		tmpl := template.Must(template.ParseFiles("static/artistes/artistes.html", "Templates/header/header.html", "Templates/groups.html"))
		tmpl.Execute(w, artiste)
	})

	http.HandleFunc("/Description", func(w http.ResponseWriter, r *http.Request) {
		descri := DescriArtiste(r.FormValue("Artiste"))
		indextest := Index(r.FormValue("Artiste"))
		tmpl := template.Must(template.ParseFiles("static/description/description.html", "Templates/header/header.html", "Templates/groups.html"))
		descri.CIndex = indextest
		tmpl.Execute(w, descri)
	})

	http.ListenAndServe(":8081", nil)
}

func DescriArtiste(id string) artiste {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	var artiste artiste
	json.Unmarshal(body, &artiste)
	artiste.Firstalbum = FormatDate(artiste.Firstalbum)

	return artiste
}

func FormatDate(rawDate string) string {
	return strings.ReplaceAll(rawDate, "-", "/")
}

func FormatLocation(location string) string {
	var ville string
	ville = strings.ReplaceAll(location, "-", " / ")
	ville = strings.ReplaceAll(ville, "_", " ")
	return strings.Title(ville)

}

func SearchArtistsByName(query string, artistes []artiste) []artiste {
	var queryArtistes []artiste
	for index := range artistes {
		if strings.HasPrefix(artistes[index].Name, query) {
			queryArtistes = append(queryArtistes, artistes[index])
		}
	}
	return queryArtistes
}
