// A simple example cli that generates a csv file using canned data.  The file
// is written to tmp and the location of the file is displayed to the user.
package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"

	"github.com/mohae/struct2csv"
)

type Person struct {
	Name    string
	Id      int
	Address map[string]Address
	Phone   map[string]Phone
	Tags    []string
}

type Address struct {
	Typ          string
	Organization string
	Street1      string
	Street2      string
	City         string
	State        string
	Code         string
}

func (a *Address) String() string {
	s := a.Organization
	if a.Street1 != "" {
		if s == "" {
			s = a.Street1
		} else {
			s = fmt.Sprintf("%s, %s", s, a.Street1)
		}
	}
	if a.Street2 != "" {
		if s == "" {
			s = a.Street2
		} else {
			s = fmt.Sprintf("%s, %s", s, a.Street2)
		}
	}
	if a.City != "" {
		if s == "" {
			s = a.City
		} else {
			s = fmt.Sprintf("%s, %s", s, a.City)
		}
	}
	if a.State != "" {
		if s == "" {
			s = a.State
		} else {
			s = fmt.Sprintf("%s, %s", s, a.State)
		}
	}
	if a.Code != "" {
		if s == "" {
			s = a.Code
		} else {
			s = fmt.Sprintf("%s, %s", s, a.Code)
		}
	}
	return s
}

type Phone struct {
	Typ        string
	NationCode string
	AreaCode   string
	Prefix     string
	Suffix     string
	Ext        string
}

func (p *Phone) String() string {
	num := fmt.Sprintf("%s:", p.Typ)
	if p.NationCode != "" {
		num = fmt.Sprintf("%s +%s", num, p.NationCode)
	}
	if p.AreaCode != "" {
		num = fmt.Sprintf("%s (%s)", num, p.AreaCode)
	}

	if p.Prefix == "" {
		goto EXT
	}
	if p.AreaCode != "" {
		num = fmt.Sprintf("%s-%s", num, p.Prefix)
	} else {
		num = fmt.Sprintf("%s %s", num, p.Prefix)
	}
	if p.Suffix != "" {
		num = fmt.Sprintf("%s-%s", num, p.Suffix)
	}
EXT:
	if p.Ext != "" {
		num = fmt.Sprintf("%s Ext. %s", num, p.Ext)
	}
	return num
}

func main() {
	people := []Person{
		Person{
			Name: "Jack Straw", Id: 420,
			Address: map[string]Address{
				"Work": Address{
					Typ:          "Work",
					Organization: "City Hall",
					Street1:      "544 N. Main",
					City:         "Wichita",
					State:        "KS",
					Code:         "67202",
				},
			},
			Phone: map[string]Phone{
				"Work": Phone{
					Typ:      "Work",
					AreaCode: "316",
					Prefix:   "942",
					Suffix:   "4482",
				},
			},
			Tags: []string{"Sante Fe", "Cheyenne", "Tuscon"},
		},
		Person{
			Name: "Sugar Magnolia", Id: 71,
			Address: map[string]Address{
				"Work": Address{
					Typ:          "Work",
					Organization: "Preservation Hall",
					Street1:      "726 St. Peters St.",
					City:         "New Orleans",
					State:        "LA",
					Code:         "70116",
				},
			},
			Phone: map[string]Phone{
				"Work": Phone{
					Typ:      "Work",
					AreaCode: "504",
					Prefix:   "522",
					Suffix:   "2841",
				},
			},
			Tags: []string{"jazz", "french quarter", "live music", "education"},
		},
	}

	enc := struct2csv.New()
	enc.SetSeparators("\"", "\"")
	data, err := enc.Marshal(people)
	// open a tmp file to write to
	f, err := ioutil.TempFile("", "CSV")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// get a new csv writer
	w := csv.NewWriter(f)
	// encode
	err = w.WriteAll(data)
	if err != nil {
		fmt.Println(err)
	}
	// output message
	fmt.Printf("Data marshaled to CSV and saved as %s\n", f.Name())
}
