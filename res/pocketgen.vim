
syn keyword	pgDeclare  type enum
syn keyword	pgKeyword  func const
syn keyword	pgType     i8 i64 i32 i64
syn keyword	pgType     u8 u16 u32 u64
syn keyword	pgType     f32 f64
syn keyword	pgType     bool char string

syn region pgCommentLine start="#" end="$"

hi def link pgKeyword     Keyword
hi def link pgType        Type
hi def link pgDeclare     Structure
hi def link pgCommentLine Comment
