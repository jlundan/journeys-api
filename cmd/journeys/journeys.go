package main

import "github.com/jlundan/journeys-api/internal/app/journeys"

func main() {
	_ = journeys.MainCommand.Execute()
}
