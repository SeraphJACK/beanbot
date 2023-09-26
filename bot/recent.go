package bot

import (
	"sort"
	"time"
)

type recentCmd struct {
	command string
	t       time.Time
}

var recentCmds []recentCmd

type recentCmdSort []recentCmd

func (r recentCmdSort) Len() int {
	return len(r)
}

func (r recentCmdSort) Less(i, j int) bool {
	return r[i].t.Before(r[j].t)
}

func (r recentCmdSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func updateRecentCommand(command string) {
	exists := false
	for _, v := range recentCmds {
		if v.command == command {
			v.t = time.Now()
			exists = true
			break
		}
	}

	if !exists {
		recentCmds = append([]recentCmd{{command: command, t: time.Now()}}, recentCmds...)
		if len(recentCmds) > 5 {
			recentCmds = recentCmds[:5]
		}
	} else {
		sort.Sort(recentCmdSort(recentCmds))
	}
}
