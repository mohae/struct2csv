# struct2csv
[![Build Status](https://travis-ci.org/mohae/struct2csv.png)](https://travis-ci.org/mohae/struct2csv)

Create CSV from a slice of structs.

The output of a marshal is a `[][]string`, which `encoding/CSV` can use.

## About
Struct2csv takes a slice of structs and creates CSV data from them.  

The field names are used as the first row of the csv data.  Use of field tags for csv header row column names is supported.  By default, struct2csv uses the looks for field tags for `csv`.  It can be configured to use the values of other field tags, e.g. `yaml` or `json`, instead, or, to ignore field tags.  

Each slice element becomes its own row.  Struct fields that are slices will be transformed to a quoted list of comma separated values.  Struct fields that are maps will be transformed into a list of quoted key:value pairs.  More complex types use separators, "(" and ")" by default, to group subtypes: e.g. a field with the type `map[string][]string` will have an output similar to `key1: (key1, slice, values)`.  

The separators to use can be set with `encoder.SetSeparators(begin, end)`. Passing empty strings, `""`, will result in no separators being used.  The separators are used for composite types with lists.

Only exported columns become part of the csv data.  Some types, like channels and funcs, are skipped.

If a non-slice is received, an error will be returned.  Any error encountered will be returned.

## Usage

    import github.com/mohae/struct2csv

    var data []MyStruct
    enc := struct2csv.New()
    rows, err := enc.Marshal(data)
    if err != nil {
        // handle error
    }

### Configuration of an Encoder
By default, an encoder will use tag fields with the tag `csv`, if they exist, as the column header value for a field. If such a tag does not exist, the column name will be used.  The encoder will also use `(` and `)` as its begin and end separator values.

The separator values can be changed with the `Encoder.SetSeparators(beginSep, endSep)` method.  If the separators are set to `""`, an empty string, nothing will be used.  This mainly applies to lists.

The tag that the encoder uses can be changed by calling `Encoder.SetTag(value)`.

Tags can be ignored by calling `Encoder.SetUseTag(false)`.  This will result in the struct field names being used as the colmn header values.

## Supported types
The following `reflect.Kind` are supported:  
```
Bool
Int
Int8
Int16
Int32
Int64
Uint
Uint8
Uint16
Uint32
Uint64
Float32
Float64
Complex64
Complex128
Array
Map
Ptr
Slice
String
Struct
```

The following `reflect.Kind` are not supported, or have not yet been implemented.  Any fields using any of these kinds will be ignored. If a map uses any of these Kinds in either its key or value, it will be ignored.
```
Chan
Func
Uintptr
Interface
UnsafePointer
```

### Embedded types
If a type is embedded, any exported fields within that struct become their own columns with the field name being the column name, unless a field tag has been defined.  The name of the embedded struct does not become part of the column header name.

### Maps, Slices, and Arrays
#### Map
Maps are a single column in the resulting CSV as maps can have a variable number of elements and there is no way to account for this within CSV.  Each map element becomes a `key:value` pair with each element seperated by a `,`.  Keys are sorted.  

    map[string]string{"Douglas Adams": "Hitchhiker's Guide to the Galaxy", "William Gibson": "Neuromancer"}

becomes:

    Douglas Adams:Hitchhiker's Guide to the Galaxy, William Gibson:Neuromancer

If the map's value is a composite type, the values of the composite type become a comma separated list surrounded by `()`.

    map[string][]string{
            "William Gibson": []string{"Neuromancer" ,"Count Zero", "Mona Lisa Overdrive"},
            "Frank Herbert": []string{"Destination Void", "Jesus Incident", "Lazurus Effect"},
    }

becomes:

    ((Frank Herbert: (Destination Void, Jesus Incident, Lazurus Effect)),
    (William Gibson:(Neuromancer, Count Zero, Mona Lisa Overdrive)))

#### Slices and Arrays
Slices and arrays are a single column in the resulting CSV as slices can have a variable number of elements and there is no way to account for this within CSV.  Arrays are treated the same as slices.  Slices become a comma separated list of values.

#### Structs
TODO add doc.

### Header row
It is possible to get the header row for a struct by calling the `GetHeaders()` func with the struct for which you want the column names

## TODO

* Add support for field tag value of `-` to explicitly skip fields.
* Add option to add names of embedded structs to the column header for its fields.
* Add support for output via encoding/csv
