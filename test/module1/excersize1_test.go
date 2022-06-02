package module1

import (
	"dongzw/dongzwhom/module1"
	"testing"
)

func TestExcersizeItemOne(t *testing.T) {
	r := module1.ItemOne()
	t.Log(r)
}

func TestExcersizeItemTwo(t *testing.T) {
	module1.ItemTwo()
}
