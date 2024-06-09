package models

import "time"

type Root struct {
    ID        int
    UserID    int
    Title     string
    Content   string
    CreatedAt time.Time
}
