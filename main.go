package main

import (
	"log"

	"github.com/bipol/bipolmud/pkg/mud"
	"github.com/gliderlabs/ssh"
)

func main() {
	client := mud.New()
	ssh.Handle(func(s ssh.Session) {
		p := client.JoinSession(s)
		client.Start(p)
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
