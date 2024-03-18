package nes

import (
	"fmt"
	"os"
)

// types of mirroring the cartridge can produce
type Mirror uint8

const (
	HORIZONTAL Mirror = iota
	VERTICAL
	ONESCREEN_LO
	ONESCREEN_HI
)

type Cartridge struct {
	//what mapper is being used
	MapperID uint8
	//what type of mirroring the cartridge uses
	Mirror Mirror
	//how many banks of memory for each type of data
	PRGBanks uint8
	CHRBanks uint8

	//stores the memory
	PRGMemory []uint8
	CHRMemory []uint8

	//the iNES header
	Header *iNESHeader

	//the mapper
	AddressMapper Mapper
}

// struct that represents the header of an iNES file
type iNESHeader struct {
	//first 16 bytes, goes from byte 0 to byte 15
	//constant header, "NES" followed by MS-DOS end of file character, 4 bytes
	name string
	//size of PRG ROM in 16 KB units, 4 bytes
	prgSize uint8
	//size of CHR ROM in 8KB units, 1 byte
	chrSize uint8
	//flags 6, mapper low nybble, mirroring, battery, trainer, 1 byte
	flag6 byte
	//flags 7, mapper high nybble, vs/playchoice, NES 1.0, 1 byte
	flag7 byte
	//flags 8, PRG-RAM size (rarely used extension), 1 byte
	flag8 byte
	//flags 9, TV system (rarely used extension), 1 byte
	flag9 byte
	//flags 10, TV system, PRG-RAM presence (unofficial, rarely used extension), 1 byte
	flag10 byte
	//unused, should be either 0s, sometimes people put their name here, 5 bytes
	pad []byte
}

func CreateCartridge(filename string) *Cartridge {
	cart := Cartridge{}

	//get header data and check if it was succesful
	header := readHeader(filename)
	if header == nil {
		return nil
	}

	cart.Header = header

	rom, err := os.Open(filename)

	//handle potential error in reading file
	if err != nil {
		return nil
	}
	//executed when the surrounding function finishes
	defer rom.Close()

	//gets the mapper ID from the header data
	cart.MapperID = (header.flag7 & 0xf0) | (header.flag6 >> 4)
	//gets how the cartridge sets up mirroring for the nametable
	//some mappers may use this header to set the mirroring in a different way, would need to overwrite this in the mapper
	//TODO make sure mappers can adjust this
	if header.flag6&0x01 == 0x01 {
		cart.Mirror = VERTICAL
	} else {
		cart.Mirror = HORIZONTAL
	}

	//where to begin reading, default is to only skip the first 16 header bytes and begin at byte
	begin := 16

	//checks if the rom file has trainer data, a depreciated mapping translation used in early NES emulators
	if header.flag6&4 == 4 {
		begin += 512
	}

	//there are types 0,1, and 2 of iNES files, but currently just uses type 1

	cart.PRGBanks = header.prgSize
	cart.CHRBanks = header.chrSize
	//set the size of the data to the size stated by the header
	cart.PRGMemory = make([]byte, int(cart.PRGBanks)*16384)
	cart.CHRMemory = make([]byte, int(cart.PRGBanks)*8192)

	//skip to the program data
	_, err = rom.Seek(int64(begin), 0)
	if err != nil {
		return nil
	}

	//read data into the virtual cartridge
	_, err = rom.Read(cart.PRGMemory)
	if err != nil {
		return nil
	}

	_, err = rom.Read(cart.CHRMemory)
	if err != nil {
		return nil
	}

	//sets the type of mapper to be used
	switch cart.MapperID {
	case 0:
		cart.AddressMapper = Mapper000{
			PRGBanks: cart.PRGBanks,
			CHRBanks: cart.CHRBanks,
		}
	}

	return &cart
}

func readHeader(filename string) *iNESHeader {
	rom, err := os.Open(filename)

	//handle potential error in reading file
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//executed when the surrounding function finishes
	defer rom.Close()

	//put the header data into a slice
	headerData := make([]byte, 16)
	_, err = rom.Read(headerData)
	if err != nil {
		return nil
	}

	head := iNESHeader{}

	//use the header data to set the values in the header struct
	head.name = string(headerData[0:4])
	head.prgSize = headerData[4]
	head.chrSize = headerData[5]
	head.flag6 = headerData[6]
	head.flag7 = headerData[7]
	head.flag8 = headerData[8]
	head.flag9 = headerData[9]
	head.flag10 = headerData[10]
	head.pad = headerData[11:]

	return &head
}

func (cart *Cartridge) CPUWrite(addr uint16, data uint8) bool {
	mapAddr, succ := cart.AddressMapper.CPUMapWrite(addr)
	//if the address was in the cartridge range, write the data and return that it was for the cartridge
	if succ {
		cart.PRGMemory[mapAddr] = data
		return true
	}

	return false
}

func (cart *Cartridge) CPURead(addr uint16, readOnly bool) (uint8, bool) {
	mapAddr, succ := cart.AddressMapper.CPUMapRead(addr)
	//if the address was in the cartridge range, return the data and return that it was for the cartridge
	if succ {
		return cart.PRGMemory[mapAddr], true
	}

	return 0x0000, false
}

// reads and writes from ppu memory
func (cart *Cartridge) PPUWrite(addr uint16, data uint8) bool {
	mapAddr, succ := cart.AddressMapper.PPUMapWrite(addr)
	//if the address was in the cartridge range, write the data and return that it was for the cartridge
	if succ {
		cart.CHRMemory[mapAddr] = data
		return true
	}

	return false
}

func (cart Cartridge) PPURead(addr uint16, readOnly bool) (uint8, bool) {
	mapAddr, succ := cart.AddressMapper.PPUMapRead(addr)
	//if the address was in the cartridge range, return the data and return that it was for the cartridge
	if succ {
		return cart.CHRMemory[mapAddr], true
	}

	return 0x0000, false
}
