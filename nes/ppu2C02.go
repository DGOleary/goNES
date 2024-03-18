package nes

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type PPU2C02 struct {
	bus       *Bus
	Cartridge *Cartridge

	//name tables, lays out background tile data, NES has 2 physical name tables but 4 logical tables due to mirroring, 1 KB each
	NameTable [2][1024]uint8
	//pattern table, defines colors for sprites, 4 KB each
	PatternTable [2][4096]uint8
	//holds the palette data
	Palette [32]uint8

	//list of SDL color structs that correspond to the colors the NES can produce
	RGBPalette [64]sdl.Color
	//SDL renderer that allows the PPU to draw
	Renderer *sdl.Renderer

	//current PPU cycle
	Cycle uint32
	//current scanline
	Scanline int
	//status of current frame
	Complete bool
	//if a NMI should be fired during Vblank
	NMI bool

	//PPU registers, what allows the CPU to control the PPU
	//comment gives the meaning of each bit in left to right order, IE starting with the MSB
	PPUCTRL   uint8 //VPHB SINN, NMI enable, PPU master/slave, sprite height, background tile select, sprite tile select, increment mode, nametable select
	PPUMASK   uint8 //BGRs bMmG, color emphasis (BRG), sprite enable, background enable, sprite left column enable, background left column enable, greyscale
	PPUSTATUS uint8 //VSO- ----, vlbank, sprite 0 hit, sprite overflow, unused(- ----)
	OAMADDR   uint8 //aaaa aaaa, OAM read/write address
	OAMDATA   uint8 //dddd dddd, OAM data read/write
	PPUSCROLL uint8 //xxxx xxxx, fine scroll position (two writes; MSB, LSB)
	PPUADDR   uint8 //aaaa aaaa, PPU read/write address
	PPUDATA   uint8 //dddd dddd, PPU data read/write
	OAMDMA    uint8 //aaaa aaaa, OAM DMA high address

	//says if the low or high byte is being written to when the PPU is written to, 1 = low byte, 0 = high byte
	AddressByte uint8
	//buffer that holds the data that is being written to the PPU because the transfer is one cycle behind
	PPUBuffer uint8
	//the "loopy" register, an internal PPU register used to keep track of progress during a scanline
	//fine Y selector (3 bits), unused (1 bit), nametable select (2 bits), coarseY (5 bits) the vertical scroll position on screen (increments by 8 pixels for every +1), coarseX (5 bits) the horizontal scroll position on the screen (increments by 8 pixels for every +1)
	//the current VRAM address
	loopyVRAM uint16
	//temporary VRAM address, thought of as the address of the top left tile
	loopyTRAM uint16
	//horizontal offset
	fineX uint8

	//variables for loading BG data
	bgNextId     uint8
	bgNextAttrib uint8
	bgNextLSB    uint8
	bgNextMSB    uint8
	//shifter registers
	bgPatternShifterHigh   uint16
	bgPatternShifterLow    uint16
	bgAttributeShifterHigh uint16
	bgAttributeShifterLow  uint16

	//TODO test val
	frame uint64
}

