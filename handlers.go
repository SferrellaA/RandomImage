package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// mainPage serves the gui at localhost:8080/main
func mainPage(w http.ResponseWriter, r *http.Request) {

	// Pick an image if there are images to pick from (otherwise 404)
	file := "assets/404.jpg"
	if len(images) > 0 {
		file = pickRandom()
	} else {
		// TODO use a different image to say there are no images
	}

	// Read the code in page.html
	html, err := ioutil.ReadFile("page.html")
	errFail(err) // it's bad if you can't read the GUI's code

	// Replace variables in the HTML with current values
	html = bytes.Replace(html, []byte("RANDOM"), []byte(file), 1)
	html = bytes.Replace(html, []byte("FILENAME"), []byte(file), 1)
	html = bytes.Replace(html, []byte("FOLDER"), []byte(folder), 1)

	// Serve the finished page up
	w.Write(html)
}

// changeFolder handles commands to change the current image folder
func changeFolder(w http.ResponseWriter, r *http.Request) {

	// Read new folder from POST message
	r.ParseForm()
	nf := r.Form["folder"]

	// TODO alert the user the folder was not set
	if !setFolder(nf[0]) {
	}

	// Send the user back to the main page
	http.Redirect(w, r, "/main", 303)
}

// serveFile checks that a given file is valid, then serves it
func serveFile(w http.ResponseWriter, r *http.Request, location, name string) bool {

	// check that file is valid
	if checkFile(location, name) {

		// serve the file and confirm success
		http.ServeFile(w, r, path.Join(location, name))
		return true
	}

	// the requested file was invalid
	return false
}

// serveAsset serves images for use in the RandomImage gui
func serveAsset(w http.ResponseWriter, r *http.Request) {

	// serve the requested asset
	if serveFile(w, r, "assets", strings.TrimPrefix(r.URL.Path, "/assets/")) {
		return
	}

	// redirect all invalid requests back to the main gui
	http.Redirect(w, r, "/main", 303)
}

// serveImage serves a specific requested image from the image folder
func serveImage(w http.ResponseWriter, r *http.Request) {

	// parse what image is requested
	file := strings.TrimPrefix(r.URL.Path, "/folder/")

	// show the folder view if no particular image is requested
	if "" == file {
		showFolder(w, r)
	}

	// Check that image was and still is valid
	for _, f := range images {
		if f == file && serveFile(w, r, folder, file) {
			return
		}
	}

	// The requested image cannot be found
	if !serveFile(w, r, "assets", "404.jpg") {

		// It's bad if the 404 image can't be found
		os.Exit(1)
	}

}

// showFolder serves the list of valid images
func showFolder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html>\n<p>" + folder + "</p><pre>\n"))
	for _, f := range images {
		w.Write([]byte("<a href=\"http://localhost:8080/folder/" + f + "\">" + f + "</a>\n"))
	}
	w.Write([]byte("</pre>\n</html>\n"))
}

// folderHandler services the multiple functions of /folder
func folderHandler(w http.ResponseWriter, r *http.Request) {

	switch {

	// Changing the folder
	case r.Method == "POST":
		changeFolder(w, r)

	// An image is requested
	case strings.HasPrefix(r.URL.Path, "/folder/"):
		serveImage(w, r)

	// Show folder contents
	default:
		showFolder(w, r)
	}
}
