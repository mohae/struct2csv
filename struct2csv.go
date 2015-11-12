// Package struct2csv creates CSV rows out of structs.
package struct2csv

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// An UnsupportedTypeError is returned when attempting to encode an
// an unsupported value type.
type UnsupportedTypeError struct {
	kind    reflect.Kind
	method  string
	message string
}

func (e UnsupportedTypeError) Error() string {
	if e.message == "" {
		return fmt.Sprintf("struct2csv:%s: unsupported type: %s", e.method, e.kind)
	}
	return fmt.Sprintf("struct2csv:%s: %s %s", e.method, e.message, e.kind)
}

// A StructRequiredError is returned when a non-struct type is received.
type StructRequiredError struct {
	Kind reflect.Kind
}

func (e StructRequiredError) Error() string {
	return fmt.Sprintf("struct2csv: a value of type struct is required: type was %s", e.Kind)
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

var (
	// ErrNilSlice occurs when the slice of structs to encode is nil.
	ErrNilSlice = errors.New("struct2csv: the slice of structs was nil")
	// ErrEmptySlice occurs when the slice of structs to encode is empty.
	ErrEmptySlice = errors.New("struct2csv: the slice of structs was empty")
)

// Below is implemented from
// https://golang.org/src/encoding/json/encode.go#L773 through L780
// This is the copyright of the original code:
// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// stringValues is a slice of reflect.Value holding *reflect.StringValue.
// It implements the methods to sort by string.
type stringValues []reflect.Value

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

// Encoder handles encoding of a CSV from a struct.
type Encoder struct {
	// Whether or not tags should be use for header (column) names; by default this is csv,
	useTags  bool
	base     int
	tag      string // The tag to use when tags are being used for headers; defaults to csv.
	sepBeg   string
	sepEnd   string
	colNames []string
}

// New returns an initialized Encoder.
func New() *Encoder {
	return &Encoder{
		useTags: true, base: 10, tag: "csv",
		sepBeg: "(", sepEnd: ")",
	}
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

// SetSeparators sets the begin and end separator values for lists.   Setting
// the separators to "", empty strings, results in no separators being added.
// By default, "(" and ")" are used as the begin and end separators,
func (e *Encoder) SetSeparators(beg, end string) {
	e.sepBeg = beg
	e.sepEnd = end
}

// SetBase sets the base for strings.FormatUint. By default, this is 10. Set
// the base if another base should be used for formatting uint values.
//
// Base 2 is the minimum value; anything less will be set to two.
func (e *Encoder) SetBase(i int) {
	if i < 2 {
		i = 2
	}
	e.base = i
}

// ColNames returns the encoder's saved column names as a copy.  The
// colNames field must be populated before using this.
func (e *Encoder) ColNames() []string {
	ret := make([]string, len(e.colNames))
	_ = copy(ret, e.colNames)
	return ret
}

// GetColNames get's the column names from the received struct.  If the
// interface is not a struct, an error will occur.
//
// Field tags are supported. By default, the column names will be the value
// of the `csv` tag, if any.  This can be changed with the SetTag(newTag)
// func; e.g. `json` to use JSON tags.  Use of field tags can be toggled with
// the the SetUseTag(bool) func.  If use of field tags is set to FALSE, the
// field's name will be used.
func (e *Encoder) GetColNames(v interface{}) ([]string, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, StructRequiredError{reflect.TypeOf(v).Kind()}
	}
	return e.getColNames(v)
}

// The private func where the work is done.  This also copies the headers
// to the Encoder.colNames field.  This is a copy of the slice contents.
func (e *Encoder) getColNames(v interface{}) ([]string, error) {
	st := reflect.TypeOf(v)
	sv := reflect.ValueOf(v)
	var cols []string
	for i := 0; i < st.NumField(); i++ {
		// skip unexported
		ft := st.Field(i)
		if ft.PkgPath != "" {
			continue
		}
		fv := sv.Field(i)
		if !supportedKind(fv.Kind()) {
			continue
		}
		switch fv.Kind() {
		case reflect.Struct:
			tmp, err := e.getColNames(fv.Interface())
			if err != nil {
				return nil, err
			}
			cols = append(cols, tmp...)
			continue
		case reflect.Array, reflect.Slice:
			_, ok := supportedSliceKind(fv.Type())
			if !ok {
				continue
			}
		case reflect.Ptr:
			_, _, ok := supportedPtrKind(fv.Type().Elem())
			if !ok {
				continue
			}
		case reflect.Map:
			_, _, ok := supportedMapKind(fv.Type())
			if !ok {
				continue
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
		cols = append(cols, name)
	}
	e.colNames = make([]string, len(cols))
	_ = copy(e.colNames, cols)
	return cols, nil
}

// GetRow get's the data from the passed struct. This only operates on
// single structs.  If you wish to transmogrify everything at once, use
// Encoder.Marshal([]T).
func (e *Encoder) GetRow(v interface{}) ([]string, error) {
	return e.marshalStruct(v)
}

// Marshal takes a slice of structs and returns a [][]byte representing CSV
// data. Each field in the struct results in a column.  Fields that are slices
// are stored in a single column as a comma separated list of values.  Fields
// that are maps are stored in a single column as a comma separted list of
// key:value pairs.
//
// If the passed data isn't a slice of structs or an error occurs during
// processing, an error will be returned.
func (e *Encoder) Marshal(v interface{}) ([][]string, error) {
	// must be a slice
	if reflect.TypeOf(v).Kind() != reflect.Slice {
		return nil, StructSliceError{Kind: reflect.TypeOf(v).Kind()}
	}
	val := reflect.ValueOf(v)
	// must be a slice of struct
	if val.IsNil() {
		return nil, ErrNilSlice
	}
	if val.Len() == 0 {
		return nil, ErrEmptySlice
	}
	var rows [][]string
	s := val.Index(0)
	switch s.Kind() {
	case reflect.Struct:
		cols, err := e.getColNames(s.Interface())
		if err != nil {
			return nil, err
		}
		rows = append(rows, cols)
	default:
		return nil, StructSliceError{Kind: reflect.Slice, SliceKind: s.Kind()}
	}
	for i := 0; i < val.Len(); i++ {
		s := val.Index(i)
		row, err := e.marshalStruct(s.Interface())
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)

	}
	return rows, nil
}

// marshal returns the marshaled struct as a slice of column values, as
// strings. Any error results a nil slice.
func (e *Encoder) marshal(val reflect.Value) (cols []string, err error) {
	var s string
	switch val.Kind() {
	case reflect.Ptr:
		cols, err = e.marshal(val.Elem())
		return cols, err
	case reflect.Struct:
		return e.marshalStruct(val.Interface())
	case reflect.Map:
		s, err = e.marshalMap(val)
		if err != nil {
			return nil, err
		}
	case reflect.Array, reflect.Slice:
		s, err = e.marshalSlice(val)
		if err != nil {
			return nil, err
		}
		s = fmt.Sprintf("%s%s%s", e.sepBeg, s, e.sepEnd)
	default:
		var ok bool
		s, ok = e.stringify(val)
		if !ok {
			return nil, UnsupportedTypeError{kind: val.Kind(), method: "marshal"}
		}
	}

	return append(cols, s), nil
}

func (e *Encoder) marshalStruct(s interface{}) ([]string, error) {
	var rows []string
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	for i := 0; i < val.NumField(); i++ {
		if typ.Field(i).PkgPath != "" {
			continue
		}
		str, err := e.marshal(val.Field(i))
		if err != nil {
			// field kind is not supported
			if _, ok := err.(UnsupportedTypeError); ok {
				continue
			}
			return nil, err
		}
		rows = append(rows, str...)
	}
	return rows, nil
}

// marshal map handles marshalling of maps.  If the map's key type is
// supported is determined by the caller.
func (e *Encoder) marshalMap(m reflect.Value) (string, error) {
	if k, v, ok := supportedMapKind(m.Type()); !ok {
		if supportedKind(k) {
			return "", UnsupportedTypeError{kind: v, method: "marshalMap"}
		}
		return "", UnsupportedTypeError{kind: k, method: "marshalMap"}
	}
	var row string
	var sv stringValues = m.MapKeys()
	// sort the map keys first
	sort.Sort(sv)
	for i, key := range sv {
		val := m.MapIndex(key)
		kk, ok := e.stringify(key)
		if !ok {
			return "", UnsupportedTypeError{kind: key.Kind(), method: "marshalMap"}
		}
		vv, ok := e.stringify(val)
		if !ok {
			return "", UnsupportedTypeError{kind: key.Kind(), method: "marshalMap"}
		}
		if i == 0 {
			row = fmt.Sprintf("%s:%s", kk, vv)
		} else {
			row = fmt.Sprintf("%s,%s:%s", row, kk, vv)
		}
	}
	return row, nil
}

// marshalSlice handles marshaling of slices. This should not receive a
// pointer. Is is assumed that any pointers to the slice have already been
// dereferenced.
func (e *Encoder) marshalSlice(val reflect.Value) (string, error) {
	var sl, str string
	var ok bool
	if !supportedKind(val.Type().Elem().Kind()) {
		return "", UnsupportedTypeError{kind: val.Type().Elem().Kind(), method: "marshalSlice"}
	}
	// if the elem is a pointer, dereference to see it's kind
	for j := 0; j < val.Len(); j++ {
		str = ""
		str, ok = e.stringify(val.Index(j))
		if !ok {
			return "", UnsupportedTypeError{kind: val.Index(j).Kind(), method: "marshalSlice"}
		}
		if j == 0 {
			sl = str
			continue
		}
		sl = fmt.Sprintf("%s,%s", sl, str)
	}
	return sl, nil
}

// stringify takes a interface and returns the value it contains as a string.
// This is not ment for composite types.  If the type is not supported, an
// error will be returned.
func (e *Encoder) stringify(v reflect.Value) (string, bool) {
	if !supportedKind(v.Kind()) {
		return "", false
	}
	switch v.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.Itoa(int(v.Int())), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(uint64(v.Uint()), e.base), true
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'E', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'E', -1, 64), true
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%g", v.Complex()), true
	case reflect.String:
		return v.String(), true
	case reflect.Ptr:
		return e.stringify(v.Elem())
	default:
		cols, err := e.marshal(v)
		if err != nil {
			return "", false
		}
		r := cols[0]
		for i := 1; i < len(cols); i++ {
			r = fmt.Sprintf("%s,%s", r, cols[i])
		}
		if strings.HasPrefix(r, "(") {
			return r, true
		}
		return fmt.Sprintf("%s%s%s", e.sepBeg, r, e.sepEnd), true
	}
}

