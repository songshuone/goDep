package service





type SecProductInfoConf struct {
	ProductId int `json:"product_id"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	Status    int `json:"status"`
	Total     int `json:"total"`
	Left      int `json:"left"`
	OnePersonBuyLimit int `json:"one_person_buy_limit"`
	BuyRate           float64 `json:"buy_rate"`
	//每秒最多能卖多少个
	SoldMaxLimit int `json:"sold_max_limit"`
}