package repl

type CLICommand struct {
	Name        string
	Description string
	Callback    func() error
}
