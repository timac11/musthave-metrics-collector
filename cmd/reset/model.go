package main

import (
	"go/ast"
)

type packageInfo struct {
	Name string
	Dir  string
	Ast  []*ast.File
}

type genInfo struct {
	Name  string
	Funcs []string
}

type structInfo struct {
	Name   string
	Fields []string
}

type structAstInfo struct {
	Name string
	Typ  *ast.StructType
}
