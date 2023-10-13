package report

import (
	"github.com/pterm/pterm"
)

type Style struct {
	Table struct {
		Bg pterm.Color
		Fg pterm.Color
	}
}

var defaultStyle = Style{
	Table: struct {
		Bg pterm.Color
		Fg pterm.Color
	}{pterm.BgCyan, pterm.FgCyan},
}

func RenderTable(header []string, data [][]string) error {
	var rows [][]string
	rows = append(rows, header)
	rows = append(rows, data...)
	pterm.DefaultBox.VerticalString = "│"
	pterm.DefaultTable.HeaderStyle = &pterm.Style{defaultStyle.Table.Fg}
	return pterm.DefaultTable.WithHasHeader().WithBoxed(true).WithData(rows).WithSeparator("  ").Render()
}
