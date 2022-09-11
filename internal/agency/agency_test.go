package agency

import (
	"bytes"
	"io"
	"testing"
)

func TestSearch(t *testing.T) {
	agencies := Agencies{"OP291": "ROST"}
	t.Run("known id", func(t *testing.T) {
		got, _ := agencies.GetNoc("OP291")
		want := Noc("ROST")
		assertNoc(t, got, want)
	})
	t.Run("unknown id", func(t *testing.T) {
		_, got := agencies.GetNoc("unknown")
		assertError(t, got, ErrNocNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Run("new word", func(t *testing.T) {
		agencies := Agencies{}
		id := AgencyId("OP291")
		noc := Noc("ROST")

		err := agencies.Add(id, noc)

		assertError(t, err, nil)
		assertNocFromId(t, agencies, id, noc)
	})
	t.Run("existing word", func(t *testing.T) {
		id := AgencyId("OP291")
		noc := Noc("ROST")
		agencies := Agencies{id: noc}
		err := agencies.Add(id, "new test")

		assertError(t, err, ErrNocExists)
		assertNocFromId(t, agencies, id, noc)
	})
}

func TestAddFromCsv(t *testing.T) {
	buffer := bytes.Buffer{}
	buffer.WriteString(`agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_noc
OP291,"Rosso","https://www.traveline.info",Europe/London,EN,"","ROST"`)
	closeableBuffer := io.NopCloser(&buffer)

	agencies := make(Agencies)

	agencies.AddFromCsv(closeableBuffer)
	got, _ := agencies.GetNoc("OP291")
	want := Noc("ROST")
	assertNoc(t, got, want)
}

func assertNoc(t testing.TB, got, want Noc) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertNocFromId(t testing.TB, agencies Agencies, id AgencyId, noc Noc) {
	t.Helper()
	got, err := agencies.GetNoc(id)
	if err != nil {
		t.Fatal("shoud find added NOC:", err)
	}
	if got != noc {
		t.Errorf("got %q want %q", got, noc)
	}
}
