package main

type Scope int

const (
	Global Scope = iota
	Local
)

func (s Scope) String() string {
	return []string{"--global", "--local"}[s]
}
