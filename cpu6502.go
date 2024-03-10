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
	addrMode addressingMode
	cycles   uint8
}

//constructor to create a cpu so it initializes the lookup table
func createCPU() *CPU6502 {
	cpu := CPU6502{}
	//ugly, gross, disgusting, bad, not good, but it initializes the entire table
	cpu.instructions = [256]Instruction{
		{name: "BRK", op: BRK, addrMode: IMP, cycles: 7}, {name: "ORA", op: ORA, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, addrMode: ZPI, cycles: 3}, {name: "ASL", op: ASL, addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "PHP", op: PHP, addrMode: IMP, cycles: 3}, {name: "ORA", op: ORA, addrMode: IMM, cycles: 2}, {name: "ASL", op: ASL, addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, addrMode: ABS, cycles: 4}, {name: "ASL", op: ASL, addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BPL", op: BPL, addrMode: REL, cycles: 2}, {name: "ORA", op: ORA, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, addrMode: ZPX, cycles: 4}, {name: "ASL", op: ASL, addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CLC", op: CLC, addrMode: IMP, cycles: 2}, {name: "ORA", op: ORA, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ORA", op: ORA, addrMode: ABX, cycles: 4}, {name: "ASL", op: ASL, addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "JSR", op: JSR, addrMode: ABS, cycles: 6}, {name: "AND", op: AND, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "BIT", op: BIT, addrMode: ZPI, cycles: 3}, {name: "AND", op: AND, addrMode: ZPI, cycles: 3}, {name: "ROL", op: ROL, addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "PLP", op: PLP, addrMode: IMP, cycles: 4}, {name: "AND", op: AND, addrMode: IMM, cycles: 2}, {name: "ROL", op: ROL, addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "BIT", op: BIT, addrMode: ABS, cycles: 4}, {name: "AND", op: AND, addrMode: ABS, cycles: 4}, {name: "ROL", op: ROL, addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BMI", op: BMI, addrMode: REL, cycles: 2}, {name: "AND", op: AND, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "AND", op: AND, addrMode: ZPX, cycles: 4}, {name: "ROL", op: ROL, addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "SEC", op: SEC, addrMode: IMP, cycles: 2}, {name: "AND", op: AND, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "AND", op: AND, addrMode: ABX, cycles: 4}, {name: "ROL", op: ROL, addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "RTI", op: RTI, addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, addrMode: ZPI, cycles: 3}, {name: "LSR", op: LSR, addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "PHA", op: PHA, addrMode: IMP, cycles: 3}, {name: "EOR", op: EOR, addrMode: IMM, cycles: 2}, {name: "LSR", op: LSR, addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "JMP", op: JMP, addrMode: ABS, cycles: 3}, {name: "EOR", op: EOR, addrMode: ABS, cycles: 4}, {name: "LSR", op: LSR, addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BVC", op: BVC, addrMode: REL, cycles: 2}, {name: "EOR", op: EOR, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, addrMode: ZPX, cycles: 4}, {name: "LSR", op: LSR, addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CLI", op: CLI, addrMode: IMP, cycles: 2}, {name: "EOR", op: EOR, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "EOR", op: EOR, addrMode: ABX, cycles: 4}, {name: "LSR", op: LSR, addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "RTS", op: RTS, addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, addrMode: ZPI, cycles: 3}, {name: "ROR", op: ROR, addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "PLA", op: PLA, addrMode: IMP, cycles: 4}, {name: "ADC", op: ADC, addrMode: IMM, cycles: 2}, {name: "ROR", op: ROR, addrMode: ACC, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "JMP", op: JMP, addrMode: IND, cycles: 5}, {name: "ADC", op: ADC, addrMode: ABS, cycles: 4}, {name: "ROR", op: ROR, addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BVS", op: BVS, addrMode: REL, cycles: 2}, {name: "ADC", op: ADC, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, addrMode: ZPX, cycles: 4}, {name: "ROR", op: ROR, addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "SEI", op: SEI, addrMode: IMP, cycles: 2}, {name: "ADC", op: ADC, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "ADC", op: ADC, addrMode: ABX, cycles: 4}, {name: "ROR", op: ROR, addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "STA", op: STA, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "STY", op: STY, addrMode: ZPI, cycles: 3}, {name: "STA", op: STA, addrMode: ZPI, cycles: 3}, {name: "STX", op: STX, addrMode: ZPI, cycles: 3}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "DEY", op: DEY, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "TXA", op: TXA, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "STY", op: STY, addrMode: ABS, cycles: 4}, {name: "STA", op: STA, addrMode: ABS, cycles: 4}, {name: "STX", op: STX, addrMode: ABS, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BCC", op: BCC, addrMode: REL, cycles: 2}, {name: "STA", op: STA, addrMode: IDY, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "STY", op: STY, addrMode: ZPX, cycles: 4}, {name: "STA", op: STA, addrMode: ZPX, cycles: 4}, {name: "STX", op: STX, addrMode: ZPY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "TYA", op: TYA, addrMode: IMP, cycles: 2}, {name: "STA", op: STA, addrMode: ABY, cycles: 5}, {name: "TXS", op: TXS, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "STA", op: STA, addrMode: ABX, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "LDY", op: LDY, addrMode: IMM, cycles: 2}, {name: "LDA", op: LDA, addrMode: IDX, cycles: 6}, {name: "LDX", op: LDX, addrMode: IMM, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, addrMode: ZPI, cycles: 3}, {name: "LDA", op: LDA, addrMode: ZPI, cycles: 3}, {name: "LDX", op: LDX, addrMode: ZPI, cycles: 3}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "TAY", op: TAY, addrMode: IMP, cycles: 2}, {name: "LDA", op: LDA, addrMode: IMM, cycles: 2}, {name: "TAX", op: TAX, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, addrMode: ABS, cycles: 4}, {name: "LDA", op: LDA, addrMode: ABS, cycles: 4}, {name: "LDX", op: LDX, addrMode: ABS, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BCS", op: BCS, addrMode: REL, cycles: 2}, {name: "LDA", op: LDA, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, addrMode: ZPX, cycles: 4}, {name: "LDA", op: LDA, addrMode: ZPX, cycles: 4}, {name: "LDX", op: LDX, addrMode: ZPY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CLV", op: CLV, addrMode: IMP, cycles: 2}, {name: "LDA", op: LDA, addrMode: ABY, cycles: 4}, {name: "TSX", op: TSX, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "LDY", op: LDY, addrMode: ABX, cycles: 4}, {name: "LDA", op: LDA, addrMode: ABX, cycles: 4}, {name: "LDX", op: LDX, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "CPY", op: CPY, addrMode: IMM, cycles: 2}, {name: "CMP", op: CMP, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CPY", op: CPY, addrMode: ZPI, cycles: 3}, {name: "CMP", op: CMP, addrMode: ZPI, cycles: 3}, {name: "DEC", op: DEC, addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "INY", op: INY, addrMode: IMP, cycles: 2}, {name: "CMP", op: CMP, addrMode: IMM, cycles: 2}, {name: "DEX", op: DEX, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CPY", op: CPY, addrMode: ABS, cycles: 4}, {name: "CMP", op: CMP, addrMode: ABS, cycles: 4}, {name: "DEC", op: DEC, addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BNE", op: BNE, addrMode: REL, cycles: 2}, {name: "CMP", op: CMP, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CMP", op: CMP, addrMode: ZPX, cycles: 4}, {name: "DEC", op: DEC, addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CLD", op: CLD, addrMode: IMP, cycles: 2}, {name: "CMP", op: CMP, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CMP", op: CMP, addrMode: ABX, cycles: 4}, {name: "DEC", op: DEC, addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "CPX", op: CPX, addrMode: IMM, cycles: 2}, {name: "SBC", op: SBC, addrMode: IDX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CPX", op: CPX, addrMode: ZPI, cycles: 3}, {name: "SBC", op: SBC, addrMode: ZPI, cycles: 3}, {name: "INC", op: INC, addrMode: ZPI, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "INX", op: INX, addrMode: IMP, cycles: 2}, {name: "SBC", op: SBC, addrMode: IMM, cycles: 2}, {name: "NOP", op: NOP, addrMode: IMP, cycles: 2}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "CPX", op: CPX, addrMode: ABS, cycles: 4}, {name: "SBC", op: SBC, addrMode: ABS, cycles: 4}, {name: "INC", op: INC, addrMode: ABS, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
		{name: "BEQ", op: BEQ, addrMode: REL, cycles: 2}, {name: "SBC", op: SBC, addrMode: IDY, cycles: 5}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "SBC", op: SBC, addrMode: ZPX, cycles: 4}, {name: "INC", op: INC, addrMode: ZPX, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "SED", op: SED, addrMode: IMP, cycles: 2}, {name: "SBC", op: SBC, addrMode: ABY, cycles: 4}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6}, {name: "SBC", op: SBC, addrMode: ABX, cycles: 4}, {name: "INC", op: INC, addrMode: ABX, cycles: 7}, {name: "NEX", op: NEX, addrMode: IMP, cycles: 6},
	}

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

