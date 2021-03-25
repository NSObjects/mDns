package registry

import (
	"context"
)

type servicesKey struct{}

//func getServiceRecords(ctx context.Context) map[string]map[string]*record {
//	memServices, ok := ctx.Value(servicesKey{}).(map[string][]*Service)
//	if !ok {
//		return nil
//	}
//
//	services := make(map[string]map[string]*record)
//
//	for name, svc := range memServices {
//		if _, ok := services[name]; !ok {
//			services[name] = make(map[string]*record)
//		}
//		// go through every version of the service
//		for _, s := range svc {
//			services[s.Name][s.Version] = serviceToRecord(s, 0)
//		}
//	}
//
//	return services
//}
//
//func serviceToRecord(s *Service, ttl time.Duration) *record {
//	metadata := make(map[string]string, len(s.Metadata))
//	for k, v := range s.Metadata {
//		metadata[k] = v
//	}
//
//	nodes := make(map[string]*node, len(s.Nodes))
//	for _, n := range s.Nodes {
//		nodes[n.Id] = &node{
//			Node:     n,
//			TTL:      ttl,
//			LastSeen: time.Now(),
//		}
//	}
//
//	endpoints := make([]*Endpoint, len(s.Endpoints))
//	for i, e := range s.Endpoints {
//		endpoints[i] = e
//	}
//
//	return &record{
//		Name:      s.Name,
//		Version:   s.Version,
//		Metadata:  metadata,
//		Nodes:     nodes,
//		Endpoints: endpoints,
//	}
//}

// Services is an option that preloads service data
func Services(s map[string][]*Service) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, servicesKey{}, s)
	}
}

type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}
