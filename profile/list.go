package profile

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// List writes a human-readable summary of all profiles to w.
func (s *Store) List(w io.Writer) {
	if len(s.Profiles) == 0 {
		fmt.Fprintln(w, "no profiles defined")
		return
	}
	names := make([]string, 0, len(s.Profiles))
	for name := range s.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		p := s.Profiles[name]
		ports := make([]string, len(p.Ports))
		for i, port := range p.Ports {
			ports[i] = fmt.Sprintf("%d", port)
		}
		fmt.Fprintf(w, "%-20s hosts=%-20s ports=%s proto=%s\n",
			p.Name,
			strings.Join(p.Hosts, ","),
			strings.Join(ports, ","),
			p.Protocol,
		)
	}
}
