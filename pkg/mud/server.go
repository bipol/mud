package mud

import (
	"fmt"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SessionTerminal struct {
	session  ssh.Session
	terminal *terminal.Terminal
}

type Server struct {
	sts map[string]SessionTerminal
}

func New() *Server {
	m := make(map[string]SessionTerminal)
	return &Server{m}
}

func PromptLogin(st SessionTerminal) (*Player, error) {
	line := ""
	_, err := io.WriteString(st.terminal, fmt.Sprint("Enter your character name and press enter.\n"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, _ = st.terminal.ReadLine()
		if line == "quit" {
			return nil, nil
		} else if line == "" {
			continue
		}
		return &Player{line, st}, nil
	}
}

type Player struct {
	name string
	st   SessionTerminal
}

func (m *Server) JoinSession(session ssh.Session) *Player {
	term := terminal.NewTerminal(session, "> ")
	st := SessionTerminal{session, term}
	player, err := PromptLogin(st)
	if err != nil {
		log.Fatal(err)
	}
	// if player isn't found they have quit i guess
	// that is probably ugly
	if player == nil {
		return nil
	}
	term.SetPrompt(fmt.Sprintf("%s%s: %s", term.Escape.Green, player.name, term.Escape.Reset))
	m.sts[player.name] = st
	return player
}

func (m *Server) LeaveSession(p *Player) {
	delete(m.sts, p.name)
}

// a user should be able to ssh onto the server
// and then login to a specific character

//so the handler is still a local context
// a bit of state that runs every time a person connects to the ssh server
// i need to split login "character creation" and terminal creation
func (m *Server) Start(p *Player) {
	line := ""
	_, err := io.WriteString(p.st.terminal, fmt.Sprint("Type 'quit' to quit\n"))
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, _ = p.st.terminal.ReadLine()
		if line == "quit" {
			m.LeaveSession(p)
			break
		} else if line == "" {
			continue
		}
		// go through all the player sessions in the world and send the message to them
		for name, st := range m.sts {
			if p.name == name {
				continue
			}
			log.Printf("writing to %s's session", st.session.User())
			_, err := io.WriteString(st.terminal, fmt.Sprintf("%s: %s\n", p.name, line))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