func supportedPtrKind(typ reflect.Type) (k, v reflect.Kind, ok bool) {
	switch typ.Kind() {
	case reflect.Ptr:
		return supportedPtrKind(typ.Elem())
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Ptr {
			return supportedPtrKind(typ.Elem())

		}
		k, ok = supportedSliceKind(typ)
		return k, reflect.Invalid, ok
	case reflect.Map:
		return supportedMapKind(typ)
	case reflect.Struct:
		fmt.Println("=============\n", "Implement supportedKind Struct\n=============")
	}
	return typ.Kind(), reflect.Invalid, supportedKind(typ.Kind())
}

func supportedMapKind(t reflect.Type) (k, v reflect.Kind, ok bool) {
	k, v, ok = mapKind(t)
	if !ok {
		return k, v, ok
	}
	if !supportedKind(k) {
		return k, v, false
	}
	if !supportedKind(v) {
		return k, v, false
	}
	return k, v, true
}

func mapKind(t reflect.Type) (k, v reflect.Kind, is bool) {
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Ptr {
			return mapKind(t.Elem())
		}
		if t.Elem().Kind() == reflect.Map {
			return t.Elem().Kind(), t.Elem().Kind(), true
		}
		return reflect.Invalid, reflect.Invalid, false
	case reflect.Map:
		k = baseKind(t.Key())
		switch t.Elem().Kind() {
		case reflect.Map:
			return mapKind(t.Elem())
		}
		v = baseKind(t.Elem())
		return k, v, true
	}
	return reflect.Invalid, reflect.Invalid, false
}

