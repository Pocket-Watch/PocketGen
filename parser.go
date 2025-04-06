package main

// Extension? .tg?
// TODO("UNIMPLEMENTED ASSERTION for every lang?")

// type Type struct {
// 	fields []Member
// 	// functions potentially?
// }
//
// type Member struct {
// 	modifiers []string
// 	name      string
// 	isList    bool
// 	typeName  string
// }

// evaluated = UNSTARTED = 0, IN-PROGRESS=1, FINISHED=2

type Parser struct {
	lexer     Lexer
	tokenNow  Token
	tokenNext Token
	structs   []TypeDecl
}

type TypeDecl struct {
	name    string
	fields  []Field
	methods []FuncDecl
}

type FuncDecl struct {
	name       string
	fields     []Field
	returnType string
}

type FieldModifier = uint32

const (
	FIELD_CONST FieldModifier = (1 << iota)
	FIELD_ARRAY
	FIELD_NULLABLE
	FIELD_PRIMITIVE
)

type Field struct {
	varName   string
	typeName  string
	modifiers []FieldModifier
}

func PlaceholderParse(parser *Parser) {

}
