// Package router provides a network routing control plane
package router

import (
	"time"

	"github.com/micro/go-micro/router/table"
)

var (
	// DefaultAddress is default router address
	DefaultAddress = ":9093"
	// DefaultName is default router service name
	DefaultName = "go.micro.router"
	// DefaultNetwork is default micro network
	DefaultNetwork = "go.micro"
	// DefaultRouter is default network router
	DefaultRouter = NewRouter()
)

// Router is an interface for a routing control plane
type Router interface {
	// Init initializes the router with options
	Init(...Option) error
	// Options returns the router options
	Options() Options
	// The routing table
	Table() table.Table
	// Advertise advertises routes to the network
	Advertise() (<-chan *Advert, error)
	// Process processes incoming adverts
	Process(*Advert) error
	// Solicit advertises the whole routing table to the network
	Solicit() error
	// Start starts the router
	Start() error
	// Status returns router status
	Status() Status
	// Stop stops the router
	Stop() error
	// Returns the router implementation
	String() string
}

// Option used by the router
type Option func(*Options)

// StatusCode defines router status
type StatusCode int

const (
	// Running means the router is up and running
	Running StatusCode = iota
	// Advertising means the router is advertising
	Advertising
	// Stopped means the router has been stopped
	Stopped
	// Error means the router has encountered error
	Error
)

func (s StatusCode) String() string {
	switch s {
	case Running:
		return "running"
	case Advertising:
		return "advertising"
	case Stopped:
		return "stopped"
	case Error:
		return "error"
	default:
		return "unknown"
	}
}

// Status is router status
type Status struct {
	// Code defines router status
	Code StatusCode
	// Error contains error description
	Error error
}

// String returns human readable status
func (s Status) String() string {
	return s.Code.String()
}

// AdvertType is route advertisement type
type AdvertType int

const (
	// Announce is advertised when the router announces itself
	Announce AdvertType = iota
	// RouteUpdate advertises route updates
	RouteUpdate
)

// String returns human readable advertisement type
func (t AdvertType) String() string {
	switch t {
	case Announce:
		return "announce"
	case RouteUpdate:
		return "update"
	default:
		return "unknown"
	}
}

// Advert contains a list of events advertised by the router to the network
type Advert struct {
	// Id is the router Id
	Id string
	// Type is type of advert
	Type AdvertType
	// Timestamp marks the time when the update is sent
	Timestamp time.Time
	// TTL is Advert TTL
	TTL time.Duration
	// Events is a list of routing table events to advertise
	Events []*table.Event
}

// Strategy is route advertisement strategy
type Strategy int

const (
	// AdvertiseAll advertises all routes to the network
	AdvertiseAll Strategy = iota
	// AdvertiseBest advertises optimal routes to the network
	AdvertiseBest
)

// String returns human readable Strategy
func (s Strategy) String() string {
	switch s {
	case AdvertiseAll:
		return "all"
	case AdvertiseBest:
		return "best"
	default:
		return "unknown"
	}
}

// NewRouter creates new Router and returns it
func NewRouter(opts ...Option) Router {
	return newRouter(opts...)
}
