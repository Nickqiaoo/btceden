package biz

type ChainTVL struct {
	GeckoId     string  `json:"gecko_id"`
	Tvl         float64 `json:"tvl"`
	TokenSymbol string  `json:"tokenSymbol"`
	CmcId       string  `json:"cmcId"`
	Name        string  `json:"name"`
}
