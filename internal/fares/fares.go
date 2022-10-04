package fares

import (
	"encoding/xml"
	"time"

	"github.com/bojanz/currency"
	"github.com/imnatgreen/busfares/internal/agency"
)

type Noc = agency.Noc
type Naptan string

// XmlTime allows timestamps in FareXChange XML to be unmarshalled into time.Time objects.
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
	Currency            string                   `xml:"dataObjects>CompositeFrame>FrameDefaults>DefaultCurrency"`
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
	Id     string `xml:"id,attr"`
	Amount string `xml:"Amount"`
}

type FareTable struct {
	TariffRef                 AttrRef              `xml:"usedIn>TariffRef"`
	UserProfileRef            AttrRef              `xml:"pricesFor>UserProfileRef"`
	SalesOfferPackageRef      AttrRef              `xml:"pricesFor>SalesOfferPackageRef"`
	PreassignedFareProductRef AttrRef              `xml:"pricesFor>PreassignedFareProductRef"`
	Columns                   []FareTableRowColumn `xml:"columns>FareTableColumn"`
	Rows                      []FareTableRowColumn `xml:"rows>FareTableRow"`
	Cells                     []FareTableCell      `xml:"includes>FareTable>cells>Cell"`
}

type FareTableRowColumn struct {
	Id           string    `xml:"id,attr"`
	Order        string    `xml:"order,attr"`
	Name         string    `xml:"Name"`
	FareZoneRefs []AttrRef `xml:"representing>FareZoneRef"` // only in columns
}

type FareTableCell struct {
	Id                         string                     `xml:"id,attr"`
	DistanceMatrixElementPrice DistanceMatrixElementPrice `xml:"DistanceMatrixElementPrice"`
	ColumnRef                  AttrRef                    `xml:"ColumnRef"`
	RowRef                     AttrRef                    `xml:"RowRef"`
}

type DistanceMatrixElementPrice struct {
	Id                           string  `xml:"id,attr"`
	GeographicalIntervalPriceRef AttrRef `xml:"GeographicalIntervalPriceRef"`
	DistanceMatrixElementRef     AttrRef `xml:"DistanceMatrixElementRef"`
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

type FareZone struct {
	Id      string                  `xml:"id,attr"`
	Name    string                  `xml:"Name"`
	Members []ScheduledStopPointRef `xml:"members>ScheduledStopPointRef"`
}

type ScheduledStopPointRef struct {
	Ref  string `xml:"ref,attr"`
	Name string `xml:",chardata"`
}

func ParseXml(data []byte) (FareObject, error) {
	var obj FareObject
	err := xml.Unmarshal(data, &obj)
	return obj, err
}

func (f *FareObject) GetFare(from, to Naptan) (fare currency.Amount, err error) {
	// get fare zones from stops
	fromRef := string("atco:" + from)
	toRef := string("atco:" + to)
	var fromZone, toZone string
	for _, fareZone := range f.FareZones {
		for _, member := range fareZone.Members {
			if member.Ref == fromRef {
				fromZone = fareZone.Id
				break
			}
			if member.Ref == toRef {
				toZone = fareZone.Id
				break
			}
		}
	}

	// get distance matrix element id from fare zones
	var distanceMatrixElementId string
	for _, d := range f.Tariffs[0].DistanceMatrixElements {
		if d.StartTariffZoneRef.Ref == fromZone && d.EndTariffZoneRef.Ref == toZone {
			distanceMatrixElementId = d.Id
			break
		}
	}

	// get price group id from distance matrix element id
	var geographicalIntervalPriceRef string
	for _, c := range f.FareTables[0].Cells {
		if c.DistanceMatrixElementPrice.DistanceMatrixElementRef.Ref == distanceMatrixElementId {
			geographicalIntervalPriceRef = c.DistanceMatrixElementPrice.GeographicalIntervalPriceRef.Ref
			break
		}
	}

	// get price from price group id
	for _, p := range f.PriceGroups {
		if p.GeographicalIntervalPrice[0].Id == geographicalIntervalPriceRef {
			fare, err = currency.NewAmount(p.GeographicalIntervalPrice[0].Amount, f.Currency)
			break
		}
	}

	return fare, err
}

func (f *FareObject) ContainsOpAndLine(op Noc, lineCode string) bool {
	for _, line := range f.Lines {
		if line.OperatorRef.Ref == "noc:"+string(op) && line.PublicCode == lineCode {
			return true
		}
	}
	return false
}

func (f *FareObject) ContainsStops(from, to Naptan) bool {
	fromRef := string("atco:" + from)
	toRef := string("atco:" + to)

	fromExists, toExists := false, false

	for _, stop := range f.ScheduledStopPoints {
		if stop.ScheduledStopPointRef == fromRef {
			fromExists = true
		}
		if stop.ScheduledStopPointRef == toRef {
			toExists = true
		}
		if fromExists && toExists {
			break
		}
	}

	return fromExists && toExists
}