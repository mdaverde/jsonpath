package jsonpath

import (
	"testing"
	"reflect"
	"encoding/json"
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
		t.Errorf("wrong get value, wanted %v, got %v", "seth", result)
	}

	result, err = Get(data, "filmography.movies[1]")
	if err != nil {
		t.Errorf("failed to get filmography.movies[1]")
	}
	if result != "Superbad" {
		t.Errorf("wrong get value, wanted %v, got %v", "Superbad", result)
	}

	result, err = Get(data, "age")
	if err != nil {
		t.Errorf("failed to get age: %v", err)
	}
	if result != 35 {
		t.Errorf("wrong get value, wanted: %v, got: %v", 35, result)
	}

	result, err = Get(data, "this.does.not[0].exist")
	if result != nil || err != DoesNotExistErr {
		t.Errorf("does not handle non-existant path correctly")
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

func TestJSON(t *testing.T) {
	test := `
{
	"pet": {
		"name": "baxter",
		"owner": {
      "name": "john doe",
      "contact": {
			  "phone": "859-289-9290"
      }
		},
		"type": "dog",
    "age": "4"
	},
	"tags": [
		12,
		true,
		{
			"hello": [
				"world"
			]
		}
	]
}
`
	var payload interface{}

	err := json.Unmarshal([]byte(test), &payload)
	if err != nil {
		t.Errorf("failed to parse: %v", err)
	}

	result, err := Get(payload, "tags[2].hello[0]")
	if result != "world" {
		t.Errorf("got wrong value from path, wanted: %v, got: %v", "world", result)
	}

	err = Set(&payload, "tags[2].hello[0]", "bobby")
	if err != nil {
		t.Errorf("failed to set: %v", err)
	}

	result, err = Get(payload, "tags[2].hello[0]")
	if result != "bobby" {
		t.Errorf("got wrong value after setting, wanted: %v, got %v", "bobby", result)
	}

	newContact := map[string]string{
		"phone": "555-555-5555",
		"email": "baxterowner@johndoe.com",
	}
	err = Set(&payload, "pet.owner.contact", newContact)
	if err != nil {
		t.Errorf("failed to set: %v", err)
	}

	contact, err := Get(&payload, "pet.owner.contact")
	if !reflect.DeepEqual(newContact, contact) {
		t.Errorf("contact set do not equal, wanted: %v, got %v", newContact, contact)
	}
}