package typesense

import (
	"fmt"
	"strings"
)

func interfaceArrayToStringArray(inputs []interface{}) []string {
	res := make([]string, len(inputs))
	for i, input := range inputs {
		res[i] = input.(string)
	}

	return res
}

func splitCollectionRelatedId(input string, resourceType string) (string, string, error) {
	eles := strings.Split(input, ".")
	if len(eles) != 2 {
		return "", "", fmt.Errorf("invalid format, format should be <collection>.<%s>", resourceType)
	}

	return eles[0], eles[1], nil
}
