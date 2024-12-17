package tribitemulator

//go:generate stringer -type Instruction
type Instruction int

const (
	// EXPECTS COMBO OPERAND
	// Perform division by taking REGISTER_A / 2^(operand)
	// Write result to register A
	INSTRUCTION_ADV Instruction = 0

	// EXPECTS LITERAL OPERAND
	// Perform bitwise XOR of register B and literal operand
	// Write result to register B
	INSTRUCTION_BXL Instruction = 1

	// EXPECTS COMBO OPERAND
	// Performs (combo operand) % 8
	// Write result to register B
	INSTRUCTION_BST Instruction = 2

	// EXPECTS LITERAL OPERAND
	// if register A is 0: noop
	// else: set instruction pointer to literal operand and instruction pointer *not incremented by 2*
	INSTRUCTION_JNZ Instruction = 3

	// EXPECTS OPERAND, IGNORED
	// Perform bitwise XOR of B and C, write result to B
	INSTRUCTION_BXC Instruction = 4

	// EXPECTS COMBO OPERAND
	// Perform (combo operand) % 8 and output value (with comma)
	INSTRUCTION_OUT Instruction = 5

	// EXPECTS COMBO OPERAND
	// Perform division by taking register A / 2^(operand)
	// Write result to register B
	INSTRUCTION_BDV Instruction = 6

	// EXPECTS COMBO OPERAND
	// Perform division by taking register A / 2^(operand)
	// Write result to register C
	INSTRUCTION_CDV Instruction = 7
)

//go:generate stringer -type ComboOperand
type ComboOperand int

const (
	COMBO_OPERAND_LITERAL_0  ComboOperand = 0
	COMBO_OPERAND_LITERAL_1  ComboOperand = 1
	COMBO_OPERAND_LITERAL_2  ComboOperand = 2
	COMBO_OPERAND_LITERAL_3  ComboOperand = 3
	COMBO_OPERAND_REGISTER_A ComboOperand = 4
	COMBO_OPERAND_REGISTER_B ComboOperand = 5
	COMBO_OPERAND_REGISTER_C ComboOperand = 6
	COMBO_OPERAND_RESERVED   ComboOperand = 7
)
