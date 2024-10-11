package encodings_test

import (
	"statelessdb/internal/helpers"
)

// SampleStruct is a sample data structure for testing.
type SampleStruct struct {
	ID      int
	Name    string
	Numbers []int
	Details map[string]string
}

//// GobEncode implements the GobEncoder interface.
//// It serializes only the exported fields, excluding private ones.
//func (s *SampleStruct) GobEncode() ([]byte, error) {
//	var buf bytes.Buffer
//	Encoder := gob.NewEncoder(&buf)
//	if err := Encoder.Encode(s.ID); err != nil { return nil, err }
//	if err := Encoder.Encode(s.Name); err != nil { return nil, err }
//	if err := Encoder.Encode(s.Numbers); err != nil { return nil, err }
//	if err := Encoder.Encode(s.Details); err != nil { return nil, err }
//	return buf.Bytes(), nil
//}
//
//// GobDecode implements the GobDecoder interface.
//// It deserializes only the exported fields, leaving private ones untouched.
//func (s *SampleStruct) GobDecode(data []byte) error {
//	buf := bytes.NewBuffer(data)
//	Decoder := gob.NewDecoder(buf)
//	if err := Decoder.Decode(&s.ID); err != nil { return err }
//	if err := Decoder.Decode(&s.Name); err != nil { return err }
//	if err := Decoder.Decode(&s.Numbers); err != nil { return err }
//	if err := Decoder.Decode(&s.Details); err != nil { return err }
//	return nil
//}

// Equals returns true if both SampleStructs are equal state. This is not counting
// calculated data like internal caches, attackMaps, etc. which should be same
// anyway if this data is equal.
func (s *SampleStruct) Equals(other *SampleStruct) bool {
	if s == other {
		return true
	}
	if other == nil {
		return false
	}
	if s.ID != other.ID ||
		s.Name != other.Name {
		return false
	}
	if !helpers.CompareSlices(s.Numbers, other.Numbers) {
		return false
	}
	if !helpers.CompareMaps(s.Details, other.Details) {
		return false
	}
	return true
}
