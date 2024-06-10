package models

import "time"

type Root struct {
    ID        int
    UserID    int
    Username  string
    Title     string
    Content   string
    CreatedAt time.Time
}