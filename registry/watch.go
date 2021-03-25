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
	var services []*registry.ServiceInstance
	select {
	case e := <-m.ch:
		txt, err := decode(e.InfoFields)
		if err != nil {
			return nil, err
		}
		if len(txt.Service) == 0 || len(txt.Version) == 0 {
			return nil, err
		}
		// Filter watch options
		// wo.Service: Only keep services we care about
		if len(m.wo.Service) > 0 && txt.Service != m.wo.Service {
			break
		}

		// skip anything without the domain we care about
		suffix := fmt.Sprintf(".%s.%s.", txt.Service, m.domain)
		if !strings.HasSuffix(e.Name, suffix) {
			break
		}

		var addr string
		if len(e.AddrV4) > 0 {
			addr = e.AddrV4.String()
		} else if len(e.AddrV6) > 0 {
			addr = "[" + e.AddrV6.String() + "]"
		} else {
			addr = e.Addr.String()
		}

		txt.Endpoints = append(txt.Endpoints, fmt.Sprintf("%s:%d", addr, e.Port))
		services = append(services, &registry.ServiceInstance{
			ID:        strings.TrimSuffix(e.Name, suffix),
			Name:      e.Name,
			Version:   txt.Version,
			Metadata:  txt.Metadata,
			Endpoints: txt.Endpoints,
		})
	case <-m.exit:
		return nil, ErrWatcherStopped
	}

	for _, v := range m.registry.watchers {
		if v.id == m.id {
			continue
		}
		e := <-v.ch

		txt, err := decode(e.InfoFields)
		if err != nil {
			return nil, err
		}
		if len(txt.Service) == 0 || len(txt.Version) == 0 {
			return nil, err
		}
		// Filter watch options
		// wo.Service: Only keep services we care about
		if len(m.wo.Service) > 0 && txt.Service != m.wo.Service {
			continue
		}
		// skip anything without the domain we care about
		suffix := fmt.Sprintf(".%s.%s.", txt.Service, m.domain)
		if !strings.HasSuffix(e.Name, suffix) {
			break
		}

		var addr string
		if len(e.AddrV4) > 0 {
			addr = e.AddrV4.String()
		} else if len(e.AddrV6) > 0 {
			addr = "[" + e.AddrV6.String() + "]"
		} else {
			addr = e.Addr.String()
		}

		txt.Endpoints = append(txt.Endpoints, fmt.Sprintf("%s:%d", addr, e.Port))
		services = append(services, &registry.ServiceInstance{
			ID:        strings.TrimSuffix(e.Name, suffix),
			Name:      e.Name,
			Version:   txt.Version,
			Metadata:  txt.Metadata,
			Endpoints: txt.Endpoints,
		})

	}

	return services, nil

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
