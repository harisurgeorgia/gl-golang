package controllers

import (
	"gl/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func ClosePeriod(c *gin.Context) {
	//c.Request.PostForm("close-period")

	month_year := models.GetCurrentPeriod()

	endDate := lastDayOfMonthDate(month_year.Month, month_year.Year)
	startDate := time.Date(month_year.Year, time.Month(month_year.Month), 1, 0, 0, 0, 0, time.UTC)
	log.Printf("Closing period from %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	//var endDate = time.Date(2025, 3, 30, 0, 0, 0, 0, time.UTC)
	err := models.ClosePeriod(startDate, endDate)
	if err != nil {

		c.String(400, "Error closing period: %v", err)
		return
	}
	c.String(200, "Close Period - To be implemented")
}

func lastDayOfMonthDate(month int, year int) time.Time {
	// Convert month to time.Month and add +1 to get the next month
	// If month is December, Go will correctly handle month rollover
	firstOfNextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.UTC)

	// Subtract 1 day to get the last day of the given month
	lastDay := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDay
}
