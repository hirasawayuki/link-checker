package linkchecker

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type IOStream struct {
	In        io.ReadCloser
	Out       io.Writer
	ErrOut    io.Writer
	indicator *spinner.Spinner
}

func New() *IOStream {
	return &IOStream{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

func (s *IOStream) StartIndicator() {
	dotStyle := spinner.CharSets[11]
	sp := spinner.New(dotStyle, 120*time.Millisecond, spinner.WithWriter(s.ErrOut), spinner.WithColor("fgCyan"))

	sp.Start()
	s.indicator = sp
}

func (s *IOStream) StopIndicator() {
	if s.indicator == nil {
		return
	}
	s.indicator.Stop()
	s.indicator = nil
}

func (s *IOStream) Font() *Font {
	color := color.New()
	return &Font{color: color}
}

type Font struct {
	color *color.Color
}

func (f *Font) Bold(t string) string {
	bold := color.New(color.Bold)
	return bold.Sprint(t)
}

func (f *Font) Green(t string) string {
	return color.GreenString(t)
}

func (f *Font) Greenf(t string, args ...interface{}) string {
	return f.Green(fmt.Sprintf(t, args...))
}

func (f *Font) Red(t string) string {
	return color.RedString(t)
}

func (f *Font) Redf(t string, args ...interface{}) string {
	return f.Red(fmt.Sprintf(t, args...))
}

func (f *Font) Yellow(t string) string {
	return color.YellowString(t)
}

func (f *Font) Yellowf(t string, args ...interface{}) string {
	return f.Yellow(fmt.Sprintf(t, args...))
}
