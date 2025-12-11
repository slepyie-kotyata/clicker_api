package ws

import (
	"clicker_api/database"
	"clicker_api/models"
	"clicker_api/service"
	"encoding/json"
	"log"
	"math"
	"time"
)

type PassiveWorker struct {
    ticker *time.Ticker
    done chan struct{}
}

var seconds uint = 3
var P = NewPassiveWorker()

func NewPassiveWorker() *PassiveWorker {
    return &PassiveWorker{
        ticker: time.NewTicker(3 * time.Second),
        done: make(chan struct{}),
    }
}

func (p *PassiveWorker) Start() {
	log.Println("passive started")
    go func() {
        for {
            select {
            case <-p.ticker.C:
				users := H.GetActiveUsers()
    			if len(users) == 0 {
					continue
    			}
				
    			for _, id := range users {
        			p.updateSessionState(id)
    			}
            case <-p.done:
                return
            }
        }
    }()
}

func (p *PassiveWorker) Stop() {
    p.ticker.Stop()
    close(p.done)
}

func (p *PassiveWorker) updateSessionState(id uint) {
    session := database.GetSessionState(id)
    if session == nil {
        return
    }

    upgrade_stats := service.CountBoostValues(service.FilterUpgrades(session, true))

	if upgrade_stats.MpS == 0 && upgrade_stats.DpS == 0 {
		log.Println("no passive updates")
		return
	}

	log.Println("init passive")

	prestigeUpgrade(session, upgrade_stats, seconds)
	passiveSellUpdate(session, upgrade_stats, seconds)
	passiveCookUpdate(session, upgrade_stats, seconds)

	log.Println("done")

	database.SaveSessionState(id, session)

	data, _ := json.Marshal(map[string]interface{}{
		"money": session.Money,
		"dishes": session.Dishes,
		"level_rank": session.LevelRank,
		"level_xp": session.LevelXP,
		"prestige_accumulated": session.PrestigeAccumulated,
	})

    H.incoming <- HubEvent{
    	Type:     BroadcastToConnection,
    	UserID:   id,
    	Session:  nil,
    	Message:  Message{
			MessageType: Response,
			RequestID:   "",
			RequestType: PassiveRequest,
			Data:        data,
		},
	}
}

func passiveSellUpdate(session *models.SessionState, upgrade_stats service.UpgradeStats, seconds uint) {
	if session.Dishes <= 0 || session.Dishes < 3 {
		return 
	}

	if upgrade_stats.MpS == 0 {
		return
	}

	prestige_boost := session.PrestigeAccumulated

	if prestige_boost == 0 {
		prestige_boost = 1
	}

	service.SetDefaults(&upgrade_stats)

	minNum := min((float64(seconds) * upgrade_stats.SpS), float64(session.Dishes))

	if session.LevelRank < 100 {
		session.LevelXP = math.Round((session.LevelXP + math.Abs(0.05 * float64(seconds) * upgrade_stats.MpS)) * 100) / 100
	}

	session.Money += uint(math.Ceil(upgrade_stats.MpS * upgrade_stats.MpM * float64(seconds) * prestige_boost * minNum))
	session.Dishes -= uint(math.Ceil(minNum))
}

func passiveCookUpdate(session *models.SessionState, upgrade_stats service.UpgradeStats, seconds uint) {
	if upgrade_stats.DpS == 0 {
		return
	}

	prestige_boost := session.PrestigeAccumulated

	if prestige_boost == 0 {
		prestige_boost = 1
	}

	service.SetDefaults(&upgrade_stats)

	if session.LevelRank < 100 {
		session.LevelXP = math.Round((session.LevelXP + math.Abs(0.2 * float64(seconds) * upgrade_stats.DpS)) * 100) / 100
	}

	session.Dishes += uint(math.Ceil(upgrade_stats.DpS * upgrade_stats.DpM * float64(seconds) * prestige_boost))
}

func prestigeUpgrade(session *models.SessionState, upgrade_stats service.UpgradeStats, seconds uint) {
	if upgrade_stats.MpS == 0 {
		return
	}

	service.SetDefaults(&upgrade_stats)

	d := upgrade_stats.MpS * upgrade_stats.MpM
	p := (d / 10000) * float64(seconds)
	p = math.Round(p * 10000) / 10000

	session.PrestigeAccumulated += p
}

