package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions are a slice of bytes
type Instructions []byte

// Returns a string representation of the Instructions
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

// Formats the instruction to a string
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// Opcode is a byte
type Opcode byte

const (
	// Special Opcodes

	OpConstant Opcode = iota // Push a constant to the stack
	OpPop                    // Pop the top of the stack
	OpNull                   // Push a null value to the stack

	// Arithmetic Opcodes

	OpAdd // Pop the top two elements of the stack, add them and push the result to the stack
	OpSub // Pop the top two elements of the stack, subtract them and push the result to the stack
	OpMul // Pop the top two elements of the stack, multiply them and push the result to the stack
	OpDiv // Pop the top two elements of the stack, divide them and push the result to the stack
	OpMod // Pop the top two elements of the stack, modulo them and push the result to the stack

	// Boolean Opcodes

	OpTrue  // Push a true value to the stack
	OpFalse // Push a false value to the stack

	// Comparison Opcodes

	OpEqual          // Pop the top two elements of the stack and compare for equality, push the result to the stack
	OpNotEqual       // Pop the top two elements of the stack and compare for inequality, push the result to the stack
	OpGreaterThan    // Pop the top two elements of the stack and compare for greater than, push the result to the stack
	OpGreaterOrEqual // Pop the top two elements of the stack and compare for greater or equal, push the result to the stack
	OpLessThan       // Pop the top two elements of the stack and compare for less than, push the result to the stack
	OpLessOrEqual    // Pop the top two elements of the stack and compare for less or equal, push the result to the stack

	// Prefix Opcodes

	OpMinus // Pop the top element of the stack, perform arithmetic negation and push the result to the stack
	OpBang  // Pop the top element of the stack, perform boolean negation and push the result to the stack

	// Jump Opcodes

	OpJumpNotTruthy // Pop the top element of the stack and jump to a specific position if it is not truthy
	OpJump          // Jump to a specific position

	// Variable Opcodes

	OpSetGlobal // Pop the top element of the stack and set it to a global scope variable
	OpGetGlobal // Push a global scope variable to the stack
	OpSetLocal  // Pop the top element of the stack and set it to a local scope variable
	OpGetLocal  // Push a local scope variable to the stack

	// Data Structure Opcodes

	OpArray // Push an array to the stack made from the n elements below it
	OpHash  // Push a hash to the stack made from the n elements below it // n is even
	OpIndex // Pop the top two elements of the stack, using the first as an index to the second, push the result to the stack

	// Function Opcodes

	OpClosure    // Push a closure to the stack
	OpGetBuiltin // Push a builtin function to the stack
	OpCall       // Call top n+1 elements of the stack as a function // n is the number of arguments // last element is the function

	OpCurrentClosure // Push the current closure as a variable // recursion
	OpSetFree        // Pop the top element of the stack and set it to a free scope variable
	OpGetFree        // Push a variable from the free scope

	OpReturnValue // Return from a function with a value
	OpReturn      // Return from a function

	// Loop Opcodes

	OpLoop  // Push a loop onto the stack // pops after the loop
	OpBreak // Break out of a loop // pops a loop
)

// Opcode definitions
type Definition struct {
	Name          string // Name of the opcode
	OperandWidths []int  // Width of the operands
}

// Mapping of Opcode to its Definition
var definitions = map[Opcode]*Definition{
	// Special Opcodes

	OpConstant: {"OpConstant", []int{2}}, // Single operand of 2 bytes, 3 bytes in total
	OpPop:      {"OpPop", []int{}},       // No operands, 1 byte in total
	OpNull:     {"OpNull", []int{}},      // No operands, 1 byte in total,

	// Arithmetic Opcodes

	OpAdd: {"OpAdd", []int{}}, // No operands, 1 byte in total
	OpSub: {"OpSub", []int{}}, // No operands, 1 byte in total
	OpMul: {"OpMul", []int{}}, // No operands, 1 byte in total
	OpDiv: {"OpDiv", []int{}}, // No operands, 1 byte in total
	OpMod: {"OpMod", []int{}}, // No operands, 1 byte in total

	// Boolean Opcodes

	OpTrue:  {"OpTrue", []int{}},  // No operands, 1 byte in total
	OpFalse: {"OpFalse", []int{}}, // No operands, 1 byte in total

	// Comparison Opcodes

	OpEqual:          {"OpEqual", []int{}},          // No operands, 1 byte in total
	OpNotEqual:       {"OpNotEqual", []int{}},       // No operands, 1 byte in total
	OpGreaterThan:    {"OpGreaterThan", []int{}},    // No operands, 1 byte in total
	OpGreaterOrEqual: {"OpGreaterOrEqual", []int{}}, // No operands, 1 byte in total
	OpLessThan:       {"OpLessThan", []int{}},       // No operands, 1 byte in total
	OpLessOrEqual:    {"OpLessOrEqual", []int{}},    // No operands, 1 byte in total

	// Prefix Opcodes

	OpMinus: {"OpMinus", []int{}}, // No operands, 1 byte in total
	OpBang:  {"OpBang", []int{}},  // No operands, 1 byte in total

	// Jump Opcodes

	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}}, // Single operand of 2 bytes, 3 bytes in total
	OpJump:          {"OpJump", []int{2}},          // Single operand of 2 bytes, 3 bytes in total

	// Variable Opcodes

	OpSetGlobal: {"OpSetGlobal", []int{2}}, // Single operand of 2 bytes, 3 bytes in total
	OpGetGlobal: {"OpGetGlobal", []int{2}}, // Single operand of 2 bytes, 3 bytes in total
	OpSetLocal:  {"OpSetLocal", []int{1}},  // Single operand of 1 byte, 2 bytes in total
	OpGetLocal:  {"OpGetLocal", []int{1}},  // Single operand of 1 byte, 2 bytes in total

	// Data Structure Opcodes

	OpArray: {"OpArray", []int{2}}, // Single operand of 2 bytes, 3 bytes in total
	OpHash:  {"OpHash", []int{2}},  // Single operand of 2 bytes, 3 bytes in total
	OpIndex: {"OpIndex", []int{}},  // No operands, 1 byte in total

	// Function Opcodes

	OpClosure:    {"OpClosure", []int{2, 1}}, // Two operands of 2 and 1 bytes, 4 bytes in total
	OpGetBuiltin: {"OpGetBuiltin", []int{1}}, // Single operand of 1 byte, 2 bytes in total
	OpCall:       {"OpCall", []int{1}},       // Single operand of 1 byte, 2 bytes in total

	OpReturnValue: {"OpReturnValue", []int{}}, // No operands, 1 byte in total
	OpReturn:      {"OpReturn", []int{}},      // No operands, 1 byte in total

	OpCurrentClosure: {"OpCurrentClosure", []int{}}, // No operands, 1 byte in total
	OpSetFree:        {"OpSetFree", []int{1}},       // Single operand of 1 byte, 2 bytes in total
	OpGetFree:        {"OpGetFree", []int{1}},       // Single operand of 1 byte, 2 bytes in total

	// Loop Opcodes

	OpLoop:  {"OpLoop", []int{2}}, // Single operand of 2 bytes, 3 bytes in total
	OpBreak: {"OpBreak", []int{}}, // No operands, 1 byte in total
}

// Returns the Definition of the opcode
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Returns the byte representation of the opcode and its operands in big-endian order
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

// Reads the operands of the opcode
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// Reads a 16-bit unsigned integer in big-endian order
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

// Reads an 8-bit unsigned integer in big-endian order
func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
