package common

import "github.com/google/uuid"

type Requester interface {
	UserId() uuid.UUID
	TokenId() uuid.UUID
	FirstName() string
	LastName() string
	Role() string
	Status() string
}

type requesterData struct {
	userId    uuid.UUID
	tid       uuid.UUID
	firstName string
	lastName  string
	role      string
	status    string
}

func (r *requesterData) UserId() uuid.UUID {
	return r.userId
}
func (r *requesterData) TokenId() uuid.UUID {
	return r.tid
}
func (r *requesterData) FirstName() string { return r.firstName }
func (r *requesterData) LastName() string  { return r.lastName }
func (r *requesterData) Role() string      { return r.role }
func (r *requesterData) Status() string    { return r.status }

func NewRequester(sub, tid uuid.UUID, firstName, lastName, role, status string) Requester {
	return &requesterData{
		userId:    sub,
		tid:       tid,
		firstName: firstName,
		lastName:  lastName,
		role:      role,
		status:    status,
	}
}
