package jsonpath

import (
	"testing"
	"reflect"
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
	result, err := Get(data, "user.firstname")
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
	err := Set(&data, "user.firstname", "chris")
	if err != nil {
		t.Errorf("failed to set user.firstname: %v", err)
	}

	firstname := reflect.ValueOf(data["user"]).MapIndex(reflect.ValueOf("firstname")).Interface()
	if firstname != "chris" {
		t.Errorf("set user.firstname to wrong value, wanted: %v, got: %v", "chris", firstname)
	}

	err = Set(&data, "filmography.movies[2]", "The Disaster Artist")
	if err != nil {
		t.Errorf("failed to set filmography.movies[2]: %v", err)
	}

	secondMovie := reflect.ValueOf(data["filmography"]).MapIndex(reflect.ValueOf("movies")).Elem().Index(2).Interface()
	if secondMovie != "The Disaster Artist" {
		t.Errorf("set filmography.movies[2] to wrong value, wanted: %v, got %v", "The Disaster Artist", secondMovie)
	}

	newUser := map[string]interface{}{
		"firstname": "james",
		"lastname": "franco",
	}

	err = Set(&data, "user", &newUser)
	if err != nil {
		t.Errorf("failed to set user: %v", err)
	}

	user := data["user"]

	if !reflect.DeepEqual(newUser, user) {
		t.Errorf("set user is not equal, wanted: %v, got %v", newUser, user)
	}
}