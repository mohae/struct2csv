package struct2csv

import (
	"bytes"
	"testing"
)

var rowTests = []struct {
	Input   [][]string
	UseCRLF bool
	Output  string
}{
	{Input: [][]string{}, UseCRLF: false, Output: ""},
	{Input: [][]string{}, UseCRLF: true, Output: ""},
	{Input: [][]string{[]string{""}}, UseCRLF: false, Output: "\n"},
	{Input: [][]string{[]string{""}}, UseCRLF: true, Output: "\r\n"},
	{Input: [][]string{[]string{"abc"}}, UseCRLF: false, Output: "abc\n"},
	{Input: [][]string{[]string{"def"}}, UseCRLF: true, Output: "def\r\n"},
	{Input: [][]string{[]string{"Col1"}, []string{"abc"}}, UseCRLF: false, Output: "Col1\nabc\n"},
	{Input: [][]string{[]string{"ColA"}, []string{"def"}}, UseCRLF: true, Output: "ColA\r\ndef\r\n"},
	{Input: [][]string{[]string{"a", "b", "c"}}, UseCRLF: false, Output: "a,b,c\n"},
	{Input: [][]string{[]string{"d", "e", "f"}}, UseCRLF: true, Output: "d,e,f\r\n"},

	{Input: [][]string{[]string{"Col1", "Col2", "Col3"}, []string{"a", "b", "c"}}, UseCRLF: false, Output: "Col1,Col2,Col3\na,b,c\n"},
	{Input: [][]string{[]string{"ColA", "ColB", "ColC"}, []string{"d", "e", "f"}}, UseCRLF: true, Output: "ColA,ColB,ColC\r\nd,e,f\r\n"},
}

func TestRows(t *testing.T) {
	for i, test := range rowTests {
		b := &bytes.Buffer{}
		w := NewWriter(b)
		w.SetUseCRLF(test.UseCRLF)
		err := w.WriteAll(test.Input)
		if err != nil {
			t.Errorf("%d: unexpected error: %s", i, err)
		}
		out := b.String()
		if out != test.Output {
			t.Errorf("%d: out=%q want %q", i, out, test.Output)
		}
	}

	for i, test := range rowTests {
		buff := &bytes.Buffer{}
		w := NewWriter(buff)
		w.SetUseCRLF(test.UseCRLF)
		for j, v := range test.Input {
			err := w.Write(v)
			if err != nil {
				t.Errorf("%d:%d: unexpected error: %s", i, j, err)
			}
		}
		w.Flush()
		out := buff.String()
		if out != test.Output {
			t.Errorf("%d: out=%q want $q", i, out, test.Output)
		}

	}
}

