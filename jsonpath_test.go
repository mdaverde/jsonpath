package jsonpath

import (
	"testing"
)

var data = map[string]interface{}{
	"user": map[string]interface{}{
		"firstname": "seth",
		"lastname":  "rogen",
	},
	"age": 35,
	"filmography": map[string]interface{}{
		"movies": []string{
			"This Is The End",
			"Superbad",
			"Neighbors",
		},
	},
}

func TestGet(t *testing.T) {
	var result interface{}
	var err error

	result, err = Get(data, "user.firstname")
	if err != nil {
		t.Errorf("failed to get user.firstname")
	}
	if result != "seth" {
		t.Errorf("wrong get value: wanted %v, got %v", "seth", result)
	}

	result, err = Get(data, "filmography.movies[1]")
	if err != nil {
		t.Errorf("failed to get filmography.movies[1]")
	}
	if result != "Superbad" {
		t.Errorf("wrong get value: wanted %v, got %v", "Superbad", result)
	}
}

func TestSet(t *testing.T) {

}
