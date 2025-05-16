# writing-interpreter-in-go

Book: Writing an Interpreter in Go by Thorsten Ball

## Notes

`objectsystem` package has the latest and full impl of Monkey lang.

The book is not on latest Go version. I migrated and upgrade to version `go 1.24.0`

## Monkey Language

* Token, Lexer, AST, Parser, Interpreter and REPL
* types: int, string, bool, array and hash
* if/else
* inflix/postfix expressions
* functions
* literals
* closures
* built-in functions

## Diego Pacheco extra notes

* I added a tracing system for better debugging
* Migrated from Go versionm 1.16 to 1.24
* Added scripts for build and testing