func supportedSliceKind(t reflect.Type) (reflect.Kind, bool) {
	k, ok := sliceKind(t)
	if !ok {
		return k, false
	}
	if supportedKind(k) {
		return k, true
	}
	return k, false
}

func sliceKind(t reflect.Type) (reflect.Kind, bool) {
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem().Kind() == reflect.Ptr {
			return sliceKind(t.Elem())
		}
		return t.Elem().Kind(), true
	case reflect.Slice, reflect.Array:
		return sliceKind(t.Elem())
	}
	if t.Kind() == reflect.Invalid {
		return t.Kind(), false
	}
	return t.Kind(), true
}

func baseKind(t reflect.Type) reflect.Kind {
	switch t.Kind() {
	case reflect.Ptr:
		return baseKind(t.Elem())
	case reflect.Slice:
		return baseKind(t.Elem())
	case reflect.Map:
		if t.Elem().Kind() == reflect.Map {
			k, _, _ := mapKind(t.Elem())
			return k
		}
		return baseKind(t.Elem())
	}
	return t.Kind()
}

// returns the actual kind of the value, ptrs will be dereference. slices will
// have their Kind returned, maps will return the Kind of the key.  This
// func assumes ptrs have been dereferenced.
func supportedKind(k reflect.Kind) bool {
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
	case reflect.Invalid:
		return false
	}
	return true
}
