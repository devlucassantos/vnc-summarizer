package converters

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"strconv"
)

func ToInt(data interface{}) (int, error) {
	number, err := strconv.ParseFloat(fmt.Sprint(data), 64)
	if err != nil {
		log.Error("Error converting data to integer: ", err.Error())
		return 0, err
	}

	return int(number), nil
}

func ToJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("Error converting data to JSON: ", err.Error())
		return nil, err
	}

	return jsonData, nil
}

func ToMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := ToJson(data)
	if err != nil {
		log.Error("ToJson(): ", err.Error())
		return nil, err
	}

	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		log.Error("Error converting data to map[string]interface{}: ", err.Error())
		return nil, err
	}

	return resultMap, err
}

func ToMapSlice(data interface{}) ([]map[string]interface{}, error) {
	jsonData, err := ToJson(data)
	if err != nil {
		log.Error("ToJson(): ", err.Error())
		return nil, err
	}

	var resultMap []map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		log.Error("Error converting data to []map[string]interface{}: ", err.Error())
		return nil, err
	}

	return resultMap, nil
}

func IntSliceToUniqueIntSlice(data []int) []int {
	var uniqueIntSlice []int
	uniqueIntMap := map[int]struct{}{}
	for _, number := range data {
		if _, exists := uniqueIntMap[number]; !exists {
			uniqueIntMap[number] = struct{}{}
			uniqueIntSlice = append(uniqueIntSlice, number)
		}
	}

	return uniqueIntSlice
}

func StringSliceToUniqueStringSlice(data []string) []string {
	var uniqueStringSlice []string
	uniqueStringMap := map[string]struct{}{}
	for _, text := range data {
		if _, exists := uniqueStringMap[text]; !exists {
			uniqueStringMap[text] = struct{}{}
			uniqueStringSlice = append(uniqueStringSlice, text)
		}
	}

	return uniqueStringSlice
}
