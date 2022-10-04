package router

import (
	"encoding/json"
	"strings"

	"github.com/bojanz/currency"
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
	Fare                 legFare
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

type legFare struct {
	Amount currency.Amount
	// TODO: user profiles, how to purchase, etc.
}

func ParseJson(data []byte) (TripPlannerResponse, error) {
	var response TripPlannerResponse
	err := json.Unmarshal(data, &response)
	return response, err
}

func (r *TripPlannerResponse) GetFares(f *fares.FareObjects, a *agency.Agencies) (err error) {
	for _, itinerary := range r.Plan.Itineraries {
		for _, leg := range itinerary.Legs {
			if leg.TransitLeg {
				agencyId := agency.AgencyId(TrimId(leg.AgencyId))
				noc, err := a.GetNoc(agencyId)
				if err != nil {
					return err
				}
				leg.GetFare(f, noc) // TODO: pass user profiles
			}
		}
	}
	return err
}

func (l *leg) GetFare(f *fares.FareObjects, n agency.Noc) (err error) {
	from := fares.Naptan(TrimId(l.From.StopId))
	to := fares.Naptan(TrimId(l.To.StopId))

	for _, obj := range f.Objects {
		if obj.ContainsOpAndLine(n, l.RouteShortName) {
			if obj.ContainsStops(from, to) {
				if obj.FareProducts[0].ProductType == "singleTrip" && obj.FareProducts[0].Id == "Trip@Single" {
					fare, err := obj.GetFare(from, to)
					if err != nil {
						return err
					}
					l.Fare.Amount = fare
					return nil
				}
			}
		}
	}
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
