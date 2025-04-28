package handlers

import (
	"clicker_api/environment"
	"clicker_api/models"
	"clicker_api/service"
	"clicker_api/utils"
	"math"
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

type ThisUpgrade struct {
	Upgrade 	models.Upgrade `json:"upgrade"`
	TimesBought uint		   `json:"times_bought"`
}

var secret string = environment.GetVariable("ACCESS_TOKEN_SECRET")

func filterUpgrades(session models.Session, is_bought bool) []ThisUpgrade {
	filtered_upgrades := make([]ThisUpgrade, 0)

	var session_upgrades []models.SessionUpgrade
	db.Where("session_id = ?", session.ID).Find(&session_upgrades)

	times_bought_map := make(map[uint]uint)
	for _, su := range session_upgrades {
		times_bought_map[su.UpgradeID] = su.TimesBought
	}

	for _, upgrade := range session.Upgrades {
		times_bought, ok := times_bought_map[upgrade.ID]
		if is_bought {
			if ok && times_bought > 0 {
				filtered_upgrades = append(filtered_upgrades, ThisUpgrade{
					Upgrade: upgrade,
					TimesBought: times_bought,
				})
			}
		} else {
			if ok && (times_bought == 0 || upgrade.UpgradeType != "dish"){
			filtered_upgrades = append(filtered_upgrades, ThisUpgrade{
				Upgrade: upgrade,
				TimesBought: times_bought,
			})	
		}		
	}
}
	sort.Slice(filtered_upgrades, func(i, j int) bool {
		return filtered_upgrades[i].Upgrade.ID < filtered_upgrades[j].Upgrade.ID
	})

	return filtered_upgrades 
}

//TODO: ДОБАВИТЬ ПОЛЕ TIMES_BOUGHT


func InitGame(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))

	var session models.Session
	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	if session.ID > 0 {
		filtered_upgrades := filterUpgrades(session, true)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "0",
			"session": session,
			"upgrades": filtered_upgrades,
		})
	}

	new_session := models.Session{
		Money: 0,
		Dishes: 0,
		UserID: id,
	}

	db.Create(&new_session)
	var upgrades []models.Upgrade
	db.Find(&upgrades)

	for _, upgrade := range upgrades {
		session_upgrade := &models.SessionUpgrade{
			SessionID: new_session.ID,
			UpgradeID: upgrade.ID,
			TimesBought: 0,
		}
		db.Create(&session_upgrade)
	}

	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&new_session)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"session": new_session,
		"upgrades": make([]ThisUpgrade, 0),
	})
}

func CookClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))

	var session models.Session

	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	filtered_upgrades := filterUpgrades(session, true)

	var (
		total_dishes_multiplier float64 = 0;
		total_dishes_per_click float64 = 0;
	)

	dish_exist := false

	for _, upgrade := range filtered_upgrades {
		if upgrade.Upgrade.UpgradeType == "dish" && dish_exist == false  {
			dish_exist = true
		}

		if upgrade.Upgrade.Boost.BoostType == "dM" {
			total_dishes_multiplier += upgrade.Upgrade.Boost.Value
		}

		if upgrade.Upgrade.Boost.BoostType == "dPc" {
			total_dishes_per_click += upgrade.Upgrade.Boost.Value
		}
	}

	if dish_exist == false {
		return c.JSON(http.StatusForbidden, map[string]string{
			"status": "5",
			"message": "can't perform action",
		})
	}

	if total_dishes_multiplier == 0 {
		total_dishes_multiplier = 1
	}

	db.Model(&session).Select("dishes").Updates(models.Session{Dishes: session.Dishes + uint(math.Ceil((1 + total_dishes_per_click) * total_dishes_multiplier))})
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
	})
}

func SellClick(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))

	var session models.Session

	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)
	filtered_upgrades := filterUpgrades(session, true)

	var (
		total_money_multiplier float64 = 0;
		total_money_per_click float64 = 0;
	)

	for _, upgrade := range filtered_upgrades {
		if upgrade.Upgrade.Boost.BoostType == "mM" {
			total_money_multiplier += upgrade.Upgrade.Boost.Value
		}

		if upgrade.Upgrade.Boost.BoostType == "mPc" {
			total_money_per_click += upgrade.Upgrade.Boost.Value
		}
	}

	if total_money_multiplier == 0 {
		total_money_multiplier = 1
	}

	if session.Dishes <= 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"status": "3",
			"message": "not enough dishes",
		})
	}

	db.Model(&session).Select("dishes", "money").Updates(models.Session{Dishes: session.Dishes - 1, Money: session.Money + uint(math.Ceil((total_money_per_click) * total_money_multiplier))})
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"dishes": session.Dishes,
		"money": session.Money,
	})
}

func BuyUpgrade(c echo.Context) error {
	user_id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))
	upgrade_id := c.Param("upgrade_id")

	var (
		session models.Session
		this_upgrade = struct {
			SessionID 	uint
			UpgradeID 	uint
			Price       uint
			PriceFactor float64
			UpgradeType models.UpgradeType
			TimesBought uint
		}{}
	)

	db.Preload("Upgrades.Boost").Where("user_id = ?", user_id).First(&session)

	query_result := db.Table("upgrades").
	Select("upgrades.price, upgrades.price_factor, upgrades.upgrade_type, session_upgrades.times_bought").
	Joins("JOIN session_upgrades ON session_upgrades.upgrade_id = upgrades.id AND session_upgrades.upgrade_id = ?", utils.StringToUint(upgrade_id)).
	Where("session_upgrades.session_id = ?", session.ID).Scan(&this_upgrade)

	if query_result.Error != nil || query_result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"status": "4",
			"message": "upgrade not found",
		})
	}

	result_price := uint(math.Ceil(float64(this_upgrade.Price) * this_upgrade.PriceFactor * (float64(this_upgrade.TimesBought) + 1)))

	if session.Money < result_price {
		return c.JSON(http.StatusConflict, map[string]string{
			"status": "3",
			"message": "not enough money",
		})
	}

	if this_upgrade.UpgradeType == "dish" && this_upgrade.TimesBought == 1 {
		return c.JSON(http.StatusConflict, map[string]string{
			"status": "3",
			"message": "already bought",
		})
	}

	db.Model(&session).Select("money").Updates(models.Session{Money: session.Money - result_price})
	db.Model(&models.SessionUpgrade{}).Where("session_id = ? AND upgrade_id = ?", session.ID, upgrade_id).Select("times_bought").Updates(models.SessionUpgrade{TimesBought: this_upgrade.TimesBought + 1})
	
	return c.JSON(http.StatusOK, map[string]string{
		"status": "0",
		"message": "success",
	})
}

func GetUpgrades(c echo.Context) error {
	id := utils.StringToUint(service.ExtractIDFromToken(c.Request().Header.Get("Authorization"), secret))

	var session models.Session
	db.Preload("Upgrades.Boost").Where("user_id = ?", id).First(&session)

	filtered_upgrades := filterUpgrades(session, false)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "0",
		"upgrades": filtered_upgrades,
	})
}