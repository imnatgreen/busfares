package main

import (
	"encoding/csv"
	"io"
)

const (
	ErrNocNotFound = AgencyErr("could not find NOC for given agency_id")
	ErrNocExists   = AgencyErr("NOC already exists, skipping")
)

type AgencyErr string

func (e AgencyErr) Error() string {
	return string(e)
}

type AgencyId string
type Noc string

type Agencies map[AgencyId]Noc

func (a Agencies) GetNoc(id AgencyId) (Noc, error) {
	definition, ok := a[id]
	if !ok {
		return "", ErrNocNotFound
	}
	return definition, nil
}

func (a Agencies) Add(id AgencyId, noc Noc) error {
	_, err := a.GetNoc(id)

	switch err {
	case ErrNocNotFound:
		a[id] = noc
	case nil:
		return ErrNocExists
	default:
		return err
	}

	return nil
}

func (a Agencies) AddFromCsv(file io.ReadCloser) error {
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Read() // skip first line of column headers
	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		a.Add(AgencyId(record[0]), Noc(record[6]))
	}
	return nil
}
