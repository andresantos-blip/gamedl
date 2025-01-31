package main

const MaxConcurrency = 2

type GameProcessReport struct {
	Err  error
	Id   string
	Year int
}
