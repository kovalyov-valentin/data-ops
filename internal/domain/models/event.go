package models

import (
	"fmt"
	"time"
)

type ClickhouseEvent struct {
	Id          int       `json:"Id"`
	ProjectId   int       `json:"ProjectId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description,omitempty"`
	Priority    int       `json:"Priority,omitempty"`
	Removed     bool      `json:"Removed,omitempty"`
	EventTime   time.Time `json:"EventTime"`
}

func (cl *ClickhouseEvent) String() string {
	return fmt.Sprintf(
		"ID:%d, ProjectId:%d, Name:%s, Description:%s, Priority:%d, Removed:%v, EventTime:%s",
		cl.Id, cl.ProjectId, cl.Name, cl.Description, cl.Priority, cl.Removed, cl.EventTime,
	)
}
