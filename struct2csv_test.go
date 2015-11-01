package struct2csv

import (
	_ "os"
	"testing"

	_ "github.com/mohae/customjson"
)

type Tst struct {
	Int          int
	Ints         []int
	ints         []int
	String       string
	Strings      []string
	StringString map[string]string
	StringInt    map[string]int
	IntInt       map[int]int
	strings      []string
}

type TstTags struct {
	Int          int   `json:"number" csv:"Number"`
	Ints         []int `json:"numbers" csv:"Numbers"`
	ints         []int
	String       string            `json:"word" csv:"Word"`
	Strings      []string          `json:"words" csv:"Words"`
	StringString map[string]string `json:"mapstringstring" csv:"MapStringString"`
	StringInt    map[string]int    `json:"mapstringint" csv:"MapStringInt"`
	IntInt       map[int]int       `json:"mapintint" csv:"MapIntInt"`
	strings      []string
}

type TstEmbed struct {
	Name string
	Location
	Notes
	Stuff
}

type Location struct {
	ID int
	Address
	Phone string
	Long  string
	Lat   string
}

type Address struct {
	Addr1 string
	Addr2 string
	City  string
	State string
	Zip   string
}

type Notes map[string]string
type Stuff []string

func TestNewTranscoder(t *testing.T) {
	tc := NewTranscoder()
	if tc.useTags != true {
		t.Errorf("expected useTags to be true got %t", tc.useTags)
	}
	if tc.tag != "csv" {
		t.Errorf("expected transcoder's tag to be \"csv\", got %q", tc.tag)
	}
	tc.SetUseTags(false)
	if tc.useTags != false {
		t.Errorf("expected useTags to be false got %t", tc.useTags)
	}
	tc.SetTag("json")
	if tc.tag != "json" {
		t.Errorf("expected transcoder's tag to be \"json\", got %q", tc.tag)
	}
	tc.SetUseTags(true)
	if tc.useTags != true {
		t.Errorf("expected useTags to be false got %t", tc.useTags)
	}
	tc.SetTag("csv")
	if tc.tag != "csv" {
		t.Errorf("expected transcoder's tag to be \"csv\", got %q", tc.tag)
	}
}
func TestGetHeaders(t *testing.T) {
	_, err := GetHeaders([]string{"a", "b", "c"})
	if err == nil {
		t.Error("expected passing of a non struct to result in an error, none received")
	} else {
		if err.Error() != "struct required: value was of type slice" {
			t.Errorf("expected error to be \"struct required: value was of type slice\", got %q", err)
		}
	}
	tc := NewTranscoder()
	tc.useTags = false
	tst := Tst{}
	expectedHeaders := []string{"Int", "Ints", "String", "Strings", "StringString", "StringInt", "IntInt"}
	hdr, err := tc.GetHeaders(tst)
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto IGNORETAG
	}
	if len(hdr) != len(expectedHeaders) {
		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
		goto IGNORETAG
	}
	for i, v := range hdr {
		if v != expectedHeaders[i] {
			t.Errorf("%d: expected %q got %q", i, expectedHeaders[i], v)
		}
	}

IGNORETAG:
	tc.useTags = false
	test := TstTags{}
	expectedHeaders = []string{"Int", "Ints", "String", "Strings", "StringString", "StringInt", "IntInt"}
	hdr, err = tc.GetHeaders(test)
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto CSVTAG
	}
	if len(hdr) != len(expectedHeaders) {
		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
		goto CSVTAG
	}
	for i, v := range hdr {
		if v != expectedHeaders[i] {
			t.Errorf("%d: expected %q got %q", i, expectedHeaders[i], v)
		}
	}

CSVTAG:
	// test using CSV tags
	tc.useTags = true
	expectedHeaders = []string{"Number", "Numbers", "Word", "Words", "MapStringString", "MapStringInt", "MapIntInt"}
	hdr, err = tc.GetHeaders(TstTags{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto JSONTAG
	}
	if len(hdr) != len(expectedHeaders) {
		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
		goto JSONTAG
	}
	for i, v := range hdr {
		if v != expectedHeaders[i] {
			t.Errorf("%d: expected %q got %q", i, expectedHeaders[i], v)
		}
	}

JSONTAG:
	// test using CSV tags
	expectedHeaders = []string{"number", "numbers", "word", "words", "mapstringstring", "mapstringint", "mapintint"}
	tc.useTags = true
	tc.tag = "json"
	hdr, err = tc.GetHeaders(TstTags{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto EMBED
	}
	if len(hdr) != len(expectedHeaders) {
		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
		goto EMBED
	}
	for i, v := range hdr {
		if v != expectedHeaders[i] {
			t.Errorf("%d: expected %q got %q", i, expectedHeaders[i], v)
		}
	}

EMBED:
	expectedHeaders = []string{"Name", "ID", "Addr1", "Addr2", "City", "State", "Zip", "Phone", "Long", "Lat", "Notes", "Stuff"}
	hdr, err = tc.GetHeaders(TstEmbed{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		return
	}
	if len(hdr) != len(expectedHeaders) {
		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
		t.Errorf("%#v\n", hdr)
		return
	}
	for i, v := range hdr {
		if v != expectedHeaders[i] {
			t.Errorf("%d: expected %q got %q", expectedHeaders[i], v)
		}
	}
}
