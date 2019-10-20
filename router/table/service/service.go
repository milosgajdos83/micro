package service

import (
	"context"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/router/table"
	pb "github.com/micro/go-micro/router/table/proto"
)

type svc struct {
	table    pb.TableService
	callOpts []client.CallOption
}

func NewTable(name string, c client.Client, callOpts []client.CallOption) table.Table {
	table := pb.NewTableService(name, c)

	return &svc{
		table:    table,
		callOpts: callOpts,
	}
}

// Create new route in the routing table
func (t *svc) Create(r table.Route) error {
	route := &pb.Route{
		Service: r.Service,
		Address: r.Address,
		Gateway: r.Gateway,
		Network: r.Network,
		Link:    r.Link,
		Metric:  int64(r.Metric),
	}

	if _, err := t.table.Create(context.Background(), route, t.callOpts...); err != nil {
		return err
	}

	return nil
}

// Delete deletes existing route from the routing table
func (t *svc) Delete(r table.Route) error {
	route := &pb.Route{
		Service: r.Service,
		Address: r.Address,
		Gateway: r.Gateway,
		Network: r.Network,
		Link:    r.Link,
		Metric:  int64(r.Metric),
	}

	if _, err := t.table.Delete(context.Background(), route, t.callOpts...); err != nil {
		return err
	}

	return nil
}

// Update updates route in the routing table
func (t *svc) Update(r table.Route) error {
	route := &pb.Route{
		Service: r.Service,
		Address: r.Address,
		Gateway: r.Gateway,
		Network: r.Network,
		Link:    r.Link,
		Metric:  int64(r.Metric),
	}

	if _, err := t.table.Update(context.Background(), route, t.callOpts...); err != nil {
		return err
	}

	return nil
}

// List returns the list of all routes in the table
func (t *svc) List() ([]table.Route, error) {
	resp, err := t.table.List(context.Background(), &pb.Request{}, t.callOpts...)
	if err != nil {
		return nil, err
	}

	routes := make([]table.Route, len(resp.Routes))
	for i, route := range resp.Routes {
		routes[i] = table.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Link:    route.Link,
			Metric:  int(route.Metric),
		}
	}

	return routes, nil
}

// Lookup looks up routes in the routing table and returns them
func (t *svc) Lookup(q ...table.QueryOption) ([]table.Route, error) {
	query := table.NewQuery(q...)

	// call the router
	resp, err := t.table.Lookup(context.Background(), &pb.LookupRequest{
		Query: &pb.Query{
			Service: query.Service,
			Gateway: query.Gateway,
			Network: query.Network,
		},
	}, t.callOpts...)

	// errored out
	if err != nil {
		return nil, err
	}

	routes := make([]table.Route, len(resp.Routes))
	for i, route := range resp.Routes {
		routes[i] = table.Route{
			Service: route.Service,
			Address: route.Address,
			Gateway: route.Gateway,
			Network: route.Network,
			Link:    route.Link,
			Metric:  int(route.Metric),
		}
	}

	return routes, nil
}

// Watch returns table watcher
func (t *svc) Watch(opts ...table.WatchOption) (table.Watcher, error) {
	rsp, err := t.table.Watch(context.Background(), &pb.WatchRequest{}, t.callOpts...)
	if err != nil {
		return nil, err
	}
	options := table.WatchOptions{
		Service: "*",
	}
	for _, o := range opts {
		o(&options)
	}
	return newWatcher(rsp, options)
}

func (t *svc) String() string {
	return "service"
}
