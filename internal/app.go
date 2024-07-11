package internal

import "github.com/praveenmahasena/aiserver/internal/listner"

func Start() error {
	p := listner.New(":42069")
	return p.Run()
}
