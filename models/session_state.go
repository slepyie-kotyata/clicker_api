package models

type SessionState struct {
    ID        				uint			`json:"id"`
    Money     				uint			`json:"money"`
    Dishes    				uint			`json:"dishes"`
    LevelRank 				uint			`json:"level_rank"`
    LevelXP   				float64			`json:"level_xp"`
    PrestigeCurrent        	float64			`json:"prestige_current"`
    PrestigeBoost          	float64			`json:"prestige_boost"`
    PrestigeAccumulated    	float64			`json:"prestige_accumulated"`
    Upgrades 				map[uint]uint	`json:"upgrades"`
}