package struct2csv

import (
	"fmt"
	"strings"
	"testing"
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
	Lat   string
	Long  string
}

type Address struct {
	Addr1 string
	Addr2 string
	City  string
	State string
	Zip   string
}

type Notes map[string]string
type Stuff map[string]string

type Basic struct {
	Name string
	List []string
}

type Structor struct {
	ValueMapMap   map[string]map[string]string
	ValueMapSlice map[string][]string
	BasicMap map[string]Basic
	BasicSlice map[string][]Basic
}

type TTypes struct {
	Bool        bool
	Bools       []bool
	ABool       [2]bool
	Int         int
	Ints        []int
	AInt        [2]int
	ints        []int
	Int8        int8
	Int8s       []int8
	AInt8       [2]int8
	Int16       int16
	Int16s      []int16
	AInt16      [2]int16
	Int32       int32
	Int32s      []int32
	AInt32      [2]int32
	Int64       int64
	Int64s      []int64
	AInt64      [2]int64
	Uint        uint
	Uints       []uint
	AUint       [2]uint
	Uint8       uint8
	Uint8s      []uint8
	AUint8      [2]uint8
	Uint16      uint16
	Uint16s     []uint16
	AUint16     [2]uint16
	Uint32      uint32
	Uint32s     []uint32
	AUint32     [2]uint32
	Uint64      uint64
	Uint64s     []uint64
	AUint64     [2]uint64
	Float32     float32
	Float32s    []float32
	AFloat32    [2]float32
	Float64     float64
	Float64s    []float64
	AFloat64    [2]float64
	Complex64   complex64
	Complex64s  []complex64
	AComplex64  [2]complex64
	Complex128  complex128
	Complex128s []complex128
	AComplex128 [2]complex128
	Chan        chan int
	Chans       []chan int
	AChan       [2]chan int
	Func        func()
	Funcs       []func()
	AFunc       [2]func()
	String      string
	Strings     []string
	AString     [2]string
	strings     []string
}

type PtrTypes struct {
	Bool        *bool
	Bools       []*bool
	Int         *int
	Ints        []*int
	Uint        *uint
	Uints       []*uint
	Float64     *float64
	Float64s    []*float64
	Complex128  *complex128
	Complex128s []*complex128
	Chan        *chan int
	Chans       []*chan int
	Func        *func()
	Funcs       []*func()
	String      *string
	Strings     []*string
	strings     []*string
	BoolM map[string]*bool
	IntM map[string]*int
	Float64M map[string]*float64
	Complex128M map[int]*complex128
	ChanM map[string]*chan string
	FuncM map[string]*func()
	StringM map[string]*string
}

func TestNew(t *testing.T) {
	tc := New()
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
	tc := New()
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
	expectedHeaders = []string{"Name", "ID", "Addr1", "Addr2", "City", "State", "Zip", "Phone", "Lat", "Long", "Notes", "Stuff"}
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
			t.Errorf("%d: expected %q got %q", i, expectedHeaders[i], v)
		}
	}
}

