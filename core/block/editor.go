package block

import (
	"context"
	"fmt"
	"time"

	"github.com/anytypeio/go-anytype-middleware/core/block/editor/table"

	"github.com/anytypeio/go-anytype-middleware/core/block/simple/link"
	"github.com/anytypeio/go-anytype-middleware/core/block/source"
	"github.com/anytypeio/go-anytype-middleware/metrics"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/schema"
	"github.com/anytypeio/go-anytype-middleware/util/internalflag"
	"github.com/anytypeio/go-anytype-middleware/util/ocache"
	ds "github.com/ipfs/go-datastore"
	"github.com/textileio/go-threads/core/thread"

	"github.com/anytypeio/go-anytype-middleware/core/block/doc"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/basic"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/bookmark"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/clipboard"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/dataview"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/file"
	_import "github.com/anytypeio/go-anytype-middleware/core/block/editor/import"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/smartblock"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/state"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/stext"
	"github.com/anytypeio/go-anytype-middleware/core/block/simple"
	"github.com/anytypeio/go-anytype-middleware/core/block/simple/text"
	"github.com/anytypeio/go-anytype-middleware/pb"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/bundle"
	coresb "github.com/anytypeio/go-anytype-middleware/pkg/lib/core/smartblock"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/localstore/objectstore"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/threads"
	"github.com/anytypeio/go-anytype-middleware/util/pbtypes"
	"github.com/gogo/protobuf/types"
)

var ErrOptionUsedByOtherObjects = fmt.Errorf("option is used by other objects")

func (s *service) MarkArchived(id string, archived bool) (err error) {
	return s.Do(id, func(b smartblock.SmartBlock) error {
		return b.SetDetails(nil, []*pb.RpcObjectSetDetailsDetail{
			{
				Key:   "isArchived",
				Value: pbtypes.Bool(archived),
			},
		}, true)
	})
}

func (s *service) SetBreadcrumbs(ctx *state.Context, req pb.RpcObjectSetBreadcrumbsRequest) (err error) {
	return s.Do(req.BreadcrumbsId, func(b smartblock.SmartBlock) error {
		if breadcrumbs, ok := b.(*editor.Breadcrumbs); ok {
			return breadcrumbs.SetCrumbs(req.Ids)
		} else {
			return ErrUnexpectedBlockType
		}
	})
}

func (s *service) CreateBlock(ctx *state.Context, req pb.RpcBlockCreateRequest) (id string, err error) {
	err = s.DoBasic(req.ContextId, func(b basic.Basic) error {
		id, err = b.Create(ctx, "", req)
		return err
	})
	return
}

func (s *service) DuplicateBlocks(ctx *state.Context, req pb.RpcBlockListDuplicateRequest) (newIds []string, err error) {
	err = s.DoBasic(req.ContextId, func(b basic.Basic) error {
		newIds, err = b.Duplicate(ctx, req)
		return err
	})
	return
}

func (s *service) UnlinkBlock(ctx *state.Context, req pb.RpcBlockListDeleteRequest) (err error) {
	return s.DoBasic(req.ContextId, func(b basic.Basic) error {
		return b.Unlink(ctx, req.BlockIds...)
	})
}

func (s *service) SetDivStyle(ctx *state.Context, contextId string, style model.BlockContentDivStyle, ids ...string) (err error) {
	return s.DoBasic(contextId, func(b basic.Basic) error {
		return b.SetDivStyle(ctx, style, ids...)
	})
}

func (s *service) SplitBlock(ctx *state.Context, req pb.RpcBlockSplitRequest) (blockId string, err error) {
	err = s.DoText(req.ContextId, func(b stext.Text) error {
		blockId, err = b.Split(ctx, req)
		return err
	})
	return
}

func (s *service) MergeBlock(ctx *state.Context, req pb.RpcBlockMergeRequest) (err error) {
	return s.DoText(req.ContextId, func(b stext.Text) error {
		return b.Merge(ctx, req.FirstBlockId, req.SecondBlockId)
	})
}

func (s *service) TurnInto(ctx *state.Context, contextId string, style model.BlockContentTextStyle, ids ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.TurnInto(ctx, style, ids...)
	})
}

func (s *service) SimplePaste(contextId string, anySlot []*model.Block) (err error) {
	var blocks []simple.Block

	for _, b := range anySlot {
		blocks = append(blocks, simple.New(b))
	}

	return s.DoBasic(contextId, func(b basic.Basic) error {
		return b.PasteBlocks(blocks)
	})
}

func (s *service) ReplaceBlock(ctx *state.Context, req pb.RpcBlockReplaceRequest) (newId string, err error) {
	err = s.DoBasic(req.ContextId, func(b basic.Basic) error {
		newId, err = b.Replace(ctx, req.BlockId, req.Block)
		return err
	})
	return
}

