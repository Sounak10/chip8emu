package main

import (
	"fmt"

	"github.com/Sounak10/chip8emu/config"
	"github.com/Sounak10/chip8emu/cpu"
	"github.com/Sounak10/chip8emu/sdlConf"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	config := config.NewConfig(64, 32, 0xFFA4FFFF, 0x000000FF, 10, true, 700)
	sdlContext, err := sdlConf.NewSDLContext(config)
	if err != nil {
		panic(err)
	}
	defer sdlContext.Destroy()

	chip8, err := cpu.NewChip8("roms/Tetris.ch8")
	if err != nil {
		panic(err)
	}

	// Initial screen clear
	sdlContext.ClearScreen(config)

	for chip8.State != cpu.QUIT {
		handleInput(chip8)
		if chip8.State == cpu.PAUSED {
			continue
		}

		before := sdl.GetPerformanceCounter()

		for range int(config.InstructionPerSec / 60) {
			chip8.EmulateInstruction(config)
		}
		after := sdl.GetPerformanceCounter()
		elapsed := float64((after-before)*1000) / float64(sdl.GetPerformanceFrequency())
		if 16.67 > elapsed {
			sdl.Delay(uint32(16.67 - elapsed))
		}

		sdlContext.UpdateScreen(chip8, config)
		chip8.UpdateTimers()
	}

}

func handleInput(chip8 *cpu.Chip8) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			chip8.State = cpu.QUIT
		case *sdl.KeyboardEvent:
			keyboardEvent := event.(*sdl.KeyboardEvent)

			if keyboardEvent.Type == sdl.KEYDOWN {
				switch keyboardEvent.Keysym.Sym {
				case sdl.K_ESCAPE:
					chip8.State = cpu.QUIT
				case sdl.K_SPACE:
					if chip8.State == cpu.PAUSED {
						chip8.State = cpu.RUNNING
					} else {
						chip8.State = cpu.PAUSED
						fmt.Println("======= PAUSED =======")
					}
				case sdl.K_1:
					chip8.KeyPad[0] = true
				case sdl.K_2:
					chip8.KeyPad[1] = true
				case sdl.K_3:
					chip8.KeyPad[2] = true
				case sdl.K_4:
					chip8.KeyPad[3] = true
				case sdl.K_q:
					chip8.KeyPad[4] = true
				case sdl.K_w:
					chip8.KeyPad[5] = true
				case sdl.K_e:
					chip8.KeyPad[6] = true
				case sdl.K_r:
					chip8.KeyPad[7] = true
				case sdl.K_a:
					chip8.KeyPad[8] = true
				case sdl.K_s:
					chip8.KeyPad[9] = true
				case sdl.K_d:
					chip8.KeyPad[10] = true
				case sdl.K_f:
					chip8.KeyPad[11] = true
				case sdl.K_z:
					chip8.KeyPad[12] = true
				case sdl.K_x:
					chip8.KeyPad[13] = true
				case sdl.K_c:
					chip8.KeyPad[14] = true
				case sdl.K_v:
					chip8.KeyPad[15] = true
				default:

				}

			} else if keyboardEvent.Type == sdl.KEYUP {
				switch keyboardEvent.Keysym.Sym {
				case sdl.K_1:
					chip8.KeyPad[0] = false
				case sdl.K_2:
					chip8.KeyPad[1] = false
				case sdl.K_3:
					chip8.KeyPad[2] = false
				case sdl.K_4:
					chip8.KeyPad[3] = false
				case sdl.K_q:
					chip8.KeyPad[4] = false
				case sdl.K_w:
					chip8.KeyPad[5] = false
				case sdl.K_e:
					chip8.KeyPad[6] = false
				case sdl.K_r:
					chip8.KeyPad[7] = false
				case sdl.K_a:
					chip8.KeyPad[8] = false
				case sdl.K_s:
					chip8.KeyPad[9] = false
				case sdl.K_d:
					chip8.KeyPad[10] = false
				case sdl.K_f:
					chip8.KeyPad[11] = false
				case sdl.K_z:
					chip8.KeyPad[12] = false
				case sdl.K_x:
					chip8.KeyPad[13] = false
				case sdl.K_c:
					chip8.KeyPad[14] = false
				case sdl.K_v:
					chip8.KeyPad[15] = false
				default:
				}

			}
		}
	}
}
