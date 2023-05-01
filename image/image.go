// Package testdata provides test images.
package image

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
)

var (
	// Dir is the path to the directory containing the test data.
	Dir = initDir()

	// Images contains all images by filename.
	Images = make(map[string]*imageserver.Image)

	// Server is an Image Server that uses filename as source.
	Server = imageserver.Server(imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
		source, err := params.GetString(imageserver_source.Param)
		if err != nil {
			return nil, err
		}
		im, err := Get(source)
		if err != nil {
			return nil, &imageserver.ParamError{Param: imageserver_source.Param, Message: err.Error()}
		}
		return im, nil
	}))
)

// Get returns an Image for a name.
func Get(name string) (*imageserver.Image, error) {
	im, ok := Images[name]
	if !ok {
		return nil, fmt.Errorf("unknown image \"%s\"", name)
	}
	return im, nil
}

func initDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFile)
}

func init() {
	loadImage()
}

// loadImage loads the test data and populates the Images map.
func loadImage() {
	files, err := ioutil.ReadDir(Dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			data, err := ioutil.ReadFile(filepath.Join(Dir, filename))
			if err != nil {
				panic(err)
			}
			im := &imageserver.Image{
				Format: getImageFormat(filename),
				Data:   data,
			}
			Images[filename] = im
			fmt.Println(filename)
		}
	}
}

func getImageFormat(filename string) string {
	switch filepath.Ext(filename) {
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".gif":
		return "gif"
	default:
		return ""
	}
}

func RefreshImage(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		fmt.Println("Refreshing image...")
		// Get the list of files in the directory
		files, err := ioutil.ReadDir(Dir)
		if err != nil {
			panic(err)
		}
		// Check if there are any new files
		newImages := false
		for _, file := range files {
			if !file.IsDir() {
				filename := file.Name()
				if _, ok := Images[filename]; !ok {
					// This is a new image
					newImages = true
					break
				}
			}
		}
		// If there are new images, reload all images
		if newImages {
			// Lock the Images map to prevent race conditions
			for key := range Images {
				delete(Images, key)
			}
			loadImage()
			fmt.Println("Image refreshed.")
		}
	}
}
