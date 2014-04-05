package main 

import (
	"ngaro"
	"flag"
	"fmt"
  "log"
	"io"
	"os"
	"os/exec"
	"github.com/kless/term/readline"
  "github.com/kless/term"
)

var Usage = func() {
	fmt.Fprint(os.Stderr, `
Gonga usage:
	gonga [options] [image file]

Gonga is the Go version of the Ngaro virtual machine.

If no image file is specified in the command line
gongaImage will be loaded, retroImage if that fails.

Options:
`)
	flag.PrintDefaults()
}

var size = flag.Int("s", 50000, "image size")
var dump = flag.String("d", "gongaImage", "image dump file")
var shrink = flag.Bool("shrink", false, "shrink image dump file")
var tty = flag.Bool("t", true, "input / output is tty")

type withFiles []*os.File

func (wf *withFiles) String() string {
	return fmt.Sprint(*wf)
}

func (wf *withFiles) Set(value string) bool {
	if f, err := os.OpenFile(value, os.O_RDONLY, 0666); err == nil {
		nwf := append(*wf, f)
		wf = &nwf
		return true
	}
	return false
}

// Default terminal
func noClear(w io.Writer)                          {}
func vt100Dimensions() (width int32, height int32) { return 80, 24 }

// Tty terminal
func ttyClear(w io.Writer) { fmt.Fprint(w, "\033[2J\033[1;1H") }
func ttyDimensions() (width int32, height int32) {
	if err := exec.Command("/bin/stty", "-F", "/dev/tty", "size").Run(); err == nil {
//		fmt.Fscan(p.Stdout, &width, &height)
		return
	}
	return vt100Dimensions()
}

// ConsoleInput is an adapter for Readline and Ngaro's byte-oriented input.
type ConsoleInput struct{
  line    *readline.Line
  pending []byte
}

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

// Read bridges the semantic gap between a raw byte-oriented reader and the line-oriented approach taken by ReadLine.
func (ci *ConsoleInput) Read(bs []byte) (int, error) {
  if len(ci.pending) == 0 {
    s, err := ci.line.Read()
    if err != nil {
      return 0, io.EOF  // Hokey, but good enough for bootstrapping purposes.
    }
    ci.pending = []byte(s)
    ci.pending = append(ci.pending, byte(13), byte(10))
  }
  n := min(len(bs), len(ci.pending))
  copy(bs[0:n], ci.pending[0:n])
  ci.pending = ci.pending[n:]
  return n, nil
}

func main() {
	var wf withFiles
	flag.Parse()

	var img []int32
	var err error

  h, err := readline.NewHistory("gonga.history")
  if err != nil {
    log.Fatal(err)
  }
  t, err := term.New()
  if err != nil {
    log.Fatal(err)
  }
  defer t.Restore()
  l, err := readline.NewLine(t, "", "", 0, h)
  if err != nil {
    log.Fatal(err)
  }
  ci := &ConsoleInput{
    line:    l,
    pending: make([]byte, 0),
  }

	switch flag.NArg() {
	case 0:
		img, err = ngaro.Load("gongaImage", *size)
		if err != nil {
			img, err = ngaro.Load("retroImage", *size)
		}
	case 1:
		img, err = ngaro.Load(flag.Arg(0), *size)
	default:
		fmt.Fprintln(os.Stderr, "too many arguments")
		flag.Usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error starting gonga: ", err)
		os.Exit(1)
	}

	// Reverse wf and add os.Stdin
	rs := make([]io.Reader, 0, len(wf)+1)
	for i, _ := range wf {
		rs = append(rs, wf[len(wf)-1-i])
	}
	input := io.MultiReader(append(rs, ci)...)
	clr := noClear
	dim := vt100Dimensions
	if *tty {
    // Gonga's original implementation set the console into raw mode with a /bin/stty call.
    // For interactive Orc emulation, this works against us.  We desperately want readline
    // capability.
		clr = ttyClear
		dim = ttyDimensions
	}

	// Run a new VM
	term := ngaro.NewTerm(clr, dim, input, os.Stdout)
	vm := ngaro.New(img, *dump, *shrink, term)
	err = vm.Run()
	if err != nil {
		os.Exit(1)
	}
}
