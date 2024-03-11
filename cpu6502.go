package goNES

//constants representing the flags for the different values the status register can have
//each constant is a different bit, so they can be combined in the status register to signal what flags are set
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

//function signatures for the different functions associated with an operation
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

//constructor to create a cpu so it initializes the lookup table
func createCPU() *CPU6502 {
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
	cpu.setFlag(U, true)

	return &cpu
}

// links the cpu to a bus, should be the bus it's contained in
func (cpu *CPU6502) connectBus(ptr *Bus) {
	cpu.bus = ptr
}

// uses the bus to attempt reads and writes
func (cpu *CPU6502) write(addr uint16, data uint8) {
	cpu.bus.write(addr, data)
}

func (cpu CPU6502) read(addr uint16, readOnly bool) uint8 {
	return cpu.bus.read(addr, readOnly)
}

//sets the flag on the status registor for one of the values, if the flag is already set it unsets it
func (cpu *CPU6502) setFlag(flag uint8, set bool) {
	if set {
		cpu.status |= flag
	} else if cpu.getFlag(flag) {
		cpu.status ^= flag
	}

}

//checks if a specific flag is set on the status register
func (cpu CPU6502) getFlag(flag uint8) bool {
	return cpu.status&flag == flag
}

//gets data for the cpu and sets the fetchedData value to what was gotten
func (cpu *CPU6502) fetchData() {
	//checks if the addressing mode itself sets the fetchedData field, if it does don't attempt to fetch data again and overwrite it
	mode := cpu.instructions[cpu.opCode].modeType
	if mode == "IMP" || mode == "IMM" || mode == "ACC" {
		return
	}
	cpu.fetchedData = cpu.read(cpu.addrAbs, false)
}

//these four functions can occur at any point in operation, and will go after the current instruction is complete
//tells the cpu to advance one clock cycle
func (cpu *CPU6502) clock() {
	if cpu.cycles == 0 {
		cpu.opCode = cpu.read(cpu.pc, false)
		cpu.pc++
		cpu.cycles = cpu.instructions[cpu.opCode].cycles + cpu.instructions[cpu.opCode].op(cpu) + cpu.instructions[cpu.opCode].addrMode(cpu)
	}

	cpu.cycles--
}

//cpu interrupts

//resets the cpu
func (cpu *CPU6502) reset() {

}

//interrupt request, can be ignored depending on the interrupt flag of the status register
func (cpu *CPU6502) irq() {

}

//non maskable interrupt request, unable to be ignored
func (cpu *CPU6502) nmi() {

}

//addressing mode

//implied addressing, address is implicit in the opcode itself so nothing is needed
func IMP(cpu *CPU6502) uint8 {
	return 0
}

//accumulator addressing, data is retrieved from the accumulator
func ACC(cpu *CPU6502) uint8 {
	cpu.fetchedData = cpu.a
	return 0
}

//immediate addressing, the second byte in the instruction is the operand
func IMM(cpu *CPU6502) uint8 {
	cpu.fetchedData = cpu.read(cpu.pc, false)
	cpu.pc++
	return 0
}

//absolute addressing, the 2nd instruction byte is the lower byte of the address, the 3rd is the high bytes, combined to allow access to any point in memory
func ABS(cpu *CPU6502) uint8 {
	low := cpu.read(cpu.pc, false)
	cpu.pc++
	hi := cpu.read(cpu.pc, false)
	cpu.pc++

	//combine into one number
	var data uint16 = uint16(hi)<<8 | uint16(low)

	cpu.addrAbs = data
	return 0
}

//zero page addressing, the second byte is the offset from the first page of the memory
//there is a glitch, replicated from original hardware where values that should go to the next page wrap back around, but shouldn't occur in this zero page mode because it only reads in one byte
func ZPI(cpu *CPU6502) uint8 {
	byte := cpu.read(cpu.pc, false)
	cpu.pc++

	//wraps address that goes to the next page back around to the zero page
	cpu.addrAbs = uint16(byte) & 255
	return 0
}

//zero page addressing with X register offset, the second byte is the offset from the first page of the memory
//there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ZPX(cpu *CPU6502) uint8 {
	byte := cpu.read(cpu.pc, false)
	cpu.pc++
	byte += cpu.x

	//wraps address that goes to the next page back around to the zero page
	cpu.addrAbs = uint16(byte) & 255
	return 0
}

//zero page addressing with Y register offset, the second byte is the offset from the first page of the memory
//there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ZPY(cpu *CPU6502) uint8 {
	byte := cpu.read(cpu.pc, false)
	cpu.pc++
	byte += cpu.y

	//wraps address that goes to the next page back around to the zero page
	cpu.addrAbs = uint16(byte) & 255
	return 0
}

