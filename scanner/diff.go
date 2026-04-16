package scanner

// DiffResult holds ports that appeared or disappeared between two scans.
type DiffResult struct {
	Opened []Result
	Closed []Result
}

// Diff compares a previous scan result set against a current one.
// It returns ports that were opened (present in current, absent in previous)
// and ports that were closed (present in previous, absent in current).
func Diff(previous, current []Result) DiffResult {
	prevMap := make(map[string]Result, len(previous))
	for _, r := range previous {
		prevMap[key(r)] = r
	}

	currMap := make(map[string]Result, len(current))
	for _, r := range current {
		currMap[key(r)] = r
	}

	var diff DiffResult

	for k, r := range currMap {
		if _, found := prevMap[k]; !found {
			diff.Opened = append(diff.Opened, r)
		}
	}

	for k, r := range prevMap {
		if _, found := currMap[k]; !found {
			diff.Closed = append(diff.Closed, r)
		}
	}

	return diff
}

func key(r Result) string {
	return r.Host + ":" + itoa(r.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
