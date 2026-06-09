package services

type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency                   string  `json:"currency"`
				Symbol                     string  `json:"symbol"`
				RegularMarketPrice         float64 `json:"regularMarketPrice"`
				RegularMarketChange        float64 `json:"regularMarketChange"`
				RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
			} `json:"meta"`
			Timestamp []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []*float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}