func TestMarshal(t *testing.T) {
	tsts := []TTypes{
		TTypes{
			Bool:        true,
			Bools:       []bool{true, false, true},
			ABool:       [2]bool{true, false},
			Int:         42,
			Ints:        []int{72, 76, 88, 19, 2},
			AInt:        [2]int{1, 2},
			ints:        []int{1, 2, 3},
			Int8:        8,
			Int8s:       []int8{8, 9, 10},
			AInt8:       [2]int8{3, 4},
			Int16:       16,
			Int16s:      []int16{16, 17, 18},
			AInt16:      [2]int16{5, 6},
			Int32:       32,
			Int32s:      []int32{32, 33, 34},
			AInt32:      [2]int32{7, 8},
			Int64:       64,
			Int64s:      []int64{64, 65, 66},
			AInt64:      [2]int64{9, 10},
			Uint:        42,
			Uints:       []uint{72, 76, 88, 19, 2},
			AUint:       [2]uint{35, 21},
			Uint8:       8,
			Uint8s:      []uint8{8, 9, 10},
			AUint8:      [2]uint8{11, 12},
			Uint16:      16,
			Uint16s:     []uint16{16, 17, 18},
			AUint16:     [2]uint16{13, 14},
			Uint32:      32,
			Uint32s:     []uint32{32, 33, 34},
			AUint32:     [2]uint32{15, 16},
			Uint64:      64,
			Uint64s:     []uint64{64, 65, 66},
			AUint64:     [2]uint64{17, 18},
			Float32:     32.42,
			Float32s:    []float32{32.42, 33, 34.4},
			AFloat32:    [2]float32{35.5, 36.754},
			Float64:     64.42,
			Float64s:    []float64{64.42, 65, 66.7},
			AFloat64:    [2]float64{69.02, 69.0132},
			Complex64:   complex64(-64 + 12i),
			Complex64s:  []complex64{complex64(-65 + 11i), complex64(66 + 10i)},
			AComplex64:  [2]complex64{complex64(-61 + 21i), complex64(76 + 8i)},
			Complex128:  complex128(-128 + 12i),
			Complex128s: []complex128{complex128(-128 + 11i), complex128(129 + 10i)},
			AComplex128: [2]complex128{complex128(-118 + 2i), complex128(131 + 5i)},
			Chan:        make(chan int),
			Chans:       []chan int{make(chan int), make(chan int)},
			AChan:       [2]chan int{make(chan int), make(chan int)},
			Func:        func() { fmt.Println("hello") },
			Funcs:       []func(){func() { fmt.Println("hola") }},
			AFunc:       [2]func(){func() { fmt.Println("adios") }, func() { fmt.Println("au revior") }},
			String:      "don't panic",
			Strings:     []string{},
			AString:     [2]string{"pangalactic", "gargleblaster"},
			strings:     []string{"hello"},
		},
		TTypes{
			Bool:        true,
			Bools:       []bool{true, true, false},
			ABool:       [2]bool{true, false},
			Int:         420,
			Ints:        []int{1, 2, 3, 4},
			AInt:        [2]int{11, 12},
			ints:        []int{1, 2, 3},
			Int8:        18,
			Int8s:       []int8{18, 19, 110},
			AInt8:       [2]int8{13, 14},
			Int16:       116,
			Int16s:      []int16{116, 117, 118},
			AInt16:      [2]int16{15, 16},
			Int32:       132,
			Int32s:      []int32{132, 133, 134},
			AInt32:      [2]int32{17, 18},
			Int64:       640,
			Int64s:      []int64{164, 165, 166},
			AInt64:      [2]int64{19, 110},
			Uint:        421,
			Uints:       []uint{121, 122, 123},
			AUint:       [2]uint{35, 21},
			Uint8:       118,
			Uint8s:      []uint8{118, 119, 110},
			AUint8:      [2]uint8{111, 112},
			Uint16:      160,
			Uint16s:     []uint16{116, 117, 118},
			AUint16:     [2]uint16{113, 114},
			Uint32:      320,
			Uint32s:     []uint32{132, 133, 134},
			AUint32:     [2]uint32{115, 116},
			Uint64:      641,
			Uint64s:     []uint64{164, 165, 166},
			AUint64:     [2]uint64{117, 118},
			Float32:     132.42,
			Float32s:    []float32{132.42, 133, 134.4},
			AFloat32:    [2]float32{135.5, 136.754},
			Float64:     164.42,
			Float64s:    []float64{164.42, 165, 166.7},
			AFloat64:    [2]float64{169.02, 169.0132},
			Complex64:   complex64(-164 + 12i),
			Complex64s:  []complex64{complex64(-165 + 11i), complex64(166 + 10i)},
			AComplex64:  [2]complex64{complex64(-161 + 21i), complex64(176 + 8i)},
			Complex128:  complex128(-124 + 12i),
			Complex128s: []complex128{complex128(-126 + 11i), complex128(229 + 10i)},
			AComplex128: [2]complex128{complex128(-116 + 2i), complex128(231 + 5i)},
			Chan:        make(chan int),
			Chans:       []chan int{make(chan int), make(chan int)},
			AChan:       [2]chan int{make(chan int), make(chan int)},
			Func:        func() { fmt.Println("hello") },
			Funcs:       []func(){func() { fmt.Println("hola") }},
			AFunc:       [2]func(){func() { fmt.Println("adios") }, func() { fmt.Println("au revior") }},
			String:      "Towel",
			Strings:     []string{},
			AString:     [2]string{"Zaphod", "Beeblebrox"},
			strings:     []string{"hello"},
		},
	}
	expected := [][]string{
		[]string{"Bool", "Bools", "ABool",
			"Int", "Ints", "AInt",
			"Int8", "Int8s", "AInt8",
			"Int16", "Int16s", "AInt16",
			"Int32", "Int32s", "AInt32",
			"Int64", "Int64s", "AInt64",
			"Uint", "Uints", "AUint",
			"Uint8", "Uint8s", "AUint8",
			"Uint16", "Uint16s", "AUint16",
			"Uint32", "Uint32s", "AUint32",
			"Uint64", "Uint64s", "AUint64",
			"Float32", "Float32s", "AFloat32",
			"Float64", "Float64s", "AFloat64",
			"Complex64", "Complex64s", "AComplex64",
			"Complex128", "Complex128s", "AComplex128",
			"String", "Strings", "AString"},
		[]string{"true", "true, false, true", "true, false",
			"42", "72, 76, 88, 19, 2", "1, 2",
			"8", "8, 9, 10", "3, 4",
			"16", "16, 17, 18", "5, 6",
			"32", "32, 33, 34", "7, 8",
			"64", "64, 65, 66", "9, 10",
			"42", "72, 76, 88, 19, 2", "35, 21",
			"8", "8, 9, 10", "11, 12",
			"16", "16, 17, 18", "13, 14",
			"32", "32, 33, 34", "15, 16",
			"64", "64, 65, 66", "17, 18",
			"3.242E+01", "3.242E+01, 3.3E+01, 3.44E+01", "3.55E+01, 3.6754E+01",
			"6.442E+01", "6.442E+01, 6.5E+01, 6.67E+01", "6.902E+01, 6.90132E+01",
			"(-64+12i)", "(-65+11i), (66+10i)", "(-61+21i), (76+8i)",
			"(-128+12i)", "(-128+11i), (129+10i)", "(-118+2i), (131+5i)",
			"don't panic", "", "pangalactic, gargleblaster"},
		[]string{"true", "true, true, false", "true, false",
			"420", "1, 2, 3, 4", "11, 12",
			"18", "18, 19, 110", "13, 14",
			"116", "116, 117, 118", "15, 16",
			"132", "132, 133, 134", "17, 18",
			"640", "164, 165, 166", "19, 110",
			"421", "121, 122, 123", "35, 21",
			"118", "118, 119, 110", "111, 112",
			"160", "116, 117, 118", "113, 114",
			"320", "132, 133, 134", "115, 116",
			"641", "164, 165, 166", "117, 118",
			"1.3242E+02", "1.3242E+02, 1.33E+02, 1.344E+02", "1.355E+02, 1.36754E+02",
			"1.6442E+02", "1.6442E+02, 1.65E+02, 1.667E+02", "1.6902E+02, 1.690132E+02",
			"(-164+12i)", "(-165+11i), (166+10i)", "(-161+21i), (176+8i)",
			"(-124+12i)", "(-126+11i), (229+10i)", "(-116+2i), (231+5i)",
			"Towel", "", "Zaphod, Beeblebrox"},
	}
	tc := New()
	data, err := tc.Marshal(Tst{})
	if err != nil {
		if err.Error() != "slice required: value was of type struct" {
			t.Errorf("Expected \"slice of struct required: value was of type struct\", got %q", err)
		}
		goto NILSLICE
	}
	if err == nil {
		t.Error("Expected an error, got none")
	}
NILSLICE:
	var sl []string
	data, err = tc.Marshal(sl)
	if err != nil {
		if err.Error() != "slice cannot be nil" {
			t.Errorf("Expected \"slice cannot be nil\", got %q", err)
		}
		goto ZEROSLICE
	}
	if err == nil {
		t.Error("Expected an error, got none")
	}
ZEROSLICE:
	sl = make([]string, 0)
	data, err = tc.Marshal(sl)
	if err != nil {
		if err.Error() != "slice must have a length of at least 1: length was 0" {
			t.Errorf("Expected \"slice must have a length of at least 1: length was 0\", got %q", err)
		}
		goto NONSTRUCT
	}
	if err == nil {
		t.Error("Expected an error, got none")
	}
NONSTRUCT:
	sl = []string{"hello", "world"}
	data, err = tc.Marshal(sl)
	if err != nil {
		if err.Error() != "slice must be of type struct; type was string" {
			t.Errorf("Expected \"slice must be of type struct; type was string\", %q", err)
		}
		goto BASIC
	}
	if err == nil {
		t.Error("Expected an error, got none")
	}
BASIC:

	data, err = tc.Marshal(tsts)
	if err != nil {
		t.Errorf("expected no error, got %q", err)
		goto EMBED
	}
	if len(data) != len(expected) {
		t.Errorf("Expected %d rows, got %d", len(expected), len(data))
		goto EMBED
	}
	for i, row := range data {
		for j, col := range row {
			//t.Errorf("%d:%d\n\t%s\n\t%s\n", i, j, expected[i][j], col)
			if col != expected[i][j] {
				t.Errorf("%d:%d: expected %q, got %q", i, j, expected[i][j], col)
			}
		}
	}
EMBED:
}

