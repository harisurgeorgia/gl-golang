package controllers

import (
	"fmt"
	"gl/models"
	"gl/redirect"
	"gl/session"
	"gl/utils"
	"gl/views"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func JournalEntry(c *gin.Context) {
	s := "/static/js/journal-entry.js"
	attr := views.LayoutAttribute{
		PageTitle:  "GL Entry",
		PageHeader: "Journal Entry",
		Script:     &s,
	}

	journal := models.Journal{JournalDate: time.Now()}

	accounts := models.GetAllAccounts()
	link := "/journal/search"
	menus, err := models.GetUserMenu("journal")
	if err != nil {
		utils.Render(c, http.StatusInternalServerError, views.Layout(nil, attr, views.ErrorPage(redirect.Error500)))
		return
	}
	navData := views.NavData{
		Menus:     menus,
		Link:      &link,
		IsTopNav:  true,
		Search:    true,
		CsrfToken: csrf.GetToken(c),
	}
	journalPageData := models.JournalPageData{
		Detail:     journal,
		Heads:      accounts,
		PostStatus: false,
	}
	token := navData.CsrfToken

	utils.Render(c, http.StatusOK, views.Layout(views.Nav(navData), attr, views.JournalEntryForm(journalPageData, token)))
}

func JournalSave(c *gin.Context) {

	menus, err := models.GetUserMenu("journal")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve user menu: %v", err)
		return
	}

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
	var journalNumberPtr *string
	if journalNumber != "" {
		journalNumberPtr = &journalNumber
	}
	description := c.PostForm("description")
	journal := models.Journal{
		JournalDate:   journalDate,
		JournalNumber: journalNumberPtr,
		Description:   description,
		Lines:         lines,
	}
	if idStr := c.PostForm("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			journal.ID = &id
		}
	}
	link := "/journal/search"
	navData := views.NavData{
		Menus:     menus,
		Link:      &link,
		IsTopNav:  true,
		Search:    true,
		CsrfToken: csrf.GetToken(c),
	}
	s := "/static/js/journal-entry.js"
	attr := views.LayoutAttribute{
		PageTitle:  "GL Entry",
		PageHeader: "Journal Entry",
		Script:     &s,
	}
	//data = views.PageData{Title: "GL", Header: "Journal Entry"}
	accounts := models.GetAllAccounts()
	journalPageData := models.JournalPageData{
		Detail:     journal,
		Heads:      accounts,
		PostStatus: false,
	}
	csrfToken := csrf.GetToken(c)
	if strings.TrimSpace(dabitBalance) != strings.TrimSpace(creditBalance) {
		utils.Render(c, http.StatusOK, views.Layout(views.Nav(navData), attr, views.JournalEntryForm(journalPageData, csrfToken)))
		return
	}

	var id *int64
	id, err = models.JournalSave(journal)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save journal entry: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/journal/edit/"+strconv.FormatInt(*id, 10))
}

func JournalList(c *gin.Context) {
	menus, err := models.GetUserMenu("journal")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve user menu: %v", err)
		return
	}
	link := "/journal/search"
	navData := views.NavData{
		Menus:     menus,
		Link:      &link,
		IsTopNav:  true,
		Search:    true,
		CsrfToken: csrf.GetToken(c),
	}
	attr := views.LayoutAttribute{
		PageTitle:  "GL Entry",
		PageHeader: "Journal Entry",
		Script:     nil,
	}

	//data = views.PageData{Title: "GL", Header: "Journal Entry"}
	filter := c.Param("filter")
	if filter == "all" {
		filter = ""
	}

	journals, err := models.GetJournals(filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journals: %v", err)
		return
	}

	utils.Render(c, http.StatusOK, views.Layout(views.Nav(navData), attr, views.JournalList(journals)))
}

