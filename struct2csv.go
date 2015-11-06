// Package struct2csv creates CSV rows out of structs.
package struct2csv

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// An UnsupportedTypeError is returned when attempting to encode an
// an unsupported value type.
type UnsupportedTypeError struct {
	Kind reflect.Kind
}

func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf("struct2csv: unsupported type: %s", e.Kind)
}

// A StructRequiredError is returned when a non-struct type is received.
type StructRequiredError struct {
	Kind reflect.Kind
}

func (e StructRequiredError) Error() string {
	return fmt.Sprintf("struct2csv: a value of type struct is required: type was %s", e.Kind)
}

// A UnsupportedStringifyTypeError is returned when stringify is asked
// to create a string out of a type it does not support.
type UnsupportedStringifyTypeError struct {
	Kind reflect.Kind
}

func (e UnsupportedStringifyTypeError) Error() string {
	return fmt.Sprintf("struct2csv: stringify encountered an unsupported type: %s", e.Kind)
}

// A StructSliceError is returned when an interface that isn't a slice of
// type struct is received.
type StructSliceError struct {
	Kind      reflect.Kind
	SliceKind reflect.Kind
}

func (e StructSliceError) Error() string {
	if e.Kind != reflect.Slice {
		return fmt.Sprintf("struct2csv: a type of slice is required: type was %s", e.Kind)
	}
	return fmt.Sprintf("struct2csv: a slice of type struct is required: slice type was %s", e.SliceKind)
}

// An UnsupportedMapKeyTypeError is returned when the type of the key is an
// unsupported value type.
type UnsupportedMapKeyTypeError struct {
	Kind reflect.Kind
}

func (e UnsupportedMapKeyTypeError) Error() string {
	return fmt.Sprintf("struct2csv: map key is an unsupported type: %s", e.Kind)
}

// An UnsupportedMapValueTypeError is returned when the type of the key is an
// unsupported value type.
type UnsupportedMapValueTypeError struct {
	Kind reflect.Kind
}

func (e UnsupportedMapValueTypeError) Error() string {
	return fmt.Sprintf("struct2csv: map value is an unsupported type: %s", e.Kind)
}

var (
	// ErrNilSlice occurs when the slice of structs to encode is nil.
	ErrNilSlice   = errors.New("struct2csv: the slice of structs was nil")
	// ErrEmptySlice occurs when the slice of structs to encode is empty.
	ErrEmptySlice = errors.New("struct2csv: the slice of structs was empty")
)

// Encoder handles encoding of a CSV from a struct.
type Encoder struct {
	// Whether or not tags should be use for header (column) names; by default this is csv,
	useTags bool
	tag     string // The tag to use when tags are being used for headers; defaults to csv.
}

// New returns an initialized Encoder.
func New() *Encoder {
	return &Encoder{useTags: true, tag: "csv"}
}

// SetTag sets the tag that the Encoder should use for header (column)
// names.  By default, this is set to 'csv'.  If the received value is an
// empty string, nothing will be done
func (e *Encoder) SetTag(s string) {
	if s == "" {
		return
	}
	e.tag = s
}

// SetUseTags sets whether or not tags should be used for header (column)
// names.
func (e *Encoder) SetUseTags(b bool) {
	e.useTags = b
}

// GetHeaders get's the column headers from the received struct.  If anything
// other than a struct is passed, an error will be returned.  If the struct
// has field tags for csv, those values will be used as the column headers,
// otherwise the field names will be used.
//
// If field tags other than the ones for csv are to be used, TODO figure out
// the struct and how to implement this comment.
func (e *Encoder) GetHeaders(v interface{}) ([]string, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, StructRequiredError{reflect.TypeOf(v).Kind()}
	}
	return e.getHeaders(v)
}

