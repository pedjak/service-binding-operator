package jsonpath

import (
	"fmt"
	"k8s.io/client-go/util/jsonpath"
)

func getValuesByJSONPath(obj map[string]interface{}, path string) ([]interface{}, error) {
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
	returnedResult := make([]interface{}, 0, len(result))
	for _, r := range result[0] {
		returnedResult = append(returnedResult, r.Interface())
	}
	return returnedResult, err
}
