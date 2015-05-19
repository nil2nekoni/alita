package alita

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

type Justify int

const (
	JustLeft Justify = iota
	JustRight
	JustCenter
)

func (j Justify) Just(width int, s string) string {
	w := runewidth.StringWidth(s)
	switch j {
	case JustLeft:
		return s + strings.Repeat(" ", width-w)
	case JustRight:
		return strings.Repeat(" ", width-w) + s
	case JustCenter:
		n := width - w
		l, r := n/2, n/2
		if n%2 != 0 {
			r += 1
		}
		return strings.Repeat(" ", l) + s + strings.Repeat(" ", r)
	}
	return s + strings.Repeat(" ", width-w)
}

var JUSTFIES_SEQUENCE = regexp.MustCompile("^[lrc]+$")

type Padding struct {
	justfies []Justify
	width    []int
}

func NewPadding() *Padding {
	return &Padding{}
}

func (p *Padding) SetJustfies(a []Justify) {
	p.justfies = a
}

func (m *Padding) String() string {
	return fmt.Sprint(*m)
}

func (p *Padding) Set(format string) error {
	if !JUSTFIES_SEQUENCE.MatchString(format) {
		return fmt.Errorf("padding: invalid format: %s", format)
	}
	a := make([]Justify, 0)
	for _, c := range format {
		switch c {
		case 'l':
			a = append(a, JustLeft)
		case 'r':
			a = append(a, JustRight)
		case 'c':
			a = append(a, JustCenter)
		}
	}
	p.SetJustfies(a)
	return nil
}

func (p *Padding) UpdateWidth(a []string) {
	if len(a) == 1 {
		return
	}
	for i, s := range a {
		w := runewidth.StringWidth(s)
		switch {
		case i == len(p.width):
			p.width = append(p.width, w)
		case w > p.width[i]:
			p.width[i] = w
		}
	}
}

func (p *Padding) justKind(i int) Justify {
	switch len(p.justfies) {
	case 0:
		return JustLeft
	case 1:
		return p.justfies[0]
	}
	j := (i-1)%(len(p.justfies)-1) + 1
	return p.justfies[j]
}

func (p *Padding) Format(a []string) []string {
	if len(a) == 1 {
		return a
	}
	for i := 0; i < len(a); i++ {
		a[i] = p.justKind(i).Just(p.width[i], a[i])
	}
	return a
}
