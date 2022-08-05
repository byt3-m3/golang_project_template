package struct_utils

import (
	"encoding/json"
	"log"
)

func SerializeStructToJSON(u interface{}) (string, error) {
	data, err := SerializeStructToBytes(u)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

func SerializeStructToBytes(u interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(u)
	if err != nil {
		return []byte{}, err
	}
	return jsonData, nil

}

func ToJson(u interface{}) string {
	data, err := SerializeStructToBytes(u)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(data)
}
