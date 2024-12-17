package tribitemulator

import (
	"log/slog"
	"math"
)

func pow2(exponent int) int {
	return int(math.Pow(2, float64(exponent)))
}

type TribitEmulator struct {
	registerA          int
	registerB          int
	registerC          int
	instructionPointer int
}

func NewTribitEmulator(initRegisterA, initRegisterB, initRegisterC int) TribitEmulator {
	return TribitEmulator{
		registerA: initRegisterA,
		registerB: initRegisterB,
		registerC: initRegisterC,
	}
}

func (emulator TribitEmulator) getComboOperand(operand ComboOperand) int {
	switch operand {
	case COMBO_OPERAND_LITERAL_0:
		return 0
	case COMBO_OPERAND_LITERAL_1:
		return 1
	case COMBO_OPERAND_LITERAL_2:
		return 2
	case COMBO_OPERAND_LITERAL_3:
		return 3
	case COMBO_OPERAND_REGISTER_A:
		return emulator.registerA
	case COMBO_OPERAND_REGISTER_B:
		return emulator.registerB
	case COMBO_OPERAND_REGISTER_C:
		return emulator.registerC
	}

	slog.Error("unexpected combo operand encountered", "operand", operand)
	return 0
}

func (emulator TribitEmulator) ExecuteProgram(program []int) []int {
	emulator.instructionPointer = 0
	output := make([]int, 0)

	// halt when instruction pointer reaches end of program
	for emulator.instructionPointer < len(program)-1 {
		instruction := Instruction(program[emulator.instructionPointer])
		operand := program[emulator.instructionPointer+1]

		// slog.Debug("program loop", "instruction pointer", emulator.instructionPointer, "instruction", instruction, "operand", operand, "emulator state", emulator)

		switch instruction {
		case INSTRUCTION_ADV:
			comboOperand := emulator.getComboOperand(ComboOperand(operand))
			emulator.registerA = emulator.registerA / pow2(comboOperand)
		case INSTRUCTION_BDV:
			comboOperand := emulator.getComboOperand(ComboOperand(operand))
			emulator.registerB = emulator.registerA / pow2(comboOperand)
		case INSTRUCTION_CDV:
			comboOperand := emulator.getComboOperand(ComboOperand(operand))
			emulator.registerC = emulator.registerA / pow2(comboOperand)
		case INSTRUCTION_BXL:
			emulator.registerB = emulator.registerB ^ operand
		case INSTRUCTION_BST:
			comboOperand := emulator.getComboOperand(ComboOperand(operand))
			emulator.registerB = comboOperand % 8
		case INSTRUCTION_JNZ:
			if emulator.registerA != 0 {
				emulator.instructionPointer = operand
				continue // Don't increment instruction pointer at end of loop
			}
		case INSTRUCTION_BXC:
			emulator.registerB = emulator.registerB ^ emulator.registerC
		case INSTRUCTION_OUT:
			comboOperand := emulator.getComboOperand(ComboOperand(operand))
			output = append(output, comboOperand%8)
		default:
			slog.Error("encountered unexpected instruction", "instruction pointer", emulator.instructionPointer, "instruction", instruction)
		}
		emulator.instructionPointer += 2
	}

	// slog.Debug("program halted", "instruction pointer", emulator.instructionPointer, "emulator state", emulator)

	return output
}
