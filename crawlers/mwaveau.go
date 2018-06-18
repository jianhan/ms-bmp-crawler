package crawlers

type mwaveau struct {
	base
}

func NewMwaveau(testMode bool) Scraper {
	b := base{
		homepageURL: "https://www.mwave.com.au/",
		name:        "Mwave Australia",
		categoryURL: "https://www.mwave.com.au/",
		currency:    "AUD",
		testMode:    testMode,
	}

	return &mwaveau{b}
}
