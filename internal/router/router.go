package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sort"
	"time"

	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/fares"
	"github.com/jackc/pgx/v5"
)

type TripPlannerResponse struct {
	Plan TripPlan `json:"plan"`
}

type TripPlan struct {
	Itineraries []Itinerary `json:"itineraries"`
}

type Itinerary struct {
	Legs []Leg `json:"legs"`
}

type Leg struct {
	StartTime            int64        `json:"startTime"`
	EndTime              int64        `json:"endTime"`
	Mode                 string       `json:"mode"`
	TransitLeg           bool         `json:"transitLeg"`
	AgencyTimeZoneOffset int          `json:"agencyTimeZoneOffset"`
	AgencyId             string       `json:"agencyId"`
	ServiceDate          string       `json:"serviceDate"`
	From                 LegVertex    `json:"from"`
	To                   LegVertex    `json:"to"`
	LegGeometry          LegGeometry  `json:"legGeometry"`
	RouteShortName       string       `json:"routeShortName"`
	Fares                []fares.Fare `json:"fares"`
}

type LegVertex struct {
	Name       string `json:"name"`
	StopId     string `json:"stopId"`
	VertexType string `json:"vertexType"`
}

type LegGeometry struct {
	Points string `json:"points"`
	Lenght int    `json:"length"`
}

func ParseJson(data []byte) (TripPlannerResponse, error) {
	var response TripPlannerResponse
	err := json.Unmarshal(data, &response)
	return response, err
}

// GetAllFares returns a list of fares for each leg of the TripPlannerResponse, following the same structure
func (r *TripPlannerResponse) GetAllFares(c *pgx.Conn, a *agency.Agencies) ([][][]fares.Fare, error) {
	var err error

	// check how many legs are in each itinerary
	itineraryLengths := []int{}
	for _, itinerary := range r.Plan.Itineraries {
		itineraryLengths = append(itineraryLengths, len(itinerary.Legs))
	}

	itineraryFares := make([][][]fares.Fare, len(r.Plan.Itineraries))

	for i, itinerary := range r.Plan.Itineraries {
		// create slice to add fares to
		itineraryFares[i] = make([][]fares.Fare, itineraryLengths[i])

		for l, leg := range itinerary.Legs {
			if leg.TransitLeg {
				agencyId := agency.AgencyId(TrimId(leg.AgencyId))
				noc, err := a.GetNoc(agencyId)
				if err != nil {
					return nil, err
				}
				legFares, err := leg.GetFares(c, noc)
				if err != nil {
					return nil, err
				}
				itineraryFares[i][l] = legFares
			}
		}
	}
	return itineraryFares, err
}

// GetFares finds the possible fares for the given transit leg and returns them
func (l *Leg) GetFares(c *pgx.Conn, n agency.Noc) (legFares []fares.Fare, err error) {
	// find all possible fares
	from := fares.Naptan(TrimId(l.From.StopId))
	to := fares.Naptan(TrimId(l.To.StopId))
	fareSlice := []fares.Fare{}

	start := time.Now()

	query := fmt.Sprintf(`select *
	from fares
	where fare_object->'Lines' @> '[{"PublicCode": "%s", "OperatorRef": {"Ref": "noc:%s"}}]'
		and fare_object->'ScheduledStopPoints' @> '[{"ScheduledStopPointRef": "atco:%s"}, {"ScheduledStopPointRef": "atco:%s"}]';`, l.RouteShortName, n, from, to)
	rows, err := c.Query(context.Background(), query)
	log.Printf("query took %s", time.Since(start))
	if err != nil {
		return legFares, err
	}

	start = time.Now()
	var uuid []byte
	var obj fares.FareObject
	pgx.ForEachRow(rows, []any{&uuid, &obj}, func() error {
		fare, err := obj.GetFare(from, to)
		if err != nil && err != fares.ErrFareNotInTable {
			return err
		}
		if err != fares.ErrFareNotInTable {
			fareSlice = append(fareSlice, fare)
		}
		return nil
	})
	log.Printf("processed query in %s", time.Since(start))

	// sort fares by type
	sort.Slice(fareSlice, func(i, j int) bool {
		return fareSlice[i].PreassignedFareProduct.Id < fareSlice[j].PreassignedFareProduct.Id
	})

	// group fares by type into separate slices
	filteredFareSlices := make([][]fares.Fare, len(fareSlice))
	prevId := ""
	currentId := ""
	j := 0
	for i, fare := range fareSlice {
		if i == 0 {
			prevId = fare.PreassignedFareProduct.Id
		} else {
			prevId = currentId
		}
		currentId = fare.PreassignedFareProduct.Id
		if prevId != currentId {
			j++
		}
		filteredFareSlices[j] = append(filteredFareSlices[j], fare)
	}

	// remove empty slices
	if len(filteredFareSlices) > 0 {
		filteredFareSlices = filteredFareSlices[:j+1]
	}

	// only keep newest fare for each type
	for _, fSlice := range filteredFareSlices {
		newestFare := fSlice[0]
		for _, fare := range fSlice {
			if fare.ValidBetween.FromDate.Time.After(newestFare.ValidBetween.FromDate.Time) {
				newestFare = fare
			}
		}
		legFares = append(legFares, newestFare)
	}

	log.Printf("found %d fares for leg %s, filtered to %d.", len(fareSlice), l.RouteShortName, len(legFares))
	return legFares, err
}

// Removes the "X:" prefix present on IDs from the OTP API (where X is the router ID)
func TrimId(id string) string {
	re := regexp.MustCompile(`^[0-9]+:`)
	trimmed := re.ReplaceAllString(id, "")
	return trimmed
}

func FilterTransitLegs(response TripPlannerResponse) []Leg {
	var transitLegs []Leg
	for _, itinerary := range response.Plan.Itineraries {
		for _, leg := range itinerary.Legs {
			if leg.TransitLeg {
				transitLegs = append(transitLegs, leg)
			}
		}
	}
	return transitLegs
}
