package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/router"
	"github.com/jackc/pgx/v5"
)

func HandleRequests(c *pgx.Conn, a *agency.Agencies) {
	http.HandleFunc("/api/", homepage)
	http.HandleFunc("/api/getfares", getFares(c, a))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(prefix string, c *pgx.Conn, a *agency.Agencies) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/", homepage)
	mux.HandleFunc(prefix+"/getfares", getFares(c, a))
	return mux
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "busfares api")
}

func getFares(c *pgx.Conn, a *agency.Agencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		tripPlan := router.TripPlannerResponse{}
		err := json.NewDecoder(r.Body).Decode(&tripPlan)
		if err != nil {
			fmt.Fprintln(w, ":/ an error occured. check logs for details.")
			log.Print(err)
		}
		tripPlan.AddFares(c, a)
		err = json.NewEncoder(w).Encode(tripPlan)
		if err != nil {
			fmt.Fprintln(w, ":/ an error occured. check logs for details.")
			log.Print(err)
		}
	}
}
