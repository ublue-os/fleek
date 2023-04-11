package fin

import (
	"github.com/pterm/pterm"
)

func Progress() *pterm.ProgressbarPrinter {
	progressBar := pterm.DefaultProgressbar
	return &progressBar
}
