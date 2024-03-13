package nes

type PPU2C02 struct {
	bus       *Bus
	Cartridge *Cartridge

	//name tables, lays out background tile data, NES has 2 physical name tables but 4 logical tables due to mirroring, 1 KB each
	NameTable [2][1024]uint8
	//holds the palette data
	Palette [32]uint8
}

func CreatePPU() *PPU2C02 {
	ppu := PPU2C02{}
	return &ppu
}

// links the cpu to a bus, should be the bus it's contained in
func (ppu *PPU2C02) ConnectBus(ptr *Bus) {
	ppu.bus = ptr
}

// connects cartridge to PPU and graphics memory
func (ppu *PPU2C02) ConnectCartridge(cart *Cartridge) {
	ppu.Cartridge = cart
}

// reads and writes from cpu memory
func (ppu *PPU2C02) CPUWrite(addr uint16, data uint8) {

}

func (ppu PPU2C02) CPURead(addr uint16, readOnly bool) uint8 {

	return 0x0000
}

// reads and writes from ppu memory
func (ppu *PPU2C02) PPUWrite(addr uint16, data uint8) {

}

func (ppu PPU2C02) PPURead(addr uint16, readOnly bool) uint8 {

	return 0x0000
}
