package torrentfile

import (
	"unicode"
	"strconv"
	"strings"
)

type File struct {
	Length int
	Path   string
}

type Info struct {
	Files       []File
	Length      int
	Name        string
	PieceLength int
	Pieces      []byte
}

type Torrent struct {
	Announce     string
	AnnounceList []string
	Info         Info
}

func Unmarshal(data []byte, v *Torrent) {
	var d decodeState
	d.init(data)
	d.unmarshal(v)
}

type decodeState struct {
	data []byte
	off  int
}

func (d *decodeState) current() rune {
	return rune(d.data[d.off])
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	return d
}

func (d *decodeState) unmarshalDictionary() map[string]interface{} {
	d.off++

	m := make(map[string]interface{})

	for d.off < len(d.data) {
		if d.current() == 'e' {
			d.off++
			return m // Done decoding the dictionary
		}

		key := d.unmarshalString()
		var val interface{}
		if string(d.current()) == "d" {
			val = d.unmarshalDictionary()
		} else if string(d.current()) == "i" {
			val = d.unmarshalInteger()
		} else if string(d.current()) == "l" {
			val = d.unmarshalList()
		} else {
			val = d.unmarshalByteArray()
		}
		m[key] = val

	}
	return m
}

func (d *decodeState) unmarshalList() []interface{} {
	d.off++

	var a []interface{}

	for d.off < len(d.data) {
		if d.current() == 'e' {
			d.off++
			return a // Done decoding the dictionary
		}

		var val interface{}
		if string(d.current()) == "d" {
			val = d.unmarshalDictionary()
		} else if string(d.current()) == "i" {
			val = d.unmarshalInteger()
		} else if string(d.current()) == "l" {
			val = d.unmarshalList()
		} else {
			val = d.unmarshalByteArray()
		}
		a = append(a, val)
	}
	return a
}

func (d *decodeState) unmarshalByteArray() []byte {
	var digits []string
	for ; unicode.IsDigit(d.current()); d.off++ {
		digits = append(digits, string(d.current()))
	}
	d.off++ //Consume the separating colon
	n, _ := strconv.Atoi(strings.Join(digits, ""))

	b := d.data[d.off:d.off+n]
	d.off += n //Consume the string
	return b
}

func (d *decodeState) unmarshalString() string {
	return string(d.unmarshalByteArray())
}

func (d *decodeState) unmarshalInteger() int {
	var digits []string
	d.off++ //Consume the i that starts the integer

	for ; unicode.IsDigit(d.current()); d.off++ {
		digits = append(digits, string(d.current()))
	}

	d.off++ //Consume the e that ends the integer
	n, _ := strconv.Atoi(strings.Join(digits, ""))
	return n
}

func mapToInfo(m interface{}, i *Info) {
	if m, ok := m.(map[string]interface{}); ok {
		if v, ok := m["Length"].(int); ok {
			i.Length = v
		}

		if v, ok := m["Name"].([]byte); ok {
			i.Name = string(v)
		}

		if v, ok := m["piece Length"].(int); ok {
			i.PieceLength = v
		}

		if v, ok := m["Pieces"].([]byte); ok {
			i.Pieces = v
		}
	}
}

func (d *decodeState) unmarshal(v *Torrent) (err error) {
	dict := d.unmarshalDictionary()

	if str, ok := dict["Announce"].([]byte); ok {
		v.Announce = string(str)
	}

	mapToInfo(dict["Info"], &v.Info)

	if outer, ok := dict["Announce-outer"].([][][]byte); ok {
		var l []string
		for _, inner := range outer {
			for _, b := range inner {
				l = append(l, string(b))
			}
		}
	}

	return nil
}