func (s *service) SetFields(ctx *state.Context, req pb.RpcBlockSetFieldsRequest) (err error) {
	return s.DoBasic(req.ContextId, func(b basic.Basic) error {
		return b.SetFields(ctx, &pb.RpcBlockListSetFieldsRequestBlockField{
			BlockId: req.BlockId,
			Fields:  req.Fields,
		})
	})
}

func (s *service) SetDetails(ctx *state.Context, req pb.RpcObjectSetDetailsRequest) (err error) {
	return s.Do(req.ContextId, func(b smartblock.SmartBlock) error {
		return b.SetDetails(ctx, req.Details, true)
	})
}

func (s *service) SetFieldsList(ctx *state.Context, req pb.RpcBlockListSetFieldsRequest) (err error) {
	return s.DoBasic(req.ContextId, func(b basic.Basic) error {
		return b.SetFields(ctx, req.BlockFields...)
	})
}

func (s *service) GetAggregatedRelations(req pb.RpcBlockDataviewRelationListAvailableRequest) (relations []*model.Relation, err error) {
	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		relations, err = b.GetAggregatedRelations(req.BlockId)
		return err
	})

	return
}

func (s *service) UpdateDataviewView(ctx *state.Context, req pb.RpcBlockDataviewViewUpdateRequest) error {
	return s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.UpdateView(ctx, req.BlockId, req.ViewId, *req.View, true)
	})
}

func (s *service) DeleteDataviewView(ctx *state.Context, req pb.RpcBlockDataviewViewDeleteRequest) error {
	return s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.DeleteView(ctx, req.BlockId, req.ViewId, true)
	})
}

func (s *service) SetDataviewActiveView(ctx *state.Context, req pb.RpcBlockDataviewViewSetActiveRequest) error {
	return s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.SetActiveView(ctx, req.BlockId, req.ViewId, int(req.Limit), int(req.Offset))
	})
}

func (s *service) SetDataviewViewPosition(ctx *state.Context, req pb.RpcBlockDataviewViewSetPositionRequest) error {
	return s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.SetViewPosition(ctx, req.BlockId, req.ViewId, req.Position)
	})
}

func (s *service) CreateDataviewView(ctx *state.Context, req pb.RpcBlockDataviewViewCreateRequest) (id string, err error) {
	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		if req.View == nil {
			req.View = &model.BlockContentDataviewView{}
		}
		view, err := b.CreateView(ctx, req.BlockId, *req.View)
		id = view.Id
		return err
	})

	return
}

func (s *service) CreateDataviewRecord(ctx *state.Context, req pb.RpcBlockDataviewRecordCreateRequest) (rec *types.Struct, err error) {
	var workspaceId string
	sbt, err := coresb.SmartBlockTypeFromID(req.ContextId)
	if err == nil && sbt == coresb.SmartBlockTypeWorkspace {
		workspaceId = req.ContextId
	} else {
		workspaceId, err = s.anytype.GetWorkspaceIdForObject(req.ContextId)
		if err != nil {
			threads.WorkspaceLogger.Debugf("cannot get workspace id for object: %v", err)
		}
	}
	if workspaceId != "" {
		if req.Record == nil {
			req.Record = &types.Struct{Fields: make(map[string]*types.Value)}
		}
		// todo: maybe this check is not needed?
		if req.Record.Fields == nil {
			req.Record.Fields = make(map[string]*types.Value)
		}
		req.Record.Fields[bundle.RelationKeyWorkspaceId.String()] = pbtypes.String(workspaceId)
	}

	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		cr, err := b.CreateRecord(ctx, req.BlockId, model.ObjectDetails{Details: req.Record}, req.TemplateId)
		if err != nil {
			return err
		}
		rec = cr.Details
		return nil
	})

	return
}

func (s *service) UpdateDataviewRecord(ctx *state.Context, req pb.RpcBlockDataviewRecordUpdateRequest) (err error) {
	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.UpdateRecord(ctx, req.BlockId, req.RecordId, model.ObjectDetails{Details: req.Record})
	})

	return
}

func (s *service) DeleteDataviewRecord(ctx *state.Context, req pb.RpcBlockDataviewRecordDeleteRequest) (err error) {
	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.DeleteRecord(ctx, req.BlockId, req.RecordId)
	})

	return
}

func (s *service) UpdateDataviewRelation(ctx *state.Context, req pb.RpcBlockDataviewRelationUpdateRequest) error {
	return s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.UpdateRelation(ctx, req.BlockId, req.RelationKey, *req.Relation, true)
	})
}

func (s *service) AddDataviewRelation(ctx *state.Context, req pb.RpcBlockDataviewRelationAddRequest) (relation *model.Relation, err error) {
	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		var err error
		relation, err = b.AddRelation(ctx, req.BlockId, *req.Relation, true)
		if err != nil {
			return err
		}
		rels, err := b.GetDataviewRelations(req.BlockId)
		if err != nil {
			return err
		}

		relation = pbtypes.GetRelation(rels, relation.Key)
		if relation.Format == model.RelationFormat_status || relation.Format == model.RelationFormat_tag {
			err = b.FillAggregatedOptions(nil)
			if err != nil {
				log.Errorf("FillAggregatedOptions failed: %s", err.Error())
			}
		}
		return nil
	})

	return
}

