package sdlConf

import (
	"github.com/Sounak10/chip8emu/config"
	"github.com/Sounak10/chip8emu/cpu"
	"github.com/veandco/go-sdl2/sdl"
)

type SDLContext struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func NewSDLContext(config *config.Config) (*SDLContext, error) {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_TIMER)
	if err != nil {
		return nil, err
	}
	window, err := sdl.CreateWindow("CHIP8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(config.WindowWidth*config.ScaleFactor), int32(config.WindowHeight*config.ScaleFactor), sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	return &SDLContext{window: window, renderer: renderer}, nil
}

func (ctx *SDLContext) Destroy() {
	if ctx.renderer != nil {
		ctx.renderer.Destroy()
	}
	if ctx.window != nil {
		ctx.window.Destroy()
	}
	sdl.Quit()
}

func (ctx *SDLContext) ClearScreen(config *config.Config) {
	// Fix color extraction to match SDL's RGBA format
	r := uint8(config.BgColor >> 24 & 0xFF)
	g := uint8(config.BgColor >> 16 & 0xFF)
	b := uint8(config.BgColor >> 8 & 0xFF)
	a := uint8(config.BgColor & 0xFF)
	ctx.renderer.SetDrawColor(r, g, b, a)
	ctx.renderer.Clear()
}

func (ctx *SDLContext) UpdateScreen(chip8 *cpu.Chip8, config *config.Config) {

	// Extract foreground color components from config
	fgR := uint8(config.FgColor >> 24 & 0xFF)
	fgG := uint8(config.FgColor >> 16 & 0xFF)
	fgB := uint8(config.FgColor >> 8 & 0xFF)
	fgA := uint8(config.FgColor & 0xFF)

	// Extract background color components from config
	bgR := uint8(config.BgColor >> 24 & 0xFF)
	bgG := uint8(config.BgColor >> 16 & 0xFF)
	bgB := uint8(config.BgColor >> 8 & 0xFF)
	bgA := uint8(config.BgColor & 0xFF)

	rect := sdl.Rect{X: 0, Y: 0, W: int32(config.ScaleFactor), H: int32(config.ScaleFactor)}

	ctx.renderer.SetDrawColor(bgR, bgG, bgB, bgA)
	ctx.renderer.Clear()

	ctx.renderer.SetDrawColor(fgR, fgG, fgB, fgA)

	// Only draw foreground pixels to improve performance
	for i := range len(chip8.Display) {
		if chip8.Display[i] {
			rect.X = int32(i%int(config.WindowWidth)) * int32(config.ScaleFactor)
			rect.Y = int32(i/int(config.WindowWidth)) * int32(config.ScaleFactor)
			ctx.renderer.FillRect(&rect)
			if config.PixelOutline {
				ctx.renderer.SetDrawColor(bgR, bgG, bgB, bgA)
				ctx.renderer.DrawRect(&rect)
				ctx.renderer.SetDrawColor(fgR, fgG, fgB, fgA)
			}
		}
	}

	ctx.renderer.Present()
}
