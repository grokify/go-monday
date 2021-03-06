package monday

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

const (
	ColumnValueTitleDate   = "Date"
	ColumnValueTitleStatus = "Status"
)

type Response struct {
	AccountID int  `json:"account_id"`
	Data      Data `json:"data"`
}

type Data struct {
	Boards []Board `json:"boards"`
}

type Board struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Columns []Column `json:"columns"`
	Owner   Owner    `json:"owner"`
	Items   []Item   `json:"items"`
	State   string   `json:"state"`
}

type Column struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Owner struct {
	ID int `json:"id"`
}

type Item struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	State        string        `json:"state"`
	ColumnValues []ColumnValue `json:"column_values"`
}

const (
	ColumnValueIDLink = "link"
)

type ColumnValue struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Value *string `json:"value"`
	Text  *string `json:"text"`
}

type ColumnValueValue struct {
	URL       string     `json:"url,omitempty"`  // "title":"Link", "id":"link",
	Text      string     `json:"text,omitempty"` // "title":"Link", "id":"link",
	ChangedAt *time.Time `json:"changed_at"`
}

func ParseColumnValueValue(data []byte) (ColumnValueValue, error) {
	var cvv ColumnValueValue
	err := json.Unmarshal(data, &cvv)
	return cvv, err
}

func (item *Item) GetColumnValue(name string, errorOnDupe bool) (ColumnValue, error) {
	cvs := []ColumnValue{}

	for _, cv := range item.ColumnValues {
		if cv.Title == name {
			cvs = append(cvs, cv)
		}
	}
	if len(cvs) == 0 {
		return ColumnValue{}, fmt.Errorf("column value not found [%s]", name)
	} else if len(cvs) > 1 && errorOnDupe {
		return ColumnValue{}, fmt.Errorf("more than one column values found for [%s] count [%d]",
			name, len(cvs))
	}
	return cvs[0], nil
}

func (item *Item) Date() (time.Time, error) {
	dateCv, err := item.GetColumnValue("Date", true)
	if err != nil {
		return time.Now(), err
	}
	/*
		"title":"Date",
		"id":"date4",
		"value":"{\"date\":\"2021-08-03\",\"icon\":null,\"changed_at\":\"2021-08-06T16:49:57.071Z\"}",
		"text":"2021-08-03"
	*/
	if dateCv.Text == nil {
		return time.Now(), errors.New("date text is nil")
	}
	return time.Parse(timeutil.RFC3339FullDate, *dateCv.Text)
}

func (item *Item) FieldsSimple() map[string]string {
	msi := map[string]string{}
	for _, cv := range item.ColumnValues {
		if cv.Text == nil {
			msi[cv.Title] = ""
		} else {
			msi[cv.Title] = *cv.Text
		}
	}
	return msi
}

func (item *Item) LastChangedAtDateStatus() (time.Time, error) {
	dates := []time.Time{}
	dtCv, err := item.GetColumnValue(ColumnValueTitleDate, true)
	if err == nil && dtCv.Value != nil {
		cvv, err := ParseColumnValueValue([]byte(*dtCv.Value))
		if err == nil && cvv.ChangedAt != nil {
			dates = append(dates, *cvv.ChangedAt)
		}
	}
	stCv, err := item.GetColumnValue(ColumnValueTitleStatus, true)
	if err == nil && stCv.Value != nil {
		cvv, err := ParseColumnValueValue([]byte(*stCv.Value))
		if err == nil && cvv.ChangedAt != nil {
			dates = append(dates, *cvv.ChangedAt)
		}
	}
	latest, err := timeutil.Latest(dates, true)
	if err != nil {
		return timeutil.TimeZeroRFC3339(), nil
	}
	return latest, nil
}
