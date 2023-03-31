package ux

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	width       = 96
	columnWidth = 30
)

var (

	// General.

	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	bold   = lipgloss.NewStyle().Bold(true)
	italic = lipgloss.NewStyle().Italic(true)
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8CC8C"))
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#DBAB79"))

	// Tabs.

	// Title.

	// Dialog.

	// List.

	list = lipgloss.NewStyle().
		MarginRight(2).
		Height(8).
		Width(columnWidth + 1)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	bullet = lipgloss.NewStyle().SetString("-").
		Foreground(special).
		PaddingRight(1).
		PaddingLeft(1).
		String()
	listBullet = func(s string) string {
		return bullet + lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	// Paragraphs/History.

	// Status Bar.

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))

	// Page.

	plainStyle = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	listStyle  = lipgloss.NewStyle().Padding(1, 0, 0, 0)
)

func OutputLipgloss() {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	w := lipgloss.Width

	statusKey := statusStyle.Render("STATUS")
	encoding := encodingStyle.Render("UTF-8")
	fishCake := fishCakeStyle.Render("🍥 Fish Cake")
	statusVal := statusText.Copy().
		Width(width - w(statusKey) - w(encoding) - w(fishCake)).
		Render("Ravishing")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	)

	doc.WriteString(statusBarStyle.Width(width).Render(bar))
	if physicalWidth > 0 {
		plainStyle = plainStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(plainStyle.Render(doc.String()))
}

func BulletListLipGloss(title string, items []string) {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	formattedList := []string{}
	formattedList = append(formattedList, listHeader(title))
	for _, item := range items {
		formattedList = append(formattedList, listBullet(item))
	}
	output := list.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			formattedList...,
		),
	)
	doc.WriteString(output)
	if physicalWidth > 0 {
		listStyle = listStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(listStyle.Render(doc.String()))
}

func ThreeColumnList(col1Title string,
	col1Items []string,
	col2Title string,
	col2Items []string,
	col3Title string,
	col3Items []string) {

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	col1List := []string{}
	col1List = append(col1List, listHeader(col1Title))
	for _, item := range col1Items {
		col1List = append(col1List, listBullet(item))
	}
	column1 := list.Copy().Width(columnWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			col1List...,
		),
	)
	col2List := []string{}
	col2List = append(col2List, listHeader(col2Title))
	for _, item := range col2Items {
		col2List = append(col2List, listBullet(item))
	}
	column2 := list.Copy().Width(columnWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			col2List...,
		),
	)
	col3List := []string{}
	col3List = append(col3List, listHeader(col3Title))
	for _, item := range col3Items {
		col3List = append(col3List, listBullet(item))
	}
	column3 := list.Copy().Width(columnWidth).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			col3List...,
		),
	)
	lists := lipgloss.JoinHorizontal(
		lipgloss.Top,
		column1,
		column2,
		column3,
	)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))

	if physicalWidth > 0 {
		listStyle = listStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(listStyle.Render(doc.String()))
}
func InfoLipGloss(title string, value string) {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	title = strings.Split(title, " ")[0]

	doc.WriteString(bold.Render(title) + bullet + italic.Render(value))
	if physicalWidth > 0 {
		plainStyle = plainStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(plainStyle.Render(doc.String()))
}
func WarningLipGloss(title string, value string) {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	title = strings.Split(title, " ")[0]

	doc.WriteString(bold.Render(title) + bullet + yellow.Render(value))
	if physicalWidth > 0 {
		plainStyle = plainStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(plainStyle.Render(doc.String()))
}
func ActionLipGloss(title string, value string) {
	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	doc := strings.Builder{}
	title = strings.Split(title, " ")[0]
	doc.WriteString(bold.Copy().Foreground(highlight).Render(title) + bullet + green.Render(value))
	if physicalWidth > 0 {
		plainStyle = plainStyle.MaxWidth(physicalWidth)
	}

	// Okay, let's print it
	fmt.Println(plainStyle.Copy().Padding(1, 0, 1).Render(doc.String()))
}
