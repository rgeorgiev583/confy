// api_test.go
package wrapper

import (
	"encoding/json"
	"testing"
)

func TestWebList(t *testing.T) {
	config, err := New("/", "", 0)
	if err != nil {
		t.Error(err)
	}

	obj, err := WebList(config, "/", nil)
	if err != nil {
		t.Error(err)
	}

	res, _ := json.Marshal(obj)
	if string(res) != "[\"/augeas\",\"/files\"]" {
		t.FailNow()
	}
}

func TestWebMatch(t *testing.T) {
	config, err := New("/", "", 0)
	if err != nil {
		t.Error(err)
	}

	obj, err := WebMatch(config, "/*", nil)
	if err != nil {
		t.Error(err)
	}

	res, _ := json.Marshal(obj)
	if string(res) != "[\"/augeas\",\"/files\"]" {
		t.FailNow()
	}
}

func TestWebGet(t *testing.T) {
	config, err := New("/", "", 0)
	if err != nil {
		t.Error(err)
	}

	obj, err := WebGet(config, "/files/etc/passwd/sid/name", nil)
	if err != nil {
		t.Error(err)
	}

	res, _ := json.Marshal(obj)
	if string(res) != "{\"value\":\"Radoslav Georgiev,,,,\"}" {
		t.FailNow()
	}
}

func TestWebLabel(t *testing.T) {
	config, err := New("/", "", 0)
	if err != nil {
		t.Error(err)
	}

	obj, err := WebGetLabel(config, "/files/etc/passwd/sid/name", nil)
	if err != nil {
		t.Error(err)
	}

	res, _ := json.Marshal(obj)
	if string(res) != "\"name\"" {
		t.FailNow()
	}
}
