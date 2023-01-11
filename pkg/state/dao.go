package state

import (
	"context"
	"time"

	"github.com/awlsring/surreal-db-client/surreal"
)

const STACK = "stack:"

func STACK_ID(stack string) string {
	return STACK + stack
}

type StateDao interface {
	Create(state Entity) error
	Read(id string) (Entity, error)
	Update(state Entity) error
	Delete(id string) error
}

type SurrealStateDao struct {
	db *surreal.Surreal
}

func New(db *surreal.Surreal) StateDao {
	return &SurrealStateDao{
		db: db,
	}
}

func (d *SurrealStateDao) Close() {
	d.db.Close()
}

func (s *SurrealStateDao) Create(state Entity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return s.db.CreateItem(ctx, state.Id, surreal.ResourceToEntry(state))
}

func (s *SurrealStateDao) Read(id string) (Entity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	blob, err := s.db.GetItem(ctx, STACK_ID(id))
	if err != nil {
		return Entity{}, err 
	}
	var state *Entity
	err = surreal.UnmarshalGet(blob, &state)
	if err != nil {
		return Entity{}, err 
	}
	return *state, nil
}

func (s *SurrealStateDao) Update(state Entity) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return s.db.UpdateItem(ctx, state.Id, surreal.ResourceToEntry(state))
}

func (s *SurrealStateDao) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return s.db.DeleteItem(ctx, STACK_ID(id))
}