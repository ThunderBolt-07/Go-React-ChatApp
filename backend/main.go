package main

import (
	files "chat/handlers/file"
	ws "chat/handlers/ws"
	"log"
	"net/http"
	//"os"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL, " hi")

	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../frontend/hello.html")
}

func main() {
	//godotenv.Load()
	err := files.SetupClient()

	if err != nil {
		log.Fatal("error occur while stttimg up cleint", err)
	}

	hub := ws.NewHub()

	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		log.Println("got request")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			log.Println("Received preflight request")
			w.WriteHeader(http.StatusOK)
			return
		}

		err := files.FileUpload(w, r, hub)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("some error occured while sneding file"))
			return

		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("sending file was a success"))

	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		log.Println("upgrade req received")

		err := ws.HandleWS(w, r, hub)

		if err != nil {
			log.Println("error while upgrading to ws connecton", err)
			w.Write([]byte("error occured, try selecting a different username"))
			return
		}

	})

	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal("unbale to start server")
	}

}
