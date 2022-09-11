// https://data.bus-data.dft.gov.uk/api/v1/fares/dataset/?noc=BPTR,HRGT,KDTR,LNUD,ROST,TPEN,YACT,YCST&api_key=BODS_API_KEY

package main

import (
	"archive/zip"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/router"
)

func main() {
	// fareDataset := fmt.Sprintf("https://data.bus-data.dft.gov.uk/api/v1/fares/dataset/?noc=BPTR,HRGT,KDTR,LNUD,ROST,TPEN,YACT,YCST&api_key=%s", os.Getenv("BODS_API_KEY"))
	// res, err := http.Get(fareDataset)
	// if err != nil {
	// 	res.Body // https://stackoverflow.com/a/31129967
	// }

	// load agencies from GTFS files
	agencies, _ := loadAgencies(os.Getenv("GTFS_DIR"))
	aTest, _ := agencies.GetNoc("OP291")
	log.Print(aTest)

	// test router
	json, _ := os.Open("internal/router/resp.json")
	defer json.Close()
	data, _ := io.ReadAll(json)
	var res router.TripPlannerResponse
	res, err := router.ParseJson(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(res)
	log.Print(res.Plan.Itineraries[1].Legs[1].To.Name)
}

func loadAgencies(gtfsDir string) (agencies agency.Agencies, err error) {
	agencies = make(agency.Agencies)

	var gtfsFiles []string
	err = filepath.WalkDir(gtfsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ".zip" {
			gtfsFiles = append(gtfsFiles, path)
		}
		return nil
	})
	if err != nil {
		return agencies, err
	}

	for _, gtfsFile := range gtfsFiles {
		log.Printf("importing from %s", gtfsFile)
		openGtfsFile, err := zip.OpenReader(gtfsFile)
		if err != nil {
			log.Fatalf("failed to open gtfs file %s", gtfsFile)
			return agencies, err
		}
		defer openGtfsFile.Close()
		for _, file := range openGtfsFile.File {
			if file.Name == "agency.txt" {
				openFile, err := file.Open()
				if err != nil {
					log.Fatalf("failed to open file %s in %s", file.Name, gtfsFile)
					return agencies, err
				}
				agencies.AddFromCsv(openFile)
			}
		}
	}
	log.Printf("imported %d agencies", len(agencies))
	return agencies, nil
}

// https://github.com/antchfx/xmlquery ? -> xpath?
// https://stackoverflow.com/questions/30256729