func (e *Encoder) getHeaders(v interface{}) ([]string, error) {
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
		if !isSupported(fv.Kind()) {
			continue
		}
		switch fv.Kind() {
		case reflect.Struct:
			tmp, err := e.getHeaders(fv.Interface())
			if err != nil {
				return nil, err
			}
			hdrs = append(hdrs, tmp...)
			continue
		case reflect.Array, reflect.Slice:
			if skipSliceField(fv) {
				continue
			}
		case reflect.Map:
			//fmt.Println(fv.MapIndex())
			switch fv.Type().Key().Kind() {
			case reflect.Func, reflect.Chan:
				continue
			case reflect.Ptr:
				k := fv.Type().Key()
				switch k.Kind() {
				case reflect.Ptr:
					if !isSupported(k.Elem().Kind()) {
						continue
					}
				case reflect.Chan, reflect.Func, reflect.Interface, reflect.Uintptr, reflect.UnsafePointer:
					continue
				}
			}
			// get the value type
			var k reflect.Type
			switch fv.Type().Kind() {
			case reflect.Ptr:
				k = fv.Type().Elem().Key()
			default:
				k = fv.Type().Key()
			}
			kv := reflect.Zero(k)
			// if it's a pointer, get what it points to
			if kv.Kind() == reflect.Ptr {
				kv = kv.Elem()
			}
			/*
				// if it's invalid, we have no way of knowing the type of the map value
				// so this column gets marshaled.
				if kv.Kind() != reflect.Invalid {
					switch fv.Type().Kind() {
					case reflect.Ptr:
						val := fv.Elem().MapIndex(kv)
					default:
						fmt.Println(fv)
						val := fv.MapIndex(kv)
					}
				}
			*/
		case reflect.Ptr:
			switch fv.Type().Elem().Kind() {
			case reflect.Func, reflect.Chan, reflect.Uintptr, reflect.Interface, reflect.UnsafePointer:
				continue
			case reflect.Ptr:
				k := fv.Type().Elem().Key()
				switch k.Kind() {
				case reflect.Ptr:
					switch k.Elem().Kind() {
					case reflect.Chan, reflect.Func, reflect.Interface, reflect.Uintptr, reflect.UnsafePointer:
						continue
					}
				case reflect.Chan, reflect.Func, reflect.Interface, reflect.Uintptr, reflect.UnsafePointer:
					continue
				}
			case reflect.Slice:
				if skipSliceField(fv) {
					continue
				}
			case reflect.Map:
				k := fv.Type().Elem().Key()
				switch k.Kind() {
				case reflect.Ptr:
					if !isSupported(k.Elem().Kind()) {
						continue
					}
				case reflect.Chan, reflect.Func, reflect.Interface, reflect.Uintptr, reflect.UnsafePointer:
					continue
				}
			}
		}
		var name string
		if e.useTags {
			name = ft.Tag.Get(e.tag)
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
func (e *Encoder) Marshal(v interface{}) ([][]string, error) {
	// must be a slice
	if reflect.TypeOf(v).Kind() != reflect.Slice {
		return nil, StructSliceError{Kind: reflect.TypeOf(v).Kind()}
	}
	// must be a slice of struct
	vv := reflect.ValueOf(v)
	if vv.IsNil() {
		return nil, ErrNilSlice
	}
	if vv.Len() == 0 {
		return nil, ErrEmptySlice
	}
	var rows [][]string
	s := vv.Index(0)
	switch s.Kind() {
	case reflect.Struct:
		hdrs, err := e.getHeaders(s.Interface())
		if err != nil {
			return nil, err
		}
		rows = append(rows, hdrs)
	default:
		return nil, StructSliceError{Kind: reflect.Slice, SliceKind: s.Kind()}
	}
	for i := 0; i < vv.Len(); i++ {
		s := vv.Index(i)
		switch s.Kind() {
		case reflect.Struct:
			row, err := e.marshalStruct(s.Interface())
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
		default:
			return nil, StructSliceError{Kind: reflect.Slice, SliceKind: s.Kind()}
		}
	}
	return rows, nil
}

func (e *Encoder) marshalStruct(v interface{}) ([]string, error) {
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
			if skipSliceField(f) {
				continue
			}
			trow, err := e.marshalSlice(f)
			if err != nil {
				return nil, err
			}
			row = append(row, trow)
		case reflect.Map:
			col, err := e.marshalMap(f)
			if err != nil {
				// skip unsupported maps
				if _, ok := err.(UnsupportedMapKeyTypeError); ok {
					continue
				}
				if _, ok := err.(UnsupportedMapValueTypeError); ok {
					continue
				}
				return nil, err
			}
			row = append(row, col)
		case reflect.String:
			row = append(row, f.String())
		case reflect.Struct:
			trow, err := e.marshalStruct(f.Interface())
			if err != nil {
				return nil, err
			}
			row = append(row, trow...)
		case reflect.Ptr:
			switch f.Type().Elem().Kind() {
			case reflect.Func, reflect.Chan, reflect.Uintptr, reflect.Interface, reflect.UnsafePointer:
				continue
			case reflect.Ptr:
				switch f.Type().Elem().Elem().Kind() {
				case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.Interface, reflect.UnsafePointer:
					continue
				}
			case reflect.Slice:
				if skipSliceField(f) {
					continue
				}
				trow, err := e.marshalSlice(f.Elem())
				if err != nil {
					fmt.Println("marshal struct ptr:slice stringify", err)
					return nil, err
				}
				row = append(row, trow)
				continue
			case reflect.Map:
				if reflect.Indirect(f).Kind() == reflect.Invalid {
					// check the key and value to see if their supported types
					k := f.Type().Elem().Key()
					if k.Kind() == reflect.Ptr {
						k = f.Type().Elem().Key().Elem()
					}
					switch k.Kind() {
					case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.Interface, reflect.UnsafePointer:
						continue
					}
					// Note, for pointers to nil maps, maps whose value types are unsupported will
					// also have an ampty string appended to the row.  This is because trying to
					// obtain the value type of a field which is a pointer to a nil map results in
					// a panic.
					// TODO handle this.
					row = append(row, "")
					continue
				}
				col, err := e.marshalMap(f.Elem())
				if err != nil {
					// skip unsupported maps
					if _, ok := err.(UnsupportedMapKeyTypeError); ok {
						continue
					}
					if _, ok := err.(UnsupportedMapValueTypeError); ok {
						continue
					}
					return nil, err
				}
				row = append(row, col)
				continue
			}
			val := f.Elem()
			if !val.IsValid() {
				continue
			}
			tmp, err := e.stringify(val)
			if err != nil {
				fmt.Println("marshal struct dflt stringify", err)
				return nil, err
			}
			row = append(row, tmp)
		default:
			return nil, UnsupportedTypeError{f.Kind()}
		}
	}
	return row, nil
}

