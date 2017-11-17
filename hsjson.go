package hld

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//HSJSONTime for converting to javascript time (in milli seconds)
type HSJSONTime struct {
	time.Time
	ts int64
}

//Valid _
func (t *HSJSONTime) Valid() bool {
	return t.ts != 0
}

//MarshalJSON _
func (t *HSJSONTime) MarshalJSON() ([]byte, error) {
	if t.Valid() {
		return []byte("\"" + strconv.FormatInt(t.Time.UnixNano()/(1000000), 10) + "\""), nil
	}
	return []byte("null"), nil
}

//UnmarshalJSON _
func (t *HSJSONTime) UnmarshalJSON(b []byte) error {

	i, err := strconv.ParseInt(string(b), 10, 0)
	if err == nil {
		t.ts = i
		t.Time = JSTimeToTime(i)
		return nil
	}
	if string(b) == "null" {
		t.ts = 0
		return nil
	}
	return err
}

//HSJsonCommand json command packet
type HSJsonCommand struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

//ToJSON Convert to json bytes
func (cmd *HSJsonCommand) ToJSON() ([]byte, error) {
	return json.Marshal(cmd)
}

//HSJsonResult json result packet
type HSJsonResult struct {
	Msg     string                     `json:"msg"`
	Success bool                       `json:"success"`
	P       float64                    `json:"p"`
	Result  map[string]json.RawMessage `json:"result"`
	Dt      int64                      `json:"dt"`
}

//SendJSONCmd Send JSON command and get the result
func SendJSONCmd(cmd *HSJsonCommand, url string) (*HSJsonResult, error) {
	var result HSJsonResult
	var err error

	httpClient := &http.Client{}
	outData, err := cmd.ToJSON()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(outData))
	defer req.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp, err := httpClient.Do(req); err == nil {
		defer resp.Body.Close()
		if data, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(data, &result)
			return &result, nil
		}
	}
	return nil, err
}
