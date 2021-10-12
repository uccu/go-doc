package doc_test

import (
	"testing"

	"github.com/uccu/go-doc"
)

func TestVV(t *testing.T) {

	err := doc.Export("doc", "github.com/uccu/go-doc/test")
	if err != nil {
		t.Error(err)
	}
}
