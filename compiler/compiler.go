package compiler

import (
	"cidoka/ast"
	"cidoka/code"
	"cidoka/object"
	"cidoka/token"
	"fmt"
	"sort"
)

var (
	loopContinuePos = []int{}
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants []object.Object

	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = symbolTable
	compiler.constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	// Statements
	case *ast.LetStatement:
		if s, ok := c.symbolTable.ResolveNoRecursion(node.Name.Value); ok && s.Scope != FunctionScope {
			return fmt.Errorf("variable %s already declared", node.Name.Value)
		}

		symbol := c.symbolTable.Define(node.Name.Value)
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		if symbol.Scope == GlobalScope {
			c.emit(code.OpDeclareGlobal, symbol.Index)
		} else {
			c.emit(code.OpDeclareLocal, symbol.Index)
		}

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LoopStatement:
		c.enterScope()

		currentContinueCount := len(loopContinuePos)

		if node.Initializer != nil {
			err := c.Compile(node.Initializer)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
			c.emit(code.OpPop)
		}

		conditionPos := len(c.currentInstructions())

		if node.Condition == nil {
			c.emit(code.OpTrue)
		} else {
			err := c.Compile(node.Condition)
			if err != nil {
				return err
			}
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumptNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		// Update any Continue statements with the correct value
		for _, pos := range loopContinuePos[currentContinueCount:] {
			c.changeOperand(pos, len(c.currentInstructions()))
		}

		if node.Update != nil {
			err = c.Compile(node.Update)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpNull)
			c.emit(code.OpPop)
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		// Update the `OpJumpNotTruthy` with the correct value
		afterUpdatePos := len(c.currentInstructions())
		c.changeOperand(jumptNotTruthyPos, afterUpdatePos)

		// Update the `OpJump` with the correct value
		c.changeOperand(jumpPos, conditionPos)

		c.emit(code.OpBreak)

		freeSymbols := c.symbolTable.FreeSymbols
		numLoc := c.symbolTable.numDefinitions
		ins := c.leaveScope()

		free := make([]object.FreeVariable, len(freeSymbols))
		for i, s := range freeSymbols {
			free[i] = object.FreeVariable{Index: s.Index, Scope: s.ScopeIndex}
		}

		compiled := &object.CompiledLoop{
			Instructions: ins,
			NumLocals:    numLoc,
			Free:         free,
		}

		loopContinuePos = loopContinuePos[:currentContinueCount]

		idx := c.addConstant(compiled)

		c.emit(code.OpLoop, idx)

	case *ast.BreakStatement:
		c.emit(code.OpBreak)

	case *ast.ContinueStatement:
		loopContinuePos = append(loopContinuePos, c.emit(code.OpJump, 9999))

	// Expressions
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

	case *ast.AssignExpression:
		var symbol Symbol

		switch left := node.Left.(type) {
		case *ast.Identifier:
			var ok bool
			symbol, ok = c.symbolTable.Resolve(left.Value)
			if !ok {
				return fmt.Errorf("undefined variable %s", left.Value)
			}

			if node.Token.Type != token.ASSIGN {
				c.loadSymbol(symbol)
			}

		case *ast.IndexExpression:
			err := c.Compile(left.Left)
			if err != nil {
				return err
			}

			err = c.Compile(left.Index)
			if err != nil {
				return err
			}
		}

		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+=":
			c.emit(code.OpAdd)
		case "-=":
			c.emit(code.OpSub)
		case "*=":
			c.emit(code.OpMul)
		case "/=":
			c.emit(code.OpDiv)
		case "%=":
			c.emit(code.OpMod)
		}

		switch node.Left.(type) {
		case *ast.Identifier:
			if symbol.Scope == GlobalScope {
				c.emit(code.OpSetGlobal, symbol.Index)
			} else if symbol.Scope == FreeScope {
				c.emit(code.OpSetFree, symbol.Index)
			} else {
				c.emit(code.OpSetLocal, symbol.Index)
			}
		case *ast.IndexExpression:
			c.emit(code.OpSetIndex)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(float))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "%":
			c.emit(code.OpMod)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		case ">=":
			c.emit(code.OpGreaterOrEqual)
		case "<":
			c.emit(code.OpLessThan)
		case "<=":
			c.emit(code.OpLessOrEqual)
		case "&&":
			c.emit(code.OpAnd)
		case "||":
			c.emit(code.OpOr)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.PostfixExpression:
		switch node.Left.(type) {
		case *ast.Identifier, *ast.IndexExpression:
			var op string
			if node.Operator == "++" {
				op = "+="
			} else {
				op = "-="
			}

			c.Compile(&ast.AssignExpression{
				Left:     node.Left,
				Operator: op,
				Token:    node.Token,
				Right:    &ast.IntegerLiteral{Value: 1},
			})
		default:
			var op string
			if node.Operator == "++" {
				op = "+"
			} else {
				op = "-"
			}

			c.Compile(&ast.InfixExpression{
				Left:     node.Left,
				Operator: op,
				Right:    &ast.IntegerLiteral{Value: 1},
			})
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumptNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumptNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.FunctionLiteral:
		c.enterScope()

		if node.Name != "" {
			c.symbolTable.DefineFunctionName(node.Name)
		}

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		fnIndex := c.addConstant(compiledFn)
		c.emit(code.OpClosure, fnIndex, len(freeSymbols))

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, arg := range node.Arguments {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpGetIndex)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) addInstruction(ins code.Instructions) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	opcode := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(opcode, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(code.OpGetFree, s.Index)
	case FunctionScope:
		c.emit(code.OpCurrentClosure)
	}
}