func (s *service) DeleteDataviewRelation(ctx *state.Context, req pb.RpcBlockDataviewRelationDeleteRequest) error {
	return s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		return b.DeleteRelation(ctx, req.BlockId, req.RelationKey, true)
	})
}

func (s *service) AddDataviewRecordRelationOption(ctx *state.Context, req pb.RpcBlockDataviewRecordRelationOptionAddRequest) (opt *model.RelationOption, err error) {
	err = s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		opt, err = b.AddRelationOption(ctx, req.BlockId, req.RecordId, req.RelationKey, *req.Option, true)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func (s *service) UpdateDataviewRecordRelationOption(ctx *state.Context, req pb.RpcBlockDataviewRecordRelationOptionUpdateRequest) error {
	err := s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		err := b.UpdateRelationOption(ctx, req.BlockId, req.RecordId, req.RelationKey, *req.Option, true)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *service) DeleteDataviewRecordRelationOption(ctx *state.Context, req pb.RpcBlockDataviewRecordRelationOptionDeleteRequest) error {
	err := s.DoDataview(req.ContextId, func(b dataview.Dataview) error {
		err := b.DeleteRelationOption(ctx, true, req.BlockId, req.RecordId, req.RelationKey, req.OptionId, true)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *service) SetDataviewSource(ctx *state.Context, contextId, blockId string, source []string) (err error) {
	return s.DoDataview(contextId, func(b dataview.Dataview) error {
		return b.SetSource(ctx, blockId, source)
	})
}

func (s *service) Copy(req pb.RpcBlockCopyRequest) (textSlot string, htmlSlot string, anySlot []*model.Block, err error) {
	err = s.DoClipboard(req.ContextId, func(cb clipboard.Clipboard) error {
		textSlot, htmlSlot, anySlot, err = cb.Copy(req)
		return err
	})

	return textSlot, htmlSlot, anySlot, err
}

func (s *service) Paste(ctx *state.Context, req pb.RpcBlockPasteRequest, groupId string) (blockIds []string, uploadArr []pb.RpcBlockUploadRequest, caretPosition int32, isSameBlockCaret bool, err error) {
	err = s.DoClipboard(req.ContextId, func(cb clipboard.Clipboard) error {
		blockIds, uploadArr, caretPosition, isSameBlockCaret, err = cb.Paste(ctx, &req, groupId)
		return err
	})

	return blockIds, uploadArr, caretPosition, isSameBlockCaret, err
}

func (s *service) Cut(ctx *state.Context, req pb.RpcBlockCutRequest) (textSlot string, htmlSlot string, anySlot []*model.Block, err error) {
	err = s.DoClipboard(req.ContextId, func(cb clipboard.Clipboard) error {
		textSlot, htmlSlot, anySlot, err = cb.Cut(ctx, req)
		return err
	})
	return textSlot, htmlSlot, anySlot, err
}

func (s *service) Export(req pb.RpcBlockExportRequest) (path string, err error) {
	err = s.DoClipboard(req.ContextId, func(cb clipboard.Clipboard) error {
		path, err = cb.Export(req)
		return err
	})
	return path, err
}

func (s *service) ImportMarkdown(ctx *state.Context, req pb.RpcObjectImportMarkdownRequest) (rootLinkIds []string, err error) {
	var rootLinks []*model.Block
	err = s.DoImport(req.ContextId, func(imp _import.Import) error {
		rootLinks, err = imp.ImportMarkdown(ctx, req)
		return err
	})
	if err != nil {
		return rootLinkIds, err
	}

	if len(rootLinks) == 1 {
		err = s.SimplePaste(req.ContextId, rootLinks)

		if err != nil {
			return rootLinkIds, err
		}
	} else {
		_, pageId, err := s.CreateLinkToTheNewObject(ctx, "", pb.RpcBlockLinkCreateWithObjectRequest{
			ContextId: req.ContextId,
			Details: &types.Struct{Fields: map[string]*types.Value{
				"name":      pbtypes.String("Import from Notion"),
				"iconEmoji": pbtypes.String("📁"),
			}},
		})

		if err != nil {
			return rootLinkIds, err
		}

		err = s.SimplePaste(pageId, rootLinks)
	}

	for _, r := range rootLinks {
		rootLinkIds = append(rootLinkIds, r.Id)
	}

	return rootLinkIds, err
}

func (s *service) SetTextText(ctx *state.Context, req pb.RpcBlockTextSetTextRequest) error {
	return s.DoText(req.ContextId, func(b stext.Text) error {
		return b.SetText(ctx, req)
	})
}

func (s *service) SetLatexText(ctx *state.Context, req pb.RpcBlockLatexSetTextRequest) error {
	return s.Do(req.ContextId, func(b smartblock.SmartBlock) error {
		return b.(basic.Basic).SetLatexText(ctx, req)
	})
}

func (s *service) SetTextStyle(ctx *state.Context, contextId string, style model.BlockContentTextStyle, blockIds ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.UpdateTextBlocks(ctx, blockIds, true, func(t text.Block) error {
			t.SetStyle(style)
			return nil
		})
	})
}

func (s *service) SetTextChecked(ctx *state.Context, req pb.RpcBlockTextSetCheckedRequest) error {
	return s.DoText(req.ContextId, func(b stext.Text) error {
		return b.UpdateTextBlocks(ctx, []string{req.BlockId}, true, func(t text.Block) error {
			t.SetChecked(req.Checked)
			return nil
		})
	})
}

func (s *service) SetTextColor(ctx *state.Context, contextId string, color string, blockIds ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.UpdateTextBlocks(ctx, blockIds, true, func(t text.Block) error {
			t.SetTextColor(color)
			return nil
		})
	})
}

