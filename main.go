// https://data.bus-data.dft.gov.uk/api/v1/fares/dataset/?noc=BPTR,HRGT,KDTR,LNUD,ROST,TPEN,YACT,YCST&api_key=BODS_API_KEY

package main

import (
	"archive/zip"
	"flag"
	"net/http"
	"time"

	"context"

	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/imnatgreen/busfares/frontend"
	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/api"
	"github.com/jackc/pgx/v5"
)

func main() {
	var err error

	var devMode bool
	var dbUrl string
	var gtfsDir string
	var nocs string
	var bodsApiKey string
	var bodsApiBase string
	flag.BoolVar(&devMode, "dev", false, "run in dev mode")
	flag.StringVar(&dbUrl, "db", os.Getenv("DATABASE_URL"), "database url")
	flag.StringVar(&gtfsDir, "gtfs", os.Getenv("GTFS_DIR"), "directory containing gtfs files")
	flag.StringVar(&nocs, "nocs", os.Getenv("NOCS"), "comma separated list of NOCS")
	flag.StringVar(&bodsApiKey, "bods-key", os.Getenv("BODS_API_KEY"), "BODS api key")
	flag.StringVar(&bodsApiBase, "bods-base", os.Getenv("BODS_API_BASE"), "BODS api base url")
	flag.Parse()

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// err = fares.GetDatasets("data/fares", nocs)
	// if err != nil {
	// 	log.Fatalf("failed to get datasets: %v", err)
	// }

	// err = fares.AddDir(conn, "data/fares")
	// if err != nil {
	// 	log.Fatalf("failed to add directory: %v", err)
	// }

	// load agencies from GTFS files
	start := time.Now()
	agencies, _ := loadAgencies(gtfsDir)
	log.Printf("loaded %d agencies from disk in %s", len(agencies), time.Since(start))

	mux := http.NewServeMux()

	mux.Handle("/api/", api.Handler("/api", conn, &agencies))

	mux.Handle("/", frontend.SvelteKitHandler("/"))

	var handler http.Handler = mux

	if devMode {
		handler = devCors(handler)
		log.Println("server running in dev mode")
	}

	log.Fatal(http.ListenAndServe(":8080", handler))

	// test router
	// jsonData, _ := os.Open("internal/router/resp.json")
	// defer jsonData.Close()
	// data, _ := io.ReadAll(jsonData)
	// var res router.TripPlannerResponse
	// res, err = router.ParseJson(data)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get fares for router response
	// start = time.Now()
	// err = res.AddFares(conn, &agencies)
	// log.Printf("added fares to response in %s", time.Since(start))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// test add fares as json
	// encoded, err := json.Marshal(res)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Print(string(encoded))
	// log.Println(currency.NewFormatter(currency.NewLocale("gb")).Format(res.Plan.Itineraries[1].Legs[1].Fares[1].Amount))

	// log.Print(res.Plan.Itineraries[1].Legs[1].Fares)

	//log.Print(res)
	//log.Print(res.Plan.Itineraries[1].Legs[1].To.Name)
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

func devCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	})
}
