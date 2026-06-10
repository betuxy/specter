package specter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	colorPurple  = "#bd93f9"
	colorCyan    = "#8be9fd"
	colorYellow  = "#f1fa8c"
	colorOrange  = "#ffb86c"
	colorPink    = "#ff79c6"
	colorComment = "#6272a4"
	colorFg      = "#f8f8f2"
)

var (
	tcellBg        = tcell.NewRGBColor(40, 42, 54)
	tcellBar       = tcell.NewRGBColor(68, 71, 90)
	tcellConnector = tcell.NewRGBColor(98, 114, 164)
)

func RunTUI(data []map[string]interface{}, expanded bool) error {
	app := tview.NewApplication()

	root := tview.NewTreeNode("root")

	if len(data) == 1 {
		addMapToTree(root, data[0], expanded)
	} else {
		for i, obj := range data {
			label := fmt.Sprintf("[%s]%s[-]", colorPurple, tview.Escape(fmt.Sprintf("[%d]", i)))
			node := tview.NewTreeNode(label).
				SetSelectable(true).
				SetExpanded(expanded)
			addMapToTree(node, obj, expanded)
			root.AddChild(node)
		}
	}

	tree := tview.NewTreeView()
	tree.SetRoot(root)
	tree.SetTopLevel(1)
	tree.SetGraphicsColor(tcellConnector)
	tree.SetBackgroundColor(tcellBg)

	if children := root.GetChildren(); len(children) > 0 {
		tree.SetCurrentNode(children[0])
	}

	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(
			"  [%s]j/k ↑/↓[-]  navigate   [%s]space/enter[-]  expand/collapse   [%s]q[-]  quit",
			colorYellow, colorYellow, colorYellow,
		))
	statusBar.SetBackgroundColor(tcellBar)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tree, 0, 1, true).
		AddItem(statusBar, 1, 0, false)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'q':
			app.Stop()
			return nil
		case ' ':
			node := tree.GetCurrentNode()
			if node != nil && len(node.GetChildren()) > 0 {
				node.SetExpanded(!node.IsExpanded())
			}
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			node := tree.GetCurrentNode()
			if node != nil && len(node.GetChildren()) > 0 {
				node.SetExpanded(!node.IsExpanded())
			}
			return nil
		}
		return event
	})

	return app.SetRoot(layout, true).Run()
}

func addMapToTree(parent *tview.TreeNode, m map[string]interface{}, expanded bool) {
	for _, k := range sortedKeys(m) {
		addValueToTree(parent, k, m[k], expanded)
	}
}

func addValueToTree(parent *tview.TreeNode, key string, value interface{}, expanded bool) {
	escapedKey := tview.Escape(key)
	switch v := value.(type) {
	case map[string]interface{}:
		label := fmt.Sprintf("[%s]%s[-]", colorPurple, escapedKey)
		node := tview.NewTreeNode(label).
			SetSelectable(true).
			SetExpanded(expanded)
		addMapToTree(node, v, expanded)
		parent.AddChild(node)
	case []interface{}:
		label := fmt.Sprintf("[%s]%s[-] [%s](%d)[-]", colorPurple, escapedKey, colorComment, len(v))
		node := tview.NewTreeNode(label).
			SetSelectable(true).
			SetExpanded(expanded)
		for i, item := range v {
			escapedIdx := tview.Escape(fmt.Sprintf("[%d]", i))
			if obj, ok := item.(map[string]interface{}); ok {
				itemLabel := fmt.Sprintf("[%s]%s[-] %s", colorCyan, escapedIdx, compactColorized(obj, 1))
				itemNode := tview.NewTreeNode(itemLabel).
					SetSelectable(true).
					SetExpanded(expanded)
				addMapToTree(itemNode, obj, expanded)
				node.AddChild(itemNode)
			} else {
				addValueToTree(node, fmt.Sprintf("[%d]", i), item, expanded)
			}
		}
		parent.AddChild(node)
	default:
		label := fmt.Sprintf("[%s]%s[-]: %s", colorPurple, escapedKey, formatValue(value))
		node := tview.NewTreeNode(label).SetSelectable(false)
		parent.AddChild(node)
	}
}

func compactColorized(m map[string]interface{}, depth int) string {
	keys := sortedKeys(m)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("[%s]%s[-]: %s", colorCyan, tview.Escape(k), compactValue(m[k], depth-1)))
	}
	sep := fmt.Sprintf("[%s], [-]", colorComment)
	return fmt.Sprintf("[%s]{[-]%s[%s]}[-]", colorComment, strings.Join(parts, sep), colorComment)
}

func compactValue(value interface{}, depth int) string {
	switch v := value.(type) {
	case map[string]interface{}:
		if depth <= 0 {
			return fmt.Sprintf("[%s]{...}[-]", colorComment)
		}
		return compactColorized(v, depth)
	case []interface{}:
		if depth <= 0 || len(v) == 0 {
			return fmt.Sprintf("[%s][%d][-]", colorComment, len(v))
		}
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, compactValue(item, depth-1))
		}
		sep := fmt.Sprintf("[%s], [-]", colorComment)
		return fmt.Sprintf("[%s][ [-]%s[%s] ][-]", colorComment, strings.Join(parts, sep), colorComment)
	default:
		return formatValue(value)
	}
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("[%s]\"%s\"[-]", colorYellow, tview.Escape(v))
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("[%s]%d[-]", colorOrange, int64(v))
		}
		return fmt.Sprintf("[%s]%g[-]", colorOrange, v)
	case bool:
		return fmt.Sprintf("[%s]%v[-]", colorPink, v)
	case nil:
		return fmt.Sprintf("[%s]null[-]", colorComment)
	default:
		return fmt.Sprintf("[%s]%v[-]", colorFg, tview.Escape(fmt.Sprintf("%v", v)))
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		iName := strings.Contains(strings.ToLower(keys[i]), "name")
		jName := strings.Contains(strings.ToLower(keys[j]), "name")
		if iName != jName {
			return iName
		}
		return keys[i] < keys[j]
	})
	return keys
}
