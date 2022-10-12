// https://data.bus-data.dft.gov.uk/api/v1/fares/dataset/?noc=BPTR,HRGT,KDTR,LNUD,ROST,TPEN,YACT,YCST&api_key=BODS_API_KEY

package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"time"

	"context"

	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/bojanz/currency"
	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/api"
	"github.com/imnatgreen/busfares/internal/fares"
	"github.com/imnatgreen/busfares/internal/persist"
	"github.com/imnatgreen/busfares/internal/router"
	"github.com/jackc/pgx/v5"
)

func main() {
	var err error

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	var fareObjects fares.FareObjects

	// err = fares.GetDatasets("data/fares", os.Getenv("NOCS"))
	// if err != nil {
	// 	log.Fatalf("failed to get datasets: %v", err)
	// }

	// err = fareObjects.AddDir("data/fares")
	// if err != nil {
	// 	log.Fatalf("failed to add directory: %v", err)
	// }

	// // save fareObjects to disk
	// start := time.Now()
	// err = persist.Save("fareObjects.json", fareObjects)
	// if err != nil {
	// 	log.Fatalf("failed to save fareObjects: %v", err)
	// }
	// log.Printf("saved fares to disk in %s", time.Since(start))

	// load fareObjects from disk
	start := time.Now()
	err = persist.Load("fareObjects.json", &fareObjects)
	if err != nil {
		log.Fatalf("failed to load fareObjects: %v", err)
	}
	log.Printf("loaded %d fareObjects from disk in %s", len(fareObjects.Objects), time.Since(start))

	// load agencies from GTFS files
	start = time.Now()
	agencies, _ := loadAgencies(os.Getenv("GTFS_DIR"))
	log.Printf("loaded %d agencies from disk in %s", len(agencies), time.Since(start))

	// test router
	jsonData, _ := os.Open("internal/router/resp.json")
	defer jsonData.Close()
	data, _ := io.ReadAll(jsonData)
	var res router.TripPlannerResponse
	res, err = router.ParseJson(data)
	if err != nil {
		log.Fatal(err)
	}

	// get fares for router response
	start = time.Now()
	err = res.AddFares(&fareObjects, &agencies)
	log.Printf("added fares to response in %s", time.Since(start))
	if err != nil {
		log.Fatal(err)
	}

	// test add fares as json
	encoded, err := json.Marshal(res)
	if err != nil {
		log.Print(err)
	}
	log.Print(string(encoded))

	log.Println(currency.NewFormatter(currency.NewLocale("gb")).Format(res.Plan.Itineraries[1].Legs[1].Fares[1].Amount))
	// log.Print(res.Plan.Itineraries[1].Legs[1].Fares)
	api.HandleRequests(&fareObjects, &agencies)
	//log.Print(res)
	//log.Print(res.Plan.Itineraries[1].Legs[1].To.Name)

	// test fares
	// xml, _ := os.Open("internal/fares/ROST_483_Outbound_Single_.xml")
	// defer xml.Close()
	// xmlData, _ := io.ReadAll(xml)
	// var fareData fares.FareObject

	// start = time.Now()
	// fareData, err = fares.ParseXml(xmlData)
	// log.Printf("parsed fares in %s", time.Since(start))

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Print(fareData.Tariffs[0].ValidBetween.FromDate.Time)
	// start = time.Now()
	// fare, err := fareData.GetFare("25001425", "2500IMG2914")
	// log.Printf("got fare in %s", time.Since(start))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(currency.NewFormatter(currency.NewLocale("gb")).Format(fare))
	// log.Print(fareData)
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
	return agencies, nil
}
