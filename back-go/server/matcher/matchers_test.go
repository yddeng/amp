package matcher

import "testing"

func TestNewMatcher(t *testing.T) {
	var ms Matchers

	ms = append(ms, NewMatcher(MatchGreatEqual, "cpu", 50))
	ms = append(ms, NewMatcher(MatchGreatEqual, "mem", 20))

	t.Log(ms.String())
	t.Log(ms.Matches(map[string]float64{
		"cpu": 50,
		"mem": 50,
	}))
}