func CreatePPU() *PPU2C02 {
	ppu := PPU2C02{}
	//another gross array fill
	//puts the RGB versions of NES colors in an array so SDL can draw to the screen
	//row 1
	ppu.RGBPalette[0] = sdl.Color{R: 98, G: 98, B: 98, A: 255}
	ppu.RGBPalette[1] = sdl.Color{R: 0, G: 31, B: 178, A: 255}
	ppu.RGBPalette[2] = sdl.Color{R: 36, G: 4, B: 200, A: 255}
	ppu.RGBPalette[3] = sdl.Color{R: 82, G: 0, B: 178, A: 255}
	ppu.RGBPalette[4] = sdl.Color{R: 115, G: 0, B: 118, A: 255}
	ppu.RGBPalette[5] = sdl.Color{R: 128, G: 0, B: 36, A: 255}
	ppu.RGBPalette[6] = sdl.Color{R: 115, G: 11, B: 0, A: 255}
	ppu.RGBPalette[7] = sdl.Color{R: 82, G: 40, B: 0, A: 255}
	ppu.RGBPalette[8] = sdl.Color{R: 36, G: 68, B: 0, A: 255}
	ppu.RGBPalette[9] = sdl.Color{R: 0, G: 87, B: 0, A: 255}
	ppu.RGBPalette[10] = sdl.Color{R: 0, G: 92, B: 0, A: 255}
	ppu.RGBPalette[11] = sdl.Color{R: 0, G: 83, B: 36, A: 255}
	ppu.RGBPalette[12] = sdl.Color{R: 0, G: 60, B: 118, A: 255}
	ppu.RGBPalette[13] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.RGBPalette[14] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.RGBPalette[15] = sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//row 2
	ppu.RGBPalette[16] = sdl.Color{R: 171, G: 171, B: 171, A: 255}
	ppu.RGBPalette[17] = sdl.Color{R: 13, G: 87, B: 255, A: 255}
	ppu.RGBPalette[18] = sdl.Color{R: 75, G: 48, B: 255, A: 255}
	ppu.RGBPalette[19] = sdl.Color{R: 138, G: 19, B: 255, A: 255}
	ppu.RGBPalette[20] = sdl.Color{R: 188, G: 8, B: 214, A: 255}
	ppu.RGBPalette[21] = sdl.Color{R: 210, G: 18, B: 105, A: 255}
	ppu.RGBPalette[22] = sdl.Color{R: 199, G: 46, B: 0, A: 255}
	ppu.RGBPalette[23] = sdl.Color{R: 157, G: 84, B: 0, A: 255}
	ppu.RGBPalette[24] = sdl.Color{R: 96, G: 123, B: 0, A: 255}
	ppu.RGBPalette[25] = sdl.Color{R: 32, G: 152, B: 0, A: 255}
	ppu.RGBPalette[26] = sdl.Color{R: 0, G: 163, B: 0, A: 255}
	ppu.RGBPalette[27] = sdl.Color{R: 0, G: 153, B: 66, A: 255}
	ppu.RGBPalette[28] = sdl.Color{R: 0, G: 125, B: 180, A: 255}
	ppu.RGBPalette[29] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.RGBPalette[30] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.RGBPalette[31] = sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//row 3
	ppu.RGBPalette[32] = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	ppu.RGBPalette[33] = sdl.Color{R: 83, G: 174, B: 255, A: 255}
	ppu.RGBPalette[34] = sdl.Color{R: 144, G: 133, B: 255, A: 255}
	ppu.RGBPalette[35] = sdl.Color{R: 211, G: 101, B: 255, A: 255}
	ppu.RGBPalette[36] = sdl.Color{R: 255, G: 87, B: 255, A: 255}
	ppu.RGBPalette[37] = sdl.Color{R: 255, G: 93, B: 207, A: 255}
	ppu.RGBPalette[38] = sdl.Color{R: 255, G: 119, B: 87, A: 255}
	ppu.RGBPalette[39] = sdl.Color{R: 250, G: 158, B: 0, A: 255}
	ppu.RGBPalette[40] = sdl.Color{R: 189, G: 199, B: 0, A: 255}
	ppu.RGBPalette[41] = sdl.Color{R: 122, G: 231, B: 0, A: 255}
	ppu.RGBPalette[42] = sdl.Color{R: 67, G: 246, B: 17, A: 255}
	ppu.RGBPalette[43] = sdl.Color{R: 38, G: 239, B: 126, A: 255}
	ppu.RGBPalette[44] = sdl.Color{R: 44, G: 213, B: 246, A: 255}
	ppu.RGBPalette[45] = sdl.Color{R: 78, G: 78, B: 78, A: 255}
	ppu.RGBPalette[46] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.RGBPalette[47] = sdl.Color{R: 0, G: 0, B: 0, A: 255}

	//row 4
	ppu.RGBPalette[48] = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	ppu.RGBPalette[49] = sdl.Color{R: 182, G: 255, B: 255, A: 255}
	ppu.RGBPalette[50] = sdl.Color{R: 206, G: 209, B: 255, A: 255}
	ppu.RGBPalette[51] = sdl.Color{R: 233, G: 195, B: 255, A: 255}
	ppu.RGBPalette[52] = sdl.Color{R: 255, G: 188, B: 255, A: 255}
	ppu.RGBPalette[53] = sdl.Color{R: 255, G: 189, B: 244, A: 255}
	ppu.RGBPalette[54] = sdl.Color{R: 255, G: 198, B: 195, A: 255}
	ppu.RGBPalette[55] = sdl.Color{R: 249, G: 210, B: 155, A: 255}
	ppu.RGBPalette[56] = sdl.Color{R: 233, G: 230, B: 129, A: 255}
	ppu.RGBPalette[57] = sdl.Color{R: 206, G: 244, B: 129, A: 255}
	ppu.RGBPalette[58] = sdl.Color{R: 182, G: 251, B: 154, A: 255}
	ppu.RGBPalette[59] = sdl.Color{R: 169, G: 250, B: 195, A: 255}
	ppu.RGBPalette[60] = sdl.Color{R: 169, G: 240, B: 244, A: 255}
	ppu.RGBPalette[61] = sdl.Color{R: 184, G: 184, B: 184, A: 255}
	ppu.RGBPalette[62] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.RGBPalette[63] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	return &ppu
}

