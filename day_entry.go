package harvest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Used for /daily API queries
type DayEntryResponse struct {
	DayEntries []*DayEntry `json:"day_entries"`
}

// Used for both /project/{id}/entries AND /user/{id}/entries
type DayEntryReport []DayEntryReportRow
type DayEntryReportRow struct {
	DayEntry *DayEntry `json:"day_entry"`
}

type DayEntry struct {
	ID                int64           `json:"id"`
	UserID            int64           `json:"user_id"`
	SpentAt           string          `json:"spent_at"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	ProjectRaw        json.RawMessage `json:"project_id"`
	ProjectID         int64
	TaskRaw           json.RawMessage `json:"task_id"`
	TaskID            int64
	Project           string  `json:"project"`
	Task              string  `json:"task"`
	Client            string  `json:"client"`
	Notes             string  `json:"notes"`
	HoursWithoutTimer int64   `json:"hours_without_timer"`
	Hours             float64 `json:"hours"`
}

// Needed to avoid recursion in UnmarshalJSON
type dayentry DayEntry

func (dayEntry *DayEntry) UnmarshalJSON(b []byte) (err error) {
	d, s, i := dayentry{}, "", float64(0.0)

	if err = json.Unmarshal(b, &d); err == nil {
		if err = json.Unmarshal(d.ProjectRaw, &s); err == nil {
			i, err = strconv.ParseFloat(s, 64)
			d.ProjectID = int64(i)
		}
		if err = json.Unmarshal(d.ProjectRaw, &i); err == nil {
			d.ProjectID = int64(i)
		}
		if err = json.Unmarshal(d.TaskRaw, &s); err == nil {
			i, err = strconv.ParseFloat(s, 64)
			d.TaskID = int64(i)
		}
		if err = json.Unmarshal(d.TaskRaw, &i); err == nil {
			d.TaskID = int64(i)
		}
		*dayEntry = DayEntry(d)
		err = nil
	}
	return
}

func (report DayEntryReport) Entries() []*DayEntry {
	results := make([]*DayEntry, len(report))
	for i, _ := range report {
		results[i] = report[i].DayEntry
	}
	return results
}

func (a *API) GetTodayEntries(args Arguments) ([]*DayEntry, error) {
	dayEntriesResponse := DayEntryResponse{}
	path := fmt.Sprintf("/daily?slim=1")
	err := a.Get(path, args, &dayEntriesResponse)
	return dayEntriesResponse.DayEntries, err
}

func (a *API) GetEntriesForProjectBetween(projectID int64, fromDate time.Time, toDate time.Time, args Arguments) ([]*DayEntry, error) {
	response := make(DayEntryReport, 0)
	from := fromDate.Format("20060102")
	to := fromDate.Format("20060102")
	path := fmt.Sprintf("/projects/%d/entries?from=%s&to=%s", projectID, from, to)
	err := a.Get(path, args, &response)
	return response.Entries(), err
}

func (a *API) GetEntriesForUserBetween(userID int64, fromDate time.Time, toDate time.Time, args Arguments) ([]*DayEntry, error) {
	response := make(DayEntryReport, 0)
	from := fromDate.Format("20060102")
	to := fromDate.Format("20060102")
	path := fmt.Sprintf("/people/%d/entries?from=%s&to=%s", userID, from, to)
	err := a.Get(path, args, &response)
	return response.Entries(), err
}
