package jobs

type updateHistory struct {
	pointer int
	history []Update
}

func newUpdateHistory(maxSize int) *updateHistory {
	return &updateHistory{
		history: make([]Update, 0, maxSize),
		pointer: 0,
	}
}

func (uh *updateHistory) push(update Update) {
	if len(uh.history) < cap(uh.history) {
		uh.history = append(uh.history, update)
	} else {
		uh.history[uh.pointer] = update
	}
	uh.pointer = (uh.pointer + 1) % cap(uh.history)
}

func (uh *updateHistory) list() []Update {
	if len(uh.history) < cap(uh.history) {
		us := make([]Update, len(uh.history))
		copy(us, uh.history)
		return us
	}

	us := make([]Update, len(uh.history))
	p := uh.pointer
	for i := len(uh.history) - 1; i >= 0; i-- {
		if p == 0 {
			p = len(uh.history) - 1
		} else {
			p--
		}
		us[i] = uh.history[p]
	}

	return us
}
