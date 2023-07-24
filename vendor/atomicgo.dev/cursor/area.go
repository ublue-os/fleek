package cursor

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// Area displays content which can be updated on the fly.
// You can use this to create live output, charts, dropdowns, etc.
type Area struct {
	height int
	writer Writer
}

// NewArea returns a new Area.
func NewArea() Area {
	return Area{
		writer: os.Stdout,
		height: 0,
	}
}

// WithWriter sets a custom writer for the Area.
func (area Area) WithWriter(writer Writer) Area {
	area.writer = writer

	return area
}

// Clear clears the content of the Area.
func (area *Area) Clear() {
	Bottom()

	if area.height > 0 {
		ClearLinesUp(area.height)
	}
}

// Update overwrites the content of the Area.
func (area *Area) Update(content string) {
	oldWriter := target

	SetTarget(area.writer) // Temporary set the target to the Area's writer so we can use the cursor functions
	area.Clear()

	lines := strings.Split(content, "\n")
	fmt.Fprintln(area.writer, strings.Repeat("\n", len(lines)-1)) // This appends space if the terminal is at the bottom
	Up(len(lines))
	SetTarget(oldWriter) // Reset the target to the old writer

	// Workaround for buggy behavior on Windows
	if runtime.GOOS == "windows" {
		for _, line := range lines {
			fmt.Fprint(area.writer, line)
			StartOfLineDown(1)
		}
	} else {
		for _, line := range lines {
			fmt.Fprintln(area.writer, line)
		}
	}

	height = 0
	area.height = len(strings.Split(content, "\n"))
}
