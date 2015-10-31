package struct2csv

import (
	_ "os"
	"testing"

	_ "github.com/mohae/customjson"
)

type TstTags struct {
	Int     int      `json:"number" csv:"Number"`
	Ints    []int    `json:"numbers" csv:"Numbers"`
  ints    []int
	String  string   `json:"word" csv:"Word"`
	Strings []string `json:"words" csv:"Words"`
  StringString  map[string]string `json:"mapstringstring" csv:"MapStringString"`
  StringInt  map[string]int `json:"mapstringInt" csv:"MapStringInt"`
  IntInt     map[int]int `json:"mapintint" csv:"MapIntInt"`
  strings []string
}

type Tst struct {
	Int     int
	Ints    []int
  ints    []int
	String  string
	Strings []string
  strings []string
  StringString  map[string]string
  StringInt  map[string]int
  IntInt     map[int]int
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
	test := Tst{}
	expectedHeaders := []string{"Int", "Ints", "String", "Strings", "StringString", "StringInt", "IntInt"}
	hdr, err := GetHeaders(test)
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
    goto TAG
	}
  if len(hdr) != len(expectedHeaders) {
		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
		goto TAG
	}
	for i, v := range hdr {
		if v != expectedHeaders[i] {
			t.Errorf("%d: expected %q got %q", expectedHeaders[i], v)
		}
	}

TAG:
   // test using CSV tags
   expectedHeaders = []string{"Number", "Numbers", "Word", "Words", "MapStringString", "MapStringInt", "MapIntInt"}
   hdr, err = GetHeaders(TstTags{})
 	if err != nil {
 		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
     return
 	}
   if len(hdr) != len(expectedHeaders) {
 		t.Errorf("Expected %d column headers, got %d", len(expectedHeaders), len(hdr))
 		return
 	}
 	for i, v := range hdr {
 		if v != expectedHeaders[i] {
 			t.Errorf("%d: expected %q got %q", expectedHeaders[i], v)
 		}
 	}
}
