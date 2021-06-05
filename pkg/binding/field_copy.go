package binding

import (
	"bytes"
	"strings"

	"k8s.io/client-go/util/jsonpath"
)

func getValuesByJSONPath(obj map[string]interface{}, path_segments ...string) (string, bool, error) {
	path := "{." + strings.Join(path_segments, ".") + "}"
	j := jsonpath.New("")
	err := j.Parse(path)
	if err != nil {
		return "", false, err
	}
	buf := new(bytes.Buffer)
	err = j.Execute(buf, obj)
	if err != nil {
		return "", false, err
	}
	return buf.String(), true, nil
}