// // puts the pattern tables to a texture to render
// // turns the pattern data into sprite textures, sprites are 2 bits per pixel, stored across 2 bitplanes, which are 8 bytes apart in memory
// // index is which pattern table you are indexing from, as there are 2 physical pattern tables
// func (ppu *PPU2C02) GetPatternTable(index uint8) {
// 	for tileY := 0; tileY < 16; tileY++ {
// 		for tileX := 0; tileX < 16; tileX++ {
// 			//transforms 2d coordinates to 1d coordinates, each tyle is 16*16 (256) bytes, so the offset is 256 ahead for each previous tile, and 16 each time the inner loop runs for each bitplane
// 			tileOffset := tileY*256 + tileX*16

// 			//goes through the 2 bitplanes of each memory section and combines them into an 8x8 sprite
// 			for row := 0; row < 8; row++ {
// 				//find the 2 bitplanes
// 				//which physical pattern table + how many tiles deep + what byte of the current tile
// 				lsb := ppu.PPURead(uint16(index)*uint16(4096)+uint16(tileOffset)+uint16(row), false)
// 				//which physical pattern table + how many tiles deep + what byte of the current tile + offset of msb from lsb (8)
// 				msb := ppu.PPURead(uint16(index)*uint16(4096)+uint16(tileOffset)+uint16(row)+uint16(8), false)

// 				//go through each bit in the bitplanes
// 				for col := 0; col < 8; col++ {
// 					//combine the bitplanes
// 					pixel := (lsb & 1) + (msb & 1)

// 					//adjust bit to check next pass
// 					lsb >>= 1
// 					msb >>= 1

// 				}
// 			}
// 		}
// 	}
// }

// returns the SDL color value for a pixel
func (ppu *PPU2C02) GetColorFromPalette(palette uint8, pixel uint8) sdl.Color {
	//multiply the palette id by 4, add the pixel as the offset into the palette, and mask it into the palette section of memory
	return ppu.RGBPalette[ppu.PPURead(0x3f00+uint16(palette<<2)+uint16(pixel), false)]
}

// links the cpu to a bus, should be the bus it's contained in
func (ppu *PPU2C02) ConnectBus(ptr *Bus) {
	ppu.bus = ptr
}

// links a renderer to the PPU
func (ppu *PPU2C02) ConnectRenderer(ptr *sdl.Renderer) {
	ppu.Renderer = ptr
}

// connects cartridge to PPU and graphics memory
func (ppu *PPU2C02) ConnectCartridge(cart *Cartridge) {
	ppu.Cartridge = cart
}

// checks the PPUCTRL register to determine if the automatic increment after a PPU read/write should be 1 or 32, 1 is the next byte horizontally, but to read vertically you need to advance 32, because the memory is sequential, it's laid out in 32 byte "rows"
func (ppu *PPU2C02) incrementMode() uint16 {
	if ppu.PPUCTRL&0x0004 == 0 {
		return uint16(1)
	} else {
		return uint16(32)
	}
}

