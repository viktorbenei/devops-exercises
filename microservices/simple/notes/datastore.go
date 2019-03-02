package main

import (
	"sync"
)

type UserID string
type NoteID string
type Note struct {
	Content string `json:"content"`
}
type Datastore struct {
	notes map[UserID]map[NoteID]Note
	mux   sync.Mutex
}

func NewDatastore() *Datastore {
	return &Datastore{
		notes: map[UserID]map[NoteID]Note{},
	}
}

// SetNote ...
func (ds *Datastore) SetNote(userID UserID, noteID NoteID, note Note) error {
	if len(userID) < 1 {
		return NewInputError("Invalid (empty) UserID")
	}
	if len(noteID) < 1 {
		return NewInputError("Invalid (empty) NoteID")
	}

	ds.mux.Lock()
	defer ds.mux.Unlock()

	if len(ds.notes[userID]) < 1 {
		ds.notes[userID] = map[NoteID]Note{}
	}

	ds.notes[userID][noteID] = note

	return nil
}

// GetNotes ...
func (ds *Datastore) GetNotes(userID UserID) (map[NoteID]Note, error) {
	if len(userID) < 1 {
		return map[NoteID]Note{}, NewInputError("Invalid (empty) UserID")
	}

	ds.mux.Lock()
	defer ds.mux.Unlock()

	if ds.notes == nil || ds.notes[userID] == nil {
		return map[NoteID]Note{}, NewNotFoundError("Not found for user")
	}

	return ds.notes[userID], nil
}
