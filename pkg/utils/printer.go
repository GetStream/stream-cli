package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
)

func PrintObject(cmd *cobra.Command, object interface{}) error {
	format, err := cmd.Flags().GetString("output-format")
	if err != nil {
		return err
	}

	switch format {
	case "json":
		return printJSONObject(cmd, object)
	case "tree":
		return printUIObject(cmd, object)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

func printJSONObject(cmd *cobra.Command, object interface{}) error {
	b, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return err
	}

	cmd.Println(string(b))

	return nil
}

func printUIObject(cmd *cobra.Command, object interface{}) error {
	var asMap map[string]interface{}
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &asMap)
	if err != nil {
		return err
	}

	_, ok := asMap["ratelimit"]
	if ok {
		delete(asMap, "ratelimit")
	}

	rootNode := &widgets.TreeNode{Nodes: []*widgets.TreeNode{}}
	for k, v := range asMap {
		addNodesRecursive(rootNode, k+": ", v)
	}

	if err := ui.Init(); err != nil {
		return err
	}
	defer ui.Close()

	renderUI(rootNode)
	return nil
}

type nodeValue string

func (nv nodeValue) String() string {
	return string(nv)
}

func addNodesRecursive(parentNode *widgets.TreeNode, prefix string, value interface{}) {
	node := &widgets.TreeNode{Value: nodeValue(prefix), Nodes: []*widgets.TreeNode{}}

	switch v := value.(type) {
	case map[string]interface{}:
		for k, val := range v {
			addNodesRecursive(node, k+": ", val)
		}
	case []interface{}:
		for i, v := range v {
			addNodesRecursive(node, fmt.Sprintf("[%v] ", i), v)
		}
	case bool:
		node = &widgets.TreeNode{Value: nodeValue(prefix + strconv.FormatBool(v))}
	case string:
		if v == "" {
			v = "\"\""
		}

		node = &widgets.TreeNode{Value: nodeValue(prefix + v)}
	case float64:
		node = &widgets.TreeNode{Value: nodeValue(prefix + strconv.FormatFloat(v, 'f', -1, 64))}
	default:
		node = &widgets.TreeNode{Value: nodeValue(prefix + "null")}
	}

	parentNode.Nodes = append(parentNode.Nodes, node)
}

func renderUI(rootNode *widgets.TreeNode) {
	// First section: instructions
	p := widgets.NewParagraph()
	p.Text = `- Press [q] or [Ctrl+C] to quit
	- Use [Arrow-Keys] to navigate
	- Press [Enter] to expand/collapse
	`
	x, y := ui.TerminalDimensions()
	p.SetRect(0, 0, x, 5)
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	// Second section: the tree
	l := widgets.NewTree()
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.SetNodes(rootNode.Nodes)
	l.CollapseAll()
	l.SetRect(0, 5, x, y)

	// Printing both sections
	ui.Render(p, l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Down>":
			l.ScrollDown()
		case "<Up>":
			l.ScrollUp()
		case "<Enter>":
			l.ToggleExpand()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 0, x, y)
		}

		ui.Render(l)
	}
}
