package main


type LoadingState int

// Enums

const (
	Phase1 LoadingState = iota
	Phase2
	Phase3
	Loaded
)