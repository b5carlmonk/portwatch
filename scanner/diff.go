package scanner

import "fmt"

// ChangeType describes what changed for a port.
type ChangeType string

const (
	ChangeOpened ChangeType = "opened"
	ChangeClosed ChangeType = "closed"
)

// PortChange represents a detected change between two scans.
type PortChange struct {
	Port     int
	Protocol string
	Change   ChangeType
	Service  string
}

func (c PortChange) String() string {
	svc := c.Service
	if svc == "" {
		svc = "unknown"
	}
	return fmt.Sprintf("port %d/%s (%s) %s", c.Port, c.Protocol, svc, c.Change)
}

// Diff compares two ScanResults and returns the list of changes.
func Diff(previous, current *ScanResult) []PortChange {
	prevMap := make(map[int]PortState, len(previous.Ports))
	for _, p := range previous.Ports {
		prevMap[p.Port] = p
	}

	currMap := make(map[int]PortState, len(current.Ports))
	for _, p := range current.Ports {
		currMap[p.Port] = p
	}

	var changes []PortChange

	for port, curr := range currMap {
		prev, existed := prevMap[port]
		if curr.Open && (!existed || !prev.Open) {
			changes = append(changes, PortChange{
				Port: port, Protocol: curr.Protocol,
				Change: ChangeOpened, Service: curr.Service,
			})
		} else if !curr.Open && existed && prev.Open {
			changes = append(changes, PortChange{
				Port: port, Protocol: prev.Protocol,
				Change: ChangeClosed, Service: prev.Service,
			})
		}
	}

	return changes
}
