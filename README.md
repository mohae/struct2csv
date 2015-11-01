# struct2csv
[![Build Status](https://travis-ci.org/mohae/struct2csv.png)](https://travis-ci.org/mohae/struct2csv)

__Under development__  
Create CSV from a slice of structs.

## About
Struct2csv takes a slice of structs and creates csv from them.  The column names are used as the first row of the csv data.  Use of field tags for csv header row column names is supported.  By default, struct2csv uses the looks for field tags for `csv`.  It can be configured to use the values of other field tags, e.g. `yaml` or `json`, instead.  Each slice element becomes its own row.  Columns that are slices will be transformed to a quoted list of comma separated values.  Columns that are maps will be transformed into a list of quoted key:value pairs.

Only exported columns become part of the csv data.  Some types, like channels and funcs, are skipped.

If a non-slice is received, an error will be returned.  Any error encountered will be returned.

This is an initial implmentation.  Not all functionality is supported yet.

## Header row
It is possible to get the header row for a struct by calling the `CSVHeaders()` func with the struct for which you want the column names

## TODO
Add support for embedded structs.
