package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-rogue/scenes"
	"log"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

type Scene interface {
	scenes.IScene
	HandleEvent(ev tcell.Event) bool
	Draw()
}

type WelcomeScene struct {
	scenes.Scene
	screen tcell.Screen
}

type GameScene struct {
	scenes.Scene
	screen tcell.Screen
}

type SettingsScene struct {
	scenes.Scene
	screen tcell.Screen
}

func (s *WelcomeScene) Draw() {
	drawText(s.screen, 0, 1, 80, 1, tcell.StyleDefault, "[N]ew Game")
	drawText(s.screen, 0, 2, 80, 2, tcell.StyleDefault, "[S]ettings")
}

func (s *GameScene) Draw() {
	drawText(s.screen, 0, 1, 80, 1, tcell.StyleDefault, "[S]ettings")
}

func (s *SettingsScene) Draw() {
	drawText(s.screen, 0, 1, 80, 1, tcell.StyleDefault, "[B]ack")
}

func (s *WelcomeScene) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Key() == tcell.KeyRune {
			switch e.Rune() {
			case 'Q', 'q':
				s.Director.ShouldQuit = true
				return true
			case 'N', 'n':
				s.screen.Clear()
				s.Director.ChangeState(NewGameScene(s.screen))
				return true
			case 'S', 's':
				s.screen.Clear()
				s.Director.PushState(NewSettingsScene(s.screen))
				return true
			}
		}
	}

	return false
}

func (s *GameScene) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Key() == tcell.KeyRune {
			switch e.Rune() {
			case 'Q', 'q':
				s.Director.ShouldQuit = true
				return true
			case 'S', 's':
				s.screen.Clear()
				s.Director.PushState(NewSettingsScene(s.screen))
				return true
			}
		}
	}

	return false
}

func (s *SettingsScene) HandleEvent(ev tcell.Event) bool {
	switch e := ev.(type) {
	case *tcell.EventKey:
		if e.Key() == tcell.KeyRune {
			switch e.Rune() {
			case 'Q', 'q':
				s.Director.ShouldQuit = true
				return true
			case 'B', 'b':
				s.screen.Clear()
				s.Director.PopState()
				return true
			}
		}
	}

	return false
}

func NewWelcomeScene(screen tcell.Screen) *WelcomeScene {
	return &WelcomeScene{
		scenes.Scene{
			Name: "Welcome",
		},
		screen,
	}
}

func NewGameScene(screen tcell.Screen) *GameScene {
	return &GameScene{
		scenes.Scene{
			Name: "Game",
		},
		screen,
	}
}

func NewSettingsScene(screen tcell.Screen) *SettingsScene {
	return &SettingsScene{
		scenes.Scene{
			Name: "Settings",
		},
		screen,
	}
}

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	screen.Clear()

	director := scenes.NewDirector(NewWelcomeScene(screen))
	for director.ShouldQuit == false {
		scene := director.PeekState().(Scene)
		scene.Draw()

		screen.Show()

		ev := screen.PollEvent()
		scene.HandleEvent(ev)
	}
}
