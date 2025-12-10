package handlers

import (
	"clicker_api/database"
	"clicker_api/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func c(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var (
		session models.Session
		level   models.LevelXP
	)

	database.DB.Preload("Level").Where("user_id = ?", id).First(&session)
	if session.Level.Rank == 100 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":       0,
			"current_rank": session.Level.Rank,
			"current_xp":   session.Level.XP,
		})
	}

	database.DB.Where("rank = ?", session.Level.Rank+1).First(&level)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       0,
		"current_rank": session.Level.Rank,
		"current_xp":   session.Level.XP,
		"needed_xp":    level.XP,
	})
}

func UpdateLevel(c echo.Context) error {
	id, _ := c.Get("id").(uint)

	var (
		level      models.Level
		next_level models.LevelXP
	)

	database.DB.Where("session_id = (?)", database.DB.Model(&models.Session{}).Select("id").Where("user_id = ?", id)).First(&level)

	if level.Rank == 100 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp":   level.XP,
		})
	}

	database.DB.Where("rank = ?", level.Rank+1).First(&next_level)

	if level.XP == float64(next_level.XP) {
		database.DB.Model(&level).Updates(map[string]interface{}{
			"xp":   0,
			"rank": gorm.Expr("rank + ?", 1),
		})
		database.DB.First(&level, level.ID)

		var new_next_level models.LevelXP
		database.DB.Where("rank = ?", next_level.Rank+1).First(&new_next_level)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp":   level.XP,
			"next_xp":      new_next_level.XP,
		})
	}

	if level.XP > float64(next_level.XP) {
		database.DB.Model(&level).Updates(map[string]interface{}{
			"xp":   gorm.Expr("ROUND(xp - ?, 2)", next_level.XP),
			"rank": gorm.Expr("rank + ?", 1),
		})
		database.DB.First(&level, level.ID)

		var new_next_level models.LevelXP
		database.DB.Where("rank = ?", next_level.Rank+1).First(&new_next_level)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_rank": level.Rank,
			"current_xp":   level.XP,
			"next_xp":      new_next_level.XP,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"current_rank": level.Rank,
		"current_xp":   level.XP,
	})
}