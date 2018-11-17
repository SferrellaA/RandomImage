package main

import (
	"github.com/zserge/webview"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	folder   string   // Holds the folder to check for images
	images   []string // Contains the list of valid images
	lastHalf []int    // Contains the index of the last so many displayed images
)

// initApp performs various setup tasks before starting app
func initApp() {

	// Most computers aren't actually random, they only pretend with math
	rand.Seed(time.Now().UnixNano())

	// See if assets/images.location.txt exists
	if _, err := os.Stat("assets/images.location.txt"); !os.IsNotExist(err) {

		// Set the initial image folder
		b, err := ioutil.ReadFile("assets/images.location.txt")
		errFail(err)

		// Verify that the folder listed within actually exists
		if !setFolder(string(b)) {
			os.Exit(1)
		}
	}

	// Specify the website resource handlers
	http.HandleFunc("/main", mainPage)
	http.HandleFunc("/folder/", folderHandler)
	http.HandleFunc("/assets/", serveAsset)
}

// main runs the app
func main() {

	// Check and set things before starting the app
	initApp()

	// Run the app's backend server forever
	go func() {
		err := http.ListenAndServe(":8080", nil)
		errFail(err)
	}()

	// Wait a second so the frontend doesn't start before the backend
	time.Sleep(1 * time.Second)

	// Open the front-end to the main page
	webview.Open("Random Image Picker", "http://localhost:8080/main", 800, 600, true)
}
