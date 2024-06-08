package models

import "time"

type Thread struct {
    ID        int
    UserID    int
    Title     string
    Content   string
    CreatedAt time.Time
}
