package models

// Scope describes the git config scope (local or global)
type Scope int

const (
	// Global defines the global git config scope
	Global Scope = iota

	// Local defines the local git config scope
	Local
)

// String returns the string representation of the scope
func (s Scope) String() string {
	return []string{"global", "local"}[s]
}

// Arg returns the argument representation of the scope
func (s Scope) Arg() string {
	return []string{"--global", "--local"}[s]
}
