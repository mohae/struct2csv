// Package struct2csv creates CSV rows out of structs.
package struct2csv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Transcoder handles transcoding from a struct to CSV
type Transcoder struct {
	// Whether or not tags should be use for header (column) names; by default this is csv,
	useTags bool
	tag     string // The tag to use when tags are being used for headers; defaults to csv.
}

// NewTranscoder returns an initialized transcoder
func NewTranscoder() *Transcoder {
	return &Transcoder{useTags: true, tag: "csv"}
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

// SetUseTags sets whether or not tags should be used for header (column)
// names.
func (t *Transcoder) SetUseTags(b bool) {
	t.useTags = b
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
	st := reflect.TypeOf(v)
	sv := reflect.ValueOf(v)
	var hdrs []string
	for i := 0; i < st.NumField(); i++ {
		// skip unexported
		ft := st.Field(i)
		if ft.PkgPath != "" {
			continue
		}
		fv := sv.Field(i)
		switch fv.Kind() {
		case reflect.Struct:
			tmp, err := t.getHeaders(fv.Interface())
			if err != nil {
				return nil, err
			}
			hdrs = append(hdrs, tmp...)
			continue
		case reflect.Chan, reflect.Func:
			continue
		case reflect.Array, reflect.Slice:
			// skip if it's a array/slice of chan or func.
			switch fv.Type().Elem().Kind() {
			case reflect.Func, reflect.Chan:
				continue
			}
		}
		var name string
		if t.useTags {
			name = ft.Tag.Get(t.tag)
		}
		// If there isn't a matching field tag, use the Field Name
		if name == "" {
			name = ft.Name
		}
		hdrs = append(hdrs, name)
	}
	return hdrs, nil
}

// Marshal takes a slice of structs and returns a [][]byte representing CSV
// data. Each field in the struct results in a column.  Fields that are slices
// are stored in a single column as a comma separated list of values.  Fields
// that are maps are stored in a single column as a comma separted list of
// key:value pairs.
//
// If the passed data isn't a slice of structs or an error occurs during
// processing, an error will be returned.
// TODO:
//    handle pointers
func (t *Transcoder) Marshal(v interface{}) ([][]string, error) {
	// must be a slice
	if reflect.TypeOf(v).Kind() != reflect.Slice {
		return nil, fmt.Errorf("slice required: value was of type %s", reflect.TypeOf(v).Kind())
	}
	// must be a slice of struct
	vv := reflect.ValueOf(v)
	if vv.IsNil() {
		return nil, fmt.Errorf("slice cannot be nil")
	}
	if vv.Len() == 0 {
		return nil, fmt.Errorf("slice must have a length of at least 1: length was 0")
	}
	var rows [][]string
	s := vv.Index(0)
	switch s.Kind() {
	case reflect.Struct:
		hdrs, err := t.getHeaders(s.Interface())
		if err != nil {
			return nil, err
		}
		rows = append(rows, hdrs)
	default:
		return nil, fmt.Errorf("slice must be of type struct; type was %s", s.Kind().String())
	}
	for i := 0; i < vv.Len(); i++ {
		s := vv.Index(i)
		switch s.Kind() {
		case reflect.Struct:
			row, err := t.marshalStruct(s.Interface())
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		default:
			return nil, fmt.Errorf("slice must be of type struct; type was %s", s.Kind().String())
		}
	}
	return rows, nil
}

func (t *Transcoder) marshalStruct(v interface{}) ([]string, error) {
	var row []string
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	for i := 0; i < val.NumField(); i++ {
		if typ.Field(i).PkgPath != "" {
			continue
		}
		f := val.Field(i)
		switch f.Kind() {
		case reflect.Bool:
			row = append(row, strconv.FormatBool(f.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			row = append(row, strconv.Itoa(int(f.Int())))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			row = append(row, strconv.FormatUint(uint64(f.Uint()), 10))
		case reflect.Float32:
			row = append(row, strconv.FormatFloat(f.Float(), 'E', -1, 32))
		case reflect.Float64:
			row = append(row, strconv.FormatFloat(f.Float(), 'E', -1, 64))
		case reflect.Complex64, reflect.Complex128:
			row = append(row, fmt.Sprintf("%g", f.Complex()))
		case reflect.Chan, reflect.Func:
			continue
		case reflect.Array, reflect.Slice:
			// skip these
			switch f.Type().Elem().Kind() {
			case reflect.Chan, reflect.Func:
				continue
			}
			trow, err := t.marshalSlice(f)
			if err != nil {
				return nil, err
			}
			row = append(row, trow)
		case reflect.Map:
			col, err := t.marshalMap(f.Interface())
			if err != nil {
				return nil, err
			}
			row = append(row, col)
		case reflect.String:
			row = append(row, f.String())
		case reflect.Struct:
			trow, err := t.marshalStruct(f.Interface())
			if err != nil {
				return nil, err
			}
			row = append(row, trow...)
		default:
			return nil, fmt.Errorf("%#v's type not supported: %s", f, f.Kind().String())
		}
	}
	return row, nil
}

// marshal map handles marshalling of maps
func (t *Transcoder) marshalMap(v interface{}) (string, error) {
	var row string
	m := reflect.ValueOf(v)
	if m.Kind() != reflect.Map {
		return "", fmt.Errorf("map expected: type was %s", m.Kind().String())
	}
	keys := m.MapKeys()
	if len(keys) == 0 {
		return "", nil
	}
	for i, key := range keys {
		val := m.MapIndex(key)
		k, err := t.stringify(key)
		if err != nil {
			return "", err
		}
		switch val.Kind() {
		case reflect.Map:
			tmp, err := t.marshalMap(val.Interface())
			if err != nil {
				return "", err
			}
			if i == 0 {
				row = fmt.Sprintf("%s:(%s)", k, tmp)
			} else {
				row = fmt.Sprintf("%s, %s:(%s)", row, k, tmp)
			}
			continue
		case reflect.Array, reflect.Slice:
			tmp, err := t.marshalSlice(val)
			if err != nil {
				return "", err
			}
			if i == 0 {
				row = fmt.Sprintf("%s:(%s)", k, tmp)
			} else {
				row = fmt.Sprintf("%s, %s:(%s)", row, k, tmp)
			}
			continue
		case reflect.Struct:
			tmp, err := t.marshalStruct(val.Interface())
			if err != nil {
				return "", err
			}
			var trow string
			for j, v := range tmp {
				// if this is a list, put it in brackets
				v = t.formatList(v)
				if j == 0 {
					trow = v
				} else {
					trow = fmt.Sprintf("%s, %s", trow, v)
				}
			}
			if i == 0 {
				row = fmt.Sprintf("%s:(%s)", k, trow)
			} else {
				row = fmt.Sprintf("%s, %s:(%s)", row, k, trow)
			}
			continue
		case reflect.Chan, reflect.Func:
			continue
		}
		v, err := t.stringify(val)
		if err != nil {
			return "", err
		}
		if i == 0 {
			row = fmt.Sprintf("%s:%s", k, v)
			continue
		}
		row = fmt.Sprintf("%s, %s:%s", row, k, v)
	}
	return row, nil
}

// marshalSlice handles marshaling of slices
func (t *Transcoder) marshalSlice(v reflect.Value) (string, error) {
	var sl, str string
	var err error
	for j := 0; j < v.Len(); j++ {
		switch v.Index(j).Kind() {
		case reflect.Struct:
			tmp, err := t.marshalStruct(v.Index(j).Interface())
			if err != nil {
				return "", err
			}
			for i, v := range tmp {
				v = t.formatList(v)
				if i == 0 {
					str = v
					continue
				}
				str = fmt.Sprintf("%s, %s", str, v)
			}
		default:
			str, err = t.stringify(v.Index(j))
			if err != nil {
				return "", err
			}
		}
		if j == 0 {
			sl = str
			continue
		}
		sl = fmt.Sprintf("%s, %s", sl, str)
	}
	return sl, nil
}

// stringify takes a interface and returns the value it contains as a string.
// This is not ment for composite types.  If the type is not supported, an
// error will be returned.
func (t *Transcoder) stringify(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(v.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(uint64(v.Uint()), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'E', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'E', -1, 64), nil
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%g", v.Complex()), nil
	case reflect.String:
		return v.String(), nil
	default:
		return "", fmt.Errorf("stringify: type not supported: %s", v.Kind().String())
	}
}

// formatList takes a string and adds brackets to the beginning and end of it
// if the string contains ", ".  Otherwise it is returned unmodified.
func (t *Transcoder) formatList(s string) string {
	if strings.Index(s, ", ") > 0 {
		return fmt.Sprintf("[%s]", s)
	}
	return s
}

// GetHeaders instantiates a Transcoder and gets the headers of the received
// struct.  If you need more control over tag processing, use NewTranscoder(),
// set accordingly, and call Transcoder's GetHeaders().
func GetHeaders(v interface{}) ([]string, error) {
	tc := NewTranscoder()
	return tc.GetHeaders(v)
}
