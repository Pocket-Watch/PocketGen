package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const EXTENSION = ".tg"

func executeCLI() {
	exec, err := os.Executable()
	if err != nil {
		exec = ""
	}
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage:")
		fmt.Printf("  %v <file path/directory> <language>\n", exec)
		fmt.Println()
		return
	}
	args = args[1:]

	path := args[0]

	info, err := getPathInfo(path)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(1)
	}

	lang := args[1]
	language := languageIdentifierToLanguage(lang)
	if language == NONE {
		fmt.Println("Unrecognized, unsupported or misspelled language identifier.")
		os.Exit(1)
	}

	var files []string
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Failed to open directory.")
			os.Exit(1)
		}
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), EXTENSION) {
				tgFile := filepath.Join(path, entry.Name())
				files = append(files, tgFile)
			}
		}
	} else {
		if strings.HasSuffix(path, EXTENSION) {
			files = append(files, path)
		}
	}

	if len(files) == 0 {
		fmt.Printf("Nothing to do. Ensure your files end with %v\n", EXTENSION)
		os.Exit(1)
	}
	start := time.Now()
	fmt.Printf("Processing %v files\n", len(files))
	for _, file := range files {
		fmt.Printf("  %v\n", file)
		parser, success := CreateParser(file)
		if !success {
			os.Exit(1)
		}

		parseResult := ParseFile(&parser)
		if !parseResult.success {
			fmt.Println(parseResult.message)
			os.Exit(1)
		}

		checkResult := TypecheckFile(&parser)
		if !checkResult.success {
			fmt.Println(checkResult.message)
			os.Exit(1)
		}

		// cat.tg -> cat.js
		newExtension := languageToExtension(language)
		generatedFile := changeExtension(file, newExtension)
		writer := openWriter(generatedFile)
		if writer == nil {
			os.Exit(1)
		}

		switch language {
		case JAVASCRIPT:
			js := JavascriptGenerator{defaultOptions()}
			js.generate(parser.structs, writer)
		case GO:
			goGen := GoGenerator{defaultOptions()}
			err = goGen.generate(parser.structs, writer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case JAVA:
			java := JavaGenerator{defaultOptions()}
			err = java.generate(parser.structs, writer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case KOTLIN:
			kotlin := KotlinGenerator{defaultOptions()}
			err = kotlin.generate(parser.structs, writer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case RUST:
			rust := RustGenerator{defaultOptions()}
			err = rust.generate(parser.structs, writer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unsupported language (coming soon).")
			os.Exit(1)
		}
	}
	end := time.Now()
	timeElapsed := end.Sub(start)
	fmt.Printf("Time elapsed processing: %v\n", timeElapsed)
}

func changeExtension(file string, newExtension string) string {
	oldExtension := filepath.Ext(file)
	if oldExtension == "" {
		return file + newExtension
	}
	return strings.Replace(file, oldExtension, newExtension, 1)
}

func getPathInfo(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return info, fmt.Errorf("the file/directory at '%v' does not exist.'", path)
		} else {
			errorMessage := err.(*fs.PathError).Err
			return info, fmt.Errorf("failed to open file/directory at '%v' because %v.'", path, errorMessage)
		}
	}
	return info, nil
}

// Language enum
type Language = int

const (
	NONE Language = iota
	GO
	JAVASCRIPT
	TYPESCRIPT
	KOTLIN
	JAVA
	C_SHARP
	RUST
)

func languageIdentifierToLanguage(lang string) Language {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return GO
	case "js", "javascript":
		return JAVASCRIPT
	case "ts", "typescript":
		return TYPESCRIPT
	case "java":
		return JAVA
	case "kt", "kotlin":
		return KOTLIN
	case "rs", "rust":
		return RUST
	case "cs", "csharp":
		return C_SHARP
	}
	return NONE
}

func languageToExtension(language Language) string {
	switch language {
	case GO:
		return ".go"
	case JAVASCRIPT:
		return ".js"
	case TYPESCRIPT:
		return ".ts"
	case JAVA:
		return ".java"
	case KOTLIN:
		return ".kt"
	case RUST:
		return ".rs"
	case C_SHARP:
		return ".cs"
	default:
		panic("language has an unmapped extension?")
	}
}
