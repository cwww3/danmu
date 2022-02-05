package repository

import (
	"danmu/model"
	"sync"
)

type LivingMap struct {
	sync.Mutex
	m map[string]*model.Room
}

var livingMap = &LivingMap{
	m: make(map[string]*model.Room),
}

func GetMemoryRepository() *LivingMap {
	return livingMap
}

func (lm *LivingMap) Save(user string, room *model.Room) {
	lm.Lock()
	defer lm.Unlock()
	lm.m[user] = room
}

func (lm *LivingMap) Delete(user string) {
	lm.Lock()
	defer lm.Unlock()
	delete(lm.m, user)
}

func (lm *LivingMap) GetList() map[string]*model.Room {
	lm.Lock()
	defer lm.Unlock()
	return lm.m
}