// marshal map handles marshalling of maps
func (e *Encoder) marshalMap(m reflect.Value) (string, error) {
	var row string
	//  TODO add support for checking unsupported key types and handle, e.g. func & chan
	// a map may be nil but we still need to know its key and value types.
	//mz := reflect.Zero(m.Type().Key())
	// check to see if the value is supported:
	switch m.Type().Key().Kind() {
	case reflect.Chan, reflect.Func:
		return "", UnsupportedMapValueTypeError{m.Type().Key().Kind()}
	case reflect.Ptr:
		k := m.Type().Key()
		switch k.Kind() {
		case reflect.Func, reflect.Chan:
			return "", UnsupportedMapValueTypeError{k.Kind()}
		case reflect.Ptr:
			switch k.Elem().Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Uintptr, reflect.UnsafePointer:
				return "", UnsupportedMapValueTypeError{k.Elem().Kind()}
			}
		}
	}
	keys := m.MapKeys()
	for i, key := range keys {
		val := m.MapIndex(key)
		k, err := e.stringify(key)
		if err != nil {
			fmt.Println("map key", err)
			return "", err
		}
		switch val.Kind() {
		case reflect.Map:
			// skip if it's a map of chan or func.
			switch val.Type().Elem().Kind() {
			case reflect.Func, reflect.Chan:
				continue
			case reflect.Ptr:
				switch val.Type().Elem().Elem().Kind() {
				case reflect.Func, reflect.Chan:
					continue
				}
			}
			tmp, err := e.marshalMap(val)
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
			tmp, err := e.marshalSlice(val)
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
			tmp, err := e.marshalStruct(val.Interface())
			if err != nil {
				return "", err
			}
			var trow string
			for j, v := range tmp {
				// if this is a list, put it in brackets
				v = e.formatList(v)
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
		case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.UnsafePointer, reflect.Interface:
			continue
		case reflect.Ptr:
			continue
		}
		v, err := e.stringify(val)
		if err != nil {
			fmt.Println("stringify val", err)
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
func (e *Encoder) marshalSlice(v reflect.Value) (string, error) {
	var sl, str string
	var err error
	// this handles *[]'s that are nil'
	if v.Kind() == reflect.Invalid {
		return "", nil
	}
	for j := 0; j < v.Len(); j++ {
		switch v.Index(j).Kind() {
		case reflect.Struct:
			tmp, err := e.marshalStruct(v.Index(j).Interface())
			if err != nil {
				return "", err
			}
			for i, v := range tmp {
				v = e.formatList(v)
				if i == 0 {
					str = v
					continue
				}
				str = fmt.Sprintf("%s, %s", str, v)
			}
		case reflect.Ptr:
			if !isSupported(v.Index(j).Elem().Kind()) {
				continue
			}
			if v.Index(j).Elem().Kind() == reflect.Ptr {
				continue
			}
			str, err = e.stringify(v.Index(j).Elem())
			if err != nil {
				fmt.Println("marshal slice ptr stringify", err)
				return "", err
			}
		default:
			str, err = e.stringify(v.Index(j))
			if err != nil {
				fmt.Println("marshal slice dflt stringify", err)
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
func (e *Encoder) stringify(v reflect.Value) (string, error) {
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
		return "", UnsupportedStringifyTypeError{v.Kind()}
	}
}

// formatList takes a string and adds brackets to the beginning and end of it
// if the string contains ", ".  Otherwise it is returned unmodified.
func (e *Encoder) formatList(s string) string {
	if strings.Index(s, ", ") > 0 {
		return fmt.Sprintf("[%s]", s)
	}
	return s
}

// GetHeaders instantiates a Encoder and gets the headers of the received
// struct.  If you need more control over tag processing, use New(),
// set accordingly, and call Encoder's GetHeaders().
func GetHeaders(v interface{}) ([]string, error) {
	tc := New()
	return tc.GetHeaders(v)
}

// func skipSliceField handles both slices and arrays
func skipSliceField(f reflect.Value) bool {
	// skip if it's not a supported type
	if !isSupported(f.Type().Elem().Kind()) {
		return true
	}
	// check to see if it's either a pointer or a slice and check it's type if it is
	if f.Type().Elem().Kind() == reflect.Ptr || f.Type().Elem().Kind() == reflect.Slice {
		if !isSupported(f.Type().Elem().Elem().Kind()) {
			return true
		}
	}
	return false

}

func isSupported(k reflect.Kind) bool {
	switch k {
	case reflect.Chan:
		return false
	case reflect.Func:
		return false
	case reflect.Interface:
		return false
	case reflect.Uintptr:
		return false
	case reflect.UnsafePointer:
		return false
	}
	return true
}
