package services

import "github.com/mksmstpck/minora-scanner/events"

type Servieces struct {
	events events.Events
}

type Serviecer interface {
}

func NewServiecs(events events.Events) Serviecer {
	return &Servieces{
		events: events,
	}
}
