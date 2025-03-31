package main

// Extension? .tg?
// TODO("UNIMPLEMENTED ASSERTION for every lang?")

type Type struct {
	fields []Member
	// functions potentially?
}

type Member struct {
	modifiers []string
	name      string
	isList    bool
	typeName  string
}

// evaluated = UNSTARTED = 0, IN-PROGRESS=1, FINISHED=2
