package handlers

import (
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"math"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ThisUpgrade struct {
	ID             uint		       			`json:"id" gorm:"primary_key"`
	Name           string		   			`json:"name"`
	IconName       string		   			`json:"icon_name"`
	UpgradeType    models.UpgradeType     	`json:"upgrade_type"`
	PriceFactor    float64		   			`json:"price_factor"`
	Price          uint			   			`json:"price"`
	AccessLevel    uint			   			`json:"access_level"`
	Boost          models.Boost    			`json:"boost"`
	TimesBought    uint		   	   			`json:"times_bought"`
}

func FilterUpgrades(session models.Session, is_bought bool) []ThisUpgrade {
	filtered_upgrades := make([]ThisUpgrade, 0)

	var session_upgrades []models.SessionUpgrade
	DB.Where("session_id = ?", session.ID).Find(&session_upgrades)

	times_bought_map := make(map[uint]uint)
	for _, su := range session_upgrades {
		times_bought_map[su.UpgradeID] = su.TimesBought
	}

	for _, upgrade := range session.Upgrades {
		times_bought, ok := times_bought_map[upgrade.ID]

		this_upgrade := ThisUpgrade {
			ID: upgrade.ID,
			Name: upgrade.Name,
			IconName: upgrade.IconName,
			UpgradeType: upgrade.UpgradeType,
			PriceFactor: upgrade.PriceFactor,
			Price: upgrade.Price,
			AccessLevel: upgrade.AccessLevel,
			Boost: upgrade.Boost,
			TimesBought: times_bought,
		}

		if is_bought {
			if ok && times_bought > 0 {
				filtered_upgrades = append(filtered_upgrades, this_upgrade)
			}
		} else {
			if ok && (times_bought == 0 || upgrade.UpgradeType != "dish") {
			filtered_upgrades = append(filtered_upgrades, this_upgrade)	
		}		
	}
}
	sort.Slice(filtered_upgrades, func(i, j int) bool {
		return filtered_upgrades[i].ID < filtered_upgrades[j].ID
	})

	return filtered_upgrades 
}

func InitGame(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session
	DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.ID > 0 {
		filtered_upgrades := FilterUpgrades(session, true)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "0",
			"session": session,
			"upgrades": filtered_upgrades,
		})
	}

	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		PrestigeValue: 0,
		UserID: id,
		Level: &models.Level{},
		Prestige: &models.Prestige{},
	}
	DB.Create(&new_session)

	var upgrades []models.Upgrade
	DB.Find(&upgrades)

	for _, upgrade := range upgrades {
		session_upgrade := &models.SessionUpgrade{
			SessionID: new_session.ID,
			UpgradeID: upgrade.ID,
			TimesBought: 0,
		}
		DB.Create(&session_upgrade)
	}

	DB.Preload("Prestige").Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
		"upgrades": make([]ThisUpgrade, 0),
	})
}

func CookClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session

	DB.Preload("Level").Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	filtered_upgrades := FilterUpgrades(session, true)

	var (
		total_dishes_multiplier float64 = 0
		total_dishes_per_click  float64 = 0
	)

	dish_exist := false

	for _, upgrade := range filtered_upgrades {
		if upgrade.UpgradeType == "dish" && dish_exist == false {
			dish_exist = true
		}

		if upgrade.Boost.BoostType == "dM" {
			total_dishes_multiplier += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}

		if upgrade.Boost.BoostType == "dPc" {
			total_dishes_per_click += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}
	}

	if dish_exist == false {
		return c.JSON(http.StatusForbidden, map[string]string{
			"status":  "5",
			"message": "can't perform action",
		})
	}

	if total_dishes_multiplier == 0 {
		total_dishes_multiplier = 1
	}

	DB.Model(&session).Update("dishes", gorm.Expr("dishes + ?", uint(math.Ceil((1 + total_dishes_per_click)*total_dishes_multiplier))))
	DB.First(&session, session.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
	})
}

func SellClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var (
		session models.Session
		level   models.Level
	)

	DB.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	DB.Where("session_id = ?", session.ID).First(&level)

	filtered_upgrades := FilterUpgrades(session, true)

	var (
		total_money_multiplier float64 = 0
		total_money_per_click  float64 = 0
	)

	for _, upgrade := range filtered_upgrades {
		if upgrade.Boost.BoostType == "mM" {
			total_money_multiplier += upgrade.Boost.Value * float64(upgrade.TimesBought)
		}

		if upgrade.Boost.BoostType == "mPc" {
			total_money_per_click += upgrade.Boost.Value
		}
	}

	if total_money_multiplier == 0 {
		total_money_multiplier = 1
	}

	if session.Dishes <= 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"status":  "3",
			"message": "not enough dishes",
		})
	}

	DB.Model(&session).Updates(map[string]interface{}{
		"money": gorm.Expr("money + ?", uint(math.Ceil((total_money_per_click)*total_money_multiplier))),
		"dishes": gorm.Expr("dishes - ?", 1),
	})
	DB.Model(&level).Update("xp", gorm.Expr("ROUND(xp + ?, 2)", 0.2))

	DB.First(&session, session.ID)
	DB.First(&level, level.ID)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"money":  session.Money,
		"xp":     level.XP,
	})
}

func BuyUpgrade(c echo.Context) error {
	user_id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))
	upgrade_id := utils.StringToUint(c.Param("upgrade_id"))

	var (
		session      models.Session
		this_upgrade ThisUpgrade
		exist        bool = false
		result_price uint = 0
	)

	DB.Preload("Upgrades.Boost").Where("user_id = ?", user_id).First(&session)

	for _, upgrade := range FilterUpgrades(session, false) {
		if upgrade.ID == upgrade_id {
			this_upgrade = upgrade
			exist = true
		}
	}

	if !exist {
		return c.JSON(http.StatusNotFound, map[string]string{
			"status":  "4",
			"message": "upgrade not found",
		})
	}

	if this_upgrade.TimesBought == 0 {
		result_price = this_upgrade.Price
	} else {
		result_price = uint(math.Ceil(float64(this_upgrade.Price) * this_upgrade.PriceFactor * float64(this_upgrade.TimesBought)))
	}

	if session.Money < result_price {
		return c.JSON(http.StatusConflict, map[string]string{
			"status":  "3",
			"message": "not enough money, you have: " + utils.IntToString(int(session.Money)) + ", you need: " + utils.IntToString(int(result_price)),
		})
	}

	DB.Model(&session).Select("money").Updates(models.Session{Money: session.Money - result_price})
	DB.Model(&models.SessionUpgrade{}).Where("session_id = ? AND upgrade_id = ?", session.ID, upgrade_id).Select("times_bought").Updates(models.SessionUpgrade{TimesBought: this_upgrade.TimesBought + 1})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"money":  session.Money,
	})
}

func GetUpgrades(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))

	var session models.Session
	DB.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	filtered_upgrades := FilterUpgrades(session, false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "0",
		"upgrades": filtered_upgrades,
	})
}

func GetLevel(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))
	var (
		session models.Session
		level   models.LevelXP
	)

	DB.Preload("Level").Where("user_id = ?", id).First(&session)
	DB.Where("rank = ?", session.Level.Rank + 1).First(&level)

	if level.ID == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": 0,
			"current_rank": session.Level.Rank,
			"current_xp": session.Level.XP,
			"needed_xp": session.Level.XP,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"current_rank": session.Level.Rank,
		"current_xp": session.Level.XP,
		"needed_xp": level.XP,
	})
}

func UpdateLevel(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), Secret))
	var (
		level   models.Level
		next_level models.LevelXP
	)

	DB.Where("session_id = (?)", DB.Model(&models.Session{}).Select("id").Where("user_id = ?", id),).First(&level)
	DB.Where("rank = ?", level.Rank + 1).First(&next_level)
	
	if next_level.ID == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp":   level.XP,
		})
	}

	if level.XP == float64(next_level.XP){
		DB.Model(&level).Select("xp","rank").Updates(map[string]interface{}{"xp": 0, "rank": level.Rank + 1})
		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp": level.XP,
		})
	}

	if level.XP > float64(next_level.XP) {
		DB.Model(&level).Select("xp","rank").Updates(map[string]interface{}{"xp": math.Round((level.XP - float64(next_level.XP)) * 100) / 100, "rank": level.Rank + 1})
		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp": level.XP,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"current_rank": level.Rank,
		"current_xp": level.XP,
	})
}