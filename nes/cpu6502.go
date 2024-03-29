package nes

import "fmt"

// constants representing the flags for the different values the status register can have
// each constant is a different bit, so they can be combined in the status register to signal what flags are set
const (
	C = (uint8)(1)   //00000001, carry bit
	Z = (uint8)(2)   //00000010, zero bit
	I = (uint8)(4)   //00000100, disable interrupts bit
	D = (uint8)(8)   //00001000, decimal mode bit (unused in the NES)
	B = (uint8)(16)  //00010000, break bit
	U = (uint8)(32)  //00100000, unused bit
	V = (uint8)(64)  //01000000, overflow bit
	N = (uint8)(128) //10000000, negative bit
)

// holds the data for the cpu
type CPU6502 struct {
	bus *Bus
	//status register, carries the different flags which are set or unset
	status uint8
	//the registers
	a uint8
	x uint8
	y uint8
	//stack pointer
	sptr uint8
	//program counter
	pc uint16
	//data received from addressing, after collecting the information from the addressing function it will be put here
	fetchedData uint8
	//addresses currently used
	//the absoulte address
	addrAbs uint16
	//the relative address, used when branching
	addrRel uint16

	//the current opcode
	opCode uint8
	//cycles for the current instruction
	cycles uint8
	//lookup table of instruction structs, the index in the table is the numerical value of the instruction
	instructions [256]Instruction
}

// function signatures for the different functions associated with an operation
type operation func(cpu *CPU6502) uint8
type addressingMode func(cpu *CPU6502) uint8

// this contains the information for each instruction/opcode
type Instruction struct {
	//name is for ease of understanding
	name     string
	op       operation
	modeType string
	addrMode addressingMode
	cycles   uint8
}

