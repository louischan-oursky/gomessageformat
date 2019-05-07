package messageformat

// doTest(t, Test{
// "You have {NUM_TASKS, plural, zero {no task} one {one task} two {two tasks} few{few tasks} many {many tasks} other {# tasks} =42 {the answer to the life, the universe and everything tasks}} remaining.",
// []Expectation{
// {map[string]interface{}{"NUM_TASKS": -1}, "You have -1 tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 0}, "You have no task remaining."},
// {map[string]interface{}{"NUM_TASKS": 1}, "You have one task remaining."},
// {map[string]interface{}{"NUM_TASKS": 2}, "You have two tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 3}, "You have few tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 6}, "You have many tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 15}, "You have 15 tasks remaining."},
// {map[string]interface{}{"NUM_TASKS": 42}, "You have the answer to the life, the universe and everything tasks remaining."},
// },
// })

func doParse(input string) (*MessageFormat, error) {
	o, err := New()
	if err != nil {
		return nil, err
	}
	mf, err := o.Parse(input, nil)
	if err != nil {
		return nil, err
	}
	return mf, nil
}
