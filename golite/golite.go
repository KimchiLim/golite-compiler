package main

import ("golite/arm" ; "golite/lexer" ; "golite/parser" ; "golite/sa" ; "flag" ; "fmt" ; "os" ; "strings")

func main() {
	lexerFlag := flag.Bool("l", false, "")
	parserFlag := flag.Bool("ast", false, "")
	llvmIRFlag := flag.String("llvm-ir", "", "")
	targetFlag := flag.String("target", "x86_64-linux-gnu", "")
	armFlag := flag.Bool("S", false, "")

	flag.Parse()
	arguments := flag.Args()

	if len(arguments) == 0 {
		fmt.Println("Error: No input file provided.\n")
		os.Exit(1)
	}

	inputSourcePath := arguments[0]

	if _, err := os.Stat(inputSourcePath); os.IsNotExist(err) {
		fmt.Printf("Error: File %s does not exist.\n", inputSourcePath)
		os.Exit(1)
	}

	lexer := lexer.NewLexer(inputSourcePath)
	if *lexerFlag { lexer.PrintTokens() ; return }

	parser := parser.NewParser(lexer)
	program := parser.Parse()

    if *parserFlag { fmt.Println(program) }
	if program == nil { os.Exit(1) }

	tables := sa.Execute(program)
	if tables == nil { os.Exit(1) }

	source := strings.TrimSuffix(inputSourcePath, ".golite")
	target := *targetFlag
	llvmProgram := program.TranslateToLLVMStack(source, target, tables)
	llvmString := llvmProgram.LLVMString(*llvmIRFlag == "stack")

	if *llvmIRFlag != "" && *llvmIRFlag != "stack" && *llvmIRFlag != "reg" {
		fmt.Println("Error: Invalid value for -llvm-ir. Expected 'stack' or 'reg'.")
		os.Exit(1)
	}

	if *llvmIRFlag == "stack" || *llvmIRFlag == "reg" {
		llvmFilename := source + ".ll"
		llvmFile, _ := os.Create(llvmFilename)
		llvmFile.WriteString(llvmString)
		llvmFile.Close()
		return
	}

	armProgram := arm.NewARMProgram(llvmProgram)

	if *armFlag {
		armFilename := source + ".s"
		armFile, _ := os.Create(armFilename)
		armFile.WriteString(armProgram.String())
		armFile.Close()
	}
}
