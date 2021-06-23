package jsonpath

import (
	"fmt"
	"k8s.io/client-go/util/jsonpath"
	"reflect"
)

func GetValuesByJSONPath(obj map[string]interface{}, path string) ([]reflect.Value, error) {
	j := jsonpath.New("")
	err := j.Parse(path)
	if err != nil {
		return nil, err
	}
	result, err := j.FindResults(obj)
	if err != nil {
		return nil, err
	}
	if len(result) > 1 {
		return nil, fmt.Errorf("Expecting max list of results, but not list of lists %v", result)
	}
	return result[0], nil
}
