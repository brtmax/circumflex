package bubble

import (
	"clx/bubble/list"
	"clx/cli"
	"clx/comment"
	"clx/core"
	"clx/hn/services/mock"
	"clx/screen"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"time"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, user, domain, url  string
	id, points, commentsCount int
	time                      int64
}

func (i item) Title() string      { return i.title }
func (i item) User() string       { return i.user }
func (i item) Domain() string     { return i.domain }
func (i item) Points() int        { return i.points }
func (i item) CommentsCount() int { return i.commentsCount }
func (i item) Time() int64        { return i.time }
func (i item) URL() string        { return i.url }
func (i item) ID() int            { return i.id }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type editorFinishedMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "e" {
			cmd := m.list.NewStatusMessageWithDuration("Test", 2*time.Second)

			return m, cmd
		}
		if msg.String() == "f" {
			cmd := m.list.NewStatusMessageWithDuration("ABCDEF", 1*time.Second)

			return m, cmd
		}
		if msg.String() == "enter" {
			if i, ok := m.list.SelectedItem().(item); ok {
				id := i.ID()
				cmd := openEditor(id)

				return m, cmd
			}

			return m, nil
		}
		if msg.String() == "u" {
			dot := spinner.Spinner{
				Frames: []string{"⣾ ", "⣷ ", "⣯ ", "⣟ ", "⡿ ", "⢿ ", "⣻ ", "⣽ "},
				FPS:    time.Second / 7, //nolint:gomnd
			}

			m.list.SetSpinner(dot)
			cmd := m.list.StartSpinner()

			return m, cmd
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func openEditor(id int) tea.Cmd {
	comments := new(mock.Service).FetchStory(id)

	screenWidth := screen.GetTerminalWidth()
	commentTree := comment.ToString(comments, core.GetConfigWithDefaults(), screenWidth, 0)

	c := cli.WrapLess(commentTree)

	return tea.Exec(tea.WrapExecCommand(c), func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func Run() {
	var items []list.Item

	service := new(mock.Service)
	stories := service.FetchStories(0, 0)

	for _, story := range stories {
		items = append(items, item{
			title:  story.Title,
			domain: story.Domain,
			user:   story.User,
			url:    story.URL,
			time:   story.Time,
			points: story.Points,
			id:     story.ID})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	cli.ClearScreen()
}
