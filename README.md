# PocketGen
Language code meta generator

## Language constructs
### Keywords
```
type
const
func
enum
```

### Primitive types
```
u8
u16
u32 
u64
i8
i16
i32
i64
f64
string
bool
char
```

## Adding custom syntax highlighting:

### Jetbrains IDEs
1. In settings navigate **Editor** | **File Types** | **Recognized file types** | **Add**
2. Write a name, description (shown as label)
3. Put `#` in **Line comment**
4. Mark selected:
 - support paired parens
 - support paired brackets
5. In **Keywords** section:
 - at number 1 copy-paste keywords
 - at number 4 copy-paste primitive types
6. In **File name patterns** associate it with `*.tg`

### Neovim / Vim
Source the syntax file editor command: 
`:so res/pocketgen.vim`

You cant also add it to your `vim.rc`/`init.lua` to source it automatically.
