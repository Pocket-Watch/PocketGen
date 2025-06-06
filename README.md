# PocketGen
Language code meta generator


## Converting tg representation to go
```tg
type Cat {
    const string name;
    u32 age;
    # This a comment
    func meow(string sound, u32 volume) string;
}
```

```bash
 tg test/cat.tg go --json
```

```go
package main

type Cat struct {
    Name string `json:"name"`
    Age  uint32 `json:"age"`
}

func (cat *Cat) meow(sound string, volume uint32) string {
    panic("TODO: Unimplemented method")
}
```
## Supported languages:
- Go
- Java
- JavaScript
- Kotlin
- Rust

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
Source the syntax file with the editor command: 
`:so res/pocketgen.vim`

You can also source it automatically by adding it to your `vim.rc`/`init.lua`.
