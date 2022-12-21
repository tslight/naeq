package jsn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// https://stackoverflow.com/a/36922225/11133327
func Valid(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func FromFile(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// fmt.Print(string(byteValue))
	if !Valid(string(byteValue)) {
		return nil, fmt.Errorf("%s is not a valid JSON file\n", path)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}
