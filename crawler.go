package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/vence722/convert"

	"github.com/PuerkitoBio/goquery"
)

const URLMatches = "http://www.goal.com/en/matches"

var CompetitionsToRetrieve = []string{
	"Premier League",
	"Primera División",
	"Serie A",
	"Bundesliga",
	"Ligue 1",
	"UEFA Champions League",
	"UEFA Europa League",
	"UEFA Nations League",
	"Friendlies",
	"FIFA Club World Cup",
}

func CrawlMatches() ([]*Match, error) {
	resp, errRequest := http.Get(URLMatches)
	if errRequest != nil {
		return nil, errRequest
	}

	data, errReadBody := ioutil.ReadAll(resp.Body)
	if errReadBody != nil {
		return nil, errReadBody
	}
	resp.Body.Close()

	matches, errParse := parseMatches(string(data))
	if errReadBody != nil {
		return nil, errParse
	}

	return matches, nil
}

func parseMatches(pageData string) ([]*Match, error) {
	var matches []*Match

	doc, errParse := goquery.NewDocumentFromReader(strings.NewReader(pageData))
	if errParse != nil {
		return nil, errParse
	}

	doc.Find("div.competition-matches").Each(func(index int, elem *goquery.Selection) {
		competitionName := elem.Find("div.competition-name").Text()
		// If it's the competition that we care
		if contains(CompetitionsToRetrieve, competitionName) {
			elem.Find("div.match-row").Each(func(index int, elem *goquery.Selection) {
				match := &Match{}
				match.Competition = competitionName

				// Retrieve match time
				value, exists := elem.Find("div.match-status > time").Attr("datetime")
				if exists {
					match.Time, _ = time.Parse("2006-01-02T15:04:05+00:00", value)
				}

				// Retrieve if match is finished
				if "FT" == elem.Find("div.match-status > span").Text() {
					match.IsFinished = true
				}

				// Retrieve match teams and score info
				match.HomeTeam = elem.Find("div.team-home > span.team-name").Text()
				match.HomeScore = convert.Str2Int(elem.Find("div.team-home > span.goals").Text())
				match.AwayTeam = elem.Find("div.team-away > span.team-name").Text()
				match.AwayScore = convert.Str2Int(elem.Find("div.team-away > span.goals").Text())

				matches = append(matches, match)
			})
		}
	})

	return matches, nil
}

func contains(slice []string, searchValue string) bool {
	for _, value := range slice {
		if value == searchValue {
			return true
		}
	}
	return false
}
