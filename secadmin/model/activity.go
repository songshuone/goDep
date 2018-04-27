package model

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type ActivityModel struct {
}

type Activity struct {
	ActivityId   int `gorm:"primary_key"`
	ActivityName string
	ProductId    int
	StartTime    int64
	EndTime      int64
	Total        int
	Status       int

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
	Speed        float64 `gorm:"column:sec_speed"`
	BuyLimit     int `gorm:"column:buy_limit"`
}

type SecProductInfoConf struct {
	ProductId         int     `json:"product_id"`
	StartTime         int64   `json:"start_time"`
	EndTime           int64   `json:"end_time"`
	Status            int     `json:"status"`
	Total             int     `json:"total"`
	Left              int     `json:"left"`
	OnePersonBuyLimit int     `json:"one_person_buy_limit"`
	BuyRate           float64 `json:"buy_rate"`
	//每秒最多能卖多少个
	SoldMaxLimit int `json:"sold_max_limit"`
}

func NewActivityModel() (*ActivityModel) {
	return &ActivityModel{}
}

func (a *ActivityModel) CreateActivity(activity *Activity) (err error) {
	Db.Create(activity)
	return
}

func (p *ActivityModel) GetActivityList() (activity []*Activity, err error) {
	err = Db.Find(&activity).Error
	return
}
