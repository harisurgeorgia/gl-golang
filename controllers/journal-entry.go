package controllers

import (
	"fmt"
	"gl/db"
	"gl/messages"
	"gl/models"
	"gl/session"
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
	data, err := getBasePageData(c, "GL Entry", "Journal Entry", "journal")
	var journal = models.Journal{JournalDate: time.Now()}
	if err != nil {

	}

	//data = views.PageData{Title: "GL Entry", Header: "Journal Entry"}

	accounts := models.GetAllAccounts(db.Conn)
	utils.Render(c, http.StatusOK, views.Layout(views.Nav(data.Menus, data.Search), data, views.JournalEntryForm(data.Header, "", journal, accounts, false)))
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
	//err, idStr := strconv.ParseInt(c.PostForm("id"), 10, 64)

	journalNumber := strings.TrimSpace(c.PostForm("journal-number"))
	description := c.PostForm("description")
	journal := models.Journal{
		JournalDate:   journalDate,
		JournalNumber: &journalNumber,
		Description:   description,
		Lines:         lines,
	}
	if idStr := c.PostForm("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			journal.ID = &id
		}
	}

	data, err := getBasePageData(c, "GL Entry", "Journal Entry", "journal")
	if err != nil {

	}
	//data = views.PageData{Title: "GL", Header: "Journal Entry"}
	accounts := models.GetAllAccounts(db.Conn)
	if strings.TrimSpace(dabitBalance) != strings.TrimSpace(creditBalance) {
		utils.Render(c, http.StatusOK, views.Layout(views.Nav(data.Menus, true), data, views.JournalEntryForm(data.Header, "", journal, accounts, false)))
		return
	}

	var id *int64
	id, err = models.JournalSave(journal, db.Conn)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save journal entry: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/journal/edit/"+strconv.FormatInt(*id, 10))
}

func JournalList(c *gin.Context) {

	data := views.PageData{Title: "Journal List", Header: "Journal Entries"}
	filter := c.Param("filter")

	journals, err := models.GetJournals(db.Conn, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journals: %v", err)
		return
	}

	utils.Render(c, http.StatusOK, views.Layout(views.Nav(data.Menus, true), data, views.JournalList(journals)))
}

// func view a journal entry
func JournalEdit(c *gin.Context) {
	data, err := getBasePageData(c, "GL Entry", "Journal Entry", "journal")
	if err != nil {

	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid journal ID")
		return
	}

	journal, err := models.GetJournalByID(id, db.Conn)
	if err != nil || journal == nil {
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, views.PageData{
			Title: "Page Not Found",
		}, views.ErrorPage(messages.Error404)))

		return
	}
	uid := session.GetSession(c, "user_id")
	user_id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journal: %v", err)
	}
	journal.CreatedBy = user_id
	journal.PostedBy = &user_id
	urole := session.GetSession(c, "user_role")

	role, err := strconv.Atoi(urole)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journal: %v", err)
		return
	}
	PostStatus := false
	if role == 4 {
		PostStatus = true
	}
	accounts := models.GetAllAccounts(db.Conn)

	utils.Render(c, http.StatusOK, views.Layout(views.Nav(data.Menus, true), data, views.JournalEntryForm(data.Header, "", *journal, accounts, PostStatus)))
}

func JournalVerify(c *gin.Context) {
}

func JournalPost(c *gin.Context) {
	c.Request.ParseForm()
	idStr := c.PostForm("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid journal ID")
		return
	}
	journal, err := models.GetJournalByID(id, db.Conn)
	if err != nil || journal == nil {
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, views.PageData{
			Title: "Page Not Found",
		}, views.ErrorPage(messages.Error404)))

		return
	}
	uid := session.GetSession(c, "user_id")
	user_id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journal: %v", err)
	}
	journal.PostedBy = &user_id
	journal.Status = "Posted"
	now := time.Now()
	journal.PostedAt = &now
	err, _ = models.JournalUpdate(*journal, db.Conn)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update journal: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/journal/edit/"+strconv.FormatInt(id, 10))
}
