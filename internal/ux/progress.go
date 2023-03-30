package ux

import (
	"github.com/pterm/pterm"
)

func Progress() *pterm.ProgressbarPrinter {
	progressBar := pterm.DefaultProgressbar
	return &progressBar
}

var (
	// Info returns a PrefixPrinter, which can be used to print text with an "info" Prefix.
	Info = pterm.PrefixPrinter{
		MessageStyle: infoMessageStyle,
		Prefix:       infoPrefix,
	}

	// Warning returns a PrefixPrinter, which can be used to print text with a "warning" Prefix.
	Warning = pterm.PrefixPrinter{
		MessageStyle: warningMessageStyle,
		Prefix:       warningPrefix,
	}

	// Success returns a PrefixPrinter, which can be used to print text with a "success" Prefix.
	Success = pterm.PrefixPrinter{
		MessageStyle: successMessageStyle,
		Prefix:       successPrefix,
	}

	// Error returns a PrefixPrinter, which can be used to print text with an "error" Prefix.
	Error = pterm.PrefixPrinter{
		MessageStyle: errorMessageStyle,
		Prefix:       errorPrefix,
	}

	// Fatal returns a PrefixPrinter, which can be used to print text with an "fatal" Prefix.
	// NOTICE: Fatal terminates the application immediately!
	Fatal = pterm.PrefixPrinter{
		MessageStyle: fatalMessageStyle,
		Prefix:       fatalPrefix,
		Fatal:        true,
	}

	// Debug Prints debug messages. By default it will only print if PrintDebugMessages is true.
	// You can change PrintDebugMessages with EnableDebugMessages and DisableDebugMessages, or by setting the variable itself.
	Debug = pterm.PrefixPrinter{
		MessageStyle: debugMessageStyle,
		Prefix:       debugPrefix,
		Debugger:     true,
	}

	// Description returns a PrefixPrinter, which can be used to print text with a "description" Prefix.
	Description = pterm.PrefixPrinter{
		MessageStyle: descriptionMessageStyle,
		Prefix:       descriptionPrefix,
	}
)
