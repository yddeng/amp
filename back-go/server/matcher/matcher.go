package matcher

import (
	"bytes"
	"fmt"
)

type MatchType int

const (
	MatchLess       MatchType = iota // 小于
	MatchLessEqual                   // 小于等于
	MatchEqual                       // 等于
	MatchNoEqual                     // 不等于
	MatchGreatEqual                  // 大于等于
	MatchGreat                       // 大于
)

var MatchType_name = map[MatchType]string{
	MatchLess:       "<",
	MatchLessEqual:  "<=",
	MatchEqual:      "==",
	MatchNoEqual:    "!=",
	MatchGreatEqual: ">=",
	MatchGreat:      ">",
}

var MatchType_value = map[string]MatchType{
	"<":  MatchLess,
	"<=": MatchLessEqual,
	"==": MatchEqual,
	"!=": MatchNoEqual,
	">=": MatchGreatEqual,
	">":  MatchGreat,
}

func (m MatchType) String() string {
	if str, ok := MatchType_name[m]; ok {
		return str
	}
	panic("unknown match type")
}

type Matcher struct {
	Type  MatchType `json:"type"`
	Name  string    `json:"name"`
	Value float64   `json:"value"`
}

func NewMatcher(t MatchType, n string, v float64) *Matcher {
	if _, ok := MatchType_name[t]; !ok {
		panic("unknown match type")
	}
	return &Matcher{
		Type:  t,
		Name:  n,
		Value: v,
	}
}

func (m *Matcher) String() string {
	return fmt.Sprintf(`%s%s%f`, m.Name, m.Type.String(), m.Value)
}

func (m *Matcher) Match(v float64) bool {
	switch m.Type {
	case MatchLess:
		return v < m.Value
	case MatchLessEqual:
		return v <= m.Value
	case MatchEqual:
		return v == m.Value
	case MatchNoEqual:
		return v != m.Value
	case MatchGreatEqual:
		return v >= m.Value
	case MatchGreat:
		return v > m.Value
	}
	return false
}

type Matchers []*Matcher

func (ms Matchers) Len() int      { return len(ms) }
func (ms Matchers) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }

func (ms Matchers) Less(i, j int) bool {
	if ms[i].Name > ms[j].Name {
		return false
	}
	if ms[i].Name < ms[j].Name {
		return true
	}
	if ms[i].Value > ms[j].Value {
		return false
	}
	if ms[i].Value < ms[j].Value {
		return true
	}
	return ms[i].Type < ms[j].Type
}

func (ms Matchers) Matches(set map[string]float64) bool {
	for _, m := range ms {
		if !m.Match(set[m.Name]) {
			return false
		}
	}
	return true
}

func (ms Matchers) String() string {
	var buf bytes.Buffer

	buf.WriteByte('{')
	for i, m := range ms {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(m.String())
	}
	buf.WriteByte('}')

	return buf.String()
}