func (s *service) ClearTextStyle(ctx *state.Context, contextId string, blockIds ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.UpdateTextBlocks(ctx, blockIds, true, func(t text.Block) error {
			t.Model().BackgroundColor = ""
			t.Model().Align = model.Block_AlignLeft
			t.Model().VerticalAlign = model.Block_VerticalAlignTop
			t.SetTextColor("")
			t.SetStyle(model.BlockContentText_Paragraph)
			return nil
		})
	})
}

func (s *service) ClearTextContent(ctx *state.Context, contextId string, blockIds ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.UpdateTextBlocks(ctx, blockIds, true, func(t text.Block) error {
			return t.SetText("", nil)
		})
	})
}

func (s *service) SetTextMark(ctx *state.Context, contextId string, mark *model.BlockContentTextMark, blockIds ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.SetMark(ctx, mark, blockIds...)
	})
}

func (s *service) SetTextIcon(ctx *state.Context, contextId, image, emoji string, blockIds ...string) error {
	return s.DoText(contextId, func(b stext.Text) error {
		return b.SetIcon(ctx, image, emoji, blockIds...)
	})
}

func (s *service) SetBackgroundColor(ctx *state.Context, contextId string, color string, blockIds ...string) (err error) {
	return s.DoBasic(contextId, func(b basic.Basic) error {
		return b.Update(ctx, func(b simple.Block) error {
			b.Model().BackgroundColor = color
			return nil
		}, blockIds...)
	})
}

func (s *service) SetLinkAppearance(ctx *state.Context, req pb.RpcBlockLinkListSetAppearanceRequest) (err error) {
	return s.DoBasic(req.ContextId, func(b basic.Basic) error {
		return b.Update(ctx, func(b simple.Block) error {
			if linkBlock, ok := b.(link.Block); ok {
				return linkBlock.SetAppearance(&model.BlockContentLink{
					IconSize:    req.IconSize,
					CardStyle:   req.CardStyle,
					Description: req.Description,
					Relations:   req.Relations,
				})
			}
			return nil
		}, req.BlockIds...)
	})
}

func (s *service) SetAlign(ctx *state.Context, contextId string, align model.BlockAlign, blockIds ...string) (err error) {
	return s.Do(contextId, func(sb smartblock.SmartBlock) error {
		return sb.SetAlign(ctx, align, blockIds...)
	})
}

func (s *service) SetVerticalAlign(ctx *state.Context, contextId string, align model.BlockVerticalAlign, blockIds ...string) (err error) {
	return s.Do(contextId, func(sb smartblock.SmartBlock) error {
		return sb.SetVerticalAlign(ctx, align, blockIds...)
	})
}

func (s *service) SetLayout(ctx *state.Context, contextId string, layout model.ObjectTypeLayout) (err error) {
	return s.Do(contextId, func(sb smartblock.SmartBlock) error {
		return sb.SetLayout(ctx, layout)
	})
}

func (s *service) FeaturedRelationAdd(ctx *state.Context, contextId string, relations ...string) error {
	return s.DoBasic(contextId, func(b basic.Basic) error {
		return b.FeaturedRelationAdd(ctx, relations...)
	})
}

func (s *service) FeaturedRelationRemove(ctx *state.Context, contextId string, relations ...string) error {
	return s.DoBasic(contextId, func(b basic.Basic) error {
		return b.FeaturedRelationRemove(ctx, relations...)
	})
}

func (s *service) UploadBlockFile(ctx *state.Context, req pb.RpcBlockUploadRequest, groupId string) (err error) {
	return s.DoFile(req.ContextId, func(b file.File) error {
		err = b.Upload(ctx, req.BlockId, file.FileSource{
			Path:    req.FilePath,
			Url:     req.Url,
			GroupId: groupId,
		}, false)
		return err
	})
}

func (s *service) UploadBlockFileSync(ctx *state.Context, req pb.RpcBlockUploadRequest) (err error) {
	return s.DoFile(req.ContextId, func(b file.File) error {
		err = b.Upload(ctx, req.BlockId, file.FileSource{
			Path: req.FilePath,
			Url:  req.Url,
		}, true)
		return err
	})
}

