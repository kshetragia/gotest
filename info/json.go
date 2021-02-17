package info

import "encoding/json"

// Json returns collected info in JSON format.
func (pinfo *FullInfo) Json() ([]byte, error) {
	return json.Marshal(pinfo)
}
