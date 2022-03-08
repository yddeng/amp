package server

import (
	"amp/back-go/server/matcher"
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

func TestNewCondition(t *testing.T) {
	c := NewCondition(5, []*matcher.Matcher{
		matcher.NewMatcher(matcher.MatchGreatEqual, "cpu", 50)})

	c.BeginFunc = func() {
		t.Log("begin")
	}
	c.EndFunc = func() {
		t.Log("end")
	}

	d, _ := json.Marshal(c)
	t.Log(string(d))

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		v := rand.Int31n(50) + 40
		t.Log(v)
		c.Match(map[string]float64{"cpu": float64(v)})
	}

}
