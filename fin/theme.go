package fin

import "github.com/pterm/pterm"

// HelpSectionPrinter is used to print the different sections in the CLIs help output.
// You can overwrite it with any Sprint function.
var HelpSectionPrinter = pterm.DefaultSection.WithLevel(2).Sprint
var TitleSectionPrinter = pterm.DefaultSection.WithLevel(1).Sprint
var DescriptionSectionPrinter = pterm.DefaultSection.WithLevel(2).Sprint
var DetailSectionPrinter = pterm.DefaultSection.WithLevel(3).Sprint
var Logger = pterm.DefaultLogger.WithTime(false).WithLevel(pterm.LogLevelWarn)

var ParagraphPrinter = pterm.DefaultParagraph.Sprint

// SetTrace sets the log level to Trace
// for trace/debugging level output
func SetTrace() {
	Logger = pterm.DefaultLogger.WithTime(false).WithLevel(pterm.LogLevelTrace)
}

// SetDebug sets the log level to Debug
// for trace/debugging level output
func SetDebug() {
	Logger = pterm.DefaultLogger.WithTime(false).WithLevel(pterm.LogLevelDebug)
}

// SetDebug sets the log level to Debug (unfortunately)
// for verbose level output
func SetVerbose() {
	Logger = pterm.DefaultLogger.WithTime(false).WithLevel(pterm.LogLevelInfo)
}

var (
	// info
	infoMessageStyle = &pterm.Style{pterm.FgLightBlue}
	infoPrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgBlue},
		Text:  "[i]",
	}
	// verbose
	verboseMessageStyle = &pterm.Style{pterm.FgLightGreen}
	verbosePrefix       = pterm.Prefix{
		Style: &pterm.Style{pterm.FgGreen},
		Text:  "[v]",
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
	debugMessageStyle = &pterm.Style{pterm.FgDefault}
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

	Paragraph = pterm.ParagraphPrinter{
		MaxWidth: 70,
	}
	// Info returns a PrefixPrinter, which can be used to print text with an "info" Prefix.
	Info = pterm.PrefixPrinter{
		MessageStyle: infoMessageStyle,
		Prefix:       infoPrefix,
	}
	// Info returns a PrefixPrinter, which can be used to print text with an "info" Prefix.
	Verbose = pterm.PrefixPrinter{
		MessageStyle: verboseMessageStyle,
		Prefix:       verbosePrefix,
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
