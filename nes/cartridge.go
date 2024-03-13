package nes

type Cartridge struct {
	//what mapper is being used
	MapperID uint8
	//how many banks of memory for each type of data
	PRGBanks uint8
	CHRBanks uint8

	//stores the memory
	PRGMemory []uint8
	CHRMemory []uint8
}

func (cart *Cartridge) CPUWrite(addr uint16, data uint8) {

}

func (cart *Cartridge) CPURead(addr uint16, readOnly bool) uint8 {

	return 0x0000
}

// reads and writes from ppu memory
func (cart *Cartridge) PPUWrite(addr uint16, data uint8) {

}

func (cart Cartridge) PPURead(addr uint16, readOnly bool) uint8 {

	return 0x0000
}
