package messageformat

import (
	"testing"
)

func TestToString(t *testing.T) {
	doTostringTest(t, toStringTest{
		`literal string !`,
		`literal string !`,
		false,
	})

	doTostringTest(t, toStringTest{
		`{   NAME   }"`,
		`{NAME}"`,
		false,
	})

	doTostringTest(t, toStringTest{
		`{   position  } liked this."`,
		`{position} liked this."`,
		false,
	})

	doTostringTest(t, toStringTest{
		`{GENDER, select, male{He} female {She} other{They}} liked this."`,
		`{GENDER, select, male{He}female{She}other{They}} liked this."`,
		true,
	})

	doTostringTest(t, toStringTest{
		"{NUM_TASKS, plural, one {a } =1 { b} other {  c  }}",
		"{NUM_TASKS, plural, =1{ b}one{a }other{  c  }}",
		true,
	})

	doTostringTest(t, toStringTest{
		"This is a { A , selectordinal , one {benchmark } other{} }",
		"This is a {A, selectordinal, one{benchmark }other{}}",
		true,
	})

	doTostringTest(t, toStringTest{
		"This is a {A, plural, offset:1     one{a}=2{b}=1{c} two{d} other{#}}",
		"This is a {A, plural, offset:1 =1{c}=2{b}one{a}other{#}}",
		true,
	})
}
