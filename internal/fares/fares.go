package fares

import (
	"encoding/xml"

	"github.com/imnatgreen/busfares/internal/agency"
)

type Noc = agency.Noc
type Naptan string
type Fare int

type FareObject struct {
	XMLName    xml.Name             `xml:"PublicationDelivery"`
	References FareObjectReferences `xml:"PublicationRequest>topics>NetworkFrameTopic>NetworkFilterByValue>objectReferences"`
}

type FareObjectReferences struct {
	Operator string `xml:"OperatorRef>ref,attr"`
	Line     string `xml:"LineRef>ref,attr"`
}

// https://www.transdevbus.co.uk/_ajax/vehicles
// https://www.transdevbus.co.uk/_ajax/stops/2500LAA15762/vehicles

func ParseFares() {
	// TODO
}

func GetFareFromFareObject(fareObject FareObject, noc Noc, from Naptan, to Naptan) (Fare, error) {
	// TODO
	return 0, nil
}
