package jsonpath

import (
	"reflect"
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
		Path    string
		Want    interface{}
		WantErr bool
	}{
		{
			Path:    ".spec.serviceName",
			Want:    "db1-svc",
		},
		{
			Path:    ".spec.template.spec.containers[0].name",
			Want:    "db1",
		},
		{
			Path:    ".spec.template.spec.containers[0].env[2].value",
			Want:    "mydb",
		},
		{
			Path:    ".spec.template.spec.containers[?(@.name==\"db1\")].env[?(@.name==\"POSTGRESQL_USER\")].value",
			Want:    "user1",
		},
		{
			Path:    ".spec.template.metadata.labels",
			Want:    []interface{}{map[string]interface{} {
			"app": "db1",

			}},
		},

	}

	var u unstructured.Unstructured
	err := u.UnmarshalJSON(json)
	if err != nil {
		t.Errorf("Error unmarshaling json input\n")
	}

	for _, test := range tests {
		t.Run(test.Path, func(t *testing.T) {
			result, err := getValuesByJSONPath(u.Object, "{"+test.Path+"}")
			if (err != nil) != test.WantErr {
				t.Errorf("Expecting err %v, got %v\n", test.WantErr, err)
			}
			if !reflect.DeepEqual(result[0],test.Want) {
				t.Errorf("Expecting %s, got %s\n", test.Want, result)
			}
		})
	}
}
