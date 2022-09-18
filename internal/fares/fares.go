package fares

import (
	"encoding/xml"

	"github.com/imnatgreen/busfares/internal/agency"
)

type Noc = agency.Noc
type Naptan string
type Fare int

type FareObject struct {
	XMLName      xml.Name             `xml:"PublicationDelivery"`
	References   FareObjectReferences `xml:"PublicationRequest>topics>NetworkFrameTopic>NetworkFilterByValue>objectReferences"`
	ServiceFrame ServiceFrame         `xml:"dataObjects>CompositeFrame>frames>ServiceFrame"`
}

type FareObjectReferences struct {
	Operator AttrRef `xml:"OperatorRef"`
	Line     AttrRef `xml:"LineRef"`
}

type AttrRef struct {
	Ref string `xml:"ref,attr"`
}

type ServiceFrame struct {
	Lines               []Line               `xml:"lines>Line"`
	ScheduledStopPoints []ScheduledStopPoint `xml:"scheduledStopPoints>ScheduledStopPoint"`
}

type Line struct {
	LineRef     string  `xml:"id,attr"`
	PublicCode  string  `xml:"PublicCode"`
	PrivateCode string  `xml:"PrivateCode"`
	OperatorRef AttrRef `xml:"OperatorRef"`
}

type ScheduledStopPoint struct {
	ScheduledStopPointRef string               `xml:"id,attr"`
	Name                  string               `xml:"Name"`
	TopographicPlaceView  TopographicPlaceView `xml:"TopographicPlaceView"`
}

type TopographicPlaceView struct {
	TopographicPlaceViewRef AttrRef `xml:"TopographicPlaceRef"`
	Name                    string  `xml:"Name"`
}

// Breakdown of Transdev single fare xml structure:
// 	CompositeFrame(resoponsibilitySetRef="tarrifs")>frames>-ServiceFrame>-lines>Line(id)>-PublicCode
// 																																											>-PrivateCode
// 																																											>-OperatorRef
// 																																			>-scheduledStopPoints>-ScheduledStopPoint(id)>-Name
// 																																																									 >-TopographicPlaceView>-TopographicPlaceViewRef(ref)
// 																																																																				 >-Name

// 																										TODO>-FareFrame(responsibilitySetRef="network_data")>fareZones>FareZone(id)>-Name
// 																												                                                                       >-members>ScheduledStopPointRef(ref)
// 																										TODO>-FareFrame(responsibilitySetRef="tarrifs")>-tarrifs
// 																																																	 >-usageParameters
// 																																																	 >-fareProducts
// 																																																	 >-salesOfferPackages
// 																										TODO>-FareFrame(responsibilitySetRef="tarrifs")>-priceGroups
// 																												                                           >-fareTables

// potential api for live vehicle data?
// https://www.transdevbus.co.uk/_ajax/vehicles
// https://www.transdevbus.co.uk/_ajax/stops/2500LAA15762/vehicles

func ParseFares() {
	// TODO
}

func ParseXml(data []byte) (FareObject, error) {
	var obj FareObject
	err := xml.Unmarshal(data, &obj)
	return obj, err
}

func GetFareFromFareObject(fareObject FareObject, noc Noc, from Naptan, to Naptan) (Fare, error) {
	// TODO
	return 0, nil
}
