package memdb

import (
	mdb "github.com/hashicorp/go-memdb"
	"github.com/micro/go-micro/router/table"
)

type memdbTable struct {
	schema *mdb.DBSchema
	db     *mdb.MemDB
}

func NewTable() *memdbTable {
	return &memdbTable{}
}

func (t *memdbTable) Create(r table.Route) error {
	return nil
}
func (t *memdbTable) Delete(r table.Route) error {
	return nil
}
func (t *memdbTable) Update(r table.Route) error {
	return nil
}

func (t *memdbTable) List() ([]table.Route, error) {
	return nil, nil
}

func (t *memdbTable) Query(q ...table.QueryOption) ([]table.Route, error) {
	return nil, nil
}

func (t *memdbTable) Watch(opts ...table.WatchOption) (table.Watcher, error) {
	return nil, nil
}