//absolute addressing with x register offset, the second byte is the offset from the first page of the memory
//there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ABX(cpu *CPU6502) uint8 {
	low := cpu.read(cpu.pc, false)
	cpu.pc++
	hi := cpu.read(cpu.pc, false)
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

//absolute addressing with Y register offset, the second byte is the offset from the first page of the memory
//there is a glitch, replicated from original hardware where values that should go to the next page wrap back around
func ABY(cpu *CPU6502) uint8 {
	low := cpu.read(cpu.pc, false)
	cpu.pc++
	hi := cpu.read(cpu.pc, false)
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

//relative addressing, the second byte is added to the program counter for branching
func REL(cpu *CPU6502) uint8 {
	offset := uint16(cpu.read(cpu.pc, false))
	cpu.pc++

	//if the number is 128 or greater unsigned, convert it to the 2's complement 16 digit version
	if offset >= 128 {
		offset |= 0xff00
	}

	cpu.addrRel = offset

	//will always increase the cycle, up to the function that calls this address mode to see if it goes to another page or not to determine the final increase
	return 1
}

//indexed indirect addressing, the second byte in the instruction is added to the X register without the carry, which is the low byte of the effective address, which is found on page zero, the following byte is the high byte, both on page zero
func IDX(cpu *CPU6502) uint8 {
	addr := cpu.read(cpu.pc, false)
	cpu.pc++

	//adds them as uint8 so the carry is discarded
	addr += cpu.x

	low := cpu.read(uint16(addr)&0x00ff, false)
	hi := cpu.read((uint16(addr)+1)&0x00ff, false)

	cpu.addrAbs = uint16(uint16(hi)<<8 | uint16(low))
	return 0
}

//indirect indexed addressing, the second byte in the instruction is added to the Y register, which is the low byte of the effective address, which is found on page zero, the following byte added with the carry of the last addition is the high byte, both on page zero
func IDY(cpu *CPU6502) uint8 {
	addr := cpu.read(cpu.pc, false)
	cpu.pc++

	low := cpu.read(uint16(addr&0x00ff), false)
	hi := cpu.read(uint16((addr+1)&0x00ff), false)

	//adds them as uint16 so the carry remains
	cpu.addrAbs = (uint16(hi)<<8 | uint16(low)) + uint16(cpu.y)

	//checks if the page increased, if so the clock cycle count needs to increase
	if uint16(hi)<<8 != cpu.addrAbs&0xff00 {
		return 1
	} else {
		return 0
	}
}

//absolute indirect addressing, the second byte in the instruction is the low byte of a memory location, the third instruction byte is the high byte, the data at that memory location is the low byte of the effective address, and the following byte is the effective high byte
//has a hardware bug where instead of going to the next page, it will wrap when forming the effective address, replicated here
func IND(cpu *CPU6502) uint8 {
	low := cpu.read(cpu.pc, false)
	cpu.pc++
	hi := cpu.read(cpu.pc, false)
	cpu.pc++

	addr := uint16(hi)<<8 | uint16(low)

	addrLow := cpu.read(addr, false)
	addrHi := cpu.read((addr + 1), false)

	//replicates the wrapping glitch, if the low byte is 0xff, then it wraps back to 0 on the same page instead of advancing
	if low == 0xff {
		addrHi = cpu.read(uint16(hi)<<8, false)
	}

	cpu.addrAbs = uint16(uint16(addrHi)<<8 | uint16(addrLow))
	return 0
}

//opcodes

//add with carry, from specified memory to accumulator
func ADC(cpu *CPU6502) uint8 {
	cpu.fetchData()

	carryFlag := 0

	if cpu.getFlag(C) {
		carryFlag = 1
	}

	res := uint16(cpu.fetchedData) + uint16(cpu.a) + uint16(carryFlag)

	//carry flag
	cpu.setFlag(C, res&0xff > 0)
	//zero flag
	cpu.setFlag(Z, res == 0)

	//overflow flag
	//checks if both values are negative and it wrapped to positive
	if (0x0080&cpu.a&cpu.fetchedData == 0x0080) && 0x0080&res == 0 {
		cpu.setFlag(V, true)
	}
	//checks if both values are positive and it wrapped to negative
	if ((^(cpu.a | cpu.fetchedData) & 0x0080) == 0x0080) && 0x0080&res == 0x0080 {
		cpu.setFlag(V, true)
	}

	//negative flag
	cpu.setFlag(N, 0x0080&res == 0x0080)

	//set the accumulator register to the new value
	cpu.a = uint8(res & 0x00ff)

	return 0
}
