package nes

//mapper interface, many different types of mappers exist on the NES, so this interface allows the use of mappers to be polymorphic
type Mapper interface {
	CPUMapRead(addr uint16) (uint32, bool)
	CPUMapWrite(addr uint16) (uint32, bool)
	PPUMapRead(addr uint16) (uint32, bool)
	PPUMapWrite(addr uint16) (uint32, bool)
}
