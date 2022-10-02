package fares

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

type FareObjects struct {
	Objects []FareObject
}

func (f *FareObjects) AddXml(file io.ReadCloser) (err error) {
	defer file.Close()
	var obj FareObject
	err = xml.NewDecoder(file).Decode(&obj)
	if err != nil {
		return err
	}
	f.Objects = append(f.Objects, obj)
	return err
}

func (f *FareObjects) AddZip(path string) (err error) {
	openZipFile, err := zip.OpenReader(path)
	if err != nil {
		log.Fatalf("failed to open zip file %s", path)
		return err
	}
	defer openZipFile.Close()
	for _, file := range openZipFile.File {
		// if filepath.Ext(file.Name) == ".xml" {
		// temporarily only add ROST fare files
		if filepath.Ext(file.Name) == ".xml" && strings.HasPrefix(file.Name, "ROST") {
			openFile, err := file.Open()
			if err != nil {
				log.Fatalf("failed to open file %s in %s", file.Name, path)
				return err
			}
			log.Printf("adding xml %s from zip", file.Name)
			f.AddXml(openFile)
		}
	}
	return err
}

func (f *FareObjects) AddDir(dir string) (err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// recursively walk through subdirectories
		// if d.IsDir() {
		// 	err = f.AddDir(path)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		if filepath.Ext(d.Name()) == ".zip" {
			log.Printf("adding zip %s", d.Name())
			f.AddZip(path)
		}
		if filepath.Ext(d.Name()) == ".xml" {
			xml, err := os.Open(path)
			if err != nil {
				return err
			}
			defer xml.Close()
			log.Printf("adding xml %s", d.Name())
			f.AddXml(xml)
		}
		return nil
	})
	return err
}

type ApiResponse struct {
	Results []ApiResult `json:"results"`
}

type ApiResult struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

func GetDatasets(dir string, nocs string) (err error) {
	// create directory if it doesn't exist
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	url := os.Getenv("BODS_API_BASE") + "/fares/dataset/?noc=" + nocs + "&status=published&api_key=" + os.Getenv("BODS_API_KEY")
	resp, err := http.Get(url)
	var apiResponse ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return err
	}

	var reqs []*grab.Request
	for _, result := range apiResponse.Results {
		req, err := grab.NewRequest(dir, result.Url)
		if err != nil {
			return err
		}
		reqs = append(reqs, req)
	}

	client := grab.NewClient()
	respch := client.DoBatch(4, reqs...)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for resp := range respch {
		for {
			select {
			case <-t.C:
				log.Printf("downloading %v (%.2f%%)", resp.Request.URL(), 100*resp.Progress())
			case <-resp.Done:
				if err := resp.Err(); err != nil {
					log.Printf("download failed: %v", err)
				} else {
					log.Printf("downloaded %v", resp.Filename)
				}
				return err
			}
		}
	}

	return err
}
