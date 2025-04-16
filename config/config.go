package config

type Config struct {
	WindowWidth       uint32
	WindowHeight      uint32
	FgColor           uint32
	BgColor           uint32
	ScaleFactor       uint32
	PixelOutline      bool
	InstructionPerSec uint32
}

func NewConfig(w uint32, h uint32, fgColor uint32, bgColor uint32, scaleFactor uint32, pixelOutline bool, instructionPerSec uint32) *Config {
	return &Config{
		WindowWidth:       w,
		WindowHeight:      h,
		FgColor:           fgColor,
		BgColor:           bgColor,
		ScaleFactor:       scaleFactor,
		PixelOutline:      pixelOutline,
		InstructionPerSec: instructionPerSec,
	}

}
