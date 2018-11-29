package exec

import "os/exec"

//Command is an interface that hides os/exec functionality
type Command interface {
	Run(string, ...string) error
}

//CommandImpl is the real implementation of COmmand interface
type CommandImpl struct{}

//Run is a method that wraps Run() from os/exec package
func (r CommandImpl) Run(command string, args ...string) error {
	return exec.Command(command, args...).Run()
}