func (ppu *PPU2C02) Reset() {
	ppu.Cycle = 0
	ppu.AddressByte = 0
	ppu.bgAttributeShifterHigh = 0
	ppu.bgAttributeShifterLow = 0
	ppu.bgNextAttrib = 0
	ppu.bgNextId = 0
	ppu.bgNextLSB = 0
	ppu.bgNextMSB = 0
	ppu.bgPatternShifterHigh = 0
	ppu.bgPatternShifterLow = 0
	ppu.fineX = 0
	ppu.loopyTRAM = 0
	ppu.loopyVRAM = 0
	ppu.PPUADDR = 0
	ppu.PPUCTRL = 0
	ppu.PPUDATA = 0
	ppu.PPUMASK = 0
	ppu.PPUSCROLL = 0
	ppu.PPUSTATUS = 0
}

// reads and writes from cpu memory
func (ppu *PPU2C02) CPUWrite(addr uint16, data uint8) {
	switch addr {
	case 0: //control
		ppu.PPUCTRL = data
		//sets the nametable bits
		//clear the bits
		ppu.loopyTRAM &= 0x73ff
		//set them to the value of the nametable bits in the PPUCTRL register
		ppu.loopyTRAM |= (uint16(ppu.PPUCTRL) & 0x0003) << 10
	case 1: //mask
		ppu.PPUMASK = data

	case 2: //status (can't be written to)

	case 3: //OAM address

	case 4: //OAM data

	case 5: //Scroll
		if ppu.AddressByte == 0 {
			//fineX is the bottom 3 bits of the data
			ppu.fineX = data & 0x0007
			//gets the coarseX
			//clears the courseX
			ppu.loopyTRAM &= 0xffe0
			//sets it
			//adjust the data 3 over to move it into place
			ppu.loopyTRAM |= uint16(data >> 3)
			ppu.AddressByte = 1
		} else {
			//fineY is the bottom 3 bits of the data
			//clear the old fineY
			ppu.loopyTRAM &= 0x1fff
			//sets the new fineY
			ppu.loopyTRAM |= (uint16(data&0x0007) << 12)
			//gets the coarseY
			//clears the coarseY
			ppu.loopyTRAM &= 0x7c1f
			//sets it
			ppu.loopyTRAM |= uint16(data>>3) << 5
			ppu.AddressByte = 0
		}
	case 6: //PPU address
		if ppu.AddressByte == 0 {
			//&'s the address with 0xff00 to clear out the old high byte
			ppu.loopyTRAM = ppu.loopyTRAM&0x00ff | uint16(data)<<8
			ppu.AddressByte = 1
		} else {
			//same as above but with low byte
			ppu.loopyTRAM = ppu.loopyTRAM&0xff00 | uint16(data)
			//update the VRAM address when 16 bits have been written to the TRAM
			ppu.loopyVRAM = ppu.loopyTRAM
			ppu.AddressByte = 0
		}
	case 7: //PPU data
		ppu.PPUWrite(ppu.loopyVRAM, data)
		//data is usually successive, so it writes automatically increments the address after a write
		ppu.loopyVRAM += ppu.incrementMode()
	}
}

func (ppu *PPU2C02) CPURead(addr uint16, readOnly bool) uint8 {
	returnData := uint8(0x00)

	switch addr {
	case 0: //control

	case 1: //mask

	case 2: //status reading from the status register does effect the value of it, it will clear bit 7 (vertical blank) along with the address latch bit (referred to in this program as the AddressByte) the value determing if the high or low byte is being written, some documentation says the unused bits is the old data buffer, no game uses that as far as I know but could be a source of bugs in the future (unlikely)
		returnData = ppu.PPUSTATUS
		ppu.PPUSTATUS &= 0x007f
		ppu.AddressByte = 0
	case 3: //OAM address

	case 4: //OAM data

	case 5: //Scroll

	case 6: //PPU address, no reason to read the address

	case 7: //PPU data
		returnData = ppu.PPUBuffer
		ppu.PPUBuffer = ppu.PPURead(ppu.loopyVRAM, false)

		//the palette memory for whatever hardware reason reads in the same clock cycle
		if addr >= 0x3f00 {
			//set the return to the current read address instead of buffering
			returnData = ppu.PPUBuffer
		}
		//also increments the address here for ease of use
		ppu.loopyVRAM += ppu.incrementMode()
	}

	return returnData
}

