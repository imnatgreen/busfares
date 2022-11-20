package colours

import (
	"encoding/csv"
	"io"
)

const (
	ErrNocNotFound  = ColoursErr("could not find NOC")
	ErrLineNotFound = ColoursErr("could not find line")
	ErrNocExists    = ColoursErr("NOC already exists, skipping")
	ErrLineExists   = ColoursErr("line already exists, skipping")
)

type ColoursErr string

func (e ColoursErr) Error() string {
	return string(e)
}

type Noc string
type Line string
type Colour string

type Colours map[Noc]map[Line]Colour

func (c Colours) Get(n Noc, l Line) (Colour, error) {
	_, ok := c[n]
	if !ok {
		return "", ErrNocNotFound
	}
	col, ok := c[n][l]
	if !ok {
		return "", ErrLineNotFound
	}
	return col, nil
}

func (c Colours) Add(n Noc, l Line, col Colour) error {
	_, err := c.Get(n, l)

	switch err {
	case ErrNocNotFound:
		c[n] = map[Line]Colour{}
		c[n][l] = col
	case ErrLineNotFound:
		c[n][l] = col
	case nil:
		return ErrNocExists
	default:
		return err
	}

	return nil
}

func (c Colours) AddFromCsv(file io.ReadCloser) error {
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

		c.Add(Noc(record[0]), Line(record[1]), Colour(record[2]))
	}
	return nil
}
