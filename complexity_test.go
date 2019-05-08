package messageformat

import (
	"testing"
)

func doTestComplexity(t *testing.T, input string, expected Complexity) {
	o, err := New()
	if err != nil {
		t.Fatalf("`%s` threw <%s>", input, err)
	}

	mf, err := o.Parse(input, nil)
	if err != nil {
		t.Fatalf("parse `%s` threw <%s>", input, err)
	}

	if c := mf.Root().Complexity(); c != expected {
		t.Fatalf("Expecting %s, but got %s: %s", expected, c, input)
	}
}

func TestComplexity(t *testing.T) {
	doTestComplexity(t,
		"{A,select,other{#, and {B,select,other{#}}}}!",
		Complex)
	doTestComplexity(t,
		"liked {NAME} this.",
		Complex)
	doTestComplexity(t,
		"liked this.",
		SingleLiteral)
	doTestComplexity(t,
		"{GENDER, select, male{He} female {She} other{They}} liked this.",
		Complex)
	doTestComplexity(t,
		"{GENDER, select, male{He} female {She} other{They}}",
		SingleRoad)
	doTestComplexity(t,
		"{A,select,other{{B,select,one{1}other{more}}}}",
		SingleRoad)
	doTestComplexity(t,
		"{A,select,one{a}other{{B,plural,one{1}other{more}}}}",
		SingleRoad)
	doTestComplexity(t,
		"{A,select,one{#a}other{{B,plural,one{1}other{more}}}}",
		Complex)
	doTestComplexity(t,
		"{A,select,one{a}other{{B,plural,one{1}other{# tasks}}}}",
		Complex)
}
