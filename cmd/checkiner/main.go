package main

import "os"

func main() {
	cmd := NewCheckinerCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
