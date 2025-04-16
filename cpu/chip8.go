package cpu

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/Sounak10/chip8emu/config"
)

const (
	QUIT = iota

	RUNNING

	PAUSED
)

type Instruction struct {
	OpCode uint16
	NNN    uint16
	NN     uint8
	N      uint8
	X      uint8
	Y      uint8
}

type Chip8 struct {
	State        int
	Ram          [4096]byte
	Display      [64 * 32]bool //Originally 64x32 pixels
	Stack        [16]uint16    //Subroutine stack
	StackPointer uint8         //Stack pointer
	V            [16]uint8     //Registers V0 to VF
	I            uint16        //Index register
	PC           uint16        //Program counter
	DelayTimer   uint8         //Delay timer decrements at 60Hz
	SoundTimer   uint8         //Sound timer decrements at 60Hz
	KeyPad       [16]bool      //Keypad state
	ROMName      string        //Name of the loaded ROM
	Instruction  Instruction   //Current instruction
	Draw         bool          //Flag to indicate if the screen needs to be redrawn
}

// NewChip8 initializes a new Chip8 instance, loads the ROM into memory, and sets up the initial state.

func NewChip8(romName string) (*Chip8, error) {
	// Initialize the random seed

	chip8 := &Chip8{}
	var entrypoint uint16 = 0x200 // Program starts at 0x200
	maxRomSize := len(chip8.Ram) - int(entrypoint)

	var font = []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, //0
		0x20, 0x60, 0x20, 0x20, 0x70, //1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
		0x90, 0x90, 0xF0, 0x10, 0x10, //4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
		0xF0, 0x10, 0x20, 0x40, 0x40, //7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
		0xF0, 0x90, 0xF0, 0x90, 0x90, //A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
		0xF0, 0x80, 0x80, 0x80, 0xF0, //C
		0xE0, 0x90, 0x90, 0x90, 0xE0, //D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
		0xF0, 0x80, 0xF0, 0x80, 0x80, //F
	}

	chip8.State = RUNNING
	chip8.PC = entrypoint
	chip8.StackPointer = 0
	chip8.ROMName = romName

	// Load font into memory
	for i := range font {
		chip8.Ram[i] = font[i]
	}

	//Load ROM into memory
	rom, err := os.ReadFile(romName)
	if err != nil {
		return nil, err
	}
	if len(rom) > maxRomSize {
		return nil, os.ErrInvalid
	}
	for i := range rom {
		chip8.Ram[entrypoint+uint16(i)] = rom[i]
	}

	return chip8, nil
}

