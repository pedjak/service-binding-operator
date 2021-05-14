package nested

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type splitPathResult struct {
	StartPath        []string
	EndPath          []string
	IsSlice          bool
	SliceIndexString string
}

func TestNestedFieldCopy(t *testing.T) {
	json := []byte(`{
		"apiVersion": "apps/v1",
		"kind": "StatefulSet",
		"metadata": {
			"name": "db1",
			"namespace": "prj1"
		},
		"spec": {
			"selector": {
				"matchLabels": {
					"app": "db1"
				}
			},
			"serviceName": "db1-svc",
			"template": {
				"metadata": {
					"labels": {
						"app": "db1"
					}
				},
				"spec": {
					"containers": [
						{
							"env": [
								{
									"name": "POSTGRESQL_USER",
									"value": "user1"
								},
								{
									"name": "POSTGRESQL_PASSWORD",
									"value": "k33p5ecret"
								},
								{
									"name": "POSTGRESQL_DATABASE",
									"value": "mydb"
								},
								{
									"name": "DUPLICATED",
									"value": "val1"
								},
								{
									"name": "DUPLICATED",
									"value": "val2"
								}
							],
							"image": "centos/postgresql-96-centos7",
							"name": "db1"
						}
					]
				}
			}
		}
	}`)
	tests := []struct {
		Path    []string
		Want    string
		WantOK  bool
		WantErr bool
	}{
		{
			Path:    []string{"spec", "serviceName"},
			Want:    "db1-svc",
			WantOK:  true,
			WantErr: false,
		},
		{
			Path:    []string{"spec", "template", "spec", "containers[0]", "name"},
			Want:    "db1",
			WantOK:  true,
			WantErr: false,
		},
		{
			Path:    []string{"spec", "template", "spec", "containers[0]", "env[2]", "value"},
			Want:    "mydb",
			WantOK:  true,
			WantErr: false,
		},
		{
			Path:    []string{"spec", "template", "spec", "containers[?(@.name==\"db1\")]", "env[?(@.name==\"POSTGRESQL_USER\")]", "value"},
			Want:    "user1",
			WantOK:  true,
			WantErr: false,
		},
		{
			Path:    []string{"spec", "template", "spec", "containers[?(@.name==\"db1\")]", "env[?(@.name==\"DUPLICATED\")]", "value"},
			Want:    "val1",
			WantOK:  true,
			WantErr: false,
		},
	}

	for _, test := range tests {
		var u unstructured.Unstructured
		err := u.UnmarshalJSON(json)
		if err != nil {
			t.Errorf("Error unmarshaling json input\n")
		}
		result, ok, err := NestedFieldCopy(u.Object, test.Path...)
		if ok != test.WantOK {
			t.Errorf("Expecting ok %v, got %v\n", test.WantOK, ok)
		}
		if (err != nil) != test.WantErr {
			t.Errorf("Expecting err %v, got %v\n", test.WantErr, err != nil)
		}
		if result != test.Want {
			t.Errorf("Expecting %s, got %s\n", test.Want, result)
		}
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		Path []string
		Want splitPathResult
	}{
		{
			Path: []string{"a", "b", "c"},
			Want: splitPathResult{
				StartPath: []string{"a", "b", "c"},
			},
		},
		{
			Path: []string{"a", "b[1]", "c"},
			Want: splitPathResult{
				StartPath:        []string{"a", "b"},
				EndPath:          []string{"c"},
				IsSlice:          true,
				SliceIndexString: "1",
			},
		},
		{
			Path: []string{"a", "b", "c[1]"},
			Want: splitPathResult{
				StartPath:        []string{"a", "b", "c"},
				IsSlice:          true,
				SliceIndexString: "1",
			},
		},
	}

	for _, test := range tests {
		start, end, is, index := splitPath(test.Path)
		if !stringSlicesEqual(start, test.Want.StartPath) || !stringSlicesEqual(end, test.Want.EndPath) || is != test.Want.IsSlice || index != test.Want.SliceIndexString {
			t.Errorf("Expecting %+v, got %v, %v, %v, %v", test.Want, start, end, is, index)
		}
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
