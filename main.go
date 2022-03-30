package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main () {

	r := mux.NewRouter()

	r.HandleFunc("/", Index).Methods("GET")

	// localhost:8080/media/12/stream

	r.HandleFunc("/media/{mediaID:[0-9]+}/stream", Stream).Methods("GET")

	// localhost:8080/media/12/index10.ts
	// localhost:8080/media/12/index9.ts
	r.HandleFunc("/media/{mediaID:[0-9]+}/{segmentName:index[0-9]+.ts}", Stream).Methods("GET")

	http.ListenAndServe(":4040", r)
}

func Index (w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func Stream (w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	mediaID, err := strconv.Atoi(params["mediaID"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	segmentName, ok := params["segmentName"]

	mediaBase := getMediaBase(mediaID)

	if !ok {
		playlist := "index.m3u8"
		serveHlsM3u8(w, r, mediaBase, playlist)
	} else {
		serveHlsTs(w, r, mediaBase, segmentName)
	}
}

func getMediaBase(mediaID int) string {

	mediaRoot := "assets/media"

	return fmt.Sprintf("%s/%d", mediaRoot, mediaID)
}

func serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, playlist string) {
	
	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, playlist)

	http.ServeFile(w, r, mediaFile)

	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segmentName string) {
	
	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, segmentName)

	http.ServeFile(w, r, mediaFile)

	w.Header().Set("Content-Type", "video/MP2T")
}