func (s *service) CreateAndUploadFile(ctx *state.Context, req pb.RpcBlockFileCreateAndUploadRequest) (id string, err error) {
	err = s.DoFile(req.ContextId, func(b file.File) error {
		id, err = b.CreateAndUpload(ctx, req)
		return err
	})
	return
}

func (s *service) UploadFile(req pb.RpcFileUploadRequest) (hash string, err error) {
	upl := file.NewUploader(s)
	if req.DisableEncryption {
		log.Errorf("DisableEncryption is deprecated and has no effect")
	}

	upl.SetStyle(req.Style)
	if req.Type != model.BlockContentFile_None {
		upl.SetType(req.Type)
	} else {
		upl.AutoType(true)
	}
	res := upl.SetFile(req.LocalPath).Upload(context.TODO())
	if res.Err != nil {
		return "", res.Err
	}
	return res.Hash, nil
}

func (s *service) DropFiles(req pb.RpcFileDropRequest) (err error) {
	return s.DoFileNonLock(req.ContextId, func(b file.File) error {
		return b.DropFiles(req)
	})
}

func (s *service) SetFileStyle(ctx *state.Context, contextId string, style model.BlockContentFileStyle, blockIds ...string) error {
	return s.DoFile(contextId, func(b file.File) error {
		return b.SetFileStyle(ctx, style, blockIds...)
	})
}

func (s *service) Undo(ctx *state.Context, req pb.RpcObjectUndoRequest) (counters pb.RpcObjectUndoRedoCounter, err error) {
	err = s.DoHistory(req.ContextId, func(b basic.IHistory) error {
		counters, err = b.Undo(ctx)
		return err
	})
	return
}

func (s *service) Redo(ctx *state.Context, req pb.RpcObjectRedoRequest) (counters pb.RpcObjectUndoRedoCounter, err error) {
	err = s.DoHistory(req.ContextId, func(b basic.IHistory) error {
		counters, err = b.Redo(ctx)
		return err
	})
	return
}

func (s *service) BookmarkFetch(ctx *state.Context, req pb.RpcBlockBookmarkFetchRequest) (err error) {
	return s.DoBookmark(req.ContextId, func(b bookmark.Bookmark) error {
		return b.Fetch(ctx, req.BlockId, req.Url, false)
	})
}

func (s *service) BookmarkFetchSync(ctx *state.Context, req pb.RpcBlockBookmarkFetchRequest) (err error) {
	return s.DoBookmark(req.ContextId, func(b bookmark.Bookmark) error {
		return b.Fetch(ctx, req.BlockId, req.Url, true)
	})
}

func (s *service) BookmarkCreateAndFetch(ctx *state.Context, req pb.RpcBlockBookmarkCreateAndFetchRequest) (id string, err error) {
	err = s.DoBookmark(req.ContextId, func(b bookmark.Bookmark) error {
		id, err = b.CreateAndFetch(ctx, req)
		return err
	})
	return
}

func (s *service) SetRelationKey(ctx *state.Context, req pb.RpcBlockRelationSetKeyRequest) error {
	return s.Do(req.ContextId, func(b smartblock.SmartBlock) error {
		rels := b.Relations()
		rel := pbtypes.GetRelation(rels, req.Key)
		if rel == nil {
			var err error
			rels, err = s.Anytype().ObjectStore().ListRelations("")
			if err != nil {
				return err
			}
			rel = pbtypes.GetRelation(rels, req.Key)
			if rel == nil {
				return fmt.Errorf("relation with provided key not found")
			}
		}

		return b.(basic.Basic).AddRelationAndSet(ctx, pb.RpcBlockRelationAddRequest{Relation: rel, BlockId: req.BlockId, ContextId: req.ContextId})
	})
}

func (s *service) AddRelationBlock(ctx *state.Context, req pb.RpcBlockRelationAddRequest) error {
	return s.DoBasic(req.ContextId, func(b basic.Basic) error {
		return b.AddRelationAndSet(ctx, req)
	})
}

func (s *service) GetDocInfo(ctx context.Context, id string) (info doc.DocInfo, err error) {
	if err = s.DoWithContext(ctx, id, func(b smartblock.SmartBlock) error {
		info, err = b.GetDocInfo()
		return err
	}); err != nil {
		return
	}
	return
}

func (s *service) Wakeup(id string) (err error) {
	return s.Do(id, func(b smartblock.SmartBlock) error {
		return nil
	})
}

func (s *service) GetRelations(objectId string) (relations []*model.Relation, err error) {
	err = s.Do(objectId, func(b smartblock.SmartBlock) error {
		relations = b.Relations()
		return nil
	})
	return
}

// ModifyExtraRelations gets and updates extra relations under the sb lock to make sure no modifications are done in the middle
func (s *service) ModifyExtraRelations(ctx *state.Context, objectId string, modifier func(current []*model.Relation) ([]*model.Relation, error)) (err error) {
	if modifier == nil {
		return fmt.Errorf("modifier is nil")
	}
	return s.Do(objectId, func(b smartblock.SmartBlock) error {
		st := b.NewStateCtx(ctx)
		rels, err := modifier(st.ExtraRelations())
		if err != nil {
			return err
		}

		return b.UpdateExtraRelations(st.Context(), rels, true)
	})
}

