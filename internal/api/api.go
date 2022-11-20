package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/colours"
	"github.com/imnatgreen/busfares/internal/router"
	"github.com/jackc/pgx/v5"
)

func HandleRequests(c *pgx.Conn, a *agency.Agencies) {
	http.HandleFunc("/api/", homepage)
	http.HandleFunc("/api/getfares", getFares(c, a))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handler(prefix string, c *pgx.Conn, a *agency.Agencies, col *colours.Colours) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/", homepage)
	mux.HandleFunc(prefix+"/getfares", getFares(c, a))
	mux.HandleFunc(prefix+"/getcolour", getColour(col, a))
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
		fares, err := tripPlan.GetAllFares(c, a)
		if err != nil {
			fmt.Fprintln(w, ":/ an error occured. check logs for details.")
			log.Print(err)
		}
		err = json.NewEncoder(w).Encode(fares)
		if err != nil {
			fmt.Fprintln(w, ":/ an error occured. check logs for details.")
			log.Print(err)
		}
	}
}

func getColour(c *colours.Colours, a *agency.Agencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		ag := r.URL.Query().Get("agency")
		var noc colours.Noc
		if ag != "" {
			aNoc, err := a.GetNoc(agency.AgencyId(ag))
			if err != nil {
				noc = colours.Noc("")
			}
			noc = colours.Noc(aNoc)
		} else {
			noc = colours.Noc(r.URL.Query().Get("noc"))
		}
		colour, err := c.Get(noc, colours.Line(r.URL.Query().Get("line")))
		if err == nil {
			fmt.Fprint(w, colour)
		}
	}
}
