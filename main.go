package main

func main() {
	err := StartTelegramBot(CrawlMatches)
	if err != nil {
		panic(err)
	}
}
