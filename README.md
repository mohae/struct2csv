# struct2csv
[![Build Status](https://travis-ci.org/mohae/struct2csv.png)](https://travis-ci.org/mohae/struct2csv)

Create CSV from a slice of structs.

The output of a marshal is a `[][]string`, which `encoding/CSV` can use.

## About
Struct2csv takes a slice of structs and creates csv from them.  The column names are used as the first row of the csv data.  Use of field tags for csv header row column names is supported.  By default, struct2csv uses the looks for field tags for `csv`.  It can be configured to use the values of other field tags, e.g. `yaml` or `json`, instead.  Each slice element becomes its own row.  Columns that are slices will be transformed to a quoted list of comma separated values.  Columns that are maps will be transformed into a list of quoted key:value pairs.

Only exported columns become part of the csv data.  Some types, like channels and funcs, are skipped.

If a non-slice is received, an error will be returned.  Any error encountered will be returned.

This is an initial implementation.

## Usage

    import github.com/mohae/struct2csv

    var data []MyStruct
    tc := struct2csv.New()
    rows, err := tc.Marshal(data)
    if err != nil {
        // handle error
    }

### Configuration of Transcoder
By default, a transcoder will use tag fields with the tag `csv`, if they exist, as the column header value for a field. If such a tag does not exist, the column name will be used.

The tag that the transcoder uses can be changed by calling `Transcoder.SetTag(value)`.

Tags can be ignored by calling `Transcoder.SetUseTag(false)`.  This will result in the struct field names being used as the colmn header values.

More customization of transcoder's behavior may be added.

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
Slice
String
Struct
```

The following `reflect.Kind` are ignored:
```
Chan
Func
```

The following `reflect.Kind` are not supported, or have not yet been implemented.  Any structs with any of these Kinds will cause an error:
```
Uintptr
Interface
Ptr
UnsafePointer
```

## Notes

Pointers have not yet been implemented.

### Embedded types
If a type is embeded, any exported fields within that struct become their own columns with the field name being the column name, unless a field tag has been defined.  The name of the embedded struct does not become part of the column header name.

### Maps, Slices, and Arrays
#### Map
Maps are a single column in the resulting CSV as maps can have a variable number of elements and there is no way to account for this within CSV.  Each map element becomes a `key:value` pair with each element seperated by a `,`.  

    map[string]string{"Douglas Adams": "Hitchhiker's Guide to the Galaxy", "William Gibson": "Neuromancer"}

becomes:

    Douglas Adams:Hitchhiker's Guide to the Galaxy, William Gibson:Neuromancer

If the map's value is a composite type, the values of the composite type become a comma separated list surrounded by `()`.

    map[string][]string{"William Gibson": []string{"Neuromancer" ,"Count Zero", "Mona Lisa Overdrive"}}

becomes:

    William Gibson:[Neuromancer, Count Zero, Mona Lisa Overdrive]

#### Slices and Arrays
Slices and arrays are a single column in the resulting CSV as slices can have a variable number of elements and there is no way to account for this within CSV.  Arrays are treated the same as slices.  Slices become a comma separated list of values.

#### Structs
If a map or slice type is a struct, the fields of the struct become a comma separated list of values for the column. Struct columns that become lists are enclosed with brackets,`[]`.  For examples see the tests.

### Header row
It is possible to get the header row for a struct by calling the `CSVHeaders()` func with the struct for which you want the column names

## TODO
Add support pointers
