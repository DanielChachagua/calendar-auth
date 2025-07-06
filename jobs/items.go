package jobs

import "google.golang.org/api/calendar/v3"

type ResultItems struct {
    Items []*calendar.Event `json:"items"`
}
