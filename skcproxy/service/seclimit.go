package service

type Limit struct {
	secLimit TimeLimit
	minLimit TimeLimit
}

type SecLimit struct {
	count   int
	curTime int64
}

func (s *SecLimit) Count(nowTime int64) (curCount int) {
	if s.curTime != nowTime {
		s.count = 1
		s.curTime = nowTime
		curCount = s.count
		return
	}
	s.count ++
	curCount = s.count
	return
}

func (s *SecLimit) Check(now int64) int {
	if s.curTime != now {
		return 0
	}
	return s.count
}
