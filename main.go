package main

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-vgo/robotgo"
)

// http://192.168.1.68:8080/ - videoctl
// http://192.168.1.68:9001/ - supervisor

//go:embed index.html
var files embed.FS

func main() {
	postVolume := func(w http.ResponseWriter, r *http.Request) {
		strVol, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			fmt.Printf("ERRO %v", err)
			return
		}

		vol, err := strconv.Atoi(string(strVol))
		if err != nil {
			http.Error(w, "Invalid body format. Expect value from 0 to 100", http.StatusBadRequest)
			fmt.Printf("ERRO %v", err)
			return
		}

		err = volume.SetVolume(vol)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			fmt.Printf("ERRO %v", err)
			return
		}
	}

	http.HandleFunc("/volume", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postVolume(w, r)
			return
		}

		level, err := volume.GetVolume()
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			fmt.Printf("ERRO %v", err)
			return
		}
		fmt.Fprintf(w, "%d", level)
	})

	http.HandleFunc("/click", func(w http.ResponseWriter, _ *http.Request) {
		robotgo.Click()
	})

	http.Handle("/", http.FileServer(http.FS(files)))
	log.Fatal(
		http.ListenAndServe("0.0.0.0:8080", nil),
	)
}
