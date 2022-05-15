package truffle

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type TruffleFile struct {
	ABI      json.RawMessage `json:"abi"`
	ByteCode string          `json:"bytecode"`
}

func LoadTruffleFile(filename string) (*TruffleFile, error) {
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var truffle TruffleFile

	json.Unmarshal(byteValue, &truffle)

	return &truffle, nil
}