// ModifyDetails performs details get and update under the sb lock to make sure no modifications are done in the middle
func (s *service) ModifyDetails(objectId string, modifier func(current *types.Struct) (*types.Struct, error)) (err error) {
	if modifier == nil {
		return fmt.Errorf("modifier is nil")
	}
	return s.Do(objectId, func(b smartblock.SmartBlock) error {
		dets, err := modifier(b.CombinedDetails())
		if err != nil {
			return err
		}

		return b.Apply(b.NewState().SetDetails(dets))
	})
}

// ModifyLocalDetails modifies local details of the object in cache, and if it is not found, sets pending details in object store
func (s *service) ModifyLocalDetails(objectId string, modifier func(current *types.Struct) (*types.Struct, error)) (err error) {
	if modifier == nil {
		return fmt.Errorf("modifier is nil")
	}
	// we set pending details if object is not in cache
	// we do this under lock to prevent races if the object is created in parallel
	// because in that case we can lose changes
	err = s.cache.DoLockedIfNotExists(objectId, func() error {
		objectDetails, err := s.objectStore.GetPendingLocalDetails(objectId)
		if err != nil && err != ds.ErrNotFound {
			return err
		}
		var details *types.Struct
		if objectDetails != nil {
			details = objectDetails.GetDetails()
		}
		modifiedDetails, err := modifier(details)
		if err != nil {
			return err
		}
		return s.objectStore.UpdatePendingLocalDetails(objectId, modifiedDetails)
	})
	if err != nil && err != ocache.ErrExists {
		return err
	}
	err = s.Do(objectId, func(b smartblock.SmartBlock) error {
		// we just need to invoke the smartblock so it reads from pending details
		// no need to call modify twice
		if err == nil {
			return nil
		}

		dets, err := modifier(b.CombinedDetails())
		if err != nil {
			return err
		}

		return b.Apply(b.NewState().SetDetails(dets))
	})
	// that means that we will apply the change later as soon as the block is loaded by thread queue
	if err == source.ErrObjectNotFound {
		return nil
	}
	return err
}

func (s *service) UpdateExtraRelations(ctx *state.Context, objectId string, relations []*model.Relation, createIfMissing bool) (err error) {
	return s.Do(objectId, func(b smartblock.SmartBlock) error {
		return b.UpdateExtraRelations(ctx, relations, createIfMissing)
	})
}

func (s *service) AddExtraRelations(ctx *state.Context, objectId string, relations []*model.Relation) (relationsWithKeys []*model.Relation, err error) {
	err = s.Do(objectId, func(b smartblock.SmartBlock) error {
		var err2 error
		relationsWithKeys, err2 = b.AddExtraRelations(ctx, relations)
		if err2 != nil {
			return err2
		}
		return nil
	})

	return
}

