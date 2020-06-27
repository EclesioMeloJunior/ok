package compiler

import (
	"fmt"
	"strings"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/lexer"
	"github.com/elliotchance/ok/vm"
)

func getBinaryInstruction(op string, left, right, result string) (vm.Instruction, string) {
	switch op {
	case "data + data":
		return &vm.Combine{Left: left, Right: right, Result: result}, "data"

	case "data += data":
		return &vm.Combine{Left: left, Right: right, Result: left}, "data"

	case "number + number":
		return &vm.Add{Left: left, Right: right, Result: result}, "number"

	case "number += number":
		return &vm.Add{Left: left, Right: right, Result: left}, "number"

	case "number - number":
		return &vm.Subtract{Left: left, Right: right, Result: result}, "number"

	case "number -= number":
		return &vm.Subtract{Left: left, Right: right, Result: left}, "number"

	case "number * number":
		return &vm.Multiply{Left: left, Right: right, Result: result}, "number"

	case "number *= number":
		return &vm.Multiply{Left: left, Right: right, Result: left}, "number"

	case "number / number":
		return &vm.Divide{Left: left, Right: right, Result: result}, "number"

	case "number /= number":
		return &vm.Divide{Left: left, Right: right, Result: left}, "number"

	case "number % number":
		return &vm.Remainder{Left: left, Right: right, Result: result}, "number"

	case "number %= number":
		return &vm.Remainder{Left: left, Right: right, Result: left}, "number"

	case "string + string":
		return &vm.Concat{Left: left, Right: right, Result: result}, "string"

	case "string += string":
		return &vm.Concat{Left: left, Right: right, Result: left}, "string"

	case "bool and bool":
		return &vm.And{Left: left, Right: right, Result: result}, "bool"

	case "bool or bool":
		return &vm.Or{Left: left, Right: right, Result: result}, "bool"

	case "bool == bool", "char == char", "data == data", "string == string":
		return &vm.Equal{Left: left, Right: right, Result: result}, "bool"

	case "number == number":
		return &vm.EqualNumber{Left: left, Right: right, Result: result}, "bool"

	case "bool != bool", "char != char", "data != data", "string != string":
		return &vm.NotEqual{Left: left, Right: right, Result: result}, "bool"

	case "number != number":
		return &vm.NotEqualNumber{Left: left, Right: right, Result: result}, "bool"

	case "number > number":
		return &vm.GreaterThanNumber{Left: left, Right: right, Result: result}, "bool"

	case "string > string":
		return &vm.GreaterThanString{Left: left, Right: right, Result: result}, "bool"

	case "number < number":
		return &vm.LessThanNumber{Left: left, Right: right, Result: result}, "bool"

	case "string < string":
		return &vm.LessThanString{Left: left, Right: right, Result: result}, "bool"

	case "number >= number":
		return &vm.GreaterThanEqualNumber{Left: left, Right: right, Result: result}, "bool"

	case "string >= string":
		return &vm.GreaterThanEqualString{Left: left, Right: right, Result: result}, "bool"

	case "number <= number":
		return &vm.LessThanEqualNumber{Left: left, Right: right, Result: result}, "bool"

	case "string <= string":
		return &vm.LessThanEqualString{Left: left, Right: right, Result: result}, "bool"
	}

	return nil, ""
}

