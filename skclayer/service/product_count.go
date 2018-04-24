package service

import "sync"

type ProductCountMgr struct {
	productCount map[int]int
	lock         sync.Mutex
}

func NewProductCountMgr() (mgr *ProductCountMgr) {
	mgr = &ProductCountMgr{productCount: make(map[int]int, 128)}
	return
}

func (mgr *ProductCountMgr) Count(productId int) (count int) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()
	count = mgr.productCount[productId]
	return
}

func (mgr *ProductCountMgr) Add(productId, count int) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	curCount, ok := mgr.productCount[productId]
	if !ok {
		curCount = count
	} else {
		curCount += count
	}
	mgr.productCount[productId] = curCount

}
