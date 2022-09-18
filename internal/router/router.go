package router

import (
	"encoding/json"
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