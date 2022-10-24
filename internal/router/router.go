package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
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

// AddFares adds a list of possible fares to each transit leg in the response
func (r *TripPlannerResponse) AddFares(c *pgx.Conn, a *agency.Agencies) (err error) {
	for i, itinerary := range r.Plan.Itineraries {
		for l, leg := range itinerary.Legs {
			if leg.TransitLeg {
				agencyId := agency.AgencyId(TrimId(leg.AgencyId))
				noc, err := a.GetNoc(agencyId)
				if err != nil {
					return err
				}
				// use i and l to update original leg
				r.Plan.Itineraries[i].Legs[l].GetFares(c, noc)
			}
		}
	}
	return err
}

// GetFares finds the possible fares for the given transit leg and returns them
func (l *Leg) GetFares(c *pgx.Conn, n agency.Noc) (err error) {
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
		return err
	}

	// start = time.Now()
	// fareObjects, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (fares.FareObject, error) {
	// 	var f fares.FareObject
	// 	var uuid []byte
	// 	startLoop := time.Now()
	// 	err := row.Scan(&uuid, &f)
	// 	log.Printf("row scan took %s", time.Since(startLoop))
	// 	return f, err
	// })
	// log.Printf("collecting rows took %s", time.Since(start))
	fareObjects := []fares.FareObject{}

	start = time.Now()
	for rows.Next() {
		var f fares.FareObject
		timeScan := time.Now()
		values := rows.RawValues()
		log.Printf("raw values took %s", time.Since(timeScan))
		if err != nil {
			return err
		}
		timeScan = time.Now()
		err = json.Unmarshal(values[1], &f)
		log.Printf("unmarshal took %s", time.Since(timeScan))
		// err := rows.Scan(&uuid, &f)
		if err != nil {
			return err
		}
		fareObjects = append(fareObjects, f)
	}
	log.Printf("scanning rows took %s", time.Since(start))
	// pgx.ForEachRow(rows, )

	for _, obj := range fareObjects {
		if obj.ContainsStops(from, to) {
			fare, err := obj.GetFare(from, to)
			if err != nil && err != fares.ErrFareNotInTable {
				return err
			}
			if err != fares.ErrFareNotInTable {
				log.Printf("importing fare %s", fare.PreassignedFareProduct.Id)
				fareSlice = append(fareSlice, fare)
			}
		}
	}

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
	filteredFareSlices = filteredFareSlices[:j+1]

	// only keep newest fare for each type
	for _, fSlice := range filteredFareSlices {
		newestFare := fSlice[0]
		for _, fare := range fSlice {
			if fare.ValidBetween.FromDate.Time.After(newestFare.ValidBetween.FromDate.Time) {
				newestFare = fare
			}
		}
		l.Fares = append(l.Fares, newestFare)
	}

	log.Printf("found %d fares for leg %s, filtered to %d.", len(fareSlice), l.RouteShortName, len(l.Fares))
	return err
}

// Removes the "2:" prefix present on IDs from the OTP API
func TrimId(id string) string {
	return strings.TrimPrefix(id, "2:")
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
