package components

import (
	"os"
	"runtime"

	"github.com/charmbracelet/lipgloss"
)

type TermCapability int

const (
	TermBasic TermCapability = iota
	TermUnicode
	TermNerdFont
)

var termCap TermCapability

type CardIconSet struct {
	CPU    string
	Mem    string
	Disk   string
	Header string
}

var cardIcons CardIconSet

func CardIcons() CardIconSet {
	return cardIcons
}

func detectTerminal() TermCapability {
	if runtime.GOOS != "windows" {
		return TermNerdFont
	}

	if os.Getenv("WT_SESSION") != "" {
		return TermNerdFont
	}

	if os.Getenv("TERM_PROGRAM") == "vscode" {
		return TermNerdFont
	}

	if os.Getenv("ConEmuPID") != "" || os.Getenv("ALACRITTY_SOCKET") != "" {
		return TermNerdFont
	}

	if os.Getenv("TERM") != "" && os.Getenv("TERM") != "dumb" {
		return TermUnicode
	}

	return TermBasic
}

func init() {
	termCap = detectTerminal()

	switch termCap {
	case TermNerdFont:
		cardIcons.CPU = "\uf0ed "    // nf-fa-microchip
		cardIcons.Mem = "\uf09d "    // nf-oct-memory
		cardIcons.Disk = "\uf084 "   // nf-fa-hdd_o
		cardIcons.Header = "\u25C6 " // diamond

	case TermUnicode:
		cardIcons.CPU = "⬡ "
		cardIcons.Mem = "▣ "
		cardIcons.Disk = "◈ "
		cardIcons.Header = "◆ "

	default:
		cardIcons.CPU = "[C] "
		cardIcons.Mem = "[M] "
		cardIcons.Disk = "[D] "
		cardIcons.Header = "* "
	}

	SparklineRunes = pickSparklineRunes()
}

func pickSparklineRunes() []rune {
	switch termCap {
	case TermNerdFont, TermUnicode:
		return []rune("▁▂▃▄▅▆▇█")
	default:
		return []rune("._-~=+#")
	}
}

func ActiveBorder() lipgloss.Border {
	switch termCap {
	case TermNerdFont, TermUnicode:
		return lipgloss.RoundedBorder()
	default:
		return lipgloss.NormalBorder()
	}
}
