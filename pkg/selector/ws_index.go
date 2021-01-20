package selector

import (
	"fmt"
	"github.com/geshido/i3-tools/pkg/i3tools"
	"go.i3wm.org/i3/v4"
	"sort"
	"strconv"
)

type wsByIndex struct {
	index int
}

func NewByIndex(index int) *wsByIndex {
	return &wsByIndex{index: index}
}

func (s wsByIndex) Select() (string, error) {
	wsList, err := i3.GetWorkspaces()
	if err != nil {
		return "", fmt.Errorf("cannot get list of workspaces: %w", err)
	}

	if s.index < 1 || s.index > len(wsList) {
		if s.index < 1 {
			return "1", nil
		}
		return strconv.Itoa(s.index), nil
	}

	sortable := i3tools.MakeSortableBy(wsList, i3tools.DefaultWorkspaceLessFn)
	sort.Sort(sortable)

	return sortable.Data[s.index-1].Name, nil
}
