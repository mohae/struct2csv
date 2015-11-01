// Package struct2csv creates CSV rows out of structs.
package struct2csv

import (
	"fmt"
	"reflect"
)

// Transcoder handles transcoding from a struct to CSV
type Transcoder struct {
	// Whether or not tags should be use for header (column) names; by default this is csv,
	useTags bool
	tag     string // The tag to use when tags are being used for headers; defaults to csv.
	csv     [][]byte
}

// NewTranscoder returns an initialized transcoder
func NewTranscoder() *Transcoder {
	return &Transcoder{useTags: true, tag: "csv"}
}

// SetUseTags sets whether or not tags should be used for header (column)
// names.
func (t *Transcoder) SetUseTags(b bool) {
	t.useTags = b
}

// SetTag sets the tag that the Transcoder should use for header (column)
// names.  By default, this is set to 'csv'.  If the received value is an
// empty string, nothing will be done
func (t *Transcoder) SetTag(s string) {
	if s == "" {
		return
	}
	t.tag = s
}

// GetHeaders get's the column headers from the received struct.  If anything
// other than a struct is passed, an error will be returned.  If the struct
// has field tags for csv, those values will be used as the column headers,
// otherwise the field names will be used.
//
// If field tags other than the ones for csv are to be used, TODO figure out
// the struct and how to implement this comment.
func (t *Transcoder) GetHeaders(v interface{}) ([]string, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, fmt.Errorf("struct required: value was of type %s", reflect.TypeOf(v).Kind())
	}
	return t.getHeaders(v)
}

func (t *Transcoder) getHeaders(v interface{}) ([]string, error) {
	s := reflect.TypeOf(v)
	var hdrs []string
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		// skip unexported
		if f.PkgPath != "" {
			continue
		}
		var name string
		if t.useTags {
			name = f.Tag.Get(t.tag)
		}
		// If there isn't a matching field tag, use the Field Name
		if name == "" {
			name = f.Name
		}
		hdrs = append(hdrs, name)
	}
	return hdrs, nil
}

// GetHeaders instantiates a Transcoder and gets the headers of the received
// struct.  If you need more control over tag processing, use NewTranscoder(),
// set accordingly, and call Transcoder's GetHeaders().
func GetHeaders(v interface{}) ([]string, error) {
	tc := NewTranscoder()
	return tc.GetHeaders(v)
}
