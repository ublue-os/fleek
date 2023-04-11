package fin

import (
	"github.com/pterm/pterm"
)

func Table() *pterm.TablePrinter {
	table := pterm.DefaultTable
	return &table
}
