#!/bin/bash

go test ./lexer
go test ./ast
go test ./parser
go test ./evaluator
go test ./object