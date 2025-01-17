package repl

import (
	"bufio"
	"dara/evaluator"
	"dara/lexer"
	"dara/parser"
	"fmt"
	"io"
	"log"
)

const PROMPT = "→ "

func Start(in io.Reader, out io.Writer) {
	var (
		scanner = bufio.NewScanner(in)
		env     = evaluator.NewEnvironment()
	)

	for {
		fmt.Print(PROMPT)
		if scanned := scanner.Scan(); !scanned {
			return
		}

		var (
			line = scanner.Text()
			l    = lexer.New(line)
			p    = parser.New(l)
		)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if _, err := io.WriteString(out, evaluated.Inspect()); err != nil {
				log.Fatalln(err)
			}

			if _, err := io.WriteString(out, "\n"); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	if _, err := io.WriteString(out, "  parser errors:\n"); err != nil {
		log.Fatalln(err)
	}
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			log.Fatalln(err)
		}
	}
}
