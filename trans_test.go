package main

import (
	"github.com/go-playground/assert/v2"
	"github.com/sonhineboy/gsadminValidator/ginValidator"
	"testing"
)

func TestTrans(t *testing.T) {
	assert.Equal(t, &ginValidator.Trans{}, ginValidator.NewTrans())
	ginValidator.NewDefaultTrans()

}

func re() int {

	return 3
}
