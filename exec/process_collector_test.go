package exec

import (
	"testing"
)

func TestProcessCollect(t *testing.T) {

	m, err := ProcessCollect(82695)
	if err != nil {
		t.Error(err)
	}
	t.Log(m)
}
