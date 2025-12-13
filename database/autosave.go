package database

import (
	"log"
	"sync"
	"time"
)

var update_interval = 3 * time.Second

type AutoSave struct {
    mu        	sync.Mutex
    changed     map[uint]bool
    done      	chan struct{}
}

func newAutoSave() *AutoSave {
	return &AutoSave{
		changed: make(map[uint]bool),
		done: make(chan struct{}),
	}
}

var A = newAutoSave()

func (a *AutoSave) Start() {
	go func() {
		ticker := time.NewTicker(update_interval)
		defer ticker.Stop()

		for {
			select {
			case <-a.done:
				return
			case <-ticker.C:
				a.save()
			}
		}
	}()
}

func (a *AutoSave) Stop() {
	close(a.done)
}

func (a *AutoSave) save() {
	log.Println("start saving")

    a.mu.Lock()
    users := make([]uint, 0, len(a.changed))
    for userID := range a.changed {
        users = append(users, userID)
    }

    a.changed = make(map[uint]bool)
    a.mu.Unlock()

    for _, userID := range users {
        session := GetSessionState(userID)
        SaveSession(session)
    }
}

func (a *AutoSave) MarkChanged(userID uint) {
    a.mu.Lock()
    a.changed[userID] = true
    a.mu.Unlock()
}