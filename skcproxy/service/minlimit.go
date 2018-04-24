package service

type TimeLimit interface {
	Count(nowTime int64) (curCount int)
	Check(now int64) int
}

type MinLimit struct {
	count   int
	curTime int64
}

//记录在同一时间段访问接口的次数
func (m *MinLimit) Count(nowTime int64) (curCount int) {
	if nowTime-m.curTime > 60 {
		m.count = 1
		curCount = m.count
		m.curTime = nowTime
		return
	}
	m.count ++
	curCount = m.count
	return
}

func (m *MinLimit) Check(now int64) int {
	if now-m.curTime > 60 {
		return 0
	}
	return m.count
}