// constructor to create a cpu so it initializes the lookup table
func CreateCPU() *CPU6502 {
	cpu := CPU6502{}
	//ugly, gross, disgusting, bad, not good, but it initializes the entire table
	cpu.instructions = [256]Instruction{
		{name: "BRK", op: BRK, modeType: "IMP", addrMode: IMP, cycles: 7}, {name: "ORA", op: ORA, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "ASL", op: ASL, modeType: "ZPI", addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "PHP", op: PHP, modeType: "IMP", addrMode: IMP, cycles: 3}, {name: "ORA", op: ORA, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "ASL", op: ASL, modeType: "ACC", addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "ASL", op: ASL, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BPL", op: BPL, modeType: "REL", addrMode: REL, cycles: 2}, {name: "ORA", op: ORA, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "ASL", op: ASL, modeType: "ZPX", addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CLC", op: CLC, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "ORA", op: ORA, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "ASL", op: ASL, modeType: "ABX", addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "JSR", op: JSR, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "AND", op: AND, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "BIT", op: BIT, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "AND", op: AND, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "ROL", op: ROL, modeType: "ZPI", addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "PLP", op: PLP, modeType: "IMP", addrMode: IMP, cycles: 4}, {name: "AND", op: AND, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "ROL", op: ROL, modeType: "ACC", addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "BIT", op: BIT, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "AND", op: AND, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "ROL", op: ROL, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BMI", op: BMI, modeType: "REL", addrMode: REL, cycles: 2}, {name: "AND", op: AND, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "AND", op: AND, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "ROL", op: ROL, modeType: "ZPX", addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "SEC", op: SEC, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "AND", op: AND, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "AND", op: AND, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "ROL", op: ROL, modeType: "ABX", addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "RTI", op: RTI, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "LSR", op: LSR, modeType: "ZPI", addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "PHA", op: PHA, modeType: "IMP", addrMode: IMP, cycles: 3}, {name: "EOR", op: EOR, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "LSR", op: LSR, modeType: "ACC", addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "JMP", op: JMP, modeType: "ABS", addrMode: ABS, cycles: 3}, {name: "EOR", op: EOR, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "LSR", op: LSR, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BVC", op: BVC, modeType: "REL", addrMode: REL, cycles: 2}, {name: "EOR", op: EOR, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "LSR", op: LSR, modeType: "ZPX", addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CLI", op: CLI, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "EOR", op: EOR, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "LSR", op: LSR, modeType: "ABX", addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "RTS", op: RTS, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "ROR", op: ROR, modeType: "ZPI", addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "PLA", op: PLA, modeType: "IMP", addrMode: IMP, cycles: 4}, {name: "ADC", op: ADC, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "ROR", op: ROR, modeType: "ACC", addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "JMP", op: JMP, modeType: "IND", addrMode: IND, cycles: 5}, {name: "ADC", op: ADC, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "ROR", op: ROR, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BVS", op: BVS, modeType: "REL", addrMode: REL, cycles: 2}, {name: "ADC", op: ADC, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "ROR", op: ROR, modeType: "ZPX", addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "SEI", op: SEI, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "ADC", op: ADC, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "ROR", op: ROR, modeType: "ABX", addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "STA", op: STA, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "STY", op: STY, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "STA", op: STA, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "STX", op: STX, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "DEY", op: DEY, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "TXA", op: TXA, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "STY", op: STY, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "STA", op: STA, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "STX", op: STX, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BCC", op: BCC, modeType: "REL", addrMode: REL, cycles: 2}, {name: "STA", op: STA, modeType: "IDY", addrMode: IDY, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "STY", op: STY, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "STA", op: STA, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "STX", op: STX, modeType: "ZPY", addrMode: ZPY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "TYA", op: TYA, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "STA", op: STA, modeType: "ABY", addrMode: ABY, cycles: 5}, {name: "TXS", op: TXS, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "STA", op: STA, modeType: "ABX", addrMode: ABX, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "LDY", op: LDY, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "LDA", op: LDA, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "LDX", op: LDX, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "LDA", op: LDA, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "LDX", op: LDX, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "TAY", op: TAY, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "LDA", op: LDA, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "TAX", op: TAX, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "LDA", op: LDA, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "LDX", op: LDX, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BCS", op: BCS, modeType: "REL", addrMode: REL, cycles: 2}, {name: "LDA", op: LDA, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "LDA", op: LDA, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "LDX", op: LDX, modeType: "ZPY", addrMode: ZPY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CLV", op: CLV, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "LDA", op: LDA, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "TSX", op: TSX, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "LDA", op: LDA, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "LDX", op: LDX, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "CPY", op: CPY, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "CMP", op: CMP, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CPY", op: CPY, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "CMP", op: CMP, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "DEC", op: DEC, modeType: "ZPI", addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "INY", op: INY, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "CMP", op: CMP, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "DEX", op: DEX, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CPY", op: CPY, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "CMP", op: CMP, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "DEC", op: DEC, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BNE", op: BNE, modeType: "REL", addrMode: REL, cycles: 2}, {name: "CMP", op: CMP, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CMP", op: CMP, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "DEC", op: DEC, modeType: "ZPX", addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CLD", op: CLD, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "CMP", op: CMP, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CMP", op: CMP, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "DEC", op: DEC, modeType: "ABX", addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "CPX", op: CPX, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "SBC", op: SBC, modeType: "IDX", addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CPX", op: CPX, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "SBC", op: SBC, modeType: "ZPI", addrMode: ZPI, cycles: 3}, {name: "INC", op: INC, modeType: "ZPI", addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "INX", op: INX, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "SBC", op: SBC, modeType: "IMM", addrMode: IMM, cycles: 2}, {name: "NOP", op: NOP, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "CPX", op: CPX, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "SBC", op: SBC, modeType: "ABS", addrMode: ABS, cycles: 4}, {name: "INC", op: INC, modeType: "ABS", addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
		{name: "BEQ", op: BEQ, modeType: "REL", addrMode: REL, cycles: 2}, {name: "SBC", op: SBC, modeType: "IDY", addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "SBC", op: SBC, modeType: "ZPX", addrMode: ZPX, cycles: 4}, {name: "INC", op: INC, modeType: "ZPX", addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "SED", op: SED, modeType: "IMP", addrMode: IMP, cycles: 2}, {name: "SBC", op: SBC, modeType: "ABY", addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6}, {name: "SBC", op: SBC, modeType: "ABX", addrMode: ABX, cycles: 4}, {name: "INC", op: INC, modeType: "ABX", addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, modeType: "IMP", addrMode: IMP, cycles: 6},
	}
	//always set to 1
	cpu.SetFlag(U, true)

	return &cpu
}

// links the cpu to a bus, should be the bus it's contained in
func (cpu *CPU6502) ConnectBus(ptr *Bus) {
	cpu.bus = ptr
}

// uses the bus to attempt reads and writes
func (cpu *CPU6502) Write(addr uint16, data uint8) {
	cpu.bus.CPUWrite(addr, data)
}

func (cpu CPU6502) Read(addr uint16, readOnly bool) uint8 {
	temp := cpu.bus.CPURead(addr, readOnly)
	return temp
}

// sets the flag on the status registor for one of the values, if the flag is already set it unsets it
func (cpu *CPU6502) SetFlag(flag uint8, set bool) {
	if set {
		cpu.status |= flag
	} else if cpu.GetFlag(flag) {
		cpu.status ^= flag
	}

}

// checks if a specific flag is set on the status register
func (cpu CPU6502) GetFlag(flag uint8) bool {
	return cpu.status&flag == flag
}

// gets data for the cpu and sets the fetchedData value to what was gotten
func (cpu *CPU6502) fetchData() {
	//checks if the addressing mode itself sets the fetchedData field, if it does don't attempt to fetch data again and overwrite it
	mode := cpu.instructions[cpu.opCode].modeType
	if mode == "IMP" || mode == "IMM" || mode == "ACC" {
		return
	}
	cpu.fetchedData = cpu.Read(cpu.addrAbs, false)
}

// these four functions can occur at any point in operation, and will go after the current instruction is complete
// tells the cpu to advance one clock cycle
func (cpu *CPU6502) Clock() {
	//TODO breakpoint for passing first frame pass
	if cpu.pc-1 >= 0xc7af {
		fmt.Printf("")
	}
	if cpu.cycles == 0 {
		cpu.opCode = cpu.Read(cpu.pc, false)
		cpu.pc++
		//TODO OPCODE PRINTOUT HERE
		fmt.Printf("%0x "+cpu.instructions[cpu.opCode].name+" %0x %0x\n", cpu.pc-1, cpu.Read(cpu.pc, false), cpu.Read(cpu.pc+1, false))
		//get how many extra cycles in the address read and set the address
		addrCycles := cpu.instructions[cpu.opCode].addrMode(cpu)
		//get how many extra cycles in the CPU action and do the action
		cpuCycles := cpu.instructions[cpu.opCode].op(cpu)
		//add the given amount of cycles with any extra cycles
		cpu.cycles = cpu.instructions[cpu.opCode].cycles + cpuCycles + addrCycles
	}

	cpu.cycles--
}

//cpu interrupts

// resets the cpu
func (cpu *CPU6502) Reset() {
	cpu.a = 0
	cpu.x = 0
	cpu.y = 0
	//resets the stack pointer
	cpu.sptr = 0xfd
	cpu.status = 0
	//make sure this is always set
	cpu.SetFlag(U, true)

	//default location to look for data when the cpu is reset
	//program counter low byte = 0xfffc, high byte = 0xfffd
	cpu.pc = (uint16(cpu.Read(0xfffd, false)) << 8) | uint16(cpu.Read(0xfffc, false))

	cpu.addrAbs = 0
	cpu.addrRel = 0
	cpu.fetchedData = 0

	//time it takes to reset
	cpu.cycles = 8
}

// interrupt request, can be ignored depending on the interrupt flag of the status register
func (cpu *CPU6502) IRQ() {
	//checks if interrupts are disabled, if so escapes
	if !cpu.GetFlag(I) {
		return
	}

	//store the program counter on the stack, both bytes, little endian
	cpu.Write(0x0100+uint16(cpu.sptr), uint8((cpu.pc&0xff00)>>8))
	cpu.sptr--
	cpu.Write(0x0100+uint16(cpu.sptr), uint8(cpu.pc&0x00ff))
	cpu.sptr--

	//break flag
	cpu.SetFlag(B, false)
	//unused, making sure its set
	cpu.SetFlag(U, true)
	//interrupt flag
	cpu.SetFlag(I, true)

	//write the new status to the stack
	cpu.Write(0x0100+uint16(cpu.sptr), cpu.status)
	cpu.sptr--

	//moves program counter to known location after interrupt
	cpu.pc = uint16(cpu.Read(0xffff, false))<<8 | uint16(cpu.Read(0xfffe, false))

	//time taken
	cpu.cycles = 7
}

// non maskable interrupt request, unable to be ignored
func (cpu *CPU6502) NMI() {
	//store the program counter on the stack, both bytes, little endian
	cpu.Write(0x0100+uint16(cpu.sptr), uint8((cpu.pc&0xff00)>>8))
	cpu.sptr--
	cpu.Write(0x0100+uint16(cpu.sptr), uint8(cpu.pc&0x00ff))
	cpu.sptr--

	//break flag
	cpu.SetFlag(B, false)
	//unused, making sure its set
	cpu.SetFlag(U, true)
	//interrupt flag
	cpu.SetFlag(I, true)

	//write the new status to the stack
	cpu.Write(0x0100+uint16(cpu.sptr), cpu.status)
	cpu.sptr--

	//moves program counter to known location after interrupt
	cpu.pc = uint16(cpu.Read(0xfffb, false))<<8 | uint16(cpu.Read(0xfffa, false))

	//time taken
	cpu.cycles = 8
}

//addressing mode

// implied addressing, address is implicit in the opcode itself so nothing is needed
func IMP(cpu *CPU6502) uint8 {
	return 0
}

// accumulator addressing, data is retrieved from the accumulator
func ACC(cpu *CPU6502) uint8 {
	cpu.fetchedData = cpu.a
	return 0
}

// immediate addressing, the second byte in the instruction is the operand
func IMM(cpu *CPU6502) uint8 {
	cpu.fetchedData = cpu.Read(cpu.pc, false)
	cpu.pc++
	return 0
}

// absolute addressing, the 2nd instruction byte is the lower byte of the address, the 3rd is the high bytes, combined to allow access to any point in memory
func ABS(cpu *CPU6502) uint8 {
	low := cpu.Read(cpu.pc, false)
	cpu.pc++
	hi := cpu.Read(cpu.pc, false)
	cpu.pc++

	//combine into one number
	var data uint16 = uint16(hi)<<8 | uint16(low)

	cpu.addrAbs = data
	return 0
}

// zero page addressing, the second byte is the offset from the first page of the memory
// there is a glitch, replicated from original hardware where values that should go to the next page wrap back around, but shouldn't occur in this zero page mode because it only reads in one byte
func ZPI(cpu *CPU6502) uint8 {
	byte := cpu.Read(cpu.pc, false)
	cpu.pc++

	//wraps address that goes to the next page back around to the zero page
	cpu.addrAbs = uint16(byte) & 255
	return 0
}

// zero page addressing with X register offset, the second byte is the offset from the first page of the memory
// there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ZPX(cpu *CPU6502) uint8 {
	byte := cpu.Read(cpu.pc, false)
	cpu.pc++
	byte += cpu.x

	//wraps address that goes to the next page back around to the zero page
	cpu.addrAbs = uint16(byte) & 255
	return 0
}

