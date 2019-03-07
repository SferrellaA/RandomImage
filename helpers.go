package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
)

// End the program if there is an error
func errFail(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

// setFolder changes the directory images are pulled from
func setFolder(newFolder string) bool {

	// Check that the new folder is actually valid
	if _, err := os.Stat(newFolder); os.IsNotExist(err) {

		// Directory is invalid
		return false
	}

	// Update the folder variable
	folder = newFolder

	// Update the location file
	err := ioutil.WriteFile("assets/images.location.txt", []byte(folder), 0644)
	errFail(err)

	// Read files in the new folder
	files, err := ioutil.ReadDir(folder)
	errFail(err)

	// Empty current lists of files
	images = []string{}
	lastHalf = []int{}

	// Go over each file in the folder
	for _, f := range files {

		// Check if image can be added to list of valid images
		if checkFile(folder, f.Name()) {
			images = append(images, f.Name())
		}
	}

	// Directory is valid
	return true
}

// checkFile checks is a file (asset or image) is valid
func checkFile(location, name string) bool {

	// A file must be named
	if name == "" {
		return false
	}

	// Prepare for later steps
	filepath := path.Join(location, name)
	x, err := os.Stat(filepath)

	// Check that the file is real, not a directory, and larger than 0 bytes
	if os.IsNotExist(err) || x.IsDir() || !(x.Size() > 0) {
		return false
	}

	// Check if the file can be read
	i, err := os.Open(filepath)
	defer i.Close()
	if nil != err {
		return false
	}

	// Prove the file can be read
	buffer := make([]byte, 512)
	_, err = i.Read(buffer)
	if nil != err {
		return false
	}

	// Check that the file is an image
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return false
	}

	// All conditions have been met
	return true
}

// pickRandom selects an image that is "perceived to be" random
func pickRandom() string {

	// Loop forever
	for {

		// Randomly pick an index in the images array
		i := rand.Int() % len(images)

		// Check if index has been in the last n/2 iterations
		if func() bool {
			for _, j := range lastHalf {
				if i == j {
					return true
				}
			}
			return false
		}() {

			// Try another random number if index is a repeat
			continue
		}

		// Record index as having been selected
		lastHalf = append(lastHalf, i)

		// Cut the least recent index selections to allow eventual repeats
		if len(lastHalf) > len(images)/2 {
			lastHalf = lastHalf[1:]
		}

		// Return the sufficiently random image
		return images[i]
	}

	// This shouldn't be reached
	return ""
}
