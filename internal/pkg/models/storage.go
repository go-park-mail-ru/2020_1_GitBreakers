package models

import "time"

type Repository struct {
	Id              int64
	OwnerID         int64
	LowerName       string
	Name            string
	Description     string
	Website         string
	DefaultBranch   string
	Size            int64
	UseCustomAvatar bool

	// Counters
	NumWatches          int
	NumStars            int
	NumForks            int
	NumIssues           int
	NumClosedIssues     int
	NumOpenIssues       int
	NumPulls            int
	NumClosedPulls      int
	NumOpenPulls        int
	NumMilestones       int
	NumClosedMilestones int
	NumOpenMilestones   int
	NumTags             int

	IsFork     bool
	ForkId     int64
	BaseRepoId int64

	Created     time.Time
	CreatedUnix int64
	Updated     time.Time
	UpdatedUnix int64
}
