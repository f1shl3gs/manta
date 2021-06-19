package manta_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/f1shl3gs/manta"
)

func TestDecodeFromString(t *testing.T) {
	var id manta.ID
	err := id.DecodeFromString("020f755c3c082000")
	if err != nil {
		t.Errorf(err.Error())
	}
	want := []byte{48, 50, 48, 102, 55, 53, 53, 99, 51, 99, 48, 56, 50, 48, 48, 48}
	got, _ := id.Encode()
	if !bytes.Equal(want, got) {
		t.Errorf("got %s not equal to wanted %s", string(got), string(want))
	}
	if id.String() != "020f755c3c082000" {
		t.Errorf("expecting string representation to contain the right value")
	}
	if !id.Valid() {
		t.Errorf("expecting ID to be a valid one")
	}
}

func TestEncode(t *testing.T) {
	var id manta.ID
	if _, err := id.Encode(); err == nil {
		t.Errorf("encoding an invalid ID should not be possible")
	}

	id.DecodeFromString("5ca1ab1eba5eba11")
	want := []byte{53, 99, 97, 49, 97, 98, 49, 101, 98, 97, 53, 101, 98, 97, 49, 49}
	got, _ := id.Encode()
	if !bytes.Equal(want, got) {
		t.Errorf("encoding error")
	}
	if id.String() != "5ca1ab1eba5eba11" {
		t.Errorf("expecting string representation to contain the right value")
	}
	if !id.Valid() {
		t.Errorf("expecting ID to be a valid one")
	}
}

func TestDecodeFromAllZeros(t *testing.T) {
	var id manta.ID
	err := id.Decode(make([]byte, manta.IDLength))
	if err == nil {
		t.Errorf("expecting all zeros ID to not be a valid ID")
	}
}

func TestDecodeFromShorterString(t *testing.T) {
	var id manta.ID
	err := id.DecodeFromString("020f75")
	if err == nil {
		t.Errorf("expecting shorter inputs to error")
	}
	if id.String() != "" {
		t.Errorf("expecting invalid ID to be serialized into empty string")
	}
}

func TestDecodeFromLongerString(t *testing.T) {
	var id manta.ID
	err := id.DecodeFromString("020f755c3c082000aaa")
	if err == nil {
		t.Errorf("expecting shorter inputs to error")
	}
	if id.String() != "" {
		t.Errorf("expecting invalid ID to be serialized into empty string")
	}
}

func TestDecodeFromEmptyString(t *testing.T) {
	var id manta.ID
	err := id.DecodeFromString("")
	if err == nil {
		t.Errorf("expecting empty inputs to error")
	}
	if id.String() != "" {
		t.Errorf("expecting invalid ID to be serialized into empty string")
	}
}

func TestMarshalling(t *testing.T) {
	var id0 manta.ID
	_, err := json.Marshal(id0)
	if err == nil {
		t.Errorf("expecting empty ID to not be a valid one")
	}

	init := "ca55e77eca55e77e"
	id1, err := manta.IDFromString(init)
	if err != nil {
		t.Errorf(err.Error())
	}

	serialized, err := json.Marshal(id1)
	if err != nil {
		t.Errorf(err.Error())
	}

	var id2 manta.ID
	json.Unmarshal(serialized, &id2)

	bytes1, _ := id1.Encode()
	bytes2, _ := id2.Encode()

	if !bytes.Equal(bytes1, bytes2) {
		t.Errorf("error marshalling/unmarshalling ID")
	}

	// When used as a map key, IDs must use their string encoding.
	// If you only implement json.Marshaller, they will be encoded with Go's default integer encoding.
	b, err := json.Marshal(map[manta.ID]int{0x1234: 5678})
	if err != nil {
		t.Error(err)
	}
	const exp = `{"0000000000001234":5678}`
	if string(b) != exp {
		t.Errorf("expected map to json.Marshal as %s; got %s", exp, string(b))
	}

	var idMap map[manta.ID]int
	if err := json.Unmarshal(b, &idMap); err != nil {
		t.Error(err)
	}
	if len(idMap) != 1 {
		t.Errorf("expected length 1, got %d", len(idMap))
	}
	if idMap[0x1234] != 5678 {
		t.Errorf("unmarshalled incorrectly; exp 0x1234:5678, got %v", idMap)
	}
}

func TestValid(t *testing.T) {
	var id manta.ID
	if id.Valid() {
		t.Errorf("expecting initial ID to be invalid")
	}

	if manta.InvalidID() != 0 {
		t.Errorf("expecting invalid ID to return a zero ID, thus invalid")
	}
}

func TestID_GoString(t *testing.T) {
	type idGoStringTester struct {
		ID manta.ID
	}
	var x idGoStringTester

	const idString = "02def021097c6000"
	if err := x.ID.DecodeFromString(idString); err != nil {
		t.Fatal(err)
	}

	sharpV := fmt.Sprintf("%#v", x)
	want := `manta_test.idGoStringTester{ID:"` + idString + `"}`
	if sharpV != want {
		t.Fatalf("bad GoString: got %q, want %q", sharpV, want)
	}
}

func BenchmarkIDEncode(b *testing.B) {
	var id manta.ID
	id.DecodeFromString("5ca1ab1eba5eba11")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b, _ := id.Encode()
		_ = b
	}
}

func BenchmarkIDDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var id manta.ID
		id.DecodeFromString("5ca1ab1eba5eba11")
	}
}

func TestName(t *testing.T) {
	fmt.Println(10e3)
}