// func view a journal entry
func JournalEdit(c *gin.Context) {

	s := "/static/js/journal-entry.js"

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		attr := views.LayoutAttribute{
			PageTitle:  "GL Entry",
			PageHeader: "Journal Entry",
			Script:     &s,
		}
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, attr, views.ErrorPage(redirect.Error400)))
		return
	}

	journal, err := models.GetJournalByID(id)
	if err != nil || journal == nil {
		attr := views.LayoutAttribute{
			PageTitle:  "GL Entry",
			PageHeader: "Journal Entry",
			Script:     &s,
		}
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, attr, views.ErrorPage(redirect.Error400)))
		return
	}
	uid := session.GetSession(c, "user_id")
	user_id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		attr := views.LayoutAttribute{
			PageTitle:  "GL Entry",
			PageHeader: "Journal Entry",
			Script:     &s,
		}
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, attr, views.ErrorPage(redirect.Error400)))
		return
	}
	journal.CreatedBy = user_id
	journal.PostedBy = &user_id
	urole := session.GetSession(c, "user_role")

	role, err := strconv.Atoi(urole)
	if err != nil {
		attr := views.LayoutAttribute{
			PageTitle:  "GL Entry",
			PageHeader: "Journal Entry",
			Script:     &s,
		}
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, attr, views.ErrorPage(redirect.Error400)))
		return
	}
	PostStatus := false
	if role == 4 && journal.Status == "closed" {
		PostStatus = true
	}
	menus, err := models.GetUserMenu("journal")
	if err != nil {
		attr := views.LayoutAttribute{
			PageTitle:  "GL Entry",
			PageHeader: "Journal Entry",
			Script:     &s,
		}
		utils.Render(c, http.StatusBadRequest, views.Layout(nil, attr, views.ErrorPage(redirect.Error400)))
		return
	}
	link := "/journal/search"
	navData := views.NavData{
		Menus:     menus,
		Link:      &link,
		IsTopNav:  true,
		Search:    true,
		CsrfToken: csrf.GetToken(c),
	}
	attr := views.LayoutAttribute{
		PageTitle:  "GL Entry",
		PageHeader: "Journal Entry",
		Script:     &s,
	}
	accounts := models.GetAllAccounts()
	journalPageData := models.JournalPageData{
		Detail:     *journal,
		Heads:      accounts,
		PostStatus: PostStatus,
	}
	csrfToken := csrf.GetToken(c)
	utils.Render(c, http.StatusOK, views.Layout(views.Nav(navData), attr, views.JournalEntryForm(journalPageData, csrfToken)))
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
	journal, err := models.GetJournalByID(id)
	menus, err := models.GetUserMenu("journal")
	s := "/static/js/journal-entry.js"
	link := "/journal/search"
	navData := views.NavData{
		Menus:     menus,
		Link:      &link,
		IsTopNav:  true,
		Search:    true,
		CsrfToken: csrf.GetToken(c),
	}
	utils.Render(c, http.StatusBadRequest,
		views.Layout(views.Nav(navData),
			views.LayoutAttribute{PageTitle: "GL Entry", PageHeader: "Journal Entry", Script: &s},
			views.ErrorPage(redirect.Error404)))

	uid := session.GetSession(c, "user_id")
	user_id, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		utils.Render(c, http.StatusBadRequest,
			views.Layout(views.Nav(navData),
				views.LayoutAttribute{PageTitle: "GL Entry", PageHeader: "Journal Entry", Script: &s},
				views.ErrorPage(redirect.Error404)))
		return
	}
	journal.PostedBy = &user_id
	journal.Status = string(models.StatusClosed)
	now := time.Now()
	journal.PostedAt = &now

	_, err = models.JournalUpdate(*journal)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update journal: %v", err)
		return
	}
	c.Redirect(http.StatusFound, "/journal/edit/"+strconv.FormatInt(id, 10))
}

func FindJournalByJournalNumber(c *gin.Context) {
	filter := strings.TrimSpace(c.Param("filter"))
	rows, err := models.FindJournalByJournalNumber(filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve journals: %v", err)
		return
	}
	menus, _ := models.GetUserMenu("journal")
	link := "/journal/search"
	navData := views.NavData{
		Menus:    menus,
		Link:     &link,
		IsTopNav: true,
		Search:   true,
	}

	utils.Render(c, http.StatusOK, views.Layout(views.Nav(navData), views.LayoutAttribute{PageTitle: "GL-list"}, views.JournalList(rows)))

}
