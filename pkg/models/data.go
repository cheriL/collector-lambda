package models

import "time"

type Data struct {
	ID         int64     `json:"id"`
	Number     int       `json:"number"`
	Type       DataType  `json:"type"`
	UserID     int64     `json:"user_id"`
	UserType   UserType  `json:"user_type"`
	UserLogin  string    `json:"user_login"`
	//State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	//UpdatedAt  time.Time `json:"updated_at"`
	//ClosedAt   time.Time `json:"closed_at"`
}

type UserType int

const (
	UserTypeContributor = iota
	UserTypePingCaper
)

type DataType int

const (
	DataTypeIssue = iota
	DataTypePr
)

type DataLabel struct {
	DataId    int
	LabelName string
}