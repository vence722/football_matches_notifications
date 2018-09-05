package main

import (
	"fmt"
	"time"
)

type Match struct {
	Competition string
	HomeTeam    string
	AwayTeam    string
	HomeScore   int
	AwayScore   int
	Time        time.Time
	IsFinished  bool
}

func (this *Match) String() string {
	return fmt.Sprintln(*this)
}