// zero page addressing with Y register offset, the second byte is the offset from the first page of the memory
// there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ZPY(cpu *CPU6502) uint8 {
	byte := cpu.Read(cpu.pc, false)
	cpu.pc++
	byte += cpu.y

	//wraps address that goes to the next page back around to the zero page
	cpu.addrAbs = uint16(byte) & 255
	return 0
}

// absolute addressing with x register offset, the second byte is the offset from the first page of the memory
// there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ABX(cpu *CPU6502) uint8 {
	low := cpu.Read(cpu.pc, false)
	cpu.pc++
	hi := cpu.Read(cpu.pc, false)
	cpu.pc++

	//combine into one number
	var data uint16 = (uint16(hi)<<8 | uint16(low)) + uint16(cpu.x)

	cpu.addrAbs = data

	//checks if the page increased, if so the clock cycle count needs to increase
	if uint16(hi)<<8 != data&0xff00 {
		return 1
	} else {
		return 0
	}
}

// absolute addressing with Y register offset, the second byte is the offset from the first page of the memory
// there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ABY(cpu *CPU6502) uint8 {
	low := cpu.Read(cpu.pc, false)
	cpu.pc++
	hi := cpu.Read(cpu.pc, false)
	cpu.pc++

	//combine into one number
	var data uint16 = (uint16(hi)<<8 | uint16(low)) + uint16(cpu.y)

	cpu.addrAbs = data

	//checks if the page increased, if so the clock cycle count needs to increase
	if uint16(hi)<<8 != data&0xff00 {
		return 1
	} else {
		return 0
	}
}

