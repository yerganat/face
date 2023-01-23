// package taskstore provides a simple in-memory "data store" for faces.
// Tasks are uniquely identified by numeric IDs.
//
// Eli Bendersky [https://eli.thegreenplace.net]
// This code is in the public domain.
package facestore

import (
	"fmt"
	"sync"
	"time"
)

type Face struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// FaceStore is a simple in-memory database of faces; FaceStore methods are
// safe to call concurrently.
type FaceStore struct {
	sync.Mutex

	faces  map[int]Face
	nextId int
}

func New() *FaceStore {
	ts := &FaceStore{}
	ts.faces = make(map[int]Face)
	ts.nextId = 0
	return ts
}

// CreateFace creates a new task in the store.
func (ts *FaceStore) CreateFace(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	face := Face{
		Id:   ts.nextId,
		Text: text,
		Due:  due}
	face.Tags = make([]string, len(tags))
	copy(face.Tags, tags)

	ts.faces[ts.nextId] = face
	ts.nextId++
	return face.Id
}

// GetFace retrieves a task from the store, by id. If no such id exists, an
// error is returned.
func (ts *FaceStore) GetFace(id int) (Face, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.faces[id]
	if ok {
		return t, nil
	} else {
		return Face{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// DeleteFace deletes the task with the given id. If no such id exists, an error
// is returned.
func (ts *FaceStore) DeleteFace(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.faces[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.faces, id)
	return nil
}

// DeleteAllFaces deletes all faces in the store.
func (ts *FaceStore) DeleteAllFaces() error {
	ts.Lock()
	defer ts.Unlock()

	ts.faces = make(map[int]Face)
	return nil
}

// GetAllFaces returns all the faces in the store, in arbitrary order.
func (ts *FaceStore) GetAllFaces() []Face {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]Face, 0, len(ts.faces))
	for _, task := range ts.faces {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

// GetFacesByTag returns all the faces that have the given tag, in arbitrary
// order.
func (ts *FaceStore) GetFacesByTag(tag string) []Face {
	ts.Lock()
	defer ts.Unlock()

	var faces []Face

taskloop:
	for _, face := range ts.faces {
		for _, faceTag := range face.Tags {
			if faceTag == tag {
				faces = append(faces, face)
				continue taskloop
			}
		}
	}
	return faces
}

// GetFacesByDueDate returns all the faces that have the given due date, in
// arbitrary order.
func (ts *FaceStore) GetFacesByDueDate(year int, month time.Month, day int) []Face {
	ts.Lock()
	defer ts.Unlock()

	var faces []Face

	for _, face := range ts.faces {
		y, m, d := face.Due.Date()
		if y == year && m == month && d == day {
			faces = append(faces, face)
		}
	}

	return faces
}