func TestMarshalStructs(t *testing.T) {
	Tsts := []TstEmbed{
		TstEmbed{
			Name: "United Center",
			Location: Location{
				ID: 1,
				Address: Address{
					Addr1: "1901 W. Madison St.",
					City:  "Chicago",
					State: "IL",
					Zip:   "60612",
				},
				Phone: "(312) 455-4500",
				Lat:   "41.8806",
				Long:  "-87.6742",
			},
			Notes: map[string]string{"NHL": "Blackhawks", "NBA": "Bulls"},
		},
		TstEmbed{
			Name: "Wrigley Field",
			Location: Location{
				ID: 1906,
				Address: Address{
					Addr1: "1060 W. Addison St.",
					Addr2: "Broadcast Booth",
					City:  "Chicago",
					State: "IL",
					Zip:   "60613",
				},
				Phone: "(773) 404-2827",
				Lat:   "41.9483",
				Long:  "-87.6556",
			},
			Notes: map[string]string{"MLB": "Cubs"},
			Stuff: map[string]string{"Jack Brickhouse": "Hey Hey", "Harry Caray": "Holy Cow"},
		},
	}
	expected := [][]string{
		[]string{"Name", "ID", "Addr1", "Addr2", "City", "State", "Zip", "Phone", "Lat", "Long", "Notes", "Stuff"},
		[]string{"United Center", "1", "1901 W. Madison St.", "", "Chicago", "IL", "60612",
			"(312) 455-4500", "41.8806", "-87.6742", "NHL:Blackhawks, NBA:Bulls", ""},
		[]string{"Wrigley Field", "1906", "1060 W. Addison St.", "Broadcast Booth", "Chicago", "IL", "60613",
			"(773) 404-2827", "41.9483", "-87.6556", "MLB:Cubs", "Jack Brickhouse:Hey Hey, Harry Caray:Holy Cow"},
	}
	tc := New()
	rows, err := tc.Marshal(Tsts)
	if err != nil {
		t.Errorf("did not expect an error: got %q", err)
		return
	}
	if len(rows) != len(expected) {
		t.Errorf("Expected %d rows of data, got %d", len(expected), len(rows))
		return
	}
	for i, row := range rows {
		if len(row) != len(expected[i]) {
			t.Errorf("Expected row to have %d columns, got %d", len(expected[i]), len(row))
			continue
		}
		for j, col := range row {
			// these are map values so the order may change
			if j == 11 || j == 12 {
				tmp := StringParts(col)
				exp := StringParts(expected[i][j])
				if len(tmp) != len(exp) {
					t.Errorf("expected values to contain: %q, got values: %q", expected[i][j], col)
				}
				for _, tv := range tmp {
					for _, xp := range exp {
						if tv == xp {
							goto FOUND
						}
					}
					t.Errorf("expected values to contain: %q, got values: %q", expected[i][j], col)
					break
				FOUND:
				}
				continue
			}
		}
	}
}

