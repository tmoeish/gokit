package jsonx_test

import (
	"testing"

	"github.com/tmoeish/gokit/jsonx"
)

func TestMarshalUnmarshal(t *testing.T) {
	type S struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	v := S{Name: "alice", Age: 30}
	s := jsonx.Marshal(v)
	if s == "" {
		t.Fatal("Marshal returned empty")
	}
	var got S
	if err := jsonx.UnmarshalStr(s, &got); err != nil || got.Name != "alice" {
		t.Fatalf("Unmarshal: %v %v", got, err)
	}
}

func TestMarshalPretty(t *testing.T) {
	s := jsonx.MarshalPretty(map[string]int{"a": 1})
	if len(s) == 0 {
		t.Fatal("MarshalPretty empty")
	}
}

func TestClone(t *testing.T) {
	type S struct{ X int }
	orig := S{X: 42}
	clone, err := jsonx.Clone(orig)
	if err != nil || clone.X != 42 {
		t.Fatalf("Clone: %v %v", clone, err)
	}
	clone.X = 99
	if orig.X != 42 {
		t.Fatal("Clone should be independent")
	}
}

func TestYAML(t *testing.T) {
	type S struct {
		Name string `yaml:"name"`
	}
	v := S{Name: "gokit"}
	y := jsonx.MarshalYAML(v)
	if y == "" {
		t.Fatal("MarshalYAML empty")
	}
	var got S
	if err := jsonx.UnmarshalYAMLStr(y, &got); err != nil || got.Name != "gokit" {
		t.Fatalf("UnmarshalYAML: %v %v", got, err)
	}
}

func TestJSONToYAML(t *testing.T) {
	y, err := jsonx.JSONToYAML(`{"name":"alice","age":30}`)
	if err != nil || y == "" {
		t.Fatalf("JSONToYAML: %q %v", y, err)
	}
}

func TestYAMLToJSON(t *testing.T) {
	j, err := jsonx.YAMLToJSON("name: alice\nage: 30\n")
	if err != nil || j == "" {
		t.Fatalf("YAMLToJSON: %q %v", j, err)
	}
}

func TestGet(t *testing.T) {
	data := []byte(`{"user":{"name":"alice","age":30}}`)
	got := jsonx.Get(data, "user", "name")
	if got != "alice" {
		t.Fatalf("Get: %v", got)
	}
}
