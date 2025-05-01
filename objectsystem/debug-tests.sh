#!/bin/bash

DEBUG=true go test ./lexer
DEBUG=true go test ./ast
DEBUG=true go test ./parser
DEBUG=true go test ./evaluator