package gnuflag

type argList struct {
	src  []string
	next int
}

func (a *argList) Next() bool {
	if a.next >= len(a.src) {
		return false
	}
	a.next++
	return true
}

func (a *argList) Value() string {
	return a.src[a.next-1]
}
