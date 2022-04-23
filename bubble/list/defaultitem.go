package list

import (
	"clx/syntax"
	"fmt"
	"github.com/nleeper/goment"
	"io"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// DefaultItemStyles defines styling for a default list item.
// See DefaultItemView for when these come into play.
type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle  lipgloss.Style
	NormalDesc   lipgloss.Style
	NormalDomain lipgloss.Style

	// The selected item state.
	SelectedTitle  lipgloss.Style
	SelectedDesc   lipgloss.Style
	SelectedDomain lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style
}

// NewDefaultItemStyles returns style definitions for a default item. See
// DefaultItemView for when these come into play.
func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle()
	//Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})
	//Padding(0, 0, 0, 2)

	s.NormalDesc = s.NormalTitle.Copy()
	//Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.NormalDomain = lipgloss.NewStyle().
		Faint(true)

	s.SelectedTitle = lipgloss.NewStyle().
		Bold(true)
	//Border(lipgloss.NormalBorder(), false, false, false, true).
	//BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
	//Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
	//Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.Copy().
		Bold(false)
	//Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	s.SelectedDomain = lipgloss.NewStyle().
		Faint(true)

	s.DimmedTitle = lipgloss.NewStyle()
	//Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
	//Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.Copy()
	//Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	return s
}

// DefaultItem describes an items designed to work with DefaultDelegate.
type DefaultItem interface {
	Item
	Title() string
	Points() int
	User() string
	Time() int64
	Domain() string
	URL() string
	ID() int
	CommentsCount() int
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
	spacing         int
}

// NewDefaultDelegate creates a new delegate with default styles.
func NewDefaultDelegate() DefaultDelegate {
	return DefaultDelegate{
		ShowDescription: true,
		Styles:          NewDefaultItemStyles(),
		spacing:         1,
	}
}

// Height returns the delegate's preferred height.
func (d DefaultDelegate) Height() int {
	if d.ShowDescription {
		return 2 //nolint:gomnd
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
		title, desc, domain string
		s                   = &d.Styles
	)

	if i, ok := item.(DefaultItem); ok {
		title = i.Title()
		domain = i.Domain()
		user := i.User()
		//fmt.Sprintf("%s is %d years old", name, age)
		desc = fmt.Sprintf("%d points by %s %s | %d comments",
			i.Points(), user, parseTime(i.Time()), i.CommentsCount())

	} else {
		return
	}

	// Prevent text from exceeding list width
	if m.width > 0 {
		textWidth := uint(m.width - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
		title = truncate.StringWithTail(title, textWidth, ellipsis)
		desc = truncate.StringWithTail(desc, textWidth, ellipsis)
	}

	// Conditions
	var (
		isSelected = index == m.Index()
	)

	if isSelected {
		//title = s.SelectedTitle.Render(title)
		title = s.SelectedTitle.Render(title)
		title = syntax.HighlightYCStartupsInHeadlines(title)
		title = syntax.HighlightYearInHeadlines(title)
		title = syntax.HighlightHackerNewsHeadlines(title)
		title = syntax.HighlightSpecialContent(title)

		title = title + " " + s.SelectedDomain.Render(domain)

		//desc = s.SelectedDesc.Render(desc)
		desc = s.SelectedDesc.Render(desc)
	} else {
		title = syntax.HighlightYCStartupsInHeadlinesNoBold(title)
		title = syntax.HighlightYearInHeadlinesNoBold(title)
		title = syntax.HighlightHackerNewsHeadlinesNoBold(title)
		title = syntax.HighlightSpecialContent(title)

		title = title + " " + s.NormalDomain.Render(domain)
		//title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc)
		return
	}
	fmt.Fprintf(w, "%s", title)
}

func parseTime(unixTime int64) string {
	moment, _ := goment.Unix(unixTime)
	now, _ := goment.New()

	return moment.From(now)
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
