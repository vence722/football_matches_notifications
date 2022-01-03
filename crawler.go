package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/vence722/convert"

	"github.com/PuerkitoBio/goquery"
)

const URLMatches = "https://www.goal.com/en/live-scores"

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

var CompetitionsNameTranslations = map[string]string{
	"Premier League":        "英超",
	"Primera División":      "西甲",
	"Serie A":               "意甲",
	"Bundesliga":            "德甲",
	"Ligue 1":               "法甲",
	"UEFA Champions League": "歐聯盃",
	"UEFA Europa League":    "歐霸盃",
	"UEFA Nations League":   "歐洲國家聯賽",
	"Friendlies":            "友誼賽",
	"FIFA Club World Cup":   "世界盃",
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
				match.Competition = CompetitionsNameTranslations[competitionName]

				// Retrieve match time
				value, exists := elem.Find("div.match-status > time").Attr("data-utc")
				if exists {
					match.Time = time.Unix(convert.Str2Int64(value)/1000, 0)
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