// reads and writes from ppu memory
func (ppu *PPU2C02) PPUWrite(addr uint16, data uint8) {

	if ppu.Cartridge.PPUWrite(addr, data) { //write into cartridge
		//the call to the function in the if statement writes the data if it's for the cartridge
	} else if addr >= 0x0000 && addr <= 0x1fff { //pattern table memory
		//usually rom but maybe could be ram
		ppu.PatternTable[(addr&0x1000)>>12][addr&0x0fff] = data
	} else if addr >= 0x2000 && addr <= 0x3eff { //name table memory
		//checks what type of mirroring is used
		//mirroring type is named after where you can find a duplicate of a physical nametable, ie horizontal means the nametable's duplicate is to it's left or right
		//"mirroring" really means duplication, the memory values are not reflected to the other side, the are exact copies, and retain changes made to the counterpart
		if ppu.Cartridge.Mirror == VERTICAL {
			//the index into the table gets masked with 0x3ff because the table has a size of 0x400 so because its inexed at 0...
			//table 0 and its mirror (top and bottom sides mirror eachother)
			if (addr >= 2000 && addr <= 0x23bf) || (addr >= 2800 && addr <= 0x2bff) {
				ppu.NameTable[0][addr&0x3ff] = data
			}
			//table 1 and its mirror
			if (addr >= 2400 && addr <= 0x27ff) || (addr >= 0x2c00 && addr <= 0x2fff) {
				ppu.NameTable[1][addr&0x3ff] = data
			}
		} else if ppu.Cartridge.Mirror == HORIZONTAL { //same thing except for a horizontal layout
			//table 0 and its mirror (left and right sides mirror eachother)
			if (addr >= 2000 && addr <= 0x23bf) || (addr >= 2400 && addr <= 0x27ff) {
				ppu.NameTable[0][addr&0x3ff] = data
			}
			//table 1 and its mirror
			if (addr >= 2800 && addr <= 0x2bff) || (addr >= 0x2c00 && addr <= 0x2fff) {
				ppu.NameTable[1][addr&0x3ff] = data
			}
		}
	} else if addr >= 0x3f00 && addr <= 0x3fff { //palette memory
		//mask the address for the palette index
		addr &= 0x001f
		//mirroring cases hardcoded in
		switch addr {
		case 0x0010:
			ppu.Palette[0x0000] = data
		case 0x0014:
			ppu.Palette[0x0004] = data
		case 0x0018:
			ppu.Palette[0x0008] = data
		case 0x001c:
			ppu.Palette[0x000c] = data
		}
	}
}

func (ppu PPU2C02) PPURead(addr uint16, readOnly bool) uint8 {
	data, read := ppu.Cartridge.PPURead(addr, readOnly)

	if read { //read into cartridge
		//if true, the data has already been read into data
	} else if addr >= 0x0000 && addr <= 0x1fff { //pattern table memory
		//gets which pattern table by checking the highest 4 bits and uses the rest of the address as the index
		data = ppu.PatternTable[(addr&0x1000)>>12][addr&0x0fff]
	} else if addr >= 0x2000 && addr <= 0x3eff { //name table memory
		//checks what type of mirroring is used
		//mirroring type is named after where you can find a duplicate of a physical nametable, ie horizontal means the nametable's duplicate is to it's left or right
		//"mirroring" really means duplication, the memory values are not reflected to the other side, the are exact copies, and retain changes made to the counterpart
		if ppu.Cartridge.Mirror == VERTICAL {
			//the index into the table gets masked with 0x3ff because the table has a size of 0x400 so because its inexed at 0...
			//table 0 and its mirror (top and bottom sides mirror eachother)
			if (addr >= 2000 && addr <= 0x23bf) || (addr >= 2800 && addr <= 0x2bff) {
				data = ppu.NameTable[0][addr&0x3ff]
			}
			//table 1 and its mirror
			if (addr >= 2400 && addr <= 0x27ff) || (addr >= 0x2c00 && addr <= 0x2fff) {
				data = ppu.NameTable[1][addr&0x3ff]
			}
		} else if ppu.Cartridge.Mirror == HORIZONTAL { //same thing except for a horizontal layout
			//table 0 and its mirror (left and right sides mirror eachother)
			if (addr >= 2000 && addr <= 0x23bf) || (addr >= 2400 && addr <= 0x27ff) {
				data = ppu.NameTable[0][addr&0x3ff]
			}
			//table 1 and its mirror
			if (addr >= 2800 && addr <= 0x2bff) || (addr >= 0x2c00 && addr <= 0x2fff) {
				data = ppu.NameTable[1][addr&0x3ff]
			}
		}
	} else if addr >= 0x3f00 && addr <= 0x3fff { //palette memory
		//mask the address for the palette index
		addr &= 0x001f
		//mirroring cases hardcoded in
		switch addr {
		case 0x0010:
			data = ppu.Palette[0x0000]
		case 0x0014:
			data = ppu.Palette[0x0004]
		case 0x0018:
			data = ppu.Palette[0x0008]
		case 0x001c:
			data = ppu.Palette[0x000c]
		}
	}

	return data
}

