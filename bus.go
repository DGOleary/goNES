package goNES

type Bus struct {
	cpu CPU6502

	//full accessable memory in all of NES
	ram [64 * 1024]uint8
}

func createBus() *Bus {
	bus := Bus{}
	bus.cpu.connectBus(&bus)
	return &bus
}

// reads and writes from memory
func (bus *Bus) write(addr uint16, data uint8) {
	if addr >= 0x0000 && addr <= 0xFFFF {
		bus.ram[addr] = data
	}
}

func (bus Bus) read(addr uint16, readOnly bool) uint8 {
	if addr >= 0x0000 && addr <= 0xFFFF {
		return bus.ram[addr]
	}

	//if there is an issue with the bounds just return 0
	return 0x0000
}
