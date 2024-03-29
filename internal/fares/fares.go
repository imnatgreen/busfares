package fares

import (
	"encoding/xml"
	"errors"
	"time"

	"github.com/bojanz/currency"
	"github.com/imnatgreen/busfares/internal/agency"
)

type Noc = agency.Noc
type Naptan string

// XmlTime allows timestamps in FareXChange XML to be unmarshalled into time.Time objects.
type XmlTime struct {
	time.Time `json:"time"`
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
	FromDate XmlTime `xml:"FromDate" json:"fromDate"`
	ToDate   XmlTime `xml:"ToDate" json:"toDate"`
}

type UserProfile struct {
	Id         string `xml:"id,attr" json:"id"`
	Name       string `xml:"Name" json:"name"`
	UserType   string `xml:"UserType" json:"userType"`
	MinimumAge string `xml:"MinimumAge" json:"minimumAge"`
	MaximumAge string `xml:"MaximumAge" json:"maximumAge"`
}

type PreassignedFareProduct struct {
	Id             string `xml:"id,attr" json:"id"`
	Name           string `xml:"Name" json:"name"`
	ChargingMoment string `xml:"ChargingMomentType" json:"chargingMoment"`
	TariffBasis    string `xml:"ConditionSummary>TariffBasis" json:"tariffBasis"`
	ProductType    string `xml:"ProductType" json:"productType"`
}

type SalesOfferPackage struct {
	Id          string `xml:"id,attr" json:"id"`
	Name        string `xml:"Name" json:"name"`
	Description string `xml:"Description" json:"description"`
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

// Type Fare is used when finding fares from an imported FareObject
type Fare struct {
	Amount                 currency.Amount        `json:"amount"`
	ValidBetween           ValidBetween           `json:"validBetween"`
	UserProfile            UserProfile            `json:"userProfile"`
	SalesOfferPackage      SalesOfferPackage      `json:"salesOfferPackage"`
	PreassignedFareProduct PreassignedFareProduct `json:"preassignedFareProduct"`
}

func ParseXml(data []byte) (FareObject, error) {
	var obj FareObject
	err := xml.Unmarshal(data, &obj)
	return obj, err
}

var ErrFareNotInTable = errors.New("fare not in table")

func (f *FareObject) GetFare(from, to Naptan) (fare Fare, err error) {
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
	// for some reason, some fare files have the ids the wrong way round, so we need to check with fare zones reversed too
	var distanceMatrixElementIds []string

	for _, d := range f.Tariffs[0].DistanceMatrixElements {
		if (d.StartTariffZoneRef.Ref == fromZone && d.EndTariffZoneRef.Ref == toZone) ||
			(d.StartTariffZoneRef.Ref == toZone && d.EndTariffZoneRef.Ref == fromZone) {
			distanceMatrixElementIds = append(distanceMatrixElementIds, d.Id)
			break
		}
	}

	// get price group id from distance matrix element id
	// also get fare related refs if price group found
	var geographicalIntervalPriceRef string
	var userProfileRef string
	var salesOfferPackageRef string
	var preassignedFareProductRef string
	var tariffRef string
	for _, d := range distanceMatrixElementIds {
		for _, t := range f.FareTables {
			for _, c := range t.Cells {
				if c.DistanceMatrixElementPrice.DistanceMatrixElementRef.Ref == d {
					geographicalIntervalPriceRef = c.DistanceMatrixElementPrice.GeographicalIntervalPriceRef.Ref
					userProfileRef = t.UserProfileRef.Ref
					salesOfferPackageRef = t.SalesOfferPackageRef.Ref
					preassignedFareProductRef = t.PreassignedFareProductRef.Ref
					tariffRef = t.TariffRef.Ref
					break
				}
			}
		}
	}
	// check if fare found, if not return empty fare
	if geographicalIntervalPriceRef == "" {
		return fare, ErrFareNotInTable
	}

	// get price from price group id
	var fareAmount currency.Amount
	for _, p := range f.PriceGroups {
		if p.GeographicalIntervalPrice[0].Id == geographicalIntervalPriceRef {
			fareAmount, err = currency.NewAmount(p.GeographicalIntervalPrice[0].Amount, f.Currency)
			break
		}
	}

	// get user profile
	var userProfile UserProfile
	for _, u := range f.UserProfiles {
		if u.Id == userProfileRef {
			userProfile = u
			break
		}
	}

	// get sales offer package
	var salesOfferPackage SalesOfferPackage
	for _, s := range f.SalesOfferPackages {
		if s.Id == salesOfferPackageRef {
			salesOfferPackage = s
			break
		}
	}

	// get preassigned fare product
	var preassignedFareProduct PreassignedFareProduct
	for _, p := range f.FareProducts {
		if p.Id == preassignedFareProductRef {
			preassignedFareProduct = p
			break
		}
	}

	// get ValidBetween
	var validBetween ValidBetween
	for _, t := range f.Tariffs {
		if t.Id == tariffRef {
			validBetween = t.ValidBetween
			break
		}
	}

	fare = Fare{
		Amount:                 fareAmount,
		ValidBetween:           validBetween,
		UserProfile:            userProfile,
		SalesOfferPackage:      salesOfferPackage,
		PreassignedFareProduct: preassignedFareProduct,
	}

	return fare, err
}