func (c *Chip8) EmulateInstruction(config *config.Config) {
	// Fetch the instruction from memory
	c.Instruction.OpCode = uint16(c.Ram[c.PC])<<8 | uint16(c.Ram[c.PC+1])
	c.PC += 2 // Increment the program counter by 2

	c.Instruction.NNN = c.Instruction.OpCode & 0x0FFF
	c.Instruction.NN = uint8(c.Instruction.OpCode & 0x0FF)
	c.Instruction.N = uint8(c.Instruction.OpCode & 0x0F)
	c.Instruction.X = uint8((c.Instruction.OpCode >> 8) & 0x0F)
	c.Instruction.Y = uint8((c.Instruction.OpCode >> 4) & 0x0F)
	// Decode the instruction
	switch (c.Instruction.OpCode >> 12) & 0x0F {
	case 0x0:
		if c.Instruction.NN == 0xE0 {

			// Clear the screen
			for i := range c.Display {
				c.Display[i] = false
			}
		} else if c.Instruction.NN == 0xEE {
			// Return from subroutine
			c.StackPointer--               // Decrement the stack pointer
			c.PC = c.Stack[c.StackPointer] // Set the program counter to the address at the top of the stack

		} else {
			// Call the machine code routine at NNN
			// This is not implemented in this emulator

		}

	case 0x01:
		// Jump to address NNN
		c.PC = c.Instruction.NNN

	case 0x02:
		// Call subroutine at NNN
		c.Stack[c.StackPointer] = c.PC // Save the current program counter
		c.StackPointer++               // Increment the stack pointer
		c.PC = c.Instruction.NNN       // Jump to the subroutine

	case 0x03:
		// Skip the next instruction if Vx == NN
		if c.V[c.Instruction.X] == c.Instruction.NN {
			c.PC += 2 // Skip the next instruction
		}

	case 0x04:
		// Skip the next instruction if Vx != NN
		if c.V[c.Instruction.X] != c.Instruction.NN {
			c.PC += 2 // Skip the next instruction
		}

	case 0x05:
		// Skip the next instruction if Vx == Vy
		if c.Instruction.N != 0 {
			fmt.Println("Unknown instruction: ", c.Instruction.OpCode)
			break
		}
		if c.V[c.Instruction.X] == c.V[c.Instruction.Y] {
			c.PC += 2 // Skip the next instruction
		}

	case 0x06:
		// Set Vx to NN
		c.V[c.Instruction.X] = c.Instruction.NN

	case 0x07:
		// Set Vx to Vx + NN
		c.V[c.Instruction.X] += c.Instruction.NN

	case 0x08:
		// Handle arithmetic operations
		switch c.Instruction.N {
		case 0x0:
			// Set Vx to Vy
			c.V[c.Instruction.X] = c.V[c.Instruction.Y]

		case 0x1:
			// Set Vx to Vx OR Vy
			c.V[c.Instruction.X] |= c.V[c.Instruction.Y]

		case 0x2:
			// Set Vx to Vx AND Vy
			c.V[c.Instruction.X] &= c.V[c.Instruction.Y]

		case 0x3:
			// Set Vx to Vx XOR Vy
			c.V[c.Instruction.X] ^= c.V[c.Instruction.Y]

		case 0x4:
			// Add Vy to Vx, set VF to 1 if there's a carry
			if c.V[c.Instruction.X] > 0xFF-c.V[c.Instruction.Y] {
				c.V[0xF] = 1 // Set carry flag
			}
			c.V[c.Instruction.X] += c.V[c.Instruction.Y]

		case 0x5:
			// Subtract Vy from Vx, set VF to 0 if there's a borrow
			if c.V[c.Instruction.X] >= c.V[c.Instruction.Y] {
				c.V[0xF] = 1 // Set no borrow
			} else {
				c.V[0xF] = 0 // Set borrow
			}
			c.V[c.Instruction.X] -= c.V[c.Instruction.Y]

		case 0x6:
			// Shift Vx right by 1, set VF to the least significant bit of Vx before the shift
			c.V[0xF] = c.V[c.Instruction.X] & 0x1 // Store the least significant bit
			c.V[c.Instruction.X] >>= 1

		case 0x7:
			// Set Vx to Vy - Vx, set VF to 0 if there's a borrow
			if c.V[c.Instruction.Y] >= c.V[c.Instruction.X] {
				c.V[0xF] = 1 // Set no borrow
			} else {
				c.V[0xF] = 0 // Set borrow
			}
			c.V[c.Instruction.X] = c.V[c.Instruction.Y] - c.V[c.Instruction.X]

		case 0xE:
			// Shift Vx left by 1, set VF to the most significant bit of Vx before the shift
			c.V[0xF] = (c.V[c.Instruction.X] & 0x80) >> 7 // Store the most significant bit
			c.V[c.Instruction.X] <<= 1

		default:
			fmt.Println("Unknown instruction: ", c.Instruction.OpCode)
		}

	case 0x09:
		// Skip the next instruction if Vx != Vy
		if c.Instruction.N != 0 {
			fmt.Println("Unknown instruction: ", c.Instruction.OpCode)
			break
		}
		if c.V[c.Instruction.X] != c.V[c.Instruction.Y] {
			c.PC += 2
		}
	case 0x0A:
		// Set I to NNN
		c.I = c.Instruction.NNN

	case 0x0B:
		// Jump to address NNN + V0
		// This is a relative jump
		c.PC = c.Instruction.NNN + uint16(c.V[0])

	case 0x0C:
		// Set Vx to a random number AND NN
		// Generate a random number between 0 and 255
		randomValue := uint8(rand.IntN(256))
		c.V[c.Instruction.X] = randomValue & c.Instruction.NN

	case 0x0D:
		// Draw sprite at (x, y) with height N
		// The sprite is stored in memory starting at location I
		// The screen is 64x32 pixels, and each sprite is 8 pixels wide
		X_pos := c.V[c.Instruction.X] % uint8(config.WindowWidth)
		Y_pos := c.V[c.Instruction.Y] % uint8(config.WindowHeight)
		c.V[0xF] = 0 // Clear the carry flag

		for i := 0; i < int(c.Instruction.N); i++ {
			sprite := c.Ram[c.I+uint16(i)]

			for j := 0; j < 8; j++ {
				// Check if current bit of sprite byte is 1
				if (sprite & (0x80 >> j)) != 0 {
					// Calculate the position in the display array
					pos := int(X_pos+uint8(j)) + (int(Y_pos+uint8(i)) * int(config.WindowWidth))

					// Make sure we don't go out of bounds
					if pos >= 0 && pos < len(c.Display) {
						// If the pixel is already set, set VF to 1 (collision)
						if c.Display[pos] {
							c.V[0xF] = 1
						}
						// XOR the pixel
						c.Display[pos] = !c.Display[pos]
					}
				}
			}
		}
		c.Draw = true

	case 0x0E:
		switch c.Instruction.NN {
		case 0x9E:
			// Skip the next instruction if the key stored in Vx is pressed
			if c.KeyPad[c.V[c.Instruction.X]] {
				c.PC += 2 // Skip the next instruction
			}

		case 0xA1:
			// Skip the next instruction if the key stored in Vx is not pressed
			if !c.KeyPad[c.V[c.Instruction.X]] {
				c.PC += 2 // Skip the next instruction
			}
		}

	case 0x0F:
		switch c.Instruction.NN {
		case 0x0A:
			// Wait for a key press and store the value in Vx
			keyPressed := false
			for i := 0; i < len(c.KeyPad); i++ {
				if c.KeyPad[i] {
					c.V[c.Instruction.X] = uint8(i) // Store the key value in Vx
					keyPressed = true
					break
				}
			}
			if !keyPressed {
				c.PC -= 2 // Decrement the program counter to wait for a key press
			}
		case 0x1E:
			// Add Vx to I
			c.I += uint16(c.V[c.Instruction.X])

		case 0x07:
			// Set Vx to the value of the delay timer
			c.V[c.Instruction.X] = c.DelayTimer

		case 0x15:
			// Set the delay timer to Vx
			c.DelayTimer = c.V[c.Instruction.X]

		case 0x18:
			// Set the sound timer to Vx
			c.SoundTimer = c.V[c.Instruction.X]

		case 0x29:
			// Set I to the location of the sprite for the character in Vx
			c.I = uint16(c.V[c.Instruction.X]) * 5

		case 0x33:
			// Store the binary-coded decimal representation of Vx in memory
			c.Ram[c.I] = c.V[c.Instruction.X] / 100
			c.Ram[c.I+1] = (c.V[c.Instruction.X] / 10) % 10
			c.Ram[c.I+2] = c.V[c.Instruction.X] % 10

		case 0x55:
			// Store V0 to Vx in memory starting at address I
			//SCHIP uses I chip8 uses I+1
			for i := 0; i <= int(c.Instruction.X); i++ {
				c.Ram[c.I+uint16(i)] = c.V[i]
			}

		case 0x65:
			// Read V0 to Vx from memory starting at address I
			for i := 0; i <= int(c.Instruction.X); i++ {
				c.V[i] = c.Ram[c.I+uint16(i)]
			}

		default:
			fmt.Println("Unknown instruction: ", c.Instruction.OpCode)

		}

	default:
		fmt.Println("Unknown instruction: ", c.Instruction.OpCode)

	}

}

func (c *Chip8) UpdateTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}
}
