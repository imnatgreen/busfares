package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/fares"
	"github.com/imnatgreen/busfares/internal/router"
)

func HandleRequests(f *fares.FareObjects, a *agency.Agencies) {
	http.HandleFunc("/api/", homepage)
	http.HandleFunc("/api/getfares", getFares(f, a))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "busfares api")
}

func getFares(f *fares.FareObjects, a *agency.Agencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tripPlan := router.TripPlannerResponse{}
		err := json.NewDecoder(r.Body).Decode(&tripPlan)
		if err != nil {
			fmt.Fprintln(w, ":/ an error occured. check logs for details.")
			log.Print(err)
		}
		tripPlan.AddFares(f, a)
		err = json.NewEncoder(w).Encode(tripPlan)
		if err != nil {
			fmt.Fprintln(w, ":/ an error occured. check logs for details.")
			log.Print(err)
		}
	}
}
