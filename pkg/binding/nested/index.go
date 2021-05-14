package nested

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	INDEX_NUMERIC = iota
	INDEX_EQUALS
	INDEX_UNKNOWN
)

type IndexValue struct {
	Type         int
	NumericValue int
	Path         []string
	CompareValue string
}

func getIndexValue(indexString string) IndexValue {
	index, err := strconv.Atoi(indexString)
	if err == nil {
		return IndexValue{
			Type:         INDEX_NUMERIC,
			NumericValue: index,
		}
	}

	regex := regexp.MustCompile(`^\?\((.*)\)$`)
	if regex.MatchString(indexString) {
		parts := regex.FindStringSubmatch(indexString)
		if len(parts) != 2 {
			return IndexValue{
				Type: INDEX_UNKNOWN,
			}
		}
		expression := parts[1]

		if strings.Contains(expression, "==") {
			keyValue := strings.Split(expression, "==")
			if len(keyValue) != 2 {
				return IndexValue{
					Type: INDEX_UNKNOWN,
				}
			}
			if !strings.HasPrefix(keyValue[0], "@") {
				return IndexValue{
					Type: INDEX_UNKNOWN,
				}
			}
			key := keyValue[0][1:]
			value := keyValue[1]
			return IndexValue{
				Type:         INDEX_EQUALS,
				Path:         fromDotted(key),
				CompareValue: value,
			}
		}
	}

	return IndexValue{
		Type: INDEX_UNKNOWN,
	}
}

func fromDotted(dottedPath string) []string {
	parts := strings.Split(dottedPath, ".")
	if parts[0] == "" {
		return parts[1:]
	}
	return parts
}
