package nes

type Bus struct {
	CPU CPU6502
	PPU PPU2C02

	//count of how many clock cycles have passed in the system
	CycleCount uint32

	//current game cartridge/pak
	Cartridge *Cartridge
	//2 KB internal ram
	CPURAM [2048]uint8
}

// uses an uppercase letter at the beginning so its exported
func CreateBus() *Bus {
	bus := Bus{}
	bus.CPU = *CreateCPU()
	bus.CPU.ConnectBus(&bus)
	bus.PPU = *CreatePPU()
	bus.PPU.ConnectBus(&bus)
	return &bus
}

// reads and writes from memory
func (bus *Bus) CPUWrite(addr uint16, data uint8) {
	//if it's for the cartridge, execute the action and exit the function
	succ := bus.Cartridge.CPUWrite(addr, data)
	if succ {
		return
	}
	//area of memory for cpu
	if addr >= 0x0000 && addr <= 0x1FFF {
		//keeps the memory in range of the RAM dedicated to the cpu, mirrored every 2 KB
		bus.CPURAM[addr&0x07ff] = data
	}
	//writes to the dedicated ppu memory
	if addr >= 0x2000 && addr <= 0x3fff {
		bus.PPU.PPUWrite(addr&0x0007, data)
	}
}

func (bus Bus) CPURead(addr uint16, readOnly bool) uint8 {
	//if it's for the cartridge, execute the action and exit the function
	data, succ := bus.Cartridge.CPURead(addr, readOnly)
	if succ {
		return data
	}
	//area of memory for cpu
	if addr >= 0x0000 && addr <= 0x1FFF {
		//keeps the memory in range of the RAM dedicated to the cpu, mirrored every 2 KB
		return bus.CPURAM[addr&0x07ff]
	}
	if addr >= 0x2000 && addr <= 0x3fff {
		//keeps the memory in range of the RAM dedicated to the cpu, mirrored every 2 KB
		bus.PPU.PPURead(addr&0x0007, false)
	}

	//if there is an issue with the bounds just return 0
	return 0x0000
}

// reads and writes from ppu memory
func (bus *Bus) PPUWrite(addr uint16, data uint8) {

}

func (bus Bus) PPURead(addr uint16, readOnly bool) uint8 {

	return 0x0000
}

func (bus *Bus) InsertCartridge(cart *Cartridge) {
	bus.Cartridge = cart
	bus.PPU.ConnectCartridge(cart)
}

func (bus *Bus) Clock() {

}

func (bus *Bus) Reset() {
	bus.CPU.Reset()
	bus.CycleCount = 0
}