func TestComplicated(t *testing.T) {
	tsts := []Structor{
		Structor{
			ValueMapMap: map[string]map[string]string{
				"Region 1": map[string]string{
					"colo1": "rack1",
					"colo2": "rack2",
				},
				"Region 2": map[string]string{
					"colo11": "rack11",
					"colo12": "rack12",
				},
			},
			ValueMapSlice: map[string][]string{
				"Canada": []string{"Alberta", "British Columbia", "Quebec"},
				"USA":    []string{"California", "Florida", "New York"},
			},
			BasicMap: map[string]Basic{
				"Gibson": Basic{Name: "William Gibson", List: []string{"Neuromancer", "Count Zero", "Mona Lisa Overdrive"}},
				"Herbert": Basic{Name: "Frank Herbert", List: []string{"Destination Void", "Jesus Incident", "Lazurus Effect"}},
			},
		BasicSlice: map[string][]Basic{
				"SciFi": []Basic{
					Basic{Name: "William Gibson", List: []string{"Neuromancer", "Count Zero", "Mona Lisa Overdrive"}},
					Basic{Name: "Frank Herbert", List: []string{"Destination Void", "Jesus Incident", "Lazurus Effect"}},
				},
			},
		},
	}
	expected := [][]string{
		[]string{"ValueMapMap", "ValueMapSlice", "BasicMap", "BasicSlice"},
		[]string{"Region 1:(colo1:rack1, colo2:rack2), Region 2:(colo11:rack11, colo12:rack12)",
			"Canada:(Alberta, British Columbia, Quebec), USA:(California, Florida, New York)",
			"Gibson:(William Gibson, [Neuromancer, Count Zero, Mona Lisa Overdrive]), Herbert:(Frank Herbert, [Destination Void, Jesus Incident, Lazurus Effect])",
		 	"SciFi:(William Gibson, [Neuromancer, Count Zero, Mona Lisa Overdrive], Frank Herbert, [Destination Void, Jesus Incident, Lazurus Effect])",
		},
	}

	tc := New()
	rows, err := tc.Marshal(tsts)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	}
	if len(rows) != len(expected) {
		t.Errorf("expected %d rows got %d", len(expected), len(rows))
		return
	}
	if len(rows[0]) != len(expected[0]) {
		t.Errorf("expected a row to have %d columns, got %d", len(expected[0]), len(rows[0]))
		return
	}
	for i, v := range rows[0] {
		if v != expected[0][i] {
			t.Errorf("Expected hdr column %d to be %q, got %q", i, expected[0][i], v)
		}
	}
	for i, v := range rows[1] {
		rvals := StringParts(v)
		evals := StringParts(expected[1][i])
		var found bool
		for _, v := range rvals {
			for _, vv := range evals {
				if vv == v {
					found = true
					goto FOUND
				}
			}
FOUND:
			if !found{
				t.Errorf("expected results to have values: %q, got %q", expected[1][i], v)
			}
			found = false
		}

	}
}

