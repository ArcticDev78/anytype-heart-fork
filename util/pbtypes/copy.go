package pbtypes

import (
	"sync"

	"github.com/gogo/protobuf/types"

	"github.com/anytypeio/go-anytype-middleware/pb"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/util/slice"
)

var bytesPool = &sync.Pool{
	New: func() interface{} {
		return []byte{}
	},
}

func CopyBlock(in *model.Block) (out *model.Block) {
	buf := bytesPool.Get().([]byte)
	size := in.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = in.MarshalToSizedBuffer(buf[:size])
	out = &model.Block{}
	_ = out.Unmarshal(buf[:size])
	bytesPool.Put(buf)
	return
}

// CopyStructMap copies pb struct map, while reusing map values' pointers
func CopyStructMap(in *types.Struct) (out *types.Struct) {
	if in == nil {
		return nil
	}
	if in.Fields == nil {
		return &types.Struct{}
	}

	out = &types.Struct{Fields: make(map[string]*types.Value, len(in.Fields))}
	for k, v := range in.Fields {
		out.Fields[k] = v
	}

	return
}

func CopyStruct(in *types.Struct) (out *types.Struct) {
	if in == nil {
		return nil
	}
	buf := bytesPool.Get().([]byte)
	size := in.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = in.MarshalToSizedBuffer(buf[:size])
	out = &types.Struct{}
	_ = out.Unmarshal(buf[:size])
	if out.Fields == nil && in.Fields != nil {
		out.Fields = make(map[string]*types.Value)
	}
	bytesPool.Put(buf)
	return
}

func CopyVal(in *types.Value) (out *types.Value) {
	if in == nil {
		return nil
	}
	buf := bytesPool.Get().([]byte)
	size := in.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = in.MarshalToSizedBuffer(buf[:size])
	out = &types.Value{}
	_ = out.Unmarshal(buf[:size])

	bytesPool.Put(buf)
	return
}

func CopyRelationLink(in *model.RelationLink) (out *model.RelationLink) {
	return &model.RelationLink{
		Key:    in.Key,
		Format: in.Format,
	}
}

func CopyRelation(in *model.Relation) (out *model.Relation) {
	if in == nil {
		return nil
	}
	buf := bytesPool.Get().([]byte)
	size := in.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = in.MarshalToSizedBuffer(buf[:size])
	out = &model.Relation{}
	_ = out.Unmarshal(buf[:size])

	bytesPool.Put(buf)
	return out
}

func CopyRelationOptions(in []*model.RelationOption) (out []*model.RelationOption) {
	out = make([]*model.RelationOption, len(in))
	for i := range in {
		out[i] = &*in[i]
	}
	return
}

func CopyLayout(in *model.Layout) (out *model.Layout) {
	return &model.Layout{Id: in.Id, Name: in.Name, RequiredRelations: CopyRelations(in.RequiredRelations)}
}

func CopyObjectType(in *model.ObjectType) (out *model.ObjectType) {
	if in == nil {
		return nil
	}

	buf := bytesPool.Get().([]byte)
	size := in.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = in.MarshalToSizedBuffer(buf[:size])
	out = &model.ObjectType{}
	_ = out.Unmarshal(buf[:size])

	bytesPool.Put(buf)
	return out
}

func CopyRelations(in []*model.Relation) (out []*model.Relation) {
	if in == nil {
		return nil
	}
	buf := bytesPool.Get().([]byte)
	inWrapped := model.Relations{Relations: in}
	size := inWrapped.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = inWrapped.MarshalToSizedBuffer(buf[:size])
	outWrapped := &model.Relations{}
	_ = outWrapped.Unmarshal(buf[:size])

	bytesPool.Put(buf)
	return outWrapped.Relations
}

func CopyOptions(in []*model.RelationOption) (out []*model.RelationOption) {
	if in == nil {
		return nil
	}

	for _, inO := range in {
		inCopy := *inO
		out = append(out, &inCopy)
	}
	return
}

func CopyRelationsToMap(in []*model.Relation) (out map[string]*model.Relation) {
	out = make(map[string]*model.Relation, len(in))
	rels := CopyRelations(in)
	for _, rel := range rels {
		out[rel.Key] = rel
	}

	return
}

func CopyFilter(in *model.BlockContentDataviewFilter) (out *model.BlockContentDataviewFilter) {
	buf := bytesPool.Get().([]byte)
	size := in.Size()
	if cap(buf) < size {
		buf = make([]byte, 0, size*2)
	}
	size, _ = in.MarshalToSizedBuffer(buf[:size])
	out = &model.BlockContentDataviewFilter{}
	_ = out.Unmarshal(buf[:size])
	bytesPool.Put(buf)
	return
}

func RelationsFilterKeys(in []*model.Relation, keys []string) (out []*model.Relation) {
	for i, inRel := range in {
		if slice.FindPos(keys, inRel.Key) >= 0 {
			out = append(out, in[i])
		}
	}
	return
}

func StructNotNilKeys(st *types.Struct) (keys []string) {
	if st == nil || st.Fields == nil {
		return nil
	}

	for k, v := range st.Fields {
		if v != nil {
			keys = append(keys, k)
		}
	}
	return
}

func EventsToSliceChange(changes []*pb.EventBlockDataviewSliceChange) []slice.Change {
	sliceOpMap := map[pb.EventBlockDataviewSliceOperation]slice.DiffOperation{
		pb.EventBlockDataview_SliceOperationNone:    slice.OperationNone,
		pb.EventBlockDataview_SliceOperationAdd:     slice.OperationAdd,
		pb.EventBlockDataview_SliceOperationMove:    slice.OperationMove,
		pb.EventBlockDataview_SliceOperationRemove:  slice.OperationRemove,
		pb.EventBlockDataview_SliceOperationReplace: slice.OperationReplace,
	}

	var res []slice.Change
	for _, eventCh := range changes {
		res = append(res, slice.Change{Op: sliceOpMap[eventCh.Op], Ids: eventCh.Ids, AfterId: eventCh.AfterId})
	}

	return res
}

func SliceChangeToEvents(changes []slice.Change) []*pb.EventBlockDataviewSliceChange {
	eventsOpMap := map[slice.DiffOperation]pb.EventBlockDataviewSliceOperation{
		slice.OperationNone:    pb.EventBlockDataview_SliceOperationNone,
		slice.OperationAdd:     pb.EventBlockDataview_SliceOperationAdd,
		slice.OperationMove:    pb.EventBlockDataview_SliceOperationMove,
		slice.OperationRemove:  pb.EventBlockDataview_SliceOperationRemove,
		slice.OperationReplace: pb.EventBlockDataview_SliceOperationReplace,
	}

	var res []*pb.EventBlockDataviewSliceChange
	for _, sliceCh := range changes {
		res = append(res, &pb.EventBlockDataviewSliceChange{Op: eventsOpMap[sliceCh.Op], Ids: sliceCh.Ids, AfterId: sliceCh.AfterId})
	}

	return res
}
