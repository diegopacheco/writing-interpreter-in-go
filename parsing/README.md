### Notes

* Parsing
* Output of a parser (AST) is abstract
* Parser Generators
	* YACC
	* BISON
	* AntLR
* Context Free Grammar (CFG) - Input (with a set of rules on how to form correct sentences on the language)
	* Backus Naur Form (BNF)
	* Extended-Backus Naur Form (EBNF)
* 2 Main strategies when writing a parser form a programming language
	* Top down parsing
	* Botton up parsing 
* Recursive Decent Parser (Top down operator precedence parser) -> Pratt Parser
* Top Down Operator Precedence
    * https://crockford.com/javascript/tdop/tdop.html

### Run

```bash
./run.sh
```

```
❯ ./run.sh
Hello diego! This is the Monkey programming language!
Feel free to type in commands
>> let x = 10;
Current char: 'l', Position: 0
{Type:LET Literal:let}
Current char: 'x', Position: 4
{Type:IDENT Literal:x}
Current char: '=', Position: 6
{Type:= Literal:=}
Current char: '1', Position: 8
{Type:INT Literal:10}
Current char: ';', Position: 10
{Type:; Literal:;}
Current char: '\x00', Position: 11
>> ^Csignal: interrupt
```

### Test

```bash
./test.sh
```