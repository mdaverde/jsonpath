package jsonpath

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var invalidObjError = errors.New("invalid object")
var pathDelimiter = "."

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

func getToken(obj interface{}, token string) (interface{}, error) {
	if reflect.TypeOf(obj) == nil {
		return nil, invalidObjError
	}

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Map:
		for _, kv := range reflect.ValueOf(obj).MapKeys() {
			if kv.String() == token {
				return reflect.ValueOf(obj).MapIndex(kv).Interface(), nil
			}
		}
		return nil, fmt.Errorf("%v not found in object", token)
	case reflect.Slice:
		idx, err := strconv.Atoi(token)
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
		return nil, fmt.Errorf("object is not a map or a slice: %v", reflect.TypeOf(obj).Kind())
	}
}

func getByTokens(data interface{}, tokens []string) (interface{}, error) {
	var err error

	child := data
	for _, token := range tokens {
		child, err = getToken(child, token)
		if err != nil {
			return nil, err
		}
	}

	if child != nil {
		return child, nil
	}

	return nil, errors.New("could not get value at path")
}

func Get(data interface{}, path string) (interface{}, error) {
	var err error
	tokens, err := tokenizePath(path)
	if err != nil {
		return nil, err
	}

	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	data = rv.Interface()

	return getByTokens(data, tokens)
}

func Set(data interface{}, path string, value interface{}) error {
	tokens, err := tokenizePath(path)
	if err != nil {
		return nil
	}
	head := tokens[:len(tokens)-1]
	last := tokens[len(tokens)-1]

	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	data = rv.Interface()

	rv = reflect.ValueOf(value)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	value = rv.Interface()

	child, err := getByTokens(data, head)

	switch reflect.TypeOf(child).Kind() {
	case reflect.Map:
		reflect.ValueOf(child).SetMapIndex(reflect.ValueOf(last),reflect.ValueOf(value))
		return nil
	case reflect.Slice:
		sliceValue := reflect.ValueOf(child)
		idx, err := strconv.Atoi(last)
		if err != nil {
			return err
		}
		sliceValue.Index(idx).Set(reflect.ValueOf(value))
		return nil
	}

	return errors.New("could not set value at path")
}