func compileBinary(compiledFunc *vm.CompiledFunc, node *ast.Binary, fns map[string]*ast.Func) (string, string, error) {
	// TokenAssign is not possible here because that is handled by an Assign
	// operation.
	if node.Op == lexer.TokenPlusAssign ||
		node.Op == lexer.TokenMinusAssign ||
		node.Op == lexer.TokenTimesAssign ||
		node.Op == lexer.TokenDivideAssign ||
		node.Op == lexer.TokenRemainderAssign {

		right, rightKind, err := compileExpr(compiledFunc, node.Right, fns)
		if err != nil {
			return "", "", err
		}

		// TODO(elliot): Check +=, etc.
		if key, ok := node.Left.(*ast.Key); ok {
			arrayOrMapResults, arrayOrMapKind, err := compileExpr(compiledFunc, key.Expr, fns)
			if err != nil {
				return "", "", err
			}

			// TODO(elliot): Check this is a sane operation.
			keyResults, _, err := compileExpr(compiledFunc, key.Key, fns)
			if err != nil {
				return "", "", err
			}

			if strings.HasPrefix(arrayOrMapKind[0], "[]") {
				ins := &vm.ArraySet{
					Array: arrayOrMapResults[0],
					Index: keyResults[0],
					Value: right[0],
				}
				compiledFunc.Append(ins)
			} else {
				ins := &vm.MapSet{
					Map:   arrayOrMapResults[0],
					Key:   keyResults[0],
					Value: right[0],
				}
				compiledFunc.Append(ins)
			}

			// TODO(elliot): Return something more reasonable here.
			return "", "", nil
		}

		variable, ok := node.Left.(*ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("cannot assign to non-variable")
		}

		// Make sure we do not assign the wrong type to an existing variable.
		if v, ok := compiledFunc.Variables[variable.Name]; ok && rightKind[0] != v {
			return "", "", fmt.Errorf(
				"cannot assign %s to variable %s (expecting %s)",
				rightKind, variable.Name, v)
		}

		returns := compiledFunc.NextRegister()

		switch node.Op {
		case lexer.TokenPlusAssign:
			switch rightKind[0] {
			case "data":
				compiledFunc.Append(&vm.Combine{
					Left:   variable.Name,
					Right:  right[0],
					Result: returns,
				})

			case "number":
				compiledFunc.Append(&vm.Add{
					Left:   variable.Name,
					Right:  right[0],
					Result: returns,
				})

			case "string":
				compiledFunc.Append(&vm.Concat{
					Left:   variable.Name,
					Right:  right[0],
					Result: returns,
				})
			}

		case lexer.TokenMinusAssign:
			compiledFunc.Append(&vm.Subtract{
				Left:   variable.Name,
				Right:  right[0],
				Result: returns,
			})

		case lexer.TokenTimesAssign:
			compiledFunc.Append(&vm.Multiply{
				Left:   variable.Name,
				Right:  right[0],
				Result: returns,
			})

		case lexer.TokenDivideAssign:
			compiledFunc.Append(&vm.Divide{
				Left:   variable.Name,
				Right:  right[0],
				Result: returns,
			})

		case lexer.TokenRemainderAssign:
			compiledFunc.Append(&vm.Remainder{
				Left:   variable.Name,
				Right:  right[0],
				Result: returns,
			})
		}

		ins := &vm.Assign{
			VariableName: variable.Name,
			Register:     returns,
		}
		compiledFunc.Append(ins)
		compiledFunc.NewVariable(variable.Name, rightKind[0])

		return variable.Name, rightKind[0], nil
	}

	_, _, returns, returnKind, err := compileComparison(compiledFunc, node, fns)

	return returns, returnKind, err
}

func compileComparison(compiledFunc *vm.CompiledFunc, node *ast.Binary, fns map[string]*ast.Func) (string, string, string, string, error) {
	left, leftKind, err := compileExpr(compiledFunc, node.Left, fns)
	if err != nil {
		return "", "", "", "", err
	}

	right, rightKind, err := compileExpr(compiledFunc, node.Right, fns)
	if err != nil {
		return "", "", "", "", err
	}

	returns := compiledFunc.NextRegister()

	op := fmt.Sprintf("%s %s %s", leftKind[0], node.Op, rightKind[0])
	if bop, kind := getBinaryInstruction(op, left[0], right[0], returns); bop != nil {
		// TODO(elliot): It would be nice to be able to evaluate expressions
		//  involving literals here. So, 1 + 1 just becomes 2.

		compiledFunc.Append(bop)

		return left[0], right[0], returns, kind, nil
	}

	return left[0], right[0], returns, "", fmt.Errorf("cannot perform %s", op)
}
