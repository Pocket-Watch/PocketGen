package main

// Metagen format grammar (TODO):
//   newline        := // The U+000A character.
//   unicode_char   := // Any unicode character except newline.
//   unicode_letter := // Any unicode character categorized as "Letter".
//
//   letter     := unicode_letter | "_"
//   identifier := letter { letter | digit }
//   string     := '"' { unicode_char } '"'
//
//   metagen := typeDecl
//   		 | funcDecl
//   		 | varDecl
//
//   typeDecl := "type" identifier "{" { varDecl | funcDecl } "}"
//
//   funcDecl := "func" identifier "(" { varDecl comma } ")"
//
//   varDecl := identifier ( "[" "]" ) identifier

func main() {
	executeCLI()
}
