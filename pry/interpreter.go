package pry

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strconv"
)

type Scope map[string]interface{}

func InterpretString(scope Scope, exprStr string) (interface{}, error) {
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, err
	}
	return InterpretExpr(scope, expr)
}

func InterpretExpr(scope Scope, expr ast.Expr) (interface{}, error) {
	//fmt.Printf("EXPR %#v\n", expr)

	switch e := expr.(type) {
	case *ast.Ident:
		builtinScope := map[string]interface{}{
			"nil":   nil,
			"true":  true,
			"false": false,
		}
		obj, exists := scope[e.Name]
		if !exists {
			obj, exists = builtinScope[e.Name]
			if !exists {
				return nil, errors.New(fmt.Sprint("Can't find EXPR ", e.Name))
			}
		}
		return obj, nil
	case *ast.SelectorExpr:
		X, err := InterpretExpr(scope, e.X)
		if err != nil {
			return nil, err
		}
		sel := e.Sel

		rVal := reflect.ValueOf(X)
		if rVal.Kind() != reflect.Struct {
			return nil, errors.New(fmt.Sprintf("%#v is not a struct and thus has no field %#v", X, sel.Name))
		}
		zero := reflect.ValueOf(nil)
		field := rVal.FieldByName(sel.Name)
		if field != zero {
			return field.Interface(), nil
		}
		method := rVal.MethodByName(sel.Name)
		if method != zero {
			return method.Interface(), nil
		}
		return nil, errors.New(fmt.Sprintf("Unknown field %#v", sel.Name))
	case *ast.CallExpr:
		fun, err := InterpretExpr(scope, e.Fun)
		if err != nil {
			return nil, err
		}
		args := e.Args
		_ = args
		funVal := reflect.ValueOf(fun)
		// TODO CALL WITH ARGS
		values := funVal.Call([]reflect.Value{})
		return ValuesToInterfaces(values)[0], nil
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT:
			return strconv.Atoi(e.Value)
		case token.FLOAT, token.IMAG:
			return strconv.ParseFloat(e.Value, 64)
		case token.CHAR:
			return e.Value[1], nil
		case token.STRING:
			return e.Value[1 : len(e.Value)-1], nil
		default:
			return nil, errors.New(fmt.Sprintf("Unknown basic literal %d", e.Kind))
		}
	case *ast.CompositeLit:
		fmt.Printf("TODO COMPLIT %#v\n", e)
		switch t := e.Type.(type) {
		case *ast.ArrayType:
			fmt.Printf("TYPE %#v\n", t)
		}
		return e, nil
	case *ast.BinaryExpr:
		x, err := InterpretExpr(scope, e.X)
		if err != nil {
			return nil, err
		}
		y, err := InterpretExpr(scope, e.Y)
		if err != nil {
			return nil, err
		}

		return ComputeBinaryOp(x, y, e.Op)
	case *ast.ParenExpr:
		return InterpretExpr(scope, e.X)
	default:
		return nil, errors.New(fmt.Sprintf("Unknown EXPR %T", e))
	}
}

func ValuesToInterfaces(vals []reflect.Value) []interface{} {
	inters := make([]interface{}, len(vals))
	for i, val := range vals {
		inters[i] = val.Interface()
	}
	return inters
}
