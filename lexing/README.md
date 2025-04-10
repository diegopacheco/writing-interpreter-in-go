### Notes

* Go lang version: go 1.24.0
* Lexing
* Source code -> Tokens -> AST
* Lexer also called tokenizer or scanner

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