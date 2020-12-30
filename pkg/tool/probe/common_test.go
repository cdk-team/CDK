package probe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTaskPortListByString(t *testing.T) {
	expect := make([]FromTo, 0)
	expect = append(expect, FromTo{
		Desc: "",
		From: 1,
		To:   1,
	})
	expect = append(expect, FromTo{
		Desc: "",
		From: 2,
		To:   2,
	})
	expect = append(expect, FromTo{
		Desc: "",
		From: 3,
		To:   5,
	})

	act, _ := GetTaskPortListByString("1,2,3-5")
	assert.Equal(t, expect, act, "That should be equal")
}