// relative addressing, the second byte is added to the program counter for branching
func REL(cpu *CPU6502) uint8 {
	offset := uint16(cpu.Read(cpu.pc, false))
	cpu.pc++

	//if the number is 128 or greater unsigned, convert it to the 2's complement 16 digit version
	//because adresses are 16 bits, it needs to be converted when the value goes into the negative
	if offset&0x0080 == 0x0080 {
		offset |= 0xff00
	}

	cpu.addrRel = offset

	//will always increase the cycle, up to the function that calls this address mode to see if it goes to another page or not to determine the final increase
	return 1
}

// indexed indirect addressing, the second byte in the instruction is added to the X register without the carry, which is the low byte of the effective address, which is found on page zero, the following byte is the high byte, both on page zero
func IDX(cpu *CPU6502) uint8 {
	addr := cpu.Read(cpu.pc, false)
	cpu.pc++

	//adds them as uint8 so the carry is discarded
	addr += cpu.x

	low := cpu.Read(uint16(addr)&0x00ff, false)
	hi := cpu.Read((uint16(addr)+1)&0x00ff, false)

	cpu.addrAbs = uint16(uint16(hi)<<8 | uint16(low))
	return 0
}

// indirect indexed addressing, the second byte in the instruction is added to the Y register, which is the low byte of the effective address, which is found on page zero, the following byte added with the carry of the last addition is the high byte, both on page zero
func IDY(cpu *CPU6502) uint8 {
	addr := cpu.Read(cpu.pc, false)
	cpu.pc++

	low := cpu.Read(uint16(addr&0x00ff), false)
	hi := cpu.Read(uint16((addr+1)&0x00ff), false)

	//adds them as uint16 so the carry remains
	cpu.addrAbs = (uint16(hi)<<8 | uint16(low)) + uint16(cpu.y)

	//checks if the page increased, if so the clock cycle count needs to increase
	if uint16(hi)<<8 != cpu.addrAbs&0xff00 {
		return 1
	} else {
		return 0
	}
}

// absolute indirect addressing, the second byte in the instruction is the low byte of a memory location, the third instruction byte is the high byte, the data at that memory location is the low byte of the effective address, and the following byte is the effective high byte
// has a hardware bug where instead of going to the next page, it will wrap when forming the effective address, replicated here
func IND(cpu *CPU6502) uint8 {
	low := cpu.Read(cpu.pc, false)
	cpu.pc++
	hi := cpu.Read(cpu.pc, false)
	cpu.pc++

	addr := uint16(hi)<<8 | uint16(low)

	addrLow := cpu.Read(addr, false)
	addrHi := cpu.Read((addr + 1), false)

	//replicates the wrapping glitch, if the low byte is 0xff, then it wraps back to 0 on the same page instead of advancing
	if low == 0xff {
		addrHi = cpu.Read(uint16(hi)<<8, false)
	}

	cpu.addrAbs = uint16(uint16(addrHi)<<8 | uint16(addrLow))
	return 0
}

//opcodes

// nex, non-existent, placeholder put in to represent unspecified opcodes, more detail can later be put in to give the unspecified opcodes that get used their functionality, for now does nothing
func NEX(cpu *CPU6502) uint8 {
	return 0
}

// adc, add with carry, from specified memory to accumulator
func ADC(cpu *CPU6502) uint8 {
	cpu.fetchData()

	carryFlag := 0

	if cpu.GetFlag(C) {
		carryFlag = 1
	}

	res := uint16(cpu.fetchedData) + uint16(cpu.a) + uint16(carryFlag)

	//carry flag
	cpu.SetFlag(C, res&0xff00 > 0)
	//zero flag
	cpu.SetFlag(Z, res == 0)

	//overflow flag
	//checks if both values are negative and it wrapped to positive
	if (0x0080&cpu.a&cpu.fetchedData == 0x0080) && 0x0080&res == 0 {
		cpu.SetFlag(V, true)
	} else if ((^(cpu.a | cpu.fetchedData) & 0x0080) == 0x0080) && 0x0080&res == 0x0080 {
		//checks if both values are positive and it wrapped to negative
		cpu.SetFlag(V, true)
	} else {
		cpu.SetFlag(V, false)
	}

	//negative flag
	cpu.SetFlag(N, 0x0080&res == 0x0080)

	//set the accumulator register to the new value
	cpu.a = uint8(res & 0x00ff)

	return 0
}

// and, operates on memory and accumulator
func AND(cpu *CPU6502) uint8 {
	cpu.fetchData()
	cpu.a &= cpu.fetchedData
	//zero flag
	cpu.SetFlag(Z, cpu.a == 0)
	//negative flag
	cpu.SetFlag(N, cpu.a&0x80 == 0x80)

	return 0
}

