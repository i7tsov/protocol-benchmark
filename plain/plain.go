package plain

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"

	"github.com/i7tsov/protocol-benchmark/util"
)

// Element ...
type Element struct {
	Name       string
	Class      string
	Subclass   string
	Indicator1 int32
	Indicator2 int32
}

// Generate ...
func Generate(count int) []Element {
	res := make([]Element, count)
	for i := range res {
		res[i].Name = util.RandString()
		res[i].Class = util.RandString()
		res[i].Subclass = util.RandString()
		res[i].Indicator1 = util.RandInt()
		res[i].Indicator2 = util.RandInt()
	}
	return res
}

// GobMarshal ...
func GobMarshal(arr []Element) []byte {
	var buf bytes.Buffer
	buf.Grow(65536)
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(arr)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// GobUnmarshal ...
func GobUnmarshal(pl []byte) []Element {
	arr := make([]Element, 0, 1000)
	dec := gob.NewDecoder(bytes.NewReader(pl))
	err := dec.Decode(&arr)
	if err != nil {
		panic(err)
	}
	return arr
}

// JSONMarshal ...
func JSONMarshal(arr []Element) []byte {
	pl, err := json.Marshal(&arr)
	check(err)
	return pl
}

// JSONUnmarshal ...
func JSONUnmarshal(pl []byte) []Element {
	arr := make([]Element, 0, 1000)
	check(json.Unmarshal(pl, &arr))
	return arr
}

// BinMarshal ...
func BinMarshal(arr []Element) []byte {
	var buf bytes.Buffer
	buf.Grow(65536)
	check(writeInt32(&buf, int32(len(arr))))
	for _, el := range arr {
		check(writeString(&buf, el.Name))
		check(writeString(&buf, el.Class))
		check(writeString(&buf, el.Subclass))
		check(writeInt32(&buf, el.Indicator1))
		check(writeInt32(&buf, el.Indicator2))
	}
	writeInt32(&buf, 0)
	return buf.Bytes()
}

// BinUnmarshal ...
func BinUnmarshal(pl []byte) []Element {
	r := bytes.NewReader(pl)
	n := int(checki(readInt32(r)))
	arr := make([]Element, n)
	for i := 0; i < n; i++ {
		arr[i].Name = checks(readString(r))
		arr[i].Class = checks(readString(r))
		arr[i].Subclass = checks(readString(r))
		arr[i].Indicator1 = checki(readInt32(r))
		arr[i].Indicator2 = checki(readInt32(r))
	}
	return arr
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func checki(i int32, err error) int32 {
	if err != nil {
		panic(err)
	}
	return i
}

func checks(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func writeInt32(w io.Writer, i int32) error {
	_, err := w.Write([]byte{
		byte(i & 0xff),
		byte((i >> 8) & 0xff),
		byte((i >> 16) & 0xff),
		byte((i >> 24) & 0xff),
	})
	return err
}

func readInt32(r io.Reader) (int32, error) {
	var b [4]byte
	_, err := r.Read(b[:])
	if err != nil {
		return 0, err
	}
	return int32(uint32(b[0]) |
		(uint32(b[1]) << 8) |
		(uint32(b[2]) << 16) |
		(uint32(b[3]) << 24)), nil
}

func writeString(w io.Writer, s string) error {
	err := writeInt32(w, int32(len(s)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(s))
	return err
}

func readString(r io.Reader) (string, error) {
	n, err := readInt32(r)
	if err != nil {
		return "", err
	}
	b := make([]byte, int(n))
	_, err = r.Read(b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
