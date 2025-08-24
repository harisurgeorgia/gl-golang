package controllers

import (
	"fmt"
	"gl/db"
	"gl/models"
	"gl/utils"
	"gl/views"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

func JournalEntry(c *gin.Context) {

	var data = views.PageData{Title: "GL Entry", Header: "Journal Entry"}

	var journal = models.Journal{JournalDate: time.Now()}

	accounts := models.GetAllAccounts(db.Conn)
	utils.Render(c, http.StatusOK, views.Layout(data, views.JournalEntryForm(data.Header, "", journal, accounts)))
}

func JournalSave(c *gin.Context) {

	c.Request.ParseForm()
	dabitBalance := c.PostForm("debit-bal")
	creditBalance := c.PostForm("credit-bal")

	accountsIDs := c.PostFormArray("accounts_id[]")
	descriptions := c.PostFormArray("line_description[]")
	debits := c.PostFormArray("debit[]")
	credits := c.PostFormArray("credit[]")

	var lines []models.JournalLine

	for i := range accountsIDs {
		fmt.Println(i)
		accountID, _ := strconv.ParseInt(accountsIDs[i], 10, 64)
		debit, _ := decimal.NewFromString(debits[i])
		credit, _ := decimal.NewFromString(credits[i])

		line := models.JournalLine{
			AccountID:   accountID,
			Description: descriptions[i],
			Debit:       debit,
			Credit:      credit,
			LineNumber:  i + 1,
		}
		lines = append(lines, line)

	}
	dateStr := c.PostForm("journal-date")
	journalDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid date format")
		return
	}

	journalNumber := strings.TrimSpace(c.PostForm("journal-number"))
	description := c.PostForm("description")
	journal := models.Journal{
		JournalDate:   journalDate,
		JournalNumber: &journalNumber,
		Description:   description,
		Lines:         lines,
	}
	var data = views.PageData{Title: "GL", Header: "Journal Entry"}
	accounts := models.GetAllAccounts(db.Conn)
	if strings.TrimSpace(dabitBalance) != strings.TrimSpace(creditBalance) {
		utils.Render(c, http.StatusOK, views.Layout(data, views.JournalEntryForm(data.Header, "", journal, accounts)))
		return
	}

	err = models.JournalSave(journal, db.Conn)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save journal entry: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/journal")
}

func JournalList(c *gin.Context) {

	var data = views.PageData{Title: "Journal List", Header: "Journal Entries"}

	journals, err := models.GetPendingJournals(db.Conn)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journals: %v", err)
		return
	}

	utils.Render(c, http.StatusOK, views.Layout(data, views.JournalList(journals)))
}
