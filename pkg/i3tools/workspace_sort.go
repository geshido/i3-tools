package i3tools

import (
	"go.i3wm.org/i3/v4"
)

func MakeSortableBy(list []i3.Workspace, lessFn func(int, int, i3.Workspace, i3.Workspace) bool) *SortableWorkspaces {
	s := SortableWorkspaces{
		Data:   list,
		lessFn: lessFn,
	}

	s.OriginalOrder = make(map[i3.WorkspaceID]int)
	for i, ws := range list {
		s.OriginalOrder[ws.ID] = i
	}

	return &s
}

type SortableWorkspaces struct {
	Data          []i3.Workspace
	OriginalOrder map[i3.WorkspaceID]int
	lessFn        func(int, int, i3.Workspace, i3.Workspace) bool
}

func (s SortableWorkspaces) Len() int {
	return len(s.Data)
}

func (s SortableWorkspaces) Less(i, j int) bool {
	return s.lessFn(s.OriginalOrder[s.Data[i].ID], s.OriginalOrder[s.Data[j].ID], s.Data[i], s.Data[j])
}

func (s SortableWorkspaces) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func DefaultWorkspaceLessFn(origid1, origid2 int, w1, w2 i3.Workspace) bool {
	if w1.Num != w2.Num {
		return w1.Num < w2.Num
	}

	return origid1 < origid2
}
