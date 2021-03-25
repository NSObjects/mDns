package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/registry"

	"github.com/NSObjects/mDns/registry/mdns"
)

type mdnsWatcher struct {
	id   string
	wo   WatchOptions
	ch   chan *mdns.ServiceEntry
	exit chan struct{}
	// the mdns domain
	domain string
	// the registry
	registry *mdnsRegistry
}

func (m *mdnsWatcher) Next() ([]*registry.ServiceInstance, error) {
	for {
		select {
		case e := <-m.ch:
			txt, err := decode(e.InfoFields)
			if err != nil {
				continue
			}

			if len(txt.Service) == 0 || len(txt.Version) == 0 {
				continue
			}

			// Filter watch options
			// wo.Service: Only keep services we care about
			if len(m.wo.Service) > 0 && txt.Service != m.wo.Service {
				continue
			}
			var action string
			if e.TTL == 0 {
				action = "delete"
			} else {
				action = "create"
			}

			service := &registry.ServiceInstance{
				Name:      txt.Service,
				Version:   txt.Version,
				Endpoints: txt.Endpoints,
			}

			// skip anything without the domain we care about
			suffix := fmt.Sprintf(".%s.%s.", service.Name, m.domain)
			if !strings.HasSuffix(e.Name, suffix) {
				continue
			}

			var addr string
			if len(e.AddrV4) > 0 {
				addr = e.AddrV4.String()
			} else if len(e.AddrV6) > 0 {
				addr = "[" + e.AddrV6.String() + "]"
			} else {
				addr = e.Addr.String()
			}

			//service.Nodes = append(service.Nodes, &Node{
			//	Id:       strings.TrimSuffix(e.Name, suffix),
			//	Address:  fmt.Sprintf("%s:%d", addr, e.Port),
			//	Metadata: txt.Metadata,
			//})

			var ss  []*registry.ServiceInstance
			for _,v :=range m.registry.watchers {
				ss = append(ss,&registry.ServiceInstance{
					ID:        ,
					Name:      "",
					Version:   "",
					Metadata:  nil,
					Endpoints: nil,
				})
			}

			return &Result{
				Action:  action,
				Service: service,
			}, nil
		case <-m.exit:
			return nil, ErrWatcherStopped
		}
	}
}

func (m *mdnsWatcher) Stop() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
		// remove self from the registry
		m.registry.mtx.Lock()
		delete(m.registry.watchers, m.id)
		m.registry.mtx.Unlock()
	}
	return nil
}

//
//func (m *mdnsWatcher) Next() (*Result, error) {
//
//}
//
//func (m *mdnsWatcher) Stop() {
//
//}

// Watch a service
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}