func (s *service) AddExtraRelationOption(ctx *state.Context, req pb.RpcObjectRelationOptionAddRequest) (opt *model.RelationOption, err error) {
	err = s.Do(req.ContextId, func(b smartblock.SmartBlock) error {
		opt, err = b.AddExtraRelationOption(ctx, req.RelationKey, *req.Option, true)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func (s *service) UpdateExtraRelationOption(ctx *state.Context, req pb.RpcObjectRelationOptionUpdateRequest) error {
	return s.Do(req.ContextId, func(b smartblock.SmartBlock) error {
		err := b.UpdateExtraRelationOption(ctx, req.RelationKey, *req.Option, true)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *service) DeleteExtraRelationOption(ctx *state.Context, req pb.RpcObjectRelationOptionDeleteRequest) error {
	objIds, err := s.anytype.ObjectStore().AggregateObjectIdsForOptionAndRelation(req.RelationKey, req.OptionId)
	if err != nil {
		return err
	}

	if !req.ConfirmRemoveAllValuesInRecords {
		for _, objId := range objIds {
			if objId != req.ContextId {
				return ErrOptionUsedByOtherObjects
			}
		}
	} else {
		for _, objId := range objIds {
			err = s.Do(objId, func(b smartblock.SmartBlock) error {
				err := b.DeleteExtraRelationOption(ctx, req.RelationKey, req.OptionId, true)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil && err != smartblock.ErrRelationOptionNotFound {
				return err
			}
		}
	}
	return nil
}

func (s *service) SetObjectTypes(ctx *state.Context, objectId string, objectTypes []string) (err error) {
	return s.Do(objectId, func(b smartblock.SmartBlock) error {
		return b.SetObjectTypes(ctx, objectTypes)
	})
}

// todo: rewrite with options
// withId may me empty
func (s *service) CreateObjectInWorkspace(ctx context.Context, workspaceId string, withId thread.ID, sbType coresb.SmartBlockType) (csm core.SmartBlock, err error) {
	startTime := time.Now()
	ev, exists := ctx.Value(ObjectCreateEvent).(*metrics.CreateObjectEvent)
	err = s.DoWithContext(ctx, workspaceId, func(b smartblock.SmartBlock) error {
		if exists {
			ev.GetWorkspaceBlockWaitMs = time.Now().Sub(startTime).Milliseconds()
		}
		workspace, ok := b.(*editor.Workspaces)
		if !ok {
			return fmt.Errorf("incorrect object with workspace id")
		}
		csm, err = workspace.CreateObject(withId, sbType)
		if exists {
			ev.WorkspaceCreateMs = time.Now().Sub(startTime).Milliseconds() - ev.GetWorkspaceBlockWaitMs
		}
		if err != nil {
			return fmt.Errorf("anytype.CreateBlock error: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return csm, nil
}

func (s *service) DeleteObjectFromWorkspace(workspaceId string, objectId string) error {
	return s.Do(workspaceId, func(b smartblock.SmartBlock) error {
		workspace, ok := b.(*editor.Workspaces)
		if !ok {
			return fmt.Errorf("incorrect object with workspace id")
		}
		return workspace.DeleteObject(objectId)
	})
}

func (s *service) CreateSet(req pb.RpcObjectCreateSetRequest) (setId string, err error) {
	req.Details = internalflag.AddToDetails(req.Details, req.InternalFlags)

	var dvContent model.BlockContentOfDataview
	var dvSchema schema.Schema
	if len(req.Source) != 0 {
		if dvContent, dvSchema, err = dataview.DataviewBlockBySource(s.anytype.ObjectStore(), req.Source); err != nil {
			return
		}
	}
	workspaceId := s.anytype.PredefinedBlocks().Account

	// TODO: here can be a deadlock if this is somehow created from workspace (as set)
	csm, err := s.CreateObjectInWorkspace(context.TODO(), workspaceId, thread.Undef, coresb.SmartBlockTypeSet)
	if err != nil {
		return "", err
	}

	setId = csm.ID()

	state := state.NewDoc(csm.ID(), nil).NewState()
	if workspaceId != "" {
		state.SetDetailAndBundledRelation(bundle.RelationKeyWorkspaceId, pbtypes.String(workspaceId))
	}

	sb, err := s.newSmartBlock(setId, &smartblock.InitContext{
		State: state,
	})
	if err != nil {
		return "", err
	}
	set, ok := sb.(*editor.Set)
	if !ok {
		return setId, fmt.Errorf("unexpected set block type: %T", sb)
	}

	name := pbtypes.GetString(req.Details, bundle.RelationKeyName.String())
	icon := pbtypes.GetString(req.Details, bundle.RelationKeyIconEmoji.String())

	if name == "" && dvSchema != nil {
		name = dvSchema.Description() + " set"
	}
	if dvSchema != nil {
		err = set.InitDataview(&dvContent, name, icon)
	} else {
		err = set.InitDataview(nil, name, icon)
	}

	return setId, err
}

func (s *service) ObjectToSet(id string, source []string) (newId string, err error) {
	var details *types.Struct
	if err = s.Do(id, func(b smartblock.SmartBlock) error {
		details = pbtypes.CopyStruct(b.Details())

		s := b.NewState()
		if layout, ok := s.Layout(); ok && layout == model.ObjectType_note {
			textBlock, err := s.GetFirstTextBlock()
			if err != nil {
				return err
			}
			if textBlock != nil {
				details.Fields[bundle.RelationKeyName.String()] = pbtypes.String(textBlock.Text.Text)
			}
		}

		return nil
	}); err != nil {
		return
	}

	details.Fields[bundle.RelationKeySetOf.String()] = pbtypes.StringList(source)
	newId, err = s.CreateSet(pb.RpcObjectCreateSetRequest{
		Source:  source,
		Details: details,
	})
	if err != nil {
		return
	}

	oStore := s.app.MustComponent(objectstore.CName).(objectstore.ObjectStore)
	res, err := oStore.GetWithLinksInfoByID(id)
	if err != nil {
		return
	}
	for _, il := range res.Links.Inbound {
		if err = s.replaceLink(il.Id, id, newId); err != nil {
			return
		}
	}
	err = s.DeleteObject(id)
	if err != nil {
		// intentionally do not return error here
		log.Errorf("failed to delete object after conversion to set: %s", err.Error())
	}

	return
}

func (s *service) RemoveExtraRelations(ctx *state.Context, objectTypeId string, relationKeys []string) (err error) {
	return s.Do(objectTypeId, func(b smartblock.SmartBlock) error {
		return b.RemoveExtraRelations(ctx, relationKeys)
	})
}

func (s *service) ListAvailableRelations(objectId string) (aggregatedRelations []*model.Relation, err error) {
	err = s.Do(objectId, func(b smartblock.SmartBlock) error {
		objType := b.ObjectType()
		aggregatedRelations = b.Relations()

		agRels, err := s.Anytype().ObjectStore().ListRelations(objType)
		if err != nil {
			return err
		}

		for _, rel := range agRels {
			if pbtypes.HasRelation(aggregatedRelations, rel.Key) {
				continue
			}
			aggregatedRelations = append(aggregatedRelations, pbtypes.CopyRelation(rel))
		}
		return nil
	})

	return
}

func (s *service) ListConvertToObjects(ctx *state.Context, req pb.RpcBlockListConvertToObjectsRequest) (linkIds []string, err error) {
	err = s.DoBasic(req.ContextId, func(b basic.Basic) error {
		linkIds, err = b.ExtractBlocksToObjects(ctx, s, req)
		return err
	})
	return
}

func (s *service) MoveBlocksToNewPage(ctx *state.Context, req pb.RpcBlockListMoveToNewObjectRequest) (linkId string, err error) {
	// 1. Create new page, link
	linkId, pageId, err := s.CreateLinkToTheNewObject(ctx, "", pb.RpcBlockLinkCreateWithObjectRequest{
		ContextId: req.ContextId,
		TargetId:  req.DropTargetId,
		Position:  req.Position,
		Details:   req.Details,
	})

	if err != nil {
		return linkId, err
	}

	// 2. Move blocks to new page
	err = s.MoveBlocks(nil, pb.RpcBlockListMoveToExistingObjectRequest{
		ContextId:       req.ContextId,
		BlockIds:        req.BlockIds,
		TargetContextId: pageId,
		DropTargetId:    "",
		Position:        0,
	})

	if err != nil {
		return linkId, err
	}

	return linkId, err
}

func (s *service) MoveBlocks(ctx *state.Context, req pb.RpcBlockListMoveToExistingObjectRequest) error {
	if req.ContextId == req.TargetContextId {
		return s.DoBasic(req.ContextId, func(b basic.Basic) error {
			return b.Move(ctx, req)
		})
	}
	return s.Do(req.ContextId, func(cb smartblock.SmartBlock) error {
		return s.Do(req.TargetContextId, func(tb smartblock.SmartBlock) error {
			cs := cb.NewState()
			blocks := basic.CutBlocks(cs, req.BlockIds)

			ts := tb.NewState()
			err := basic.PasteBlocks(ts, blocks)
			if err != nil {
				return fmt.Errorf("paste blocks: %w", err)
			}

			err = tb.Apply(ts)
			if err != nil {
				return fmt.Errorf("apply target block state: %w", err)
			}
			return cb.Apply(cs)
		})
	})
}

func (s *service) CreateTableBlock(ctx *state.Context, req pb.RpcBlockTableCreateRequest) (id string, err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		id, err = t.TableCreate(st, req)
		return err
	})
	return
}

func (s *service) TableRowCreate(ctx *state.Context, req pb.RpcBlockTableRowCreateRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.RowCreate(st, req)
	})
	return
}

func (s *service) TableColumnCreate(ctx *state.Context, req pb.RpcBlockTableColumnCreateRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.ColumnCreate(st, req)
	})
	return
}

func (s *service) TableRowDelete(ctx *state.Context, req pb.RpcBlockTableRowDeleteRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.RowDelete(st, req)
	})
	return
}

func (s *service) TableColumnDelete(ctx *state.Context, req pb.RpcBlockTableColumnDeleteRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.ColumnDelete(st, req)
	})
	return
}

func (s *service) TableColumnMove(ctx *state.Context, req pb.RpcBlockTableColumnMoveRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.ColumnMove(st, req)
	})
	return
}

