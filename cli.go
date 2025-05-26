package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const EXTENSION = ".tg"

func executeCLI() {
	args := os.Args[1:]
	if len(args) < 2 {
		printHelp()
		return
	}
	path := args[0]

	info, err := getPathInfo(path)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(1)
	}

	lang := args[1]
	language := languageIdentifierToLanguage(lang)
	if language == NONE {
		fmt.Printf("ERROR: Unrecognized, unsupported or misspelled language identifier: %v\n", lang)
		os.Exit(1)
	}

	options := parseArguments(args[2:])

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

		codeBuffer := bytes.Buffer{}

		switch language {
		case JAVASCRIPT:
			js := JavascriptGenerator{options}
			js.generate(parser.structs, &codeBuffer)
		case GO:
			goGen := GoGenerator{options}
			err = goGen.generate(parser.structs, &codeBuffer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case JAVA:
			java := JavaGenerator{options}
			err = java.generate(parser.structs, &codeBuffer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case KOTLIN:
			kotlin := KotlinGenerator{options}
			err = kotlin.generate(parser.structs, &codeBuffer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		case RUST:
			rust := RustGenerator{options}
			err = rust.generate(parser.structs, &codeBuffer, parser.filepath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Println("Unsupported language (coming soon).")
			os.Exit(1)
		}

		// cat.tg -> cat.js
		newExtension := languageToExtension(language)
		generatedFile := changeExtension(file, newExtension)
		writer := openWriter(generatedFile)
		if writer == nil {
			os.Exit(1)
		}
		_, writeErr := writer.Write(codeBuffer.Bytes())
		if writeErr != nil {
			fmt.Println("ERROR writing contents to file:", writeErr)
			os.Exit(1)
		}
		writer.Flush()
	}
	end := time.Now()
	timeElapsed := end.Sub(start)
	fmt.Printf("Time elapsed processing: %v\n", timeElapsed)
}

// Method parseArguments parses arguments starting from index 0, returns generator options.
// On error exits with code 1. If help is passed as argument it's displayed and the program exits.
func parseArguments(args []string) GeneratorOptions {
	options := defaultOptions()
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h", "--help":
			printHelp()
			os.Exit(0)
		case "--json":
			options.jsonAnnotations = true
		case "--indent":
			if i+1 >= len(args) {
				fmt.Println("ERROR: No argument passed for indentation")
				os.Exit(1)
			}
			indent, err := strconv.Atoi(args[i+1])
			if err != nil || indent < 1 {
				fmt.Printf("ERROR: Invalid indentation: %v\n", args[i+1])
				os.Exit(1)
			}
			options.indent = indent
			i++
		default:
			fmt.Println("WARN: Unknown option", args[i])
		}
	}

	return options
}

func printHelp() {
	exec, err := os.Executable()
	if err == nil {
		exec = filepath.Base(exec)
	} else {
		exec = ""
	}
	fmt.Println("Usage:")
	fmt.Printf("  %v <file path/directory> <language> [options...]\n", exec)
	fmt.Println("Options:")
	fmt.Println("    --json                Generate JSON-annotations")
	fmt.Println("    --indent [number]     Specify code indentation level")
	fmt.Println("    -h, --help            Display this help message")
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
