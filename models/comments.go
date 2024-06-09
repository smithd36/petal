package models

import "time"

type Comment struct {
    ID        int
    RootID  int
    UserID    int
    Content   string
    CreatedAt time.Time
}
