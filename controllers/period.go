package controllers

import (
	"gl/models"
	"gl/utils"
	"gl/views"
	"log"
	"net/http"
	"strconv"
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

func PeriodAddEdit(c *gin.Context) {
	idStr := c.Param("id")
	if idStr != "" {
		parsedID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "invalid id",
			})
			return
		}
		id := parsedID
		period, err := models.GetPeriod(id)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"error": "invalid id",
			})
			return
		}
		utils.Render(c, http.StatusOK, views.Period(*period))
		return
	}

	endDate, err := GetLastPeriodEndDate()
	if err != nil {
		log.Fatal("Error finding end_date")
	}
	if endDate == nil {
		t := time.Now().UTC()

		// First day of current month
		firstOfMonth := time.Date(
			t.Year(), t.Month(), 1,
			0, 0, 0, 0,
			time.UTC,
		)

		// Last day of previous month
		lastPrevMonth := firstOfMonth.Add(-time.Nanosecond)

		endDate = &lastPrevMonth
	}
	//startData = endDate plus one day
	startDate := endDate.AddDate(0, 0, 1)
	periodEnd := utils.LastDayOfMonth(startDate)
	month := startDate.Month().String()

	//log.Printf("Last date of new period %s", periodEnd)
	periodEndYead := strconv.Itoa(startDate.Year())
	var newPeriod models.Period
	newPeriod.PeriodName = month + " - " + periodEndYead
	newPeriod.StartDate = startDate
	newPeriod.EndDate = periodEnd
	newPeriod.Status = models.StatusPending

	data, err := getBasePageData(c, "Period", "New Period", "", nil)
	if err != nil {

	}
	utils.Render(c, http.StatusOK, views.Layout(views.Nav(data.Menus, data.Search), data, views.Period(newPeriod)))
}

func GetLastPeriodEndDate() (*time.Time, error) {
	date, err := models.GetLastPeriodEndDate()
	if err != nil {
		log.Fatal("Error capturing date")
	}
	return date, nil
}

func PeriodSave(c *gin.Context) {

	var p models.Period
	p.PeriodName = c.PostForm("period_name")
	startDate, err := time.Parse("2006-01-02", c.PostForm("start_date"))
	if err != nil {
		log.Fatal("Invaild Date")
	}
	p.StartDate = startDate
	endDate, err := time.Parse("2006-01-02", c.PostForm("end_date"))
	if err != nil {
		log.Fatal("Invaild Date")
	}
	p.EndDate = endDate

	p.Status = models.StatusPending

	models.SavePeriod(p)

}
