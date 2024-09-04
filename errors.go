package stripe

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code    string `json:"code"`
	DocURL  string `json:"doc_url"`
	Message string `json:"message"`
	Param   string `json:"param"`
	Type    string `json:"type"`
}

func (err Error) Error() string {
	return fmt.Sprintf("%s See %s for more information.", err.Code, err.DocURL)
}

func (err *Error) UnmarshalJSON(bytes []byte) error {
	var tmp struct {
		Error struct {
			Code    string `json:"code"`
			DocURL  string `json:"doc_url"`
			Message string `json:"message"`
			Param   string `json:"param"`
			Type    string `json:"type"`
		} `json:"error"`
	}

	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return err
	}
	*err = tmp.Error
	return nil
}

func (err Error) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Error struct {
			Code    string `json:"code"`
			DocURL  string `json:"doc_url"`
			Message string `json:"message"`
			Param   string `json:"param"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	tmp.Error = err
	return json.Marshal(tmp)

}
