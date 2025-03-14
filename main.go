//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	timer         timer.Model
	userText      textinput.Model
	timeout       bool
	levels        [4]string
	level         int
	selectedLevel bool
	restart       bool
}

var (
	badChar = 0
	str     string
	texts   = []string{
		"lorem ipsum dolor , dolor sit amet , dolor sit amet neque , dolor sit amet neque nisi , dolor sit amet neque nisi nisi",
		"the quick brown fox jumps over the lazy dog and the lazy dog jumps ; this is the last thing that we will do together ; i will be so rgatefulll to save myself",
		"pack my box with five dozen liquor jugs , ninkynanka : hello wuddy what would you like to do wzap me : guys l'est pas nothing to oyu an pertex revive em all you should do is finish with this stupid thing",
		"l'est pas nothing to oyu an pertex revive em all you should do is finish with this stupid thing yokal wzap me : guys l'est pas nothing to oyu an pertex revive em all you ? :) aller au marche pas nothing to oyu an pertex revive em all you should do is finish with this stupid thing",
	}
	levelBox = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#100a8c")).Padding(1, 2).Width(71)
)

const timeout = time.Second * 30

var startTime time.Time

func initModel() model {
	userText := textinput.New()

	return model{
		timer:         timer.NewWithInterval(timeout, time.Millisecond),
		userText:      userText,
		timeout:       false,
		levels:        [4]string{"Beginner", "Medium", "Advanced", "Guuru"},
		level:         0,
		selectedLevel: false,
		restart:       false,
	}
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if strings.Compare(m.userText.Value(), "") == -1 && strings.Compare(m.userText.Value(), texts[m.level]) == 0 && m.level != -1 {
		if m.level < len(texts)-1 {
			m.userText.SetValue("")
			m.userText.CharLimit = len(texts[m.level])
			m.timeout = false
			m.timer = timer.NewWithInterval(timeout, time.Millisecond)
			return m, m.timer.Init()
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.timeout = true
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if m.timer.Timedout() {

				m.restart = true
				m.userText.SetValue("")
				m.selectedLevel = false
				m.level = 0
				m.timer = timer.NewWithInterval(timeout, time.Millisecond)
				return m, cmd
			}

		case "up":
			if !m.selectedLevel {

				if m.level > 0 {
					m.level--
				} else {
					m.level = len(m.levels) - 1
				}
			}

		case "down":
			if !m.selectedLevel {

				if m.level < len(m.levels)-1 {
					m.level++
				} else {
					m.level = 0
				}
			}
		case "enter":
			m.selectedLevel = true
			startTime = time.Now()
		case "ctrl+c":
			return m, tea.Quit

		}
	}
	m.userText.Focus()
	m.userText, cmd = m.userText.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.restart {
		return m.getLevelDisplay()
	}

	if !m.selectedLevel {
		return m.getLevelDisplay()
	}

	if !m.timeout {
		str, badChar = m.getText()
		return str
	}

	return m.getResult(badChar, str)
}

func (m model) getLevelDisplay() string {
	s := ""

	selectedBox := lipgloss.NewStyle().Foreground(lipgloss.Color("#09e322"))

	for i, level := range m.levels {
		if i == m.level {
			s += selectedBox.Render(fmt.Sprintf("> %s", level))
		} else {
			s += level
		}

		s += "\n"
	}

	return levelBox.Render(s)

}

func (m model) getText() (string, int) {
	s := ""

	currentText := strings.Split(m.userText.Value(), "")
	correctText := lipgloss.NewStyle().Foreground(lipgloss.Color("#0bfc03"))
	badText := lipgloss.NewStyle().Foreground(lipgloss.Color("#A44043"))
	greyText := lipgloss.NewStyle().Foreground(lipgloss.Color("#CFE4E8"))
	TextBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#100a8c")).Padding(1, 2).Width(71)
	TimerBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderBottomForeground(lipgloss.Color("#09e322")).Padding(0, 1).Width(10)
	RedTimerBox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderBottomForeground(lipgloss.Color("#e31409")).Padding(0, 1).Width(10)

	if time.Since(startTime) < time.Second*20 {
		s += TimerBox.Render(m.timer.View()) + "\n\n"
	} else {
		s += RedTimerBox.Render(m.timer.View()) + "\n\n"
	}

	badChar := 0

	for i, char := range strings.Split(texts[m.level], "") {
		if i < len(currentText) {
			if strings.Compare(char, currentText[i]) == 0 {
				s += correctText.Render(char)
			} else {
				s += badText.Render(char)
				badChar++
			}
		} else {
			s += greyText.Render(char)
		}
	}

	s += "\n\n"

	if !m.timeout {
		s += m.userText.View()
	}

	return TextBox.Render(s), badChar
}

func (m model) getResult(badChar int, str string) string {
	s :=
		levelBox.Render(fmt.Sprintf(`
	%s  ----------------------------- %d
	%s  ----------------------------- %s
	`, "Bad Character type", badChar, "Level ", m.levels[m.level]))

	s += "\n\n"

	s += str

	s += "\n"

	s += " Tapez a pour recommencer  ----------------------------- Tapez x pour quitter\n"

	return s
}

func main() {
	fmt.Println(`
,-----.          ,--.                 ,--------.                       
|  |) /_  ,--,--.|  |,-. ,---. ,-----.'--.  .--',--. ,--.,---.  ,---.  
|  .-.  \' ,-.  ||     /| .-. :'-----'   |  |    \  '  /| .-. || .-. : 
|  '--' /\ '-'  ||  \  \\   --.          |  |     \   ' | '-' '\   --. 
"------'  "--"--'"--'"--'"----'          "--'   .-'  /  |  |-'  "----' 
                                                "---'   "--'           
    `)

	prog := tea.NewProgram(initModel())

	if _, err := prog.Run(); err != nil {
		fmt.Println("Failed to initialize" + err.Error())
	}
}
