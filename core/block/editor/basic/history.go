package basic

import (
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/smartblock"
	"github.com/anytypeio/go-anytype-middleware/core/session"
	"github.com/anytypeio/go-anytype-middleware/pb"
	"github.com/anytypeio/go-anytype-middleware/util/pbtypes"
)

type IHistory interface {
	Undo(*session.Context) (counters pb.RpcObjectUndoRedoCounter, err error)
	Redo(*session.Context) (counters pb.RpcObjectUndoRedoCounter, err error)
}

func NewHistory(sb smartblock.SmartBlock) IHistory {
	return &history{sb}
}

type history struct {
	smartblock.SmartBlock
}

func (h *history) Undo(ctx *session.Context) (counters pb.RpcObjectUndoRedoCounter, err error) {
	s := h.NewStateCtx(ctx)
	action, err := h.History().Previous()
	if err != nil {
		return
	}

	for _, b := range action.Add {
		s.Unlink(b.Model().Id)
	}
	for _, b := range action.Remove {
		s.Set(b.Copy())
	}
	for _, b := range action.Change {
		s.Set(b.Before.Copy())
	}
	if action.ObjectTypes != nil {
		ot := make([]string, len(action.ObjectTypes.Before))
		copy(ot, action.ObjectTypes.Before)
		s.SetObjectTypes(ot)
	}
	if action.Relations != nil {
		s.SetExtraRelations(pbtypes.CopyRelations(action.Relations.Before))
	}
	if action.Details != nil {
		s.SetDetails(pbtypes.CopyStruct(action.Details.Before))
	}
	if err = h.Apply(s, smartblock.NoHistory, smartblock.NoRestrictions); err != nil {
		return
	}
	counters.Undo, counters.Redo = h.History().Counters()
	return
}

func (h *history) Redo(ctx *session.Context) (counters pb.RpcObjectUndoRedoCounter, err error) {
	s := h.NewStateCtx(ctx)
	action, err := h.History().Next()
	if err != nil {
		return
	}

	for _, b := range action.Add {
		s.Set(b.Copy())
	}
	for _, b := range action.Remove {
		s.Unlink(b.Model().Id)
	}
	for _, b := range action.Change {
		s.Set(b.After.Copy())
	}
	if action.ObjectTypes != nil {
		ot := make([]string, len(action.ObjectTypes.After))
		copy(ot, action.ObjectTypes.After)
		s.SetObjectTypes(ot)
	}
	if action.Relations != nil {
		s.SetExtraRelations(pbtypes.CopyRelations(action.Relations.After))
	}
	if action.Details != nil {
		s.SetDetails(pbtypes.CopyStruct(action.Details.After))
	}
	if err = h.Apply(s, smartblock.NoHistory, smartblock.NoRestrictions); err != nil {
		return
	}
	counters.Undo, counters.Redo = h.History().Counters()
	return
}
