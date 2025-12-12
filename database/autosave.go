package database

import (
	"sync"
	"time"
)

var update_interval = 5 * time.Second

type AutoSave struct {
    mu        	sync.Mutex
    changed     map[uint]struct{}
    done      	chan struct{}
}

func newAutoSave() *AutoSave {
	return &AutoSave{
		changed: make(map[uint]struct{}),
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
    a.mu.Lock()
    users := make([]uint, 0, len(a.changed))
    for userID := range a.changed {
        users = append(users, userID)
    }

    a.changed = make(map[uint]struct{})
    a.mu.Unlock()

    for _, userID := range users {
        session := GetSessionState(userID)
        SaveSession(session)
    }
}

func (a *AutoSave) MarkChanged(userID uint) {
    a.mu.Lock()
    a.changed[userID] = struct{}{}
    a.mu.Unlock()
}