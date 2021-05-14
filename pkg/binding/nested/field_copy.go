package nested

import (
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func NestedFieldCopy(obj map[string]interface{}, path ...string) (interface{}, bool, error) {

	startPath, endPath, isslice, sliceIndexString := splitPath(path)

	if !isslice {
		return unstructured.NestedFieldCopy(obj, startPath...)
	}

	indexValue := getIndexValue(sliceIndexString)

	switch indexValue.Type {
	case INDEX_NUMERIC:
		val, ok, err := unstructured.NestedSlice(obj, startPath...)
		if !ok || err != nil {
			return nil, ok, err
		}
		if len(val) <= indexValue.NumericValue {
			return nil, false, nil
		}

		subObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&val[indexValue.NumericValue])
		if err != nil {
			return nil, false, err
		}
		return NestedFieldCopy(subObject, endPath...)

	case INDEX_EQUALS:
		val, ok, err := unstructured.NestedSlice(obj, startPath...)
		if !ok || err != nil {
			return nil, ok, err
		}
		for _, entry := range val {
			cast, ok := entry.(map[string]interface{})
			if !ok {
				return nil, false, nil
			}

			compareVal, ok, err := unstructured.NestedString(cast, indexValue.Path...)
			if !ok || err != nil {
				return nil, ok, err
			}
			if "\""+compareVal+"\"" == indexValue.CompareValue {
				return NestedFieldCopy(cast, endPath...)
			}

		}
		return nil, false, nil

	default:
		return nil, false, nil
	}
}

func splitPath(path []string) (startPath []string, endPath []string, isSlice bool, sliceIndexString string) {
	i := 0
	for _, subpath := range path {
		if strings.HasSuffix(subpath, "]") {
			parts := strings.Split(subpath[:len(subpath)-1], "[")
			startPath = append(startPath, parts[0])
			isSlice = true
			sliceIndexString = parts[1]
			endPath = path[i+1:]
			break
		}
		startPath = append(startPath, subpath)
		i++
	}
	return
}
