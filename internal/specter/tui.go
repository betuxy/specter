package specter

import (
	"fmt"
	"regexp"
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
	selectedStyle  = tcell.StyleDefault.
			Background(tcell.NewRGBColor(68, 71, 90)).
			Foreground(tcell.NewRGBColor(248, 248, 242))
	tagRe = regexp.MustCompile(`\[[^\[\]]*\]`)
)

func RunTUI(data []map[string]interface{}, expanded bool, cfg *Config) error {
	kb := cfg.Keybindings
	app := tview.NewApplication()

	root := tview.NewTreeNode("root")

	if len(data) == 1 {
		addMapToTree(root, data[0], expanded)
	} else {
		for i, obj := range data {
			label := fmt.Sprintf("[%s]%s[-]", colorPurple, tview.Escape(fmt.Sprintf("[%d]", i)))
			node := tview.NewTreeNode(label).
				SetSelectable(true).
				SetExpanded(expanded).
				SetSelectedTextStyle(selectedStyle).
				SetReference(stripColorTags(label))
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

	// Search state
	var (
		matches  []*tview.TreeNode
		matchIdx int
	)

	statusBar := tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf(
			"  [%s]%s/%s ↑/↓[-]  navigate   [%s]%s/enter[-]  expand/collapse   [%s]%s[-]  expand/collapse all   [%s]/[-]  search   [%s]%s[-]  quit",
			colorYellow, keyName(kb.Up), keyName(kb.Down),
			colorYellow, keyName(kb.ExpandCollapse),
			colorYellow, keyName(kb.ExpandCollapseAll),
			colorYellow,
			colorYellow, keyName(kb.Quit),
		))
	statusBar.SetBackgroundColor(tcellBar)

	searchInput := tview.NewInputField().
		SetLabel(" / ").
		SetLabelColor(tcell.NewRGBColor(241, 250, 140)).
		SetFieldBackgroundColor(tcellBg).
		SetFieldTextColor(tcell.NewRGBColor(248, 248, 242))
	searchInput.SetBackgroundColor(tcellBar)

	bottomPages := tview.NewPages()
	bottomPages.AddPage("bar", statusBar, true, true)
	bottomPages.AddPage("search", searchInput, true, false)

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tree, 0, 1, true).
		AddItem(bottomPages, 1, 0, false)

	jumpToMatch := func(idx int) {
		if len(matches) == 0 {
			return
		}
		for _, child := range root.GetChildren() {
			setExpandedRecursive(child, false)
		}
		node := matches[idx]
		expandToNode(root, node)
		tree.SetCurrentNode(node)
		searchInput.SetLabel(fmt.Sprintf(" [%d/%d] / ", idx+1, len(matches)))
	}

	searchInput.SetChangedFunc(func(text string) {
		if text == "" {
			matches = nil
			matchIdx = 0
			searchInput.SetLabel(" / ")
			return
		}
		matches = findMatches(root, strings.ToLower(text))
		matchIdx = 0
		if len(matches) == 0 {
			searchInput.SetLabel(" [no match] / ")
			return
		}
		jumpToMatch(0)
	})

	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			searchInput.SetText("")
			matches = nil
			matchIdx = 0
			searchInput.SetLabel(" / ")
		}
		bottomPages.SwitchToPage("bar")
		app.SetFocus(tree)
	})

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case kb.Down:
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case kb.Up:
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case kb.Quit:
			app.Stop()
			return nil
		case kb.ExpandCollapse:
			node := tree.GetCurrentNode()
			if node != nil && len(node.GetChildren()) > 0 {
				node.SetExpanded(!node.IsExpanded())
			}
			return nil
		case kb.ExpandCollapseAll:
			node := tree.GetCurrentNode()
			if node != nil && len(node.GetChildren()) > 0 {
				setExpandedRecursive(node, !node.IsExpanded())
			}
			return nil
		case '/':
			bottomPages.SwitchToPage("search")
			app.SetFocus(searchInput)
			return nil
		case 'n':
			if len(matches) > 0 {
				matchIdx = (matchIdx + 1) % len(matches)
				jumpToMatch(matchIdx)
			}
			return nil
		case 'N':
			if len(matches) > 0 {
				matchIdx = (matchIdx - 1 + len(matches)) % len(matches)
				jumpToMatch(matchIdx)
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

func findMatches(root *tview.TreeNode, query string) []*tview.TreeNode {
	var matches []*tview.TreeNode
	collectMatches(root, query, &matches)
	return matches
}

func collectMatches(node *tview.TreeNode, query string, matches *[]*tview.TreeNode) {
	if ref, ok := node.GetReference().(string); ok {
		if strings.Contains(strings.ToLower(ref), query) {
			*matches = append(*matches, node)
		}
	}
	for _, child := range node.GetChildren() {
		collectMatches(child, query, matches)
	}
}

// expandToNode expands all ancestors of target so it becomes visible.
func expandToNode(node, target *tview.TreeNode) bool {
	if node == target {
		return true
	}
	for _, child := range node.GetChildren() {
		if expandToNode(child, target) {
			node.SetExpanded(true)
			return true
		}
	}
	return false
}

func stripColorTags(s string) string {
	s = strings.ReplaceAll(s, "[[]", "\x00")
	s = tagRe.ReplaceAllString(s, "")
	return strings.ReplaceAll(s, "\x00", "[")
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
		label := fmt.Sprintf("[%s]%s[-] [%s](%d)[-]", colorPurple, escapedKey, colorComment, len(v))
		node := tview.NewTreeNode(label).
			SetSelectable(true).
			SetExpanded(expanded).
			SetSelectedTextStyle(selectedStyle).
			SetReference(stripColorTags(label))
		addMapToTree(node, v, expanded)
		parent.AddChild(node)
	case []interface{}:
		label := fmt.Sprintf("[%s]%s[-] [%s](%d)[-]", colorPurple, escapedKey, colorComment, len(v))
		node := tview.NewTreeNode(label).
			SetSelectable(true).
			SetExpanded(expanded).
			SetSelectedTextStyle(selectedStyle).
			SetReference(stripColorTags(label))
		for i, item := range v {
			escapedIdx := tview.Escape(fmt.Sprintf("[%d]", i))
			if obj, ok := item.(map[string]interface{}); ok {
				itemLabel := fmt.Sprintf("[%s]%s[-] %s", colorCyan, escapedIdx, compactColorized(obj, 1))
				itemNode := tview.NewTreeNode(itemLabel).
					SetSelectable(true).
					SetExpanded(expanded).
					SetSelectedTextStyle(selectedStyle).
					SetReference(fmt.Sprintf("[%d]", i))
				addMapToTree(itemNode, obj, expanded)
				node.AddChild(itemNode)
			} else {
				addValueToTree(node, fmt.Sprintf("[%d]", i), item, expanded)
			}
		}
		parent.AddChild(node)
	default:
		label := fmt.Sprintf("[%s]%s[-]: %s", colorPurple, escapedKey, formatValue(value))
		node := tview.NewTreeNode(label).
			SetSelectable(false).
			SetReference(stripColorTags(label))
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

func keyName(r rune) string {
	if r == ' ' {
		return "space"
	}
	return string(r)
}

func setExpandedRecursive(node *tview.TreeNode, expanded bool) {
	node.SetExpanded(expanded)
	for _, child := range node.GetChildren() {
		setExpandedRecursive(child, expanded)
	}
}
