package compiler

import (
	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/vm"
)

func compileUnary(compiledFunc *vm.CompiledFunc, e *ast.Unary, file *Compiled) (vm.Register, string, error) {
	returns1, kinds, err := compileExpr(compiledFunc, e.Expr, file)
	if err != nil {
		return "", "", err
	}

	var ins vm.Instruction
	switch e.Op {
	case "not":
		returns2 := compiledFunc.NextRegister()
		ins = &vm.Not{
			Left:   returns1[0],
			Result: returns2,
		}
		compiledFunc.Append(ins)

		return returns2, kinds[0], nil

	case "-":
		zeroAt := compiledFunc.NextRegister()
		compiledFunc.Append(&vm.Assign{
			VariableName: zeroAt,
			Value:        asttest.NewLiteralNumber("0"),
		})

		returns2 := compiledFunc.NextRegister()
		ins = &vm.Subtract{
			Left:   zeroAt,
			Right:  returns1[0],
			Result: returns2,
		}
		compiledFunc.Append(ins)

		return returns2, kinds[0], nil

	case "++":
		oneAt := compiledFunc.NextRegister()
		compiledFunc.Append(&vm.Assign{
			VariableName: oneAt,
			Value:        asttest.NewLiteralNumber("1"),
		})

		ins = &vm.Add{
			Left:   vm.Register(returns1[0]),
			Right:  vm.Register(oneAt),
			Result: vm.Register(returns1[0]),
		}
		compiledFunc.Append(ins)

		return returns1[0], kinds[0], nil

	case "--":
		oneAt := compiledFunc.NextRegister()
		compiledFunc.Append(&vm.Assign{
			VariableName: oneAt,
			Value:        asttest.NewLiteralNumber("1"),
		})

		ins = &vm.Subtract{
			Left:   returns1[0],
			Right:  oneAt,
			Result: returns1[0],
		}
		compiledFunc.Append(ins)

		return returns1[0], kinds[0], nil
	}

	panic(e.Op)
}
