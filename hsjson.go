package hld

import (
	"fmt"
	"time"
)

//HSJSONTime for converting to javascript time (in milli seconds)
type HSJSONTime time.Time

//MarshalJSON _
func (t HSJSONTime) MarshalJSON() ([]byte, error) {
	stamp := "\"" + fmt.Sprint(time.Time(t).UnixNano()/(1000000)) + "\""
	return []byte(stamp), nil
}
