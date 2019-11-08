package messageformat

import (
	"testing"
)

func TestToString(t *testing.T) {
	doTostringTest(t, toStringTest{
		`{   NAME   }"`,
		`{NAME}"`,
	})

	doTostringTest(t, toStringTest{
		`{   position  } liked this."`,
		`{position} liked this."`,
	})

	doTostringTest(t, toStringTest{
		`{GENDER, select, male{He} female {She} other{They}} liked this."`,
		`{GENDER, select, male{He}female{She}other{They}} liked this."`,
	})

	doTostringTest(t, toStringTest{
		"{NUM_TASKS, plural, one {a } =1 { b} other {  c  }}",
		"{NUM_TASKS, plural, =1{ b}one{a }other{  c  }}",
	})

	doTostringTest(t, toStringTest{
		"This is a { A , plural , =2 {benchmark } other{} }",
		"This is a {A, plural, =2{benchmark }other{}}",
	})

	doTostringTest(t, toStringTest{
		"This is a {A, plural, offset:1     one{a}=2{b}=1{c}other{}}",
		"This is a {A, plural, offset:1 =1{c}=2{b}one{a}other{}}",
	})
}
