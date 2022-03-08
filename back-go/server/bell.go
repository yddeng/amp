package server

import "amp/back-go/server/matcher"

type Condition struct {
	Matchers  matcher.Matchers `json:"ms"`
	Duration  int64            `json:"duration"` // 持续时间
	BeginFunc func()           `json:"_"`
	EndFunc   func()           `json:"_"`
	startTime int64
	beginCall bool
}

func NewCondition(dur int64, ms []*matcher.Matcher) *Condition {
	return &Condition{
		Matchers: ms,
		Duration: dur,
	}
}

func (c *Condition) Match(set map[string]float64) {
	ok := c.Matchers.Matches(set)
	if ok {
		nowUnix := NowUnix()
		if c.startTime == 0 {
			c.startTime = nowUnix
		} else if nowUnix >= c.startTime+c.Duration && !c.beginCall {
			if c.BeginFunc != nil {
				c.BeginFunc()
			}
			c.beginCall = true
		}
	} else {
		if c.beginCall {
			if c.EndFunc != nil {
				c.EndFunc()
			}
		}
		c.beginCall = false
		c.startTime = 0
	}
}
