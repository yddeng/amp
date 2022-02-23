package exec

import (
	"testing"
)

func TestProcessCollect(t *testing.T) {

	m, err := ProcessCollect(82696)
	if err != nil {
		t.Error(err)
	}
	t.Log(m)
}
