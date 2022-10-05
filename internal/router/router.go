package router

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/imnatgreen/busfares/internal/agency"
	"github.com/imnatgreen/busfares/internal/fares"
)

type TripPlannerResponse struct {
	Plan tripPlan `json:"plan"`
}

type tripPlan struct {
	Itineraries []itinerary `json:"itineraries"`
}

type itinerary struct {
	Legs []leg `json:"legs"`
}

type leg struct {
	StartTime            int64       `json:"startTime"`
	EndTime              int64       `json:"endTime"`
	Mode                 string      `json:"mode"`
	TransitLeg           bool        `json:"transitLeg"`
	AgencyTimeZoneOffset int         `json:"agencyTimeZoneOffset"`
	AgencyId             string      `json:"agencyId"`
	ServiceDate          string      `json:"serviceDate"`
	From                 legVertex   `json:"from"`
	To                   legVertex   `json:"to"`
	LegGeometry          legGeometry `json:"legGeometry"`
	RouteShortName       string      `json:"routeShortName"`
	Fares                []fares.Fare
}

type legVertex struct {
	Name       string `json:"name"`
	StopId     string `json:"stopId"`
	VertexType string `json:"vertexType"`
}

type legGeometry struct {
	Points string `json:"points"`
	Lenght int    `json:"length"`
}

func ParseJson(data []byte) (TripPlannerResponse, error) {
	var response TripPlannerResponse
	err := json.Unmarshal(data, &response)
	return response, err
}

// AddFares adds a list of possible fares to each transit leg in the response
func (r *TripPlannerResponse) AddFares(f *fares.FareObjects, a *agency.Agencies) (err error) {
	for i, itinerary := range r.Plan.Itineraries {
		for l, leg := range itinerary.Legs {
			if leg.TransitLeg {
				agencyId := agency.AgencyId(TrimId(leg.AgencyId))
				noc, err := a.GetNoc(agencyId)
				if err != nil {
					return err
				}
				// use i and l to update original leg
				r.Plan.Itineraries[i].Legs[l].GetFares(f, noc)
			}
		}
	}
	return err
}

// GetFares finds the possible fares for the given transit leg and returns them
func (l *leg) GetFares(f *fares.FareObjects, n agency.Noc) (err error) {
	from := fares.Naptan(TrimId(l.From.StopId))
	to := fares.Naptan(TrimId(l.To.StopId))
	fareSlice := []fares.Fare{}
	for _, obj := range f.Objects {
		if obj.ContainsOpAndLine(n, l.RouteShortName) {
			if obj.ContainsStops(from, to) {
				fare, err := obj.GetFare(from, to)
				if err != nil && err != fares.ErrFareNotInTable {
					return err
				}
				if err != fares.ErrFareNotInTable {
					log.Printf("importing fare %s", fare.PreassignedFareProduct.Id)
					l.Fares = append(l.Fares, fare)
					fareSlice = append(fareSlice, fare)
				}
			}
		}
	}
	log.Printf("found %d fares for leg %s", len(fareSlice), l.RouteShortName)
	return err
}

// Removes the "2:" prefix present on IDs from the OTP API
func TrimId(id string) string {
	return strings.TrimPrefix(id, "2:")
}

func FilterTransitLegs(response TripPlannerResponse) []leg {
	var transitLegs []leg
	for _, itinerary := range response.Plan.Itineraries {
		for _, leg := range itinerary.Legs {
			if leg.TransitLeg {
				transitLegs = append(transitLegs, leg)
			}
		}
	}
	return transitLegs
}
