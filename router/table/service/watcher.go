package service

import (
	"io"
	"sync"
	"time"

	"github.com/micro/go-micro/router/table"
	pb "github.com/micro/go-micro/router/table/proto"
)

type watcher struct {
	sync.RWMutex
	opts    table.WatchOptions
	resChan chan *table.Event
	done    chan struct{}
}

func newWatcher(rsp pb.Table_WatchService, opts table.WatchOptions) (*watcher, error) {
	w := &watcher{
		opts:    opts,
		resChan: make(chan *table.Event),
		done:    make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-w.done:
				return
			default:
				if err := w.watch(rsp); err != nil {
					w.Stop()
					return
				}
			}
		}
	}()

	return w, nil
}

// watchRouter watches router and send events to all registered watchers
func (w *watcher) watch(stream pb.Table_WatchService) error {
	defer stream.Close()

	var watchErr error

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				watchErr = err
			}
			break
		}

		route := table.Route{
			Service: resp.Route.Service,
			Address: resp.Route.Address,
			Gateway: resp.Route.Gateway,
			Network: resp.Route.Network,
			Link:    resp.Route.Link,
			Metric:  int(resp.Route.Metric),
		}

		event := &table.Event{
			Type:      table.EventType(resp.Type),
			Timestamp: time.Unix(0, resp.Timestamp),
			Route:     route,
		}

		select {
		case w.resChan <- event:
		case <-w.done:
		}
	}

	return watchErr
}

// Next is a blocking call that returns watch result
func (w *watcher) Next() (*table.Event, error) {
	for {
		select {
		case res := <-w.resChan:
			switch w.opts.Service {
			case res.Route.Service, "*":
				return res, nil
			default:
				continue
			}
		case <-w.done:
			return nil, table.ErrWatcherStopped
		}
	}
}

// Chan returns event channel
func (w *watcher) Chan() (<-chan *table.Event, error) {
	return w.resChan, nil
}

// Stop stops watcher
func (w *watcher) Stop() {
	w.Lock()
	defer w.Unlock()

	select {
	case <-w.done:
		return
	default:
		close(w.done)
	}
}