var structTests = []struct {
	Input   []Basic
	UseTags bool
	Tag     string
	Output  string
	Error   string
}{
	{[]Basic{}, true, "csv", "\n", "struct2csv: the slice of structs was empty"},
	{[]Basic{}, true, "json", "\n", "struct2csv: the slice of structs was empty"},
	{[]Basic{}, true, "", "\n", "struct2csv: the slice of structs was empty"},
	{[]Basic{}, false, "", "\n", "struct2csv: the slice of structs was empty"},
	{[]Basic{
		Basic{Name: "Fyodor Dostoyevsky", List: []string{"Brothers Karamazov", "Crime and Punishment"}},
	}, true, "csv", "Nom,Liste\nFyodor Dostoyevsky,\"(Brothers Karamazov,Crime and Punishment)\"\n", ""},
	{[]Basic{
		Basic{Name: "Fyodor Dostoyevsky", List: []string{"Brothers Karamazov", "Crime and Punishment"}},
	}, true, "json", "name,list\nFyodor Dostoyevsky,\"(Brothers Karamazov,Crime and Punishment)\"\n", ""},
	{[]Basic{
		Basic{Name: "Fyodor Dostoyevsky", List: []string{"Brothers Karamazov", "Crime and Punishment"}},
	}, true, "", "Nom,Liste\nFyodor Dostoyevsky,\"(Brothers Karamazov,Crime and Punishment)\"\n", ""},
	{[]Basic{
		Basic{Name: "Fyodor Dostoyevsky", List: []string{"Brothers Karamazov", "Crime and Punishment"}},
	}, false, "", "Name,List\nFyodor Dostoyevsky,\"(Brothers Karamazov,Crime and Punishment)\"\n", ""},
	{[]Basic{
		Basic{Name: "Anatoly Rybakov", List: []string{"Children of the Arbat", "Fear", "Dust and Ashes"}},
		Basic{Name: "Vladimir Nabokov", List: []string{"Lolita", "Pnin", "Pale Fire"}},
	}, true, "csv",
		"Nom,Liste\nAnatoly Rybakov,\"(Children of the Arbat,Fear,Dust and Ashes)\"\nVladimir Nabokov,\"(Lolita,Pnin,Pale Fire)\"\n", ""},
	{[]Basic{
		Basic{Name: "Anatoly Rybakov", List: []string{"Children of the Arbat", "Fear", "Dust and Ashes"}},
		Basic{Name: "Vladimir Nabokov", List: []string{"Lolita", "Pnin", "Pale Fire"}},
	}, true, "json",
		"name,list\nAnatoly Rybakov,\"(Children of the Arbat,Fear,Dust and Ashes)\"\nVladimir Nabokov,\"(Lolita,Pnin,Pale Fire)\"\n", ""},
	{[]Basic{
		Basic{Name: "Anatoly Rybakov", List: []string{"Children of the Arbat", "Fear", "Dust and Ashes"}},
		Basic{Name: "Vladimir Nabokov", List: []string{"Lolita", "Pnin", "Pale Fire"}},
	}, true, "",
		"Nom,Liste\nAnatoly Rybakov,\"(Children of the Arbat,Fear,Dust and Ashes)\"\nVladimir Nabokov,\"(Lolita,Pnin,Pale Fire)\"\n", ""},
	{[]Basic{
		Basic{Name: "Anatoly Rybakov", List: []string{"Children of the Arbat", "Fear", "Dust and Ashes"}},
		Basic{Name: "Vladimir Nabokov", List: []string{"Lolita", "Pnin", "Pale Fire"}},
	}, false, "",
		"Name,List\nAnatoly Rybakov,\"(Children of the Arbat,Fear,Dust and Ashes)\"\nVladimir Nabokov,\"(Lolita,Pnin,Pale Fire)\"\n", ""},
}

func TestStructs(t *testing.T) {
	for i, test := range structTests {
		buff := &bytes.Buffer{}
		w := NewWriter(buff)
		w.SetTag(test.Tag)
		w.SetUseTags(test.UseTags)
		err := w.WriteStructs(test.Input)
		if err != nil {
			if err.Error() != test.Error {
				t.Errorf("%d: expected error %q got %q", i, test.Error, err)
			}
			continue
		}
		s := buff.String()
		if s != test.Output {
			t.Errorf("%d: got %s, want %s", i, s, test.Output)
		}
	}

	for i, test := range structTests {
		if len(test.Input) == 0 {
			continue
		}
		buff := &bytes.Buffer{}
		w := NewWriter(buff)
		w.SetTag(test.Tag)
		w.SetUseTags(test.UseTags)
		err := w.WriteColNames(test.Input[0])
		if err != nil {
			t.Errorf("%d: unexpected error: %s", i, err)
			continue
		}
		for j, row := range test.Input {
			err = w.WriteStruct(row)
			if err != nil {
				t.Errorf("%d:%d unexpected error: %s", i, j, err)
				continue
			}
		}
		w.Flush()
		s := buff.String()
		if s != test.Output {
			t.Errorf("%d: got %s, want %s", i, s, test.Output)
		}
	}
}
