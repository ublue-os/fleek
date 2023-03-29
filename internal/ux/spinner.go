package ux

import (
	"github.com/pterm/pterm"
)

func Spinner() *pterm.SpinnerPrinter {
	fleekSpinner := pterm.DefaultSpinner
	// Replace the InfoPrinter with a custom "NOCHG" one
	fleekSpinner.InfoPrinter = &pterm.PrefixPrinter{
		MessageStyle: infoMessageStyle,
		Prefix:       infoPrefix,
	}
	fleekSpinner.WarningPrinter = &pterm.PrefixPrinter{
		MessageStyle: warningMessageStyle,
		Prefix:       warningPrefix,
	}
	fleekSpinner.FailPrinter = &pterm.PrefixPrinter{
		MessageStyle: failMessageStyle,
		Prefix:       failPrefix,
	}
	fleekSpinner.SuccessPrinter = &pterm.PrefixPrinter{
		MessageStyle: successMessageStyle,
		Prefix:       successPrefix,
	}

	return &fleekSpinner
}
