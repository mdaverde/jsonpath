package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var invalidObjError = errors.New("invalid object")
var pathDelimiter string = "."

func tokenizePath(path string) ([]string, error) {
	var tokens []string
	for _, stem := range strings.Split(path, pathDelimiter) {
		if !strings.Contains(stem, "[") {
			tokens = append(tokens, stem)
			continue
		}
		firstBracketIndex := strings.Index(stem, "[")
		lastBracketIndex := strings.LastIndex(stem, "]")
		if lastBracketIndex < 0 {
			return nil, fmt.Errorf("invalid path: %v", path)
		}
		tokens = append(tokens, stem[0:firstBracketIndex])
		innerText := stem[firstBracketIndex+1 : lastBracketIndex]
		tokens = append(tokens, innerText)
	}
	return tokens, nil
}

func getKey(obj interface{}, key string) (interface{}, error) {
	if reflect.TypeOf(obj) == nil {
		return nil, invalidObjError
	}

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Map:
		for _, kv := range reflect.ValueOf(obj).MapKeys() {
			if kv.String() == key {
				return reflect.ValueOf(obj).MapIndex(kv).Interface(), nil
			}
		}
		return nil, fmt.Errorf("%v not found in object", key)
	case reflect.Slice:
		idx, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		length := reflect.ValueOf(obj).Len()
		if idx > -1 {
			if idx >= length {
				return nil, fmt.Errorf("index out of range: %v len: %v", idx, length)
			}
			return reflect.ValueOf(obj).Index(idx).Interface(), nil
		}
		return nil, fmt.Errorf("%v not found in object", idx)
	default:
		fmt.Println(reflect.TypeOf(obj).Kind())
		return nil, fmt.Errorf("object is not a map or a slice")
	}
}

func Get(data interface{}, path string) (interface{}, error) {
	var err error
	tokens, err := tokenizePath(path)

	child := data
	for _, token := range tokens {
		child, err = getKey(child, token)
		if err != nil {
			return nil, err
		}
	}

	if child != nil {
		return child, nil
	}

	return nil, errors.New("could not get value at path")
}

func Set(data interface{}, path string, value interface{}) error {
	return errors.New("could not set value at path")
}
