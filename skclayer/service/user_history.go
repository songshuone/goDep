package service

import "sync"

type UserBuyHistory struct {
	history map[int]int
	lock    sync.Mutex
}

func (h *UserBuyHistory) GetProductBuyCount(productId int) int {
	h.lock.Lock()
	defer h.lock.Unlock()
	count, _ := h.history[productId]
	return count
}

func (h *UserBuyHistory) Add(productId, count int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	cur, ok := h.history[productId]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	h.history[productId] = cur
}