func (ppu *PPU2C02) Clock() {
	//anonymous functions for the drawing operations

	//this controls the increment for the X registers, making sure they wrap to the right values and increment to the next table when necessary
	scrollXIncrement := func() {
		//check if rendering the sprites is enabled
		if ppu.PPUMASK&0x0008 == 0x0008 || ppu.PPUMASK&0x0010 == 0x0010 {
			//check if a name table, which is 32 tiles long, needs to wrap around
			if ppu.loopyVRAM&0b11111 == 31 {
				//resets the coarseX to 0 and flips the nametable bit so it goes to the other table
				//checks if the lower bit of the nametable flags is 1
				if ppu.loopyVRAM&0x0400 == 0x0400 {
					ppu.loopyVRAM &= 0x7be0
				} else {
					//set it to 1
					ppu.loopyVRAM |= 0x0400
				}
			} else {
				//gets the value at the coarseX register
				coarseX := (ppu.loopyVRAM & 31) + 1
				//clears the coarseX
				ppu.loopyVRAM &= 0xffe0
				//sets it to the new value
				ppu.loopyVRAM |= coarseX
			}
		}
	}

	//same as above, but here for the Y registers
	scrollYIncrement := func() {
		//check if rendering the sprites is enabled
		if ppu.PPUMASK&0x0008 == 0x0008 || ppu.PPUMASK&0x0010 == 0x0010 {
			//check if fineY can be incremented
			//get the fineY and move it over to the right so you can see the actual value
			fineY := ((ppu.loopyVRAM & 0x7000) >> 12)
			if fineY < 7 {
				//add 1 to the number and move it back into the correct position
				fineY++
				//reset the fineY back to 0
				ppu.loopyVRAM &= 0x0fff
				//put in the fineY
				ppu.loopyVRAM |= fineY
			} else {
				//get courseY in a number
				coarseY := (ppu.loopyVRAM & 0x03e0) >> 5
				//if the it goes over the height limit of 8 pixels, increment to the next

				//reset the fineY back to 0
				ppu.loopyVRAM &= 0x0fff

				//table is 32x30 so when the coarseY reaches 29 it needs to roll over
				if coarseY == 29 {
					//reset coarseY to 0
					ppu.loopyVRAM &= 0x7c1f
					//flip the upper nametable bit
					//if it's 1
					if ppu.loopyVRAM&0x0800 == 0x0800 {
						//set the bit to 0
						ppu.loopyVRAM &= 0x77ff
					} else {
						//set it to 1
						ppu.loopyVRAM |= 0x0800
					}
				} else if coarseY == 31 {
					//30 and 31 are attribute memory, reset it if it goes into here
					//reset coarseY to 0
					ppu.loopyVRAM &= 0x7c1f
				} else {
					//no wrapping, so just increment
					ppu.loopyVRAM |= ((coarseY + 1) << 5)
				}
			}
		}
	}

	//move the temporary values into the actual VRAM registers
	transferXRegisters := func() {
		//check if rendering the sprites is enabled
		if ppu.PPUMASK&0x0008 == 0x0008 || ppu.PPUMASK&0x0010 == 0x0010 {
			//reset the lower nametable bit to 0
			ppu.loopyVRAM &= 0x7bff
			//get the lower nametable bit from the TRAM
			ppu.loopyVRAM |= (ppu.loopyTRAM & 0x0400)
			//clears the coarseX
			ppu.loopyVRAM &= 0xffe0
			//gets the coarseX from the TRAM
			ppu.loopyVRAM |= (ppu.loopyTRAM & 31)
		}
	}

	transferYRegisters := func() {
		//check if rendering the sprites is enabled
		if ppu.PPUMASK&0x0008 == 0x0008 || ppu.PPUMASK&0x0010 == 0x0010 {
			//reset the fineY back to 0
			ppu.loopyVRAM &= 0x0fff
			//get the TRAM fineY
			ppu.loopyVRAM |= (ppu.loopyTRAM & 0x7000)
			//reset the nametable bit to 0
			ppu.loopyVRAM &= 0x77ff
			//get it from the TRAM
			ppu.loopyVRAM |= (ppu.loopyTRAM & 0x0800)
			//set the coarseY to 0
			ppu.loopyVRAM &= 0x7c1f
			//get it from the TRAM
			ppu.loopyVRAM |= (ppu.loopyTRAM & 0x03e0)
		}
	}

	loadBackgroundShifters := func() {
		//each PPU tick one pixel is drawn, these shifters move along with the current pixel and the 15 next, the MSB is being drawn, but the fineX register can utilize other pixels so the entire register needs to be updated
		//moves the current tiles into the MSB and the next tile into the LSB
		ppu.bgPatternShifterLow = (ppu.bgPatternShifterLow & 0xff00) | uint16(ppu.bgNextLSB)
		ppu.bgPatternShifterHigh = (ppu.bgPatternShifterHigh & 0x00ff) | uint16(ppu.bgNextMSB)

		//attribute bits change every 8 pixels, but to synchronize them they can be blown up to an entire byte of the value
		if ppu.bgNextAttrib&1 == 1 {
			ppu.bgAttributeShifterLow = (ppu.bgAttributeShifterLow & 0xff00) | 0x00ff
		} else {
			ppu.bgAttributeShifterLow = (ppu.bgAttributeShifterLow & 0xff00) | 0x0000
		}

		if ppu.bgNextAttrib&2 == 2 {
			ppu.bgAttributeShifterHigh = (ppu.bgAttributeShifterHigh & 0xff00) | 0x00ff
		} else {
			ppu.bgAttributeShifterHigh = (ppu.bgAttributeShifterHigh & 0xff00) | 0x0000
		}
	}

	//every tick the shifters move theri contents by 1 to draw the next pixel
	updateShifters := func() {
		//check if render background is enabled
		if ppu.PPUMASK&0x0008 == 0x0008 {
			ppu.bgPatternShifterLow <<= 1
			ppu.bgPatternShifterHigh <<= 1

			ppu.bgAttributeShifterLow <<= 1
			ppu.bgAttributeShifterHigh <<= 1
		}
	}

	//actions that apply to most scanlines
	if ppu.Scanline >= -1 && ppu.Scanline < 240 {
		//skip this frame
		if ppu.Scanline == 0 && ppu.Cycle == 0 {
			ppu.Cycle = 1
		}

		//this is when Vblank ends, back to the top left of the screen
		if ppu.Scanline == -1 && ppu.Cycle == 1 {
			//sets the Vblank bit to off
			ppu.PPUSTATUS &= 0x7f
		}

		if (ppu.Cycle >= 2 && ppu.Cycle < 258) || (ppu.Cycle >= 321 && ppu.Cycle < 338) {
			updateShifters()

			switch (ppu.Cycle - 1) % 8 {
			case 0:
				//when a tile has looped, load new shifters
				loadBackgroundShifters()
				ppu.bgNextId = ppu.PPURead(0x2000|(ppu.loopyVRAM&0x0fff), false)
			case 2:
				ppu.bgNextAttrib = ppu.PPURead(0x23c0|((ppu.loopyVRAM&2)<<11)|((ppu.loopyVRAM&1)<<10)|(((ppu.loopyVRAM&0b1111100000)>>2)<<3)|((ppu.loopyVRAM&0b11111)>>2), false)
				if (ppu.loopyVRAM&0b1111100000)&0x0002 == 0x0002 {
					ppu.bgNextAttrib >>= 4
				}
				if (ppu.loopyVRAM&0b11111)&0x0002 == 0x0002 {
					ppu.bgNextAttrib >>= 2
				}
				ppu.bgNextAttrib &= 0x03
			case 4:
				ppu.bgNextLSB = ppu.PPURead((uint16(ppu.PPUCTRL&0b10000)<<12)+(uint16(ppu.bgNextId)<<4)+(ppu.loopyVRAM&0x7000), false)
			case 6:
				ppu.bgNextMSB = ppu.PPURead((uint16(ppu.PPUCTRL&0b10000)<<12)+(uint16(ppu.bgNextId)<<4)+((ppu.loopyVRAM&0x7000)+8), false)
			case 7:
				//done with 8 pixels, go to next 8
				scrollXIncrement()
			}
		}
		//done with a row, go to next Y
		if ppu.Cycle == 256 {
			scrollYIncrement()
		}

		//reset X registers
		if ppu.Cycle == 256 {
			transferXRegisters()
		}

		//ready for new frame
		if ppu.Scanline == -1 && ppu.Cycle >= 280 && ppu.Cycle < 305 {
			transferYRegisters()
		}
	}

	if ppu.Scanline == 240 {
		//here to signify that scanline 240 takes no actions
	}

	//this is when Vblank starts
	if ppu.Scanline == 241 && ppu.Cycle == 1 {
		//sets the Vblank bit to on
		ppu.PPUSTATUS |= 0x80
		//checks if the NMI bit is on, if so trigger it to be fired
		if ppu.PPUCTRL&0x80 == 0x80 {
			ppu.NMI = true
		}
	}

	//the 2 bit pixel being rendered
	bgPixel := uint8(0)
	//the 3 bit index of the palette
	bgPalette := uint8(0)

	//check if background rendering is enabled
	if ppu.PPUMASK&0x0008 == 0x0008 {
		//selects the relevant bit using fineX offset
		bit := uint16(0x8000) >> ppu.fineX

		//select the bitplane bixesl by using the pattern shifter
		pixelPlane0 := 0
		pixelPlane1 := 0
		if (ppu.bgPatternShifterLow & bit) > 0 {
			pixelPlane0 = 1
		}

		if (ppu.bgPatternShifterHigh & bit) > 0 {
			pixelPlane1 = 1
		}

		bgPalette = (uint8(pixelPlane1) << 1) | uint8(pixelPlane0)

		if ppu.frame >= 4 {
			fmt.Println(ppu.frame)
		}
	}

	//fmt.Println(bgPalette)
	//fmt.Println(bgPixel)
	pixelColor := ppu.GetColorFromPalette(bgPalette, bgPixel)
	//fmt.Println(pixelColor)
	//fmt.Println(ppu.Cycle - 1)
	//fmt.Println(ppu.Scanline)
	ppu.Renderer.SetDrawColor(pixelColor.R, pixelColor.G, pixelColor.B, pixelColor.A)
	ppu.Renderer.DrawPoint(int32(ppu.Cycle)-1, int32(ppu.Scanline))
	//ppu.Renderer.Present()

	// if ppu.Scanline > 0 && (ppu.Scanline%50 == 0 && ppu.Cycle%50 == 0) {
	// 	temp := 0
	// 	temp++
	// }
	ppu.Cycle++
	//each scanline lasts for 341 PPU cycles
	if ppu.Cycle >= 341 {
		ppu.Cycle = 0
		ppu.Scanline++

		//how many scanlines per frame
		if ppu.Scanline >= 261 {
			ppu.frame++
			ppu.Scanline = -1
			ppu.Complete = true
			//TODO make accurate time
		}
	}
}