func (s *service) TableRowDuplicate(ctx *state.Context, req pb.RpcBlockTableRowDuplicateRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.RowDuplicate(st, req)
	})
	return
}

func (s *service) TableColumnDuplicate(ctx *state.Context, req pb.RpcBlockTableColumnDuplicateRequest) (id string, err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		id, err = t.ColumnDuplicate(st, req)
		return err
	})
	return id, err
}

func (s *service) TableExpand(ctx *state.Context, req pb.RpcBlockTableExpandRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.Expand(st, req)
	})
	return err
}

func (s *service) TableRowListFill(ctx *state.Context, req pb.RpcBlockTableRowListFillRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.RowListFill(st, req)
	})
	return err
}

func (s *service) TableRowListClean(ctx *state.Context, req pb.RpcBlockTableRowListCleanRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.RowListClean(st, req)
	})
	return err
}

func (s *service) TableRowSetHeader(ctx *state.Context, req pb.RpcBlockTableRowSetHeaderRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.RowSetHeader(st, req)
	})
	return err
}

func (s *service) TableSort(ctx *state.Context, req pb.RpcBlockTableSortRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.Sort(st, req)
	})
	return err
}

func (s *service) TableColumnListFill(ctx *state.Context, req pb.RpcBlockTableColumnListFillRequest) (err error) {
	err = s.DoTable(req.ContextId, ctx, func(st *state.State, t table.Editor) error {
		return t.ColumnListFill(st, req)
	})
	return err
}
