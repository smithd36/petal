package models

import "time"

type Comment struct {
    ID        int
    RootID    int
    UserID    int
    Username  string
    Content   string
    CreatedAt time.Time
}