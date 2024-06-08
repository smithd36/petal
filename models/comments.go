package models

import "time"

type Comment struct {
    ID        int
    ThreadID  int
    UserID    int
    Content   string
    CreatedAt time.Time
}
