// Package struct2csv creates CSV rows out of structs.
package struct2csv

import (
  "fmt"
  "reflect"
)

// GetHeaders get's the column headers from the received struct.  If anything
// other than a struct is passed, an error will be returned.  If the struct
// has field tags for csv, those values will be used as the column headers,
// otherwise the field names will be used.
//
// If field tags other than the ones for csv are to be used, TODO figure out
// the struct and how to implement this comment.
func GetHeaders(v interface{}) ([]string, error) {
  if reflect.TypeOf(v).Kind() != reflect.Struct {
    return nil, fmt.Errorf("struct required: value was of type %s", reflect.TypeOf(v).Kind())
  }
  s := reflect.TypeOf(v)
  var hdrs []string
  for i := 0; i < s.NumField(); i++ {
    f := s.Field(i)
    // skip unexported
    if f.PkgPath != "" {
      continue
    }
    name := f.Tag.Get("csv")
    // If there isn't a matching field tag, use the Field Name
    if name == "" {
      name = f.Name
    }
    hdrs = append(hdrs, name)
  }
  return hdrs, nil
}
