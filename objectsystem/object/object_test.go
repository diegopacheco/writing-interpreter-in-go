package object

import (
	"testing"
)

type StringHashKey string

func (s StringHashKey) HashKey() string {
	return string(s)
}

func TestStringHashKey(t *testing.T) {
	hello1 := StringHashKey("hello")
	hello2 := StringHashKey("hello")
	diff1 := StringHashKey("world")
	diff2 := StringHashKey("world")

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("Expected hash keys to be equal, got %v and %v", hello1.HashKey(), hello2.HashKey())
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("Expected hash keys to be equal, got %v and %v", diff1.HashKey(), diff2.HashKey())
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("Expected hash keys to be different, got %v and %v", hello1.HashKey(), diff1.HashKey())
	}

}
