package nes

//mapper for NES-NROM-128, NES-NROM-256

type Mapper000 struct {
	//how many banks of memory for each type of data
	PRGBanks uint8
	CHRBanks uint8
}

//accesses the CPU memory
//in mapper 0, if its a 32k PRG ram it maps CPU 0x8000 - 0xfff to ROM 0x0000 - 0x7fff, if it's 16k it mirrors itself to the other half, so CPU 0x8000 - 0xfff to ROM 0x0000 - 0x3fff
func (mapper Mapper000) CPUMapRead(addr uint16) (uint32, bool) {
	if addr >= 0x8000 && addr <= 0xffff {
		if mapper.PRGBanks > 1 {
			addr &= 0x7fff
		} else {
			addr &= 0x3fff
		}
		return uint32(addr), true
	}

	return 0x0000, false
}

func (mapper Mapper000) CPUMapWrite(addr uint16) (uint32, bool) {
	if addr >= 0x8000 && addr <= 0xffff {
		if mapper.PRGBanks > 1 {
			addr &= 0x7fff
		} else {
			addr &= 0x3fff
		}
		return uint32(addr), true
	}

	return 0x0000, false
}

//accesses the PPU memory
//in mapper 0 the CHR ROM is only 8K, so it's within the limits of the NES hardware already and doesn't need to be mapped
func (mapper Mapper000) PPUMapRead(addr uint16) (uint32, bool) {
	if addr >= 0x0000 && addr <= 0x1fff {
		return uint32(addr), true
	}

	return 0x0000, false
}

//since it is CHR ROM only here, it doesn't let you write
func (mapper Mapper000) PPUMapWrite(addr uint16) (uint32, bool) {

	return 0x0000, false
}
