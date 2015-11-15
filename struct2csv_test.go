package struct2csv

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type Tags struct {
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

var (
	expectedTagCols     = []string{"Int", "Ints", "String", "Strings", "StringString", "StringInt", "IntInt"}
	expectedTagJSONCols = []string{"number", "numbers", "word", "words", "mapstringstring", "mapstringint", "mapintint"}
	expectedTagCSVCols  = []string{"Number", "Numbers", "Word", "Words", "MapStringString", "MapStringInt", "MapIntInt"}
)

type Embedded struct {
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

var (
	expectedEmbedCols = []string{"Name", "ID", "Addr1", "Addr2", "City", "State", "Zip", "Phone", "Lat", "Long", "Notes", "Stuff"}
	EmbeddedTests     = []Embedded{
		Embedded{
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
		Embedded{
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
)

var expectedEmbedded = [][]string{
	[]string{"Name", "ID", "Addr1", "Addr2", "City", "State", "Zip", "Phone", "Lat", "Long", "Notes", "Stuff"},
	[]string{"United Center", "1", "1901 W. Madison St.", "", "Chicago", "IL", "60612",
		"(312) 455-4500", "41.8806", "-87.6742", "NBA:Bulls,NHL:Blackhawks", ""},
	[]string{"Wrigley Field", "1906", "1060 W. Addison St.", "Broadcast Booth", "Chicago", "IL", "60613",
		"(773) 404-2827", "41.9483", "-87.6556", "MLB:Cubs", "Harry Caray:Holy Cow,Jack Brickhouse:Hey Hey"},
}

type Basic struct {
	Name string   `json:"name" csv:"Nom"`
	List []string `json:"list" csv:"Liste"`
}

type BaseSliceTypes struct {
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
	Bool          *bool
	Bools         []*bool
	BoolPs        *[]bool
	Int           *int
	Ints          []*int
	IntPs         *[]int
	Uint          *uint
	Uints         []*uint
	UintPs        *[]int
	Float64       *float64
	Float64s      []*float64
	Float64Ps     *[]float64
	Complex128    *complex128
	Complex128s   []*complex128
	Complex128Ps  *[]complex128
	Chan          *chan int
	Chans         []*chan int
	ChanPs        *[]chan string
	Func          *func()
	Funcs         []*func()
	FuncPs        *[]func()
	String        *string
	Strings       []*string
	StringPs      *[]string
	strings       []*string
	BoolM         map[bool]*bool
	IntM          map[int]*int
	Float64M      map[float64]*float64
	Complex128M   map[complex128]*complex128
	ChanM         map[*chan string]*chan string
	FuncM         map[int]func()
	FuncMP        map[int]*func()
	StringM       map[string]*string
	KBoolM        map[*bool]*int
	KIntM         map[*int]bool
	KFloat64M     map[*float64]float64
	KComplex128M  map[*complex128]complex128
	KChanM        map[*chan int]string
	KFuncM        map[*func()]func()
	KStringM      map[*string]string
	PBoolM        *map[bool]*bool
	PIntM         *map[int]*int
	PFloat64M     *map[float64]*float64
	PComplex128M  *map[complex128]*complex128
	PChanM        *map[chan int]*chan string
	PStringM      *map[string]*string
	PKBoolM       *map[*bool]bool
	PKIntM        *map[*int]int
	PKFloat64M    *map[*float64]float64
	PKComplex128M *map[*complex128]complex128
	PKChanM       *map[*chan int]string
	PKFuncM       *map[*func()]func()
	PKStringM     *map[*string]string
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

func TestGetColNames(t *testing.T) {
	tc := New()
	_, err := tc.GetColNames([]string{"a", "b", "c"})
	if err == nil {
		t.Error("expected passing of a non struct to result in an error, none received")
	} else {
		if err.Error() != "struct2csv: a value of type struct is required: type was slice" {
			t.Errorf("expected error to be \"struct2csv: a value of type struct is required: type was slice\", got %q", err)
		}
	}

	tc.useTags = false
	hdr, err := tc.GetColNames(Tags{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto CSVTAG
	}
	if len(hdr) != len(expectedTagCols) {
		t.Errorf("Expected %d column ColNames, got %d", len(expectedTagCols), len(hdr))
		goto CSVTAG
	}
	for i, v := range hdr {
		if v != expectedTagCols[i] {
			t.Errorf("%d: expected %q got %q", i, expectedTagCols[i], v)
		}
	}
CSVTAG:
	// test using CSV tags
	tc.useTags = true
	hdr, err = tc.GetColNames(Tags{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto JSONTAG
	}
	if len(hdr) != len(expectedTagCSVCols) {
		t.Errorf("Expected %d column ColNames, got %d", len(expectedTagCSVCols), len(hdr))
		goto JSONTAG
	}
	for i, v := range hdr {
		if v != expectedTagCSVCols[i] {
			t.Errorf("%d: expected %q got %q", i, expectedTagCSVCols[i], v)
		}
	}
JSONTAG:
	// test using CSV tags
	tc.useTags = true
	tc.tag = "json"
	hdr, err = tc.GetColNames(Tags{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		goto EMBED
	}
	if len(hdr) != len(expectedTagJSONCols) {
		t.Errorf("Expected %d column ColNames, got %d", len(expectedTagJSONCols), len(hdr))
		goto EMBED
	}
	for i, v := range hdr {
		if v != expectedTagJSONCols[i] {
			t.Errorf("%d: expected %q got %q", i, expectedTagJSONCols[i], v)
		}
	}
EMBED:
	hdr, err = tc.GetColNames(Embedded{})
	if err != nil {
		t.Errorf("unexpected error getting header information from Tst{}: %q", err)
		return
	}
	if len(hdr) != len(expectedEmbedCols) {
		t.Errorf("Expected %d column ColNames, got %d", len(expectedEmbedCols), len(hdr))
		t.Errorf("%#v\n", hdr)
		return
	}
	for i, v := range hdr {
		if v != expectedEmbedCols[i] {
			t.Errorf("%d: expected %q got %q", i, expectedEmbedCols[i], v)
		}
	}
}

func TestMarshal(t *testing.T) {
	tsts := []BaseSliceTypes{
		BaseSliceTypes{
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
		BaseSliceTypes{
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
		[]string{"true", "true,false,true", "true,false",
			"42", "72,76,88,19,2", "1,2",
			"8", "8,9,10", "3,4",
			"16", "16,17,18", "5,6",
			"32", "32,33,34", "7,8",
			"64", "64,65,66", "9,10",
			"42", "72,76,88,19,2", "35,21",
			"8", "8,9,10", "11,12",
			"16", "16,17,18", "13,14",
			"32", "32,33,34", "15,16",
			"64", "64,65,66", "17,18",
			"3.242E+01", "3.242E+01,3.3E+01,3.44E+01", "3.55E+01,3.6754E+01",
			"6.442E+01", "6.442E+01,6.5E+01,6.67E+01", "6.902E+01,6.90132E+01",
			"(-64+12i)", "(-65+11i),(66+10i)", "(-61+21i),(76+8i)",
			"(-128+12i)", "(-128+11i),(129+10i)", "(-118+2i),(131+5i)",
			"don't panic", "", "pangalactic,gargleblaster"},
		[]string{"true", "true,true,false", "true,false",
			"420", "1,2,3,4", "11,12",
			"18", "18,19,110", "13,14",
			"116", "116,117,118", "15,16",
			"132", "132,133,134", "17,18",
			"640", "164,165,166", "19,110",
			"421", "121,122,123", "35,21",
			"118", "118,119,110", "111,112",
			"160", "116,117,118", "113,114",
			"320", "132,133,134", "115,116",
			"641", "164,165,166", "117,118",
			"1.3242E+02", "1.3242E+02,1.33E+02,1.344E+02", "1.355E+02,1.36754E+02",
			"1.6442E+02", "1.6442E+02,1.65E+02,1.667E+02", "1.6902E+02,1.690132E+02",
			"(-164+12i)", "(-165+11i),(166+10i)", "(-161+21i),(176+8i)",
			"(-124+12i)", "(-126+11i),(229+10i)", "(-116+2i),(231+5i)",
			"Towel", "", "Zaphod,Beeblebrox"},
	}
	tc := New()
	data, err := tc.Marshal(Tags{})
	if err != nil {
		if err.Error() != "struct2csv: a type of slice is required: type was struct" {
			t.Errorf("Expected \"struct2csv: a type of slice is required: type was struct\", got %q", err)
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
		if err.Error() != "struct2csv: the slice of structs was nil" {
			t.Errorf("Expected \"struct2csv: the slice of structs was nil\", got %q", err)
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
		if err.Error() != "struct2csv: the slice of structs was empty" {
			t.Errorf("Expected \"struct2csv: the slice of structs was empty\", got %q", err)
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
		if err.Error() != "struct2csv: a slice of type struct is required: slice type was string" {
			t.Errorf("Expected \"struct2csv: a slice of type struct is required: slice type was string\", %q", err)
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
	return
	if len(data) != len(expected) {
		t.Errorf("Expected %d rows, got %d", len(expected), len(data))
		goto EMBED
	}
	for i, row := range data {
		for j, col := range row {
			if col != expected[i][j] {
				t.Errorf("%d:%d: expected %q, got %q", i, j, expected[i][j], col)
			}
		}
	}
EMBED:
}

func TestMarshalStructs(t *testing.T) {
	tc := New()
	rows, err := tc.Marshal(EmbeddedTests)
	if err != nil {
		t.Errorf("did not expect an error: got %q", err)
		return
	}
	if len(rows) != len(expectedEmbedded) {
		t.Errorf("Expected %d rows of data, got %d", len(expectedEmbedded), len(rows))
		return
	}
	for i, row := range rows {
		if len(row) != len(expectedEmbedded[i]) {
			t.Errorf("%d: expected row to have %d columns, got %d", i, len(expectedEmbedded[i]), len(row))
			continue
		}
		for j, col := range row {
			if col != expectedEmbedded[i][j] {
				t.Errorf("%d:%d: expected %v got %v", i, j, expectedEmbedded[i][j], col)
			}
		}
	}
	// test GetRow
	for i, tst := range EmbeddedTests {
		row, err := tc.GetRow(tst)
		if err != nil {
			t.Errorf("Unexpected error")
		}
		for j, col := range row {
			if col != expectedEmbedded[i+1][j] {
				t.Errorf("%d:%d: expected %v, got %v", i, j, expectedEmbedded[i+1][j], col)
			}
		}
	}
}

type MapPtr struct {
	MapBasicP       map[string]*Basic
	MapBasicSliceP  map[string]*[]Basic
	MapPBasicSlice  map[string][]*Basic
	PMapPBasicSlice *map[string][]*Basic
}

func TestPtrStructs(t *testing.T) {
	tsts := []MapPtr{
		MapPtr{},
		MapPtr{
			MapBasicP:       map[string]*Basic{"MapBasicP": &Basic{Name: "Jason Bourne", List: []string{"keystone"}}},
			MapBasicSliceP:  map[string]*[]Basic{"MapBasicSliceP": new([]Basic)},
			MapPBasicSlice:  map[string][]*Basic{"MapPBasicSlice": []*Basic{&Basic{Name: "Foo", List: []string{"bar", "baz"}}}},
			PMapPBasicSlice: new(map[string][]*Basic),
		},
	}
	expected := [][]string{
		[]string{"MapBasicP", "MapBasicSliceP", "MapPBasicSlice", "PMapPBasicSlice"},
		[]string{"", "", "", "", ""},
		[]string{
			"MapBasicP:(Jason Bourne,(keystone))",
			"MapBasicSliceP:()",
			"MapPBasicSlice:((Foo,(bar,baz)))",
			"parks:((Wyoming,(Yellowstone,Grand Tetons)))",
		},
	}
	bsc := &Basic{Name: "Wyoming", List: []string{"Yellowstone", "Grand Tetons"}}
	m1 := map[string][]*Basic{"parks": []*Basic{bsc}}
	tsts[1].PMapPBasicSlice = &m1
	enc := New()
	rows, err := enc.Marshal(tsts)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	}
	for i, row := range rows {
		for j, col := range row {
			if col != expected[i][j] {
				t.Errorf("%d:%d: expected %q got %q", i, j, expected[i][j], col)
			}
		}
	}
}

type ComplexMap struct {
	MapMap          map[string]map[string]string
	MapSlice        map[string][]string
	Map2DSlice      map[string][][]string
	MapBasic        map[string]Basic
	MapBasicSlice   map[string][]Basic
	MapBasic2DSlice map[string][][]Basic
}

var complexTests = []ComplexMap{
	ComplexMap{
		MapMap: map[string]map[string]string{
			"Region 1": map[string]string{
				"colo1": "rack1",
				"colo2": "rack2",
			},
			"Region 2": map[string]string{
				"colo11": "rack11",
				"colo12": "rack12",
			},
		},
		MapSlice: map[string][]string{
			"Canada": []string{"Alberta", "British Columbia", "Quebec"},
			"USA":    []string{"California", "Florida", "New York"},
		},
		Map2DSlice: map[string][][]string{},
		MapBasic: map[string]Basic{
			"Gibson":  Basic{Name: "William Gibson", List: []string{"Neuromancer", "Count Zero", "Mona Lisa Overdrive"}},
			"Herbert": Basic{Name: "Frank Herbert", List: []string{"Destination Void", "Jesus Incident", "Lazurus Effect"}},
		},
		MapBasicSlice: map[string][]Basic{
			"SciFi": []Basic{
				Basic{Name: "William Gibson", List: []string{"Neuromancer", "Count Zero", "Mona Lisa Overdrive"}},
				Basic{Name: "Frank Herbert", List: []string{"Destination Void", "Jesus Incident", "Lazurus Effect"}},
			},
		},
		MapBasic2DSlice: map[string][][]Basic{
			"Sci-Fi": [][]Basic{
				[]Basic{
					Basic{Name: "William Gibson", List: []string{"Neuromancer", "Count Zero", "Mona Lisa Overdrive"}},
					Basic{Name: "Frank Herbert", List: []string{"Destination Void", "Jesus Incident", "Lazurus Effect"}},
				},
				[]Basic{
					Basic{Name: "Douglas Adams", List: []string{"Restaurant at the End of the Universe"}},
				},
			},
		},
	},
}

var complexExpected = [][]string{
	[]string{"MapMap", "MapSlice", "Map2DSlice", "MapBasic", "MapBasicSlice", "MapBasic2DSlice"},
	[]string{
		"Region 1:(colo1:rack1,colo2:rack2),Region 2:(colo11:rack11,colo12:rack12)",
		"Canada:(Alberta,British Columbia,Quebec),USA:(California,Florida,New York)",
		"Canada:((Calgary,Edmonton,Fort McMurray),(Winnipeg)),USA:((San Diego,Los Angeles))",
		"Gibson:(William Gibson,(Neuromancer,Count Zero,Mona Lisa Overdrive)),Herbert:(Frank Herbert,(Destination Void,Jesus Incident,Lazurus Effect))",
		"SciFi:((William Gibson,(Neuromancer,Count Zero,Mona Lisa Overdrive)),(Frank Herbert,(Destination Void,Jesus Incident,Lazurus Effect)))",
		"Sci-Fi:(((William Gibson,(Neuromancer,Count Zero,Mona Lisa Overdrive)),(Frank Herbert,(Destination Void,Jesus Incident,Lazurus Effect))),((Douglas Adams,(Restaurant at the End of the Universe))))",
	},
}

func TestComplicated(t *testing.T) {
	complexTests[0].Map2DSlice["Canada"] = [][]string{[]string{"Calgary", "Edmonton", "Fort McMurray"}, []string{"Winnipeg"}}
	complexTests[0].Map2DSlice["USA"] = [][]string{[]string{"San Diego", "Los Angeles"}}
	tc := New()
	rows, err := tc.Marshal(complexTests)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
		return
	}
	if len(rows) != len(complexExpected) {
		t.Errorf("expected %d rows got %d", len(complexExpected), len(rows))
		return
	}
	if len(rows[0]) != len(complexExpected[0]) {
		t.Errorf("expected a row to have %d columns, got %d", len(complexExpected[0]), len(rows[0]))
		return
	}
	for i, row := range rows {
		for j, col := range row {
			if col != complexExpected[i][j] {
				t.Errorf("%d:%d: expected %v got %v", i, j, complexExpected[i][j], col)
			}
		}
	}
}

func TestPtrs(t *testing.T) {
	tst := []PtrTypes{
		PtrTypes{},
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
		[]string{"Bool", "Bools", "BoolPs",
			"Int", "Ints", "IntPs",
			"Uint", "Uints", "UintPs",
			"Float64", "Float64s", "Float64Ps",
			"Complex128", "Complex128s", "Complex128Ps",
			"String", "Strings", "StringPs",
			"BoolM", "IntM", "Float64M", "Complex128M", "StringM",
			"KBoolM", "KIntM", "KFloat64M", "KComplex128M", "KStringM",
			"PBoolM", "PIntM", "PFloat64M", "PComplex128M", "PStringM",
			"PKBoolM", "PKIntM", "PKFloat64M", "PKComplex128M", "PKStringM",
		},
		[]string{
			"", "", "",
			"", "", "",
			"", "", "",
			"", "", "",
			"", "", "",
			"", "", "",
			"", "", "", "", "",
			"", "", "", "", "",
			"", "", "", "", "",
			"", "", "", "", "",
		},
		[]string{
			"true", "true,true", "false,false",
			"0", "0,0", "0",
			"0", "0,0", "0",
			"0+E00", "0+E00,0+E00", "",
			"0+0i", "0+0i,0+0i", "",
			"", ",", "",
			"", "", "", "", "",
			"", "", "", "", "",
			"", "", "", "", "",
			"", "", "", "", "",
		},
	}
	pbool(tst[1].Bool, true)
	pbool(tst[1].Bools[0], true)
	pbool(tst[1].Bools[1], true)
	pchan(tst[1].Chan)
	pchan(tst[1].Chans[0])
	pchan(tst[1].Chans[1])
	pfunc(tst[1].Func)
	pfunc(tst[1].Funcs[0])
	tmp := make([]bool, 2)
	tst[0].BoolPs = &tmp
	m := make(map[bool]*bool)
	tst[0].PBoolM = &m
	tc := New()
	data, err := tc.Marshal(tst)
	if err != nil {
		t.Errorf("Unexpected error %q", err)
	}
	//t.Errorf("\n%#v\n%#v\n", expected, data)
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

func TestPrivateMarshal(t *testing.T) {
	tsts := []struct {
		val     interface{}
		expKey  reflect.Kind
		expVal  reflect.Kind
		expCols []string
		isMap   bool
	}{
		{12, reflect.Int, reflect.Invalid, []string{"12"}, false},
		{12.4, reflect.Float64, reflect.Invalid, []string{"1.24E+01"}, false},
		{"Ni", reflect.String, reflect.Invalid, []string{"Ni"}, false},
		{(8 + 5i), reflect.Complex128, reflect.Invalid, []string{"(8+5i)"}, false},
		{[]chan int{}, reflect.Chan, reflect.Invalid, []string{}, false},
		{[]func(){}, reflect.Func, reflect.Invalid, []string{}, false},
		{[]string{"a", "b", "c"}, reflect.String, reflect.Invalid, []string{"a,b,c"}, false},
		{[]*int{}, reflect.Int, reflect.Invalid, []string{""}, false},
		{new([]float32), reflect.Float32, reflect.Invalid, []string{""}, false},
		{new([]*float64), reflect.Float64, reflect.Invalid, []string{""}, false},

		{[][]string{[]string{"a", "b", "c"}, []string{"d", "e", "f"}},
			reflect.String, reflect.Invalid, []string{"(a,b,c),(d,e,f)"}, false},
		{[][]*int{[]*int{}, []*int{}}, reflect.Int, reflect.Invalid, []string{"(2,2),(2)"}, false},
		{map[string]int{"a": 1, "b": 2}, reflect.String, reflect.Int, []string{"a:1,b:2"}, true},
		{map[int]string{1: "a", 2: "b", 3: "c"}, reflect.Int, reflect.String, []string{"1:a,2:b,3:c"}, true},
		{new(map[int]bool), reflect.Int, reflect.Bool, []string{""}, true},
		{map[*func()]int{}, reflect.Func, reflect.Int, []string{}, true},
		{map[int]func(){}, reflect.Int, reflect.Func, []string{}, true},
		{map[chan int]string{}, reflect.Chan, reflect.String, []string{}, true},
		{map[string]chan int{}, reflect.String, reflect.Chan, []string{}, true},
		{map[string]map[int]string{}, reflect.String, reflect.Int, []string{""}, true},

		{map[int]map[int]string{1: map[int]string{11: "a"}, 2: map[int]string{21: "A", 22: "B"}},
			reflect.Int, reflect.Map, []string{"1:(11:a),2:(21:A,22:B)"}, true},
		{map[int]map[chan string]int{}, reflect.Int, reflect.Chan, []string{}, true},
	}
	var i = new(int) //*int
	pint(i, 2)
	tsts[11].val = [][]*int{[]*int{i, i}, []*int{i}}
	enc := New()
	for i, tst := range tsts {
		cols, _ := enc.marshal(reflect.ValueOf(tst.val), false)
		if len(cols) != len(tst.expCols) {
			t.Errorf("%d: expected marshal to result in %d rows, got %d", i, len(tst.expCols), len(cols))
			continue
		}
		for j, col := range cols {
			if tst.isMap {
				cparts := splitCleanSort(col)
				eparts := splitCleanSort(tst.expCols[j])
				eq := true
				for i := 0; i < len(cparts); i++ {
					if cparts[i] != eparts[i] {
						eq = false
						goto EQ
					}
				}
			EQ:
				if !eq {
					t.Errorf("%d:%d: expected elements to have %q, got %q", i, j, tst.expCols[j], col)
				}
				continue
			}
			if tst.expCols[j] != col {
				t.Errorf("%d:%d: expected %q, got %q", i, j, tst.expCols[j], col)
			}
		}
	}
}

func TestIsSupportedKind(t *testing.T) {
	tsts := []struct {
		kind reflect.Kind
		ok   bool
	}{
		{reflect.Invalid, false},
		{reflect.Bool, true},
		{reflect.Int, true},
		{reflect.Int8, true},
		{reflect.Int16, true},
		{reflect.Int32, true},
		{reflect.Int64, true},
		{reflect.Uint, true},
		{reflect.Uint8, true},
		{reflect.Uint16, true},
		{reflect.Uint32, true},
		{reflect.Uint64, true},
		{reflect.Uintptr, false},
		{reflect.Float32, true},
		{reflect.Float64, true},
		{reflect.Complex64, true},
		{reflect.Complex128, true},
		{reflect.Array, true},
		{reflect.Chan, false},
		{reflect.Func, false},
		{reflect.Interface, false},
		{reflect.Map, true},
		{reflect.Ptr, true},
		{reflect.Slice, true},
		{reflect.String, true},
		{reflect.Struct, true},
		{reflect.UnsafePointer, false},
	}
	for i, tst := range tsts {
		ok := isSupportedKind(tst.kind)
		if ok != tst.ok {
			t.Errorf("%d: expected %s to be supported == %t; got %t", i, tst.kind, tst.ok, ok)
		}
	}
}

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
	*p = func() { fmt.Println("Hello") }
}

func splitCleanSort(s string) []string {
	s = strings.Replace(s, "(", ":", -1)
	s = strings.Replace(s, ")", ":", -1)
	s = strings.Replace(s, ",", ":", -1)
	s = strings.Replace(s, " ", ":", -1)
	parts := strings.Split(s, ":")
	sort.Strings(parts)
	return parts
}
