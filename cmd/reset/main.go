package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pkgs, err := parsePackages(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		structs := findResettableStructs(pkg)

		if len(structs) == 0 {
			continue
		}

		var buf bytes.Buffer

		err = fileTemplate.Execute(&buf, genInfo{
			Name:  pkg.Name,
			Funcs: generateRestFuncs(structs),
		})
		if err != nil {
			log.Fatalf("template execution error for package %s: %v", pkg.Name, err)
		}

		outPath := filepath.Join(pkg.Dir, "reset.gen.go")
		if err = os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func generateRestFuncs(structs []*structAstInfo) []string {
	var funcs []string
	for _, st := range structs {
		funcs = append(funcs, generateFunc(st))
	}
	return funcs
}

func generateFunc(st *structAstInfo) string {
	var buf bytes.Buffer
	err := methodTemplate.Execute(&buf, structInfo{
		Name:   st.Name,
		Fields: genResetFieldLines(st),
	})
	if err != nil {
		log.Fatalf("failed to render method for %s: %v", st.Name, err)
	}
	return buf.String()
}

func genResetFieldLines(st *structAstInfo) []string {
	var lines []string
	for _, f := range st.Typ.Fields.List {
		fieldType := exprToString(f.Type)
		for _, nameIdent := range f.Names {
			lines = append(lines, genResetFieldLine("s."+nameIdent.Name, fieldType))
		}
	}
	return lines
}

func genResetFieldLine(fieldName, typ string) string {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "float32", "float64", "complex64", "complex128":
		return fmt.Sprintf("%s = 0", fieldName)
	case "string":
		return fmt.Sprintf("%s = \"\"", fieldName)
	case "bool":
		return fmt.Sprintf("%s = false", fieldName)

	default:
		switch {
		case strings.HasPrefix(typ, "[]"):
			return fieldName + " = " + fieldName + "[:0]"
		case typ == "map[interface{}]interface{}":
			return "clear(" + fieldName + ")"
		case strings.HasPrefix(typ, "*"):
			return handlePointerType(typ, fieldName)
		default:
			return handleResetableType(fieldName)
		}
	}
}

func handlePointerType(typ, fieldName string) string {
	baseType := strings.TrimPrefix(typ, "*")
	if isPrimitive(baseType) {
		return fmt.Sprintf("if %s != nil { *%s = 0 }", fieldName, fieldName)
	}
	return fmt.Sprintf("if r, ok := any(%s).(interface{ Reset() }); ok && %s != nil { r.Reset() }",
		fieldName, fieldName)
}

func handleResetableType(fieldName string) string {
	return fmt.Sprintf("if r, ok := any(%s).(interface{ Reset() }); ok { r.Reset() }", fieldName)
}

func parsePackages(root string) ([]*packageInfo, error) {
	var pkgs []*packageInfo
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasPrefix(info.Name(), "_") {
			if strings.Contains(path, "/vendor/") || strings.Contains(path, "/testdata/") {
				return nil
			}

			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			var pkg *packageInfo
			for _, _pkg := range pkgs {
				if filepath.Dir(path) == _pkg.Dir {
					pkg = _pkg
					break
				}
			}
			if pkg == nil {
				pkg = &packageInfo{
					Name: f.Name.Name,
					Dir:  filepath.Dir(path),
				}
				pkgs = append(pkgs, pkg)
			}
			pkg.Ast = append(pkg.Ast, f)
		}
		return nil
	})
	return pkgs, err
}

func findResettableStructs(pkg *packageInfo) []*structAstInfo {
	var res []*structAstInfo

	for _, f := range pkg.Ast {
		for _, decl := range f.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			comment := ""
			if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
				comment = strings.TrimSpace(genDecl.Doc.Text())
			}

			if !strings.Contains(comment, "generate:reset") {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				res = append(res, &structAstInfo{
					Name: typeSpec.Name.Name,
					Typ:  structType,
				})
			}
		}
	}
	return res
}

func exprToString(expr ast.Expr) string {
	var b strings.Builder
	if err := format.Node(&b, token.NewFileSet(), expr); err != nil {
		return ""
	}
	return b.String()
}

var primitives = map[string]bool{
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	"float32":    true,
	"float64":    true,
	"complex64":  true,
	"complex128": true,
	"string":     true,
	"bool":       true,
}

func isPrimitive(typ string) bool {
	return primitives[typ]
}
