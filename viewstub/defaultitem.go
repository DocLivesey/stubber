package viewstub

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	// The Normal state.
	NormalName     lipgloss.Style
	NormalPath     lipgloss.Style
	NormalStateOn  lipgloss.Style
	NormalStateOff lipgloss.Style
	NormalPid      lipgloss.Style

	// The selected item state.
	SelectedName     lipgloss.Style
	SelectedPath     lipgloss.Style
	SelectedStateOn  lipgloss.Style
	SelectedStateOff lipgloss.Style
	SelectedPid      lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedName     lipgloss.Style
	DimmedPath     lipgloss.Style
	DimmedStateOn  lipgloss.Style
	DimmedStateOff lipgloss.Style
	DimmedPid      lipgloss.Style

	// Charcters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalName = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2)

	s.NormalStateOn = s.NormalName.Copy().Foreground(lipgloss.Color("#0DE291")).PaddingLeft(2)
	s.NormalStateOff = s.NormalName.Copy().Foreground(lipgloss.Color("#E71616")).PaddingLeft(2)

	s.NormalPath = s.NormalName.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.NormalPid = s.NormalPath.Copy()

	s.SelectedName = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
		Padding(0, 0, 0, 1)

	s.SelectedStateOn = s.SelectedName.Copy().Foreground(lipgloss.Color("#0DE291")).PaddingLeft(1)
	s.SelectedStateOff = s.SelectedName.Copy().Foreground(lipgloss.Color("#e71616")).PaddingLeft(1)

	s.SelectedPath = s.SelectedName.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	s.SelectedPid = s.SelectedPath.Copy()

	s.DimmedName = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2)

	s.DimmedStateOn = s.DimmedName.Copy().Foreground(lipgloss.Color("#7BE2BB")).PaddingLeft(2)
	s.DimmedStateOff = s.DimmedName.Copy().Foreground(lipgloss.Color("#DB7272")).PaddingLeft(2)

	s.DimmedPath = s.DimmedName.Copy().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.DimmedPid = s.DimmedPath.Copy()

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

// DefaultItem describes an items designed to work with DefaultDelegate.
type DefaultItem interface {
	Item
	Jar() string
	Path() string
	State() string
	Pid() string
}

// DefaultDelegate is a standard delegate designed to work in lists. It's
// styled by DefaultItemStyles, which can be customized as you like.
//
// The description line can be hidden by setting Description to false, which
// renders the list as single-line-items. The spacing between items can be set
// with the SetSpacing method.
//
// Setting UpdateFunc is optional. If it's set it will be called when the
// ItemDelegate called, which is called when the list's Update function is
// invoked.
//
// Settings ShortHelpFunc and FullHelpFunc is optional. They can can be set to
// include items in the list's default short and full help menus.
type DefaultDelegate struct {
	ShowDescription bool
	Styles          DefaultItemStyles
	UpdateFunc      func(tea.Msg, *Model) tea.Cmd
	ShortHelpFunc   func() []key.Binding
	FullHelpFunc    func() [][]key.Binding
	height          int
	spacing         int
}

// NewDefaultDelegate creates a new delegate with default styles.
func NewDefaultDelegate() DefaultDelegate {
	return DefaultDelegate{
		ShowDescription: true,
		Styles:          NewDefaultItemStyles(),
		height:          2,
		spacing:         1,
	}
}

// SetHeight sets delegate's preferred height.
func (d *DefaultDelegate) SetHeight(i int) {
	d.height = i
}

// Height returns the delegate's preferred height.
// This has effect only if ShowDescription is true,
// otherwise height is always 1.
func (d DefaultDelegate) Height() int {
	if d.ShowDescription {
		return d.height
	}
	return 1
}

// SetSpacing set the delegate's spacing.
func (d *DefaultDelegate) SetSpacing(i int) {
	d.spacing = i
}

// Spacing returns the delegate's spacing.
func (d DefaultDelegate) Spacing() int {
	return d.spacing
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (d DefaultDelegate) Update(msg tea.Msg, m *Model) tea.Cmd {
	if d.UpdateFunc == nil {
		return nil
	}
	return d.UpdateFunc(msg, m)
}

// Render prints an item.
func (d DefaultDelegate) Render(w io.Writer, m Model, index int, item Item) {
	var (
		name, path, state, pid string
		matchedRunes           []int
		s                      = &d.Styles
	)

	if i, ok := item.(DefaultItem); ok {
		name = i.Jar()
		path = i.Path()
		state = i.State()
		pid = "Pid:" + i.Pid()
	} else {
		return
	}

	if m.width <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := uint(m.width - s.NormalName.GetPaddingLeft() - s.NormalName.GetPaddingRight())
	name = truncate.StringWithTail(name, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(path, "\n") {
			if i >= d.height-1 {
				break
			}
			lines = append(lines, truncate.StringWithTail(line, textwidth, ellipsis))
		}
		path = strings.Join(lines, "\n")
	}

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == Filtering || m.FilterState() == FilterApplied
	)

	if isFiltered && index < len(m.filteredItems) {
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	if emptyFilter {
		name = s.DimmedName.Render(name)
		path = s.DimmedPath.Render(path)
		if state == ("On") {
			state = s.DimmedStateOn.Render(state)
		} else {
			state = s.DimmedStateOff.Render(state)
		}
		pid = s.DimmedPid.Render(pid)
	} else if isSelected && m.FilterState() != Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedName.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			name = lipgloss.StyleRunes(name, matchedRunes, matched, unmatched)
		}
		name = s.SelectedName.Render(name)
		path = s.SelectedPath.Render(path)
		if state == ("On") {
			state = s.SelectedStateOn.Render(state)
		} else {
			state = s.SelectedStateOff.Render(state)
		}
		pid = s.SelectedPid.Render(pid)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalName.Inline(true)
			matched := unmatched.Copy().Inherit(s.FilterMatch)
			name = lipgloss.StyleRunes(name, matchedRunes, matched, unmatched)
		}
		name = s.NormalName.Render(name)
		path = s.NormalPath.Render(path)
		if state == ("On") {
			state = s.NormalStateOn.Render(state)
		} else {
			state = s.NormalStateOff.Render(state)
		}
		pid = s.NormalPid.Render(pid)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s\n%s  %s", name, state, path, pid)
		return
	}
	fmt.Fprintf(w, "%s", name)
}

// ShortHelp returns the delegate's short help.
func (d DefaultDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

// FullHelp returns the delegate's full help.
func (d DefaultDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
