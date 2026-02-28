package utils

import (
	"fmt"
	"strings"
)

// PrintTableWithBorders prints a table with borders
func PrintTableWithBorders(headers []string, rows [][]string, colWidths []int) {
	// Calculate total width
	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w + 3 // +3 for " | "
	}
	totalWidth += 1 // for final "|"

	// Top border
	fmt.Println("┌" + strings.Repeat("─", totalWidth-2) + "┐")

	// Headers
	fmt.Print("│ ")
	for i, header := range headers {
		fmt.Printf("%-*s", colWidths[i], header)
		if i < len(headers)-1 {
			fmt.Print(" │ ")
		}
	}
	fmt.Println(" │")

	// Header separator
	fmt.Print("├")
	for i := range headers {
		fmt.Print(strings.Repeat("─", colWidths[i]+2))
		if i < len(headers)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")

	// Rows
	for _, row := range rows {
		fmt.Print("│ ")
		for i, cell := range row {
			// Handle colored text - don't count ANSI codes in width
			displayWidth := colWidths[i]
			// If cell contains ANSI codes, adjust padding
			if strings.Contains(cell, "\033[") {
				// Count visible characters (excluding ANSI codes)
				visibleLen := len(stripANSI(cell))
				padding := colWidths[i] - visibleLen
				fmt.Print(cell)
				if padding > 0 {
					fmt.Print(strings.Repeat(" ", padding))
				}
			} else {
				fmt.Printf("%-*s", displayWidth, cell)
			}

			if i < len(row)-1 {
				fmt.Print(" │ ")
			}
		}
		fmt.Println(" │")
	}

	// Bottom border
	fmt.Println("└" + strings.Repeat("─", totalWidth-2) + "┘")
}

// stripANSI removes ANSI color codes from a string for length calculation
func stripANSI(s string) string {
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			result += string(r)
		}
	}
	return result
}
