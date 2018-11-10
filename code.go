package main

import (
	"bytes"
	"fmt"
	"github.com/zserge/webview"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var (
	folder string   // Holds the folder to check for images
	images []string // Contains the list of valid images
)

// End the program if there is an error
func errFail(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

// Check that a folder is valid to search for images
func checkFolder(newFolder string) bool {
	if _, err := os.Stat(newFolder); os.IsNotExist(err) {
		fmt.Printf("\"%s\" is not a directory\n", newFolder)
		//os.Exit(1)
		return false
	}
	return true
}

// Change to a new folder of images
func setFolder(newFolder string) {

	// Empty current list of files
	images = []string{}

	// Check that you can read files from the folder

	if checkFolder(newFolder) {
		folder = newFolder
	} else {
		return
	}

	files, err := ioutil.ReadDir(folder)
	errFail(err)

	// Go over each file in the folder
	for _, f := range files {

		// Convenience
		filepath := path.Join(folder, f.Name())

		// Check that the file is not a folder and not empty
		x, err := os.Stat(filepath)
		errFail(err)
		if x.IsDir() && x.Size() > 0 {
			continue
		}

		// Check if the file can be read
		i, err := os.Open(filepath)
		errFail(err)
		defer i.Close()

		// Prove the file can be read
		buffer := make([]byte, 512)
		_, err = i.Read(buffer)
		errFail(err)

		// Check if the file is an image
		contentType := http.DetectContentType(buffer)
		if strings.HasPrefix(contentType, "image/") {

			// Add the file to the list of valid images
			images = append(images, f.Name())
		}
	}
}

// Serve page code at localhost:8080/random
func mainPage(w http.ResponseWriter, r *http.Request) {

	// Pick an image if there are images to pick from (otherwise 404)
	file := "assets/404.jpg"
	if len(images) > 0 {
		file = images[rand.Int()%len(images)]
	}

	// Read the code in page.htmland fill in variables as needed
	html, err := ioutil.ReadFile("page.html")
	errFail(err)
	html = bytes.Replace(html, []byte("RANDOM"), []byte(file), 1)
	html = bytes.Replace(html, []byte("FILENAME"), []byte(file), 1)
	html = bytes.Replace(html, []byte("FOLDER"), []byte(folder), 1)

	// Serve the finished page up
	w.Write(html)
}

// Read when the user picks a different folder to read images from
func viewFolder(w http.ResponseWriter, r *http.Request) {
	// This page should get PUT requests, so be cheeky for everything else
	if r.Method != "POST" {
		w.Write([]byte("\"" + folder + "\"\n\n"))
		for _, f := range images {
			w.Write([]byte(f + "\n"))
		}
		return
	}

	// Change the current folder
	r.ParseForm()
	nf := r.Form["folder"]
	if checkFolder(nf[0]) {
		setFolder(nf[0])
	}
	// TODO make a way to alert the user that an invalid directory is invalid

	// Send the user back to the main page
	http.Redirect(w, r, "main", 303)
}

// Serve a single image (TODO maybe merge with manage folder?)
func showImage(w http.ResponseWriter, r *http.Request) {

	// What does the URL say?
	file := r.URL.Path[8:]

	// No image is specified
	if file == "" {
		http.Redirect(w, r, "/folder", 303)
		return
	}

	// Check that the image is a real one
	for _, f := range images {
		if f == file {

			// Serve the file
			http.ServeFile(w, r, path.Join(folder, file))
			return
		}
	}

	// An invalid file was requested
	http.ServeFile(w, r, "assets/404.jpg")
}

// Self explanatory
func runApp() {

	// Specify the website resources
	http.HandleFunc("/main", mainPage)
	http.HandleFunc("/folder", viewFolder)
	http.HandleFunc("/folder/", showImage)

	// TODO make this just return the 404 image when the image cannot be found
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	// Run the server forever
	err := http.ListenAndServe(":8080", nil)
	errFail(err)
}

// Check that required assets are available before starting
func initApp() {
	// See if assets/images.location.txt exists
	if _, err := os.Stat("./assets/images.location.txt"); !os.IsNotExist(err) {
		b, err := ioutil.ReadFile("assets/images.location.txt")
		errFail(err)

		// Verify that the folder listed within actually exists
		if checkFolder(string(b)) {
			setFolder(string(b))
		} else {
			os.Exit(1)
		}
	}

	// Computers aren't random, they only pretend with math
	rand.Seed(time.Now().UnixNano())
}

// Run everything
func main() {

	// Check and set things before starting the app
	initApp()

	// Run the app's backend
	go runApp()

	// Wait a second so the frontend doesn't start before the backend
	time.Sleep(1 * time.Second)

	// Open the front-end to the main page
	webview.Open("Random Image Picker", "http://localhost:8080/main", 800, 600, true)
}
