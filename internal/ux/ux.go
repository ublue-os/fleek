package ux

import (
	"github.com/pterm/pterm"
)

var (
	// info
	infoMessageStyle = &pterm.Style{pterm.FgLightBlue}
	infoPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgBlue},
		Text:  "[i]",
	}
	// warning
	warningMessageStyle = &pterm.Style{pterm.FgLightYellow}
	warningPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgYellow},
		Text:  "[!]",
	}
	// fail
	failMessageStyle = &pterm.Style{pterm.FgLightRed}
	failPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgRed},
		Text:  "[!]",
	}
	// error
	errorMessageStyle = &pterm.Style{pterm.FgLightRed}
	errorPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgRed},
		Text:  "[!]",
	}
	// fatal
	fatalMessageStyle = &pterm.Style{pterm.FgLightMagenta}
	fatalPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgMagenta},
		Text:  "[x]",
	}
	// success
	successMessageStyle = &pterm.Style{pterm.FgLightGreen}
	successPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgGreen},
		Text:  "[✓]",
	}
	// debug
	debugMessageStyle = &pterm.Style{pterm.FgLightGreen}
	debugPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgGreen},
		Text:  "[d]",
	}
	// description
	descriptionMessageStyle = &pterm.Style{pterm.FgDefault}
	descriptionPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgDefault},
		Text:  "[*]",
	}
)
