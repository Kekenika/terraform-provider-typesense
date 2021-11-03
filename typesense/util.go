package typesense

func interfaceArrayToStringArray(inputs []interface{}) []string {
	res := make([]string, len(inputs))
	for i, input := range inputs {
		res[i] = input.(string)
	}

	return res
}