func TestPtrs(t *testing.T) {
	tst := []PtrTypes{
		PtrTypes{
			Bool:        new(bool),
			Bools:       []*bool{new(bool), new(bool)},
			Int:         new(int),
			Ints:        []*int{new(int), new(int)},
			Uint:        new(uint),
			Uints:       []*uint{new(uint), new(uint)},
			Float64:     new(float64),
			Float64s:    []*float64{new(float64), new(float64)},
			Complex128:  new(complex128),
			Complex128s: []*complex128{new(complex128), new(complex128)},
			Chan:        new(chan int),
			Chans:       []*chan int{new(chan int), new(chan int)},
			Func:        new(func()),
			Funcs:       []*func(){new(func())},
			String:      new(string),
			Strings:     []*string{new(string), new(string)},
			strings:     []*string{new(string)},
		},
	}
	expected := [][]string{
		[]string{"Bool", "Bools", "Int", "Ints", "Uint", "Uints", "Float64", "Float64s",
			"Complex128", "Complex128s", "String", "Strings",
			"BoolM", "IntM", "Float64M", "Complex128M", "StringM"},
		[]string{"false", "", "0", "", "0", "", "0+E00", "", "(0+0i)", "", "", "", "", "", "", ""},
	}
//	pbool(tst[0].Bool, true)
//	pbool(tst[0].Bools[0], true)
//	pbool(tst[0].Bools[1], true)
//	pchan(tst[0].Chan)
//	pchan(tst[0].Chans[0])
//	pchan(tst[0].Chans[1])
//	pfunc(tst[0].Func)
//	pfunc(tst[0].Funcs[0])
	tc := New()
	data, err := tc.Marshal(tst)
	t.Errorf("\n%#v\n%#v\n", expected, data)
	if err != nil {
		t.Errorf("unexpected err: %q", err)
	}

	for i, v := range data {
		if len(v) != len(expected[i]) {
			t.Errorf("%d: expected row to have %d cols, got %d", i, len(expected[i]), len(v))
		}
		for j, c := range v {
			if c != expected[i][j] {
				t.Errorf("Expected col value to be %q, got %q", expected[i][j], c)
			}
		}
		if i == 0 {
			break
		}
	}
}
// takes a string and returns it's parts:
// e.g. key1:("value1", "value2"), key2:("value11", "value12")
// would result in the following slice:
// []string("key1", "key1value1", "key1value1",
//           "key2", "key2value11", "key2value12")
// this makes it easier to make sure the results are as expected for strings
// built from maps and slices
func StringParts(s string) []string {
	if s == "" {
		return nil
	}
	var parts []string
	var key string
	tmp := strings.Split(s, "), ")
	tmp[len(tmp) -1] = strings.TrimSuffix(tmp[len(tmp)-1], ")")
	// get the key, which is followed by a :
	for _, v := range tmp {
		vals := strings.Split(v, ":")
		if len(vals) > 1 {
			key = vals[0]
			vals = append(vals[1:])
			parts = append(parts, key)
		}
		vals[0] = strings.TrimPrefix(vals[0], "(")
		for _, item := range vals {
			items := strings.Split(item, ", ")
			for _, vv := range items {
				parts = append(parts, fmt.Sprintf("%s%s", key, vv))
			}
		}
	}
	return parts

}
/*
func TestGetValueSliceType(t * testing.T) {
	v := interface{}(map[string]int{})
	val := reflect.ValueOf(v)
	keys := val.MapKeys()
	t.Errorf("%#v\n", getKeyType(keys).String())

}
*/
// funcs to set pointers
func pbool(p *bool, v bool) {
	*p = v
}

func pint(p *int, v int) {
	*p = v
}

func puint(p *uint, v uint) {
	*p = v
}

func pfloat64(p *float64, v float64) {
	*p = v
}

func pcomplex128(p *complex128, v complex128) {
	*p = v
}

func pstring(p *string, v string) {
	*p = v
}

func pchan(p *chan int) {
	*p = make(chan int)
}

func pfunc(p *func()) {
	*p = func(){fmt.Println("Hello")}
}
