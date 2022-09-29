package fares

import (
	"encoding/xml"
	"time"

	"github.com/imnatgreen/busfares/internal/agency"
)

type Noc = agency.Noc
type Naptan string
type Fare int

type XmlTime struct {
	time.Time
}

func (x *XmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	error := d.DecodeElement(&v, &start)
	if error != nil {
		return error
	}
	parse, error := time.Parse("2006-01-02T15:04:05", v)
	*x = XmlTime{parse}
	return error
}

type FareObject struct {
	XMLName             xml.Name                 `xml:"PublicationDelivery"`
	References          FareObjectReferences     `xml:"PublicationRequest>topics>NetworkFrameTopic>NetworkFilterByValue>objectReferences"`
	Operators           []Operator               `xml:"dataObjects>CompositeFrame>frames>ResourceFrame>organisations>Operator"`
	Lines               []Line                   `xml:"dataObjects>CompositeFrame>frames>ServiceFrame>lines>Line"`
	ScheduledStopPoints []ScheduledStopPoint     `xml:"dataObjects>CompositeFrame>frames>ServiceFrame>scheduledStopPoints>ScheduledStopPoint"`
	FareZones           []FareZone               `xml:"dataObjects>CompositeFrame>frames>FareFrame>fareZones>FareZone"`
	Tariffs             []Tariff                 `xml:"dataObjects>CompositeFrame>frames>FareFrame>tariffs>Tariff"`
	UserProfiles        []UserProfile            `xml:"dataObjects>CompositeFrame>frames>FareFrame>usageParameters>UserProfile"`
	FareProducts        []PreassignedFareProduct `xml:"dataObjects>CompositeFrame>frames>FareFrame>fareProducts>PreassignedFareProduct"`
	SalesOfferPackages  []SalesOfferPackage      `xml:"dataObjects>CompositeFrame>frames>FareFrame>salesOfferPackages>SalesOfferPackage"`
	PriceGroups         []PriceGroup             `xml:"dataObjects>CompositeFrame>frames>FareFrame>priceGroups>PriceGroup"`
	FareTables          []FareTable              `xml:"dataObjects>CompositeFrame>frames>FareFrame>fareTables>FareTable"`
}

type Operator struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"Name"`
}
type FareObjectReferences struct {
	Operator AttrRef `xml:"OperatorRef"`
	Line     AttrRef `xml:"LineRef"`
}

type AttrRef struct {
	Ref string `xml:"ref,attr"`
}

type Tariff struct {
	Id                     string                  `xml:"id,attr"`
	Name                   string                  `xml:"Name"`
	OperatorRef            AttrRef                 `xml:"OperatorRef"`
	LineRef                AttrRef                 `xml:"LineRef"`
	TariffBasis            string                  `xml:"TariffBasis"`
	DistanceMatrixElements []DistanceMatrixElement `xml:"fareStructureElements>FareStructureElement>distanceMatrixElements>DistanceMatrixElement"`
	TripType               string                  `xml:"fareStructureElements>FareStructureElement>GenericParameterAssignment>limitations>RoundTrip>TripType"`
	ValidBetween           ValidBetween            `xml:"validityConditions>ValidBetween"`
}

type DistanceMatrixElement struct {
	Id                 string  `xml:"id,attr"`
	StartTariffZoneRef AttrRef `xml:"StartTariffZoneRef"`
	EndTariffZoneRef   AttrRef `xml:"EndTariffZoneRef"`
}

type ValidBetween struct {
	FromDate XmlTime `xml:"FromDate"`
	ToDate   XmlTime `xml:"ToDate"`
}

type UserProfile struct {
	Id         string `xml:"id,attr"`
	Name       string `xml:"Name"`
	UserType   string `xml:"UserType"`
	MinimumAge string `xml:"MinimumAge"`
	MaximumAge string `xml:"MaximumAge"`
}

type PreassignedFareProduct struct {
	Id             string `xml:"id,attr"`
	Name           string `xml:"Name"`
	ChargingMoment string `xml:"ChargingMomentType"`
	TariffBasis    string `xml:"ConditionSummary>TariffBasis"`
	ProductType    string `xml:"ProductType"`
}

type SalesOfferPackage struct {
	Id          string `xml:"id,attr"`
	Name        string `xml:"Name"`
	Description string `xml:"Description"`
}

type PriceGroup struct {
	Id                        string                      `xml:"id,attr"`
	GeographicalIntervalPrice []GeographicalIntervalPrice `xml:"members>GeographicalIntervalPrice"`
}

type GeographicalIntervalPrice struct {
	Id     string  `xml:"id,attr"`
	Amount float32 `xml:"Amount"`
}

type FareTable struct {
	TariffRef                 AttrRef `xml:"usedIn>TariffRef"`
	UserProfileRef            AttrRef `xml:"pricesFor>UserProfileRef"`
	SalesOfferPackageRef      AttrRef `xml:"pricesFor>SalesOfferPackageRef"`
	PreassignedFareProductRef AttrRef `xml:"pricesFor>PreassignedFareProductRef"`
}

// type ServiceFrame struct {
// 	Lines               []Line               `xml:"lines>Line"`
// 	ScheduledStopPoints []ScheduledStopPoint `xml:"scheduledStopPoints>ScheduledStopPoint"`
// }

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

// type FareFrame struct {
// 	// Frame 1
// 	// ResponsibilitySetRef string     `xml:"responsibilitySetRef,attr"`
// 	// TypeOfFrameRef       AttrRef    `xml:"TypeOfFrameRef"`
// 	FareZones            []FareZone `xml:"fareZones>FareZone"`
// 	// Frame 2
// 	Tarrifs            []Tarrif                 `xml:"tarrifs>Tarrif"`
// 	UsageParameters    []UserProfile            `xml:"usageParameters>UserProfile"`
// 	FareProducts       []PreassignedFareProduct `xml:"fareProducts>PreassignedFareProduct"`
// 	SalesOfferPackages []SalesOfferPackage      `xml:"salesOfferPackages>SalesOfferPackage"`
// 	// Frame 3
//   PriceGroups []PriceGroup `xml:"priceGroups>PriceGroup"`
// 	FareTables  []FareTable  `xml:"fareTables>FareTable"`
// }

// https://github.com/jclgoodwin/bustimes.org/blob/main/fares/management/commands/import_netex_fares.py
// https://github.com/antchfx/xmlquery

type FareZone struct {
	Id      string                  `xml:"id,attr"`
	Name    string                  `xml:"Name"`
	Members []ScheduledStopPointRef `xml:"members>ScheduledStopPointRef"`
}

type ScheduledStopPointRef struct {
	Ref  string `xml:"ref,attr"`
	Name string `xml:",chardata"`
}

//    TODO>-FareFrame(responsibilitySetRef="tarrifs")>-tarrifs
// 																									 >-usageParameters
// 					  																			 >-fareProducts
// 																									 >-salesOfferPackages
// 		TODO>-FareFrame(responsibilitySetRef="tarrifs")>-priceGroups
// 				                                           >-fareTables

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
