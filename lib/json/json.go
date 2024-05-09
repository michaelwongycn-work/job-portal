package json

import (
	"encoding/json"
	"os"
)

func ReadJSON(fname string, cfg any) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return err
	}

	return nil
}