// asl, arithmetic shift left, shifts left one bit, memory or accumulator
func ASL(cpu *CPU6502) uint8 {
	//if it is in accumulator mode, write to it, otherwise write the shift to the memory location
	if cpu.instructions[cpu.opCode].modeType == "ACC" {
		//carry flag
		cpu.SetFlag(C, cpu.a&0x0080 == 0x0080)
		cpu.a = cpu.a << 1
		//zero flag
		cpu.SetFlag(Z, uint8(cpu.a&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(cpu.a&0xff)&0x80 == 0x80)
	} else {
		//memory mode
		cpu.fetchData()
		mem := cpu.fetchedData

		//carry flag
		cpu.SetFlag(C, mem&0x0080 == 0x0080)

		mem = mem << 1

		//zero flag
		cpu.SetFlag(Z, uint8(mem&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(mem&0xff)&0x80 == 0x80)

		cpu.Write(cpu.addrAbs, mem)
	}

	return 0
}

// bcc, branch if carry clear, if the carry flag is clear, branch to the location as specified in the instruction
func BCC(cpu *CPU6502) uint8 {
	if !cpu.GetFlag(C) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// bcs, branch if carry set, if the carry flag is set, branch to the location as specified in the instruction
func BCS(cpu *CPU6502) uint8 {
	if cpu.GetFlag(C) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// beq, branch if equal/branch if zero, if the zero flag is set, branch to the location as specified in the instruction
func BEQ(cpu *CPU6502) uint8 {
	if cpu.GetFlag(Z) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// bit, bit test, the accumulator is ANDed with the supplied memory value and sets/clears the N and V flags
func BIT(cpu *CPU6502) uint8 {
	cpu.fetchData()

	and := cpu.fetchedData & cpu.a

	//zero flag
	cpu.SetFlag(Z, and == 0)
	//overflow flag
	cpu.SetFlag(V, and&0x40 == 0x40)
	//negative flag
	cpu.SetFlag(N, and&0x80 == 0x80)

	return 0
}

// bmi, branch if minus, if the negative flag is set branch to location
func BMI(cpu *CPU6502) uint8 {
	if cpu.GetFlag(N) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// bne, branch not equal/branch if not zero, if the zero flag is not set, branch to the location as specified in the instruction
func BNE(cpu *CPU6502) uint8 {
	if !cpu.GetFlag(Z) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// bpl, branch if positive, if the negative flag is not set, branch to the location as specified in the instruction
func BPL(cpu *CPU6502) uint8 {
	if !cpu.GetFlag(N) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// brk, break, generates an interrupt, stores the pc to the stack and sets the pc to 0xfffe and 0xffff and sets the break flag
func BRK(cpu *CPU6502) uint8 {
	cpu.pc++
	cpu.SetFlag(I, true)
	//store the program counter on the stack, both bytes, little endian
	cpu.Write(0x0100+uint16(cpu.sptr), uint8((cpu.pc&0xff00)>>8))
	cpu.sptr--
	cpu.Write(0x0100+uint16(cpu.sptr), uint8(cpu.pc&0x00ff))
	cpu.sptr--

	//break flag
	cpu.SetFlag(B, true)

	//write the new status to the stack
	cpu.Write(0x0100+uint16(cpu.sptr), cpu.status)
	cpu.sptr--

	//break flag set to 0 after status was pushed
	cpu.SetFlag(B, false)

	//moves program counter to known location after interrupt
	cpu.pc = uint16(cpu.Read(0xfffe, false))<<8 | uint16(cpu.Read(0xffff, false))

	return 0
}

// bvc, branch if overflow is clear, if the overflow flag is not set branch to location
func BVC(cpu *CPU6502) uint8 {
	if !cpu.GetFlag(V) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// bvs, branch if overflow is set, if the overflow flag is set branch to location
func BVS(cpu *CPU6502) uint8 {
	if cpu.GetFlag(V) {
		org := cpu.pc
		cpu.pc += cpu.addrRel

		//checks if the offset went to another page, if so return the other cycle increment
		if (org & 0xff00) != (cpu.pc & 0xff00) {
			return 1
		} else {
			return 0
		}
	}

	return 0
}

// clc, clears the carry flag
func CLC(cpu *CPU6502) uint8 {
	cpu.SetFlag(C, false)
	return 0
}

// cld, clears the decimal mode flag, should be set to 0 anyways
func CLD(cpu *CPU6502) uint8 {
	cpu.SetFlag(D, false)
	return 0
}

// cli, clears the interrupt disable flag
func CLI(cpu *CPU6502) uint8 {
	cpu.SetFlag(I, false)
	return 0
}

// clv, clears the overflow flag
func CLV(cpu *CPU6502) uint8 {
	cpu.SetFlag(V, false)
	return 0
}

// cmp, compare, sets the zero and carry flags appropriately with a compare between the accumulator and a memory value
func CMP(cpu *CPU6502) uint8 {
	cpu.fetchData()

	//carry flag, set if a >= m
	cpu.SetFlag(C, cpu.a >= cpu.fetchedData)

	//zero, set if a - m == 0, 8 bit
	cpu.SetFlag(Z, cpu.a-cpu.fetchedData == 0)

	//negative flag, set if a - m < 0
	cpu.SetFlag(N, (cpu.a-cpu.fetchedData)&0x80 == 0x80)
	return 0
}

// cpx, compare with x register, sets the zero and carry flags appropriately with a compare between the x register and a memory value
func CPX(cpu *CPU6502) uint8 {
	cpu.fetchData()

	//carry flag, set if a >= m
	cpu.SetFlag(C, cpu.x >= cpu.fetchedData)

	//zero, set if a - m == 0, 8 bit
	cpu.SetFlag(Z, cpu.x-cpu.fetchedData == 0)

	//negative flag, set if a - m < 0
	cpu.SetFlag(N, (cpu.x-cpu.fetchedData)&0x80 == 0x80)
	return 0
}

// cpy, compare with y register, sets the zero and carry flags appropriately with a compare between the y register and a memory value
func CPY(cpu *CPU6502) uint8 {
	cpu.fetchData()

	//carry flag, set if a >= m
	cpu.SetFlag(C, cpu.y >= cpu.fetchedData)

	//zero, set if a - m == 0, 8 bit
	cpu.SetFlag(Z, cpu.y-cpu.fetchedData == 0)

	//negative flag, set if a - m < 0
	cpu.SetFlag(N, (cpu.y-cpu.fetchedData)&0x80 == 0x80)
	return 0
}

// dec, decrement memory, decrement the memory value at the specific location by 1 and set the zero and negative flags if needed
func DEC(cpu *CPU6502) uint8 {
	cpu.fetchData()

	dec := cpu.fetchedData - 1

	//zero
	cpu.SetFlag(Z, dec == 0)

	//negative flag
	cpu.SetFlag(N, dec&0x80 == 0x80)

	cpu.Write(cpu.addrAbs, dec)
	return 0
}

// dex, decrement x register, decrement the x register by 1 and set the zero and negative flags if needed
func DEX(cpu *CPU6502) uint8 {
	cpu.x -= 1

	//zero
	cpu.SetFlag(Z, cpu.x == 0)

	//negative flag
	cpu.SetFlag(N, cpu.x&0x80 == 0x80)

	return 0
}

// dey, decrement y register, decrement the x register by 1 and set the zero and negative flags if needed
func DEY(cpu *CPU6502) uint8 {
	cpu.y -= 1

	//zero
	cpu.SetFlag(Z, cpu.y == 0)

	//negative flag
	cpu.SetFlag(N, cpu.y&0x80 == 0x80)

	return 0
}

// eor, exclusive or, xor's the accumulator with the provided memory value and set the zero and negative flags if needed
func EOR(cpu *CPU6502) uint8 {
	cpu.fetchData()

	cpu.a ^= cpu.fetchedData

	//zero
	cpu.SetFlag(Z, cpu.a == 0)

	//negative flag
	cpu.SetFlag(N, cpu.a&0x80 == 0x80)

	return 0
}

// inc, increment memory, increment the memory value at the specific location by 1 and set the zero and negative flags if needed
func INC(cpu *CPU6502) uint8 {
	cpu.fetchData()

	dec := cpu.fetchedData + 1

	//zero
	cpu.SetFlag(Z, dec == 0)

	//negative flag
	cpu.SetFlag(N, dec&0x80 == 0x80)

	cpu.Write(cpu.addrAbs, dec)
	return 0
}

// inc, increment x register, increment x by 1 and set the zero and negative flags if needed
func INX(cpu *CPU6502) uint8 {
	cpu.x += 1

	//zero
	cpu.SetFlag(Z, cpu.x == 0)

	//negative flag
	cpu.SetFlag(N, cpu.x&0x80 == 0x80)

	return 0
}

// inc, increment y register, increment y by 1 and set the zero and negative flags if needed
func INY(cpu *CPU6502) uint8 {
	cpu.y += 1

	//zero
	cpu.SetFlag(Z, cpu.y == 0)

	//negative flag
	cpu.SetFlag(N, cpu.y&0x80 == 0x80)

	return 0
}

// jmp, jump, jumps to the value specified by the operand by moving the program counter
// wraparound glitch is handled in the addressing function
func JMP(cpu *CPU6502) uint8 {
	cpu.pc = cpu.addrAbs
	return 0
}

// jsr, jump to subroutine, pushes the address minus one of the current point to the stack and then jumps to the value specified by the operand by moving the program counter
func JSR(cpu *CPU6502) uint8 {
	cpu.pc--

	cpu.Write(0x0100+uint16(cpu.sptr), uint8((cpu.pc&0xff00)>>8))
	cpu.sptr--

	cpu.Write(0x0100+uint16(cpu.sptr), uint8(cpu.pc&0x00ff))
	cpu.sptr--

	cpu.pc = cpu.addrAbs
	return 0
}

// lda, load accumulator, loads byte of memory into accumulator, sets zero and negative flags if necessary
func LDA(cpu *CPU6502) uint8 {
	cpu.fetchData()

	cpu.a = cpu.fetchedData
	//zero
	cpu.SetFlag(Z, cpu.a == 0)

	//negative flag
	cpu.SetFlag(N, cpu.a&0x80 == 0x80)

	return 0
}

// ldx, load x register, loads byte of memory into x register, sets zero and negative flags if necessary
func LDX(cpu *CPU6502) uint8 {
	cpu.fetchData()

	cpu.x = cpu.fetchedData

	//zero
	cpu.SetFlag(Z, cpu.x == 0)

	//negative flag
	cpu.SetFlag(N, cpu.x&0x80 == 0x80)

	return 0
}

// ldy, load y register, loads byte of memory into y register, sets zero and negative flags if necessary
func LDY(cpu *CPU6502) uint8 {
	cpu.fetchData()

	cpu.y = cpu.fetchedData

	//zero
	cpu.SetFlag(Z, cpu.y == 0)

	//negative flag
	cpu.SetFlag(N, cpu.y&0x80 == 0x80)

	return 0
}

// lsr, logical shift right, the accumulator or the given memory location is shifted to the right one bit, the 0 bit being but into the carry flag, the addressing mode specifies what is shifted
func LSR(cpu *CPU6502) uint8 {
	//if it is in accumulator mode, write to it, otherwise write the shift to the memory location
	if cpu.instructions[cpu.opCode].modeType == "ACC" {
		//carry flag
		cpu.SetFlag(C, cpu.a&1 == 1)
		cpu.a = cpu.a >> 1
		//zero flag
		cpu.SetFlag(Z, uint8(cpu.a&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(cpu.a&0xff)&0x80 == 0x80)
	} else {
		//memory mode
		cpu.fetchData()

		mem := cpu.fetchedData

		//carry flag
		cpu.SetFlag(C, mem&1 == 1)
		mem = mem >> 1
		//zero flag
		cpu.SetFlag(Z, uint8(mem&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(mem&0xff)&0x80 == 0x80)

		cpu.Write(cpu.addrAbs, mem)
	}

	return 0
}

// nop, no operation, simply passes and lets the clock function increment the program counter, unspecified opcodes can cause a nop to have slightly different behavior, but currently unimplemented
func NOP(cpu *CPU6502) uint8 {
	return 0
}

// ora, logical inclusive or, or's the accumulator with the specified byte in memory, sets zero and negative flags if needed
func ORA(cpu *CPU6502) uint8 {
	cpu.fetchData()

	cpu.a |= cpu.fetchedData

	//zero
	cpu.SetFlag(Z, cpu.a == 0)

	//negative flag
	cpu.SetFlag(N, cpu.a&0x80 == 0x80)

	return 0
}

// pha, push accumulator, pushes accumulator onto the stack
func PHA(cpu *CPU6502) uint8 {
	cpu.Write(0x0100+uint16(cpu.sptr), cpu.a)
	cpu.sptr--

	return 0
}

// php, push processor status, pushes status onto the stack
// some sources claim it will always push the B flag as 1? if errors happen this could be a potential source
func PHP(cpu *CPU6502) uint8 {
	cpu.Write(0x0100+uint16(cpu.status), cpu.a)
	cpu.sptr--

	return 0
}

// pla, pull accumulator, pulls the top value from the stack onto the accumulator, sets the zero and negative flags if needed
func PLA(cpu *CPU6502) uint8 {
	cpu.sptr++
	cpu.a = cpu.Read(0x0100+uint16(cpu.sptr), false)

	//zero
	cpu.SetFlag(Z, cpu.a == 0)

	//negative flag
	cpu.SetFlag(N, cpu.a&0x80 == 0x80)

	return 0
}

// plp, pull processor status, pulls the top value from the stack onto the status register, sets the zero and negative flags if needed
func PLP(cpu *CPU6502) uint8 {
	cpu.sptr++
	cpu.status = cpu.Read(0x0100+uint16(cpu.sptr), false)

	//zero
	cpu.SetFlag(Z, cpu.status == 0)

	//negative flag
	cpu.SetFlag(N, cpu.status&0x80 == 0x80)

	return 0
}

// rol, rotate left, rotates the accumulator or memory value one left, the old most significant bit becoming the carry, and the current carry becoming the least significant bit of the new value
func ROL(cpu *CPU6502) uint8 {
	//if it is in accumulator mode, write to it, otherwise write the shift to the memory location
	if cpu.instructions[cpu.opCode].modeType == "ACC" {
		//holds the value of the old carry
		car := uint8(0)
		if cpu.GetFlag(C) {
			car = 1
		}

		//carry flag
		cpu.SetFlag(C, cpu.a&0x0080 == 0x0080)
		cpu.a = cpu.a << 1
		cpu.a |= car
		//zero flag
		cpu.SetFlag(Z, uint8(cpu.a&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(cpu.a&0xff)&0x80 == 0x80)
	} else {
		//memory mode
		cpu.fetchData()
		mem := cpu.fetchedData

		//holds the value of the old carry
		car := uint8(0)
		if cpu.GetFlag(C) {
			car = 1
		}

		//carry flag
		cpu.SetFlag(C, mem&0x0080 == 0x0080)
		mem = mem << 1
		mem |= car
		//zero flag
		cpu.SetFlag(Z, uint8(mem&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(mem&0xff)&0x80 == 0x80)

		cpu.Write(cpu.addrAbs, mem)
	}

	return 0
}

// ror, rotate right, rotates the accumulator or memory value one right, the old least significant bit becoming the carry, and the current carry becoming the most significant bit of the new value
func ROR(cpu *CPU6502) uint8 {
	//if it is in accumulator mode, write to it, otherwise write the shift to the memory location
	if cpu.instructions[cpu.opCode].modeType == "ACC" {
		//holds the value of the old carry
		car := uint8(0)
		if cpu.GetFlag(C) {
			car = 0x0080
		}

		//carry flag
		cpu.SetFlag(C, cpu.a&1 == 1)
		cpu.a = cpu.a >> 1
		cpu.a |= car
		//zero flag
		cpu.SetFlag(Z, uint8(cpu.a&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(cpu.a&0xff)&0x80 == 0x80)
	} else {
		//memory mode
		cpu.fetchData()
		mem := cpu.fetchedData

		//holds the value of the old carry
		car := uint8(0)
		if cpu.GetFlag(C) {
			car = 0x0080
		}

		//carry flag
		cpu.SetFlag(C, mem&1 == 1)
		mem = mem >> 1
		mem |= car
		//zero flag
		cpu.SetFlag(Z, uint8(mem&0xff) == 0)
		//negative flag
		cpu.SetFlag(N, uint8(mem&0xff)&0x80 == 0x80)

		cpu.Write(cpu.addrAbs, mem)
	}

	return 0
}

// rti, return from interrupt, used at the end of an interrupt, retrieves the status and the program counter from the stack
func RTI(cpu *CPU6502) uint8 {
	cpu.sptr++
	cpu.status = cpu.Read(0x0100+uint16(cpu.sptr), false)
	cpu.sptr++
	//little endian, so low byte is on the smaller address
	low := cpu.Read(0x0100+uint16(cpu.sptr), false)
	cpu.sptr++
	hi := cpu.Read(0x0100+uint16(cpu.sptr), false)

	cpu.pc = uint16(hi)<<8 | uint16(low)

	return 0
}

// rts, return from subroutine, retreives the program counter from the stack, and advances to the next instruction
func RTS(cpu *CPU6502) uint8 {
	cpu.sptr++
	//little endian, so low byte is on the smaller address
	low := cpu.Read(0x0100+uint16(cpu.sptr), false)
	cpu.sptr++
	hi := cpu.Read(0x0100+uint16(cpu.sptr), false)

	cpu.pc = uint16(hi)<<8 | uint16(low)

	cpu.pc++
	if cpu.sptr != 255 {
		return 0
	}
	return 0
}

// sbc, subtract with carry, from specified memory to accumulator
func SBC(cpu *CPU6502) uint8 {
	cpu.fetchData()

	carryFlag := 0

	if cpu.GetFlag(C) {
		carryFlag = 1
	}

	//subtracts the supplied memory byte from the accumulator
	res := uint16(cpu.a) - uint16(cpu.fetchedData) - ^uint16(carryFlag)

	//carry flag
	cpu.SetFlag(C, res&0xff00 > 0)
	//zero flag
	cpu.SetFlag(Z, res == 0)

	//overflow flag
	//checks if both values are negative and it wrapped to positive
	if (0x0080&cpu.a&cpu.fetchedData == 0x0080) && 0x0080&res == 0 {
		cpu.SetFlag(V, true)
	}
	//checks if both values are positive and it wrapped to negative
	if ((^(cpu.a | cpu.fetchedData) & 0x0080) == 0x0080) && 0x0080&res == 0x0080 {
		cpu.SetFlag(V, true)
	}

	//negative flag
	cpu.SetFlag(N, 0x0080&res == 0x0080)

	//set the accumulator register to the new value
	cpu.a = uint8(res & 0x00ff)

	return 0
}

// sec, set carry flag, sets carry flag
func SEC(cpu *CPU6502) uint8 {
	cpu.SetFlag(C, true)

	return 0
}

// sed, set decimal flag, sets decimal flag
func SED(cpu *CPU6502) uint8 {
	cpu.SetFlag(D, true)

	return 0
}

// sei, set interrupt flag, sets interrupt flag
func SEI(cpu *CPU6502) uint8 {
	cpu.SetFlag(I, true)

	return 0
}

// sta, store accumulator, stores accumulator to memory
func STA(cpu *CPU6502) uint8 {
	if cpu.a != 0 {
		fmt.Printf("")
	}
	cpu.Write(cpu.addrAbs, cpu.a)

	return 0
}

// stx, store x register, stores x register to memory
func STX(cpu *CPU6502) uint8 {
	cpu.Write(cpu.addrAbs, cpu.x)

	return 0
}

// sty, store y register, stores y register to memory
func STY(cpu *CPU6502) uint8 {
	cpu.Write(cpu.addrAbs, cpu.y)

	return 0
}

// tax, transfer accumulator to x, copies the accumulator to the x register, sets zero and negative flags if needed
func TAX(cpu *CPU6502) uint8 {
	cpu.x = cpu.a

	//zero flag
	cpu.SetFlag(Z, uint8(cpu.x&0xff) == 0)
	//negative flag
	cpu.SetFlag(N, uint8(cpu.x&0xff)&0x80 == 0x80)

	return 0
}

// tay, transfer accumulator to y, copies the accumulator to the y register, sets zero and negative flags if needed
func TAY(cpu *CPU6502) uint8 {
	cpu.y = cpu.a

	//zero flag
	cpu.SetFlag(Z, uint8(cpu.y&0xff) == 0)
	//negative flag
	cpu.SetFlag(N, uint8(cpu.y&0xff)&0x80 == 0x80)

	return 0
}

// tsx, transfer stack pointer to x, copies the stack pointer to the x register, sets zero and negative flags if needed
func TSX(cpu *CPU6502) uint8 {
	cpu.x = cpu.sptr

	//zero flag
	cpu.SetFlag(Z, uint8(cpu.x&0xff) == 0)
	//negative flag
	cpu.SetFlag(N, uint8(cpu.x&0xff)&0x80 == 0x80)

	return 0
}

// txa, transfer x to accumulator, copies the x register to the accumulator, sets zero and negative flags if needed
func TXA(cpu *CPU6502) uint8 {
	cpu.a = cpu.x

	//zero flag
	cpu.SetFlag(Z, uint8(cpu.a&0xff) == 0)
	//negative flag
	cpu.SetFlag(N, uint8(cpu.a&0xff)&0x80 == 0x80)

	return 0
}

// txs, transfer x to stack pointer, copies the x register to the stack pointer
func TXS(cpu *CPU6502) uint8 {
	cpu.sptr = cpu.x

	return 0
}

// tya, transfer y to accumulator, copies the y register to the accumulator, sets zero and negative flags if needed
func TYA(cpu *CPU6502) uint8 {
	cpu.a = cpu.y

	//zero flag
	cpu.SetFlag(Z, uint8(cpu.a&0xff) == 0)
	//negative flag
	cpu.SetFlag(N, uint8(cpu.a&0xff)&0x80 == 0x80)

	return 0
}
