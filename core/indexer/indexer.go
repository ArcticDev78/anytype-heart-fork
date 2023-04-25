package indexer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/anytypeio/any-sync/app"
	"github.com/gogo/protobuf/types"
	ds "github.com/ipfs/go-datastore"
	"github.com/textileio/go-threads/core/thread"
	"golang.org/x/exp/slices"

	"github.com/anytypeio/go-anytype-middleware/core/anytype/config"
	"github.com/anytypeio/go-anytype-middleware/core/block/doc"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/state"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/template"
	"github.com/anytypeio/go-anytype-middleware/core/block/source"
	"github.com/anytypeio/go-anytype-middleware/core/relation"
	"github.com/anytypeio/go-anytype-middleware/metrics"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/bundle"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core/smartblock"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/database"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/localstore/addr"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/localstore/ftsearch"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/localstore/objectstore"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/logging"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/space/typeprovider"
	"github.com/anytypeio/go-anytype-middleware/util/pbtypes"
	"github.com/anytypeio/go-anytype-middleware/util/slice"
)

const (
	CName = "indexer"

	// ### Increasing counters below will trigger existing account to reindex their

	// ForceThreadsObjectsReindexCounter reindex thread-based objects
	ForceThreadsObjectsReindexCounter int32 = 8
	// ForceFilesReindexCounter reindex ipfs-file-based objects
	ForceFilesReindexCounter int32 = 9 //
	// ForceBundledObjectsReindexCounter reindex objects like anytypeProfile
	ForceBundledObjectsReindexCounter int32 = 5 // reindex objects like anytypeProfile
	// ForceIdxRebuildCounter erases localstore indexes and reindex all type of objects
	// (no need to increase ForceThreadsObjectsReindexCounter & ForceFilesReindexCounter)
	ForceIdxRebuildCounter int32 = 37
	// ForceFulltextIndexCounter  performs fulltext indexing for all type of objects (useful when we change fulltext config)
	ForceFulltextIndexCounter int32 = 4
	// ForceFilestoreKeysReindexCounter reindex filestore keys in all objects
	ForceFilestoreKeysReindexCounter int32 = 2
)

var log = logging.Logger("anytype-doc-indexer")

var (
	ftIndexInterval         = time.Minute
	ftIndexForceMinInterval = time.Second * 10
)

func New() Indexer {
	return &indexer{}
}

type Indexer interface {
	ForceFTIndex()
	app.ComponentRunnable
}

type ThreadLister interface {
	Threads() (thread.IDSlice, error)
}

type Hasher interface {
	Hash() string
}

type reindexFlags uint64

const cacheTimeout = 4 * time.Second

const (
	reindexBundledTypes reindexFlags = 1 << iota
	removeAllIndexedObjects
	reindexBundledRelations
	eraseIndexes
	reindexThreadObjects
	reindexFileObjects
	reindexFulltext
	reindexBundledTemplates
	reindexBundledObjects
	reindexFileKeys
)

type indexer struct {
	store objectstore.ObjectStore
	// todo: move logstore to separate component?
	anytype         core.Service
	source          source.Service
	relationService relation.Service

	doc         doc.Service
	quit        chan struct{}
	mu          sync.Mutex
	btHash      Hasher
	archivedMap map[string]struct{}
	favoriteMap map[string]struct{}
	newAccount  bool
	forceFt     chan struct{}

	relationBulkMigration relation.BulkMigration
	relationMigratorMu    sync.Mutex
	typeProvider          typeprovider.ObjectTypeProvider
}

func (i *indexer) Init(a *app.App) (err error) {
	i.newAccount = a.MustComponent(config.CName).(*config.Config).NewAccount
	i.anytype = a.MustComponent(core.CName).(core.Service)
	i.store = a.MustComponent(objectstore.CName).(objectstore.ObjectStore)
	i.relationService = a.MustComponent(relation.CName).(relation.Service)
	i.typeProvider = a.MustComponent(typeprovider.CName).(typeprovider.ObjectTypeProvider)
	i.source = a.MustComponent(source.CName).(source.Service)
	i.btHash = a.MustComponent("builtintemplate").(Hasher)
	i.doc = a.MustComponent(doc.CName).(doc.Service)
	i.quit = make(chan struct{})
	i.archivedMap = make(map[string]struct{}, 100)
	i.favoriteMap = make(map[string]struct{}, 100)
	i.forceFt = make(chan struct{})
	return
}

func (i *indexer) Name() (name string) {
	return CName
}

func (i *indexer) saveLatestChecksums() error {
	// todo: add layout indexing when needed
	checksums := model.ObjectStoreChecksums{
		BundledObjectTypes:         bundle.TypeChecksum,
		BundledRelations:           bundle.RelationChecksum,
		BundledTemplates:           i.btHash.Hash(),
		ObjectsForceReindexCounter: ForceThreadsObjectsReindexCounter,
		FilesForceReindexCounter:   ForceFilesReindexCounter,

		IdxRebuildCounter:                ForceIdxRebuildCounter,
		FulltextRebuild:                  ForceFulltextIndexCounter,
		BundledObjects:                   ForceBundledObjectsReindexCounter,
		FilestoreKeysForceReindexCounter: ForceFilestoreKeysReindexCounter,
	}
	return i.store.SaveChecksums(&checksums)
}

func (i *indexer) saveLatestCounters() error {
	// todo: add layout indexing when needed
	checksums := model.ObjectStoreChecksums{
		BundledObjectTypes:               bundle.TypeChecksum,
		BundledRelations:                 bundle.RelationChecksum,
		BundledTemplates:                 i.btHash.Hash(),
		ObjectsForceReindexCounter:       ForceThreadsObjectsReindexCounter,
		FilesForceReindexCounter:         ForceFilesReindexCounter,
		IdxRebuildCounter:                ForceIdxRebuildCounter,
		FulltextRebuild:                  ForceFulltextIndexCounter,
		BundledObjects:                   ForceBundledObjectsReindexCounter,
		FilestoreKeysForceReindexCounter: ForceFilestoreKeysReindexCounter,
	}
	return i.store.SaveChecksums(&checksums)
}

func (i *indexer) Run(context.Context) (err error) {
	if ftErr := i.ftInit(); ftErr != nil {
		log.Errorf("can't init ft: %v", ftErr)
	}
	err = i.reindexIfNeeded()
	if err != nil {
		return err
	}
	i.migrateRemoveNonindexableObjects()
	go i.ftLoop()
	return
}

func (i *indexer) ForceFTIndex() {
	select {
	case i.forceFt <- struct{}{}:
	default:
	}
}

func (i *indexer) migrateRemoveNonindexableObjects() {
	ids, err := i.getIdsForTypes(
		smartblock.SmartblockTypeMarketplaceType, smartblock.SmartblockTypeMarketplaceRelation,
		smartblock.SmartblockTypeMarketplaceTemplate, smartblock.SmartBlockTypeDate, smartblock.SmartBlockTypeBreadcrumbs,
	)
	if err != nil {
		log.Errorf("migrateRemoveNonindexableObjects: failed to get ids: %s", err.Error())
	}

	for _, id := range ids {
		err = i.store.DeleteDetails(id)
		if err != nil {
			log.Errorf("migrateRemoveNonindexableObjects: failed to get ids: %s", err.Error())
		}
	}
}

func (i *indexer) reindexIfNeeded() error {
	var (
		err       error
		checksums *model.ObjectStoreChecksums
		reindex   reindexFlags
	)
	checksums, err = i.store.GetChecksums()
	if err != nil && err != ds.ErrNotFound {
		return err
	}
	if checksums == nil {
		checksums = &model.ObjectStoreChecksums{
			// do no add bundled relations checksums, because we want to index them for new accounts
			ObjectsForceReindexCounter:       ForceThreadsObjectsReindexCounter,
			FilesForceReindexCounter:         ForceFilesReindexCounter,
			IdxRebuildCounter:                ForceIdxRebuildCounter,
			FilestoreKeysForceReindexCounter: ForceFilestoreKeysReindexCounter,
		}
	}

	if checksums.BundledRelations != bundle.RelationChecksum {
		reindex = reindex | reindexBundledRelations
	}
	if checksums.BundledObjectTypes != bundle.TypeChecksum {
		reindex = reindex | reindexBundledTypes
	}
	if checksums.ObjectsForceReindexCounter != ForceThreadsObjectsReindexCounter {
		reindex = reindex | reindexThreadObjects
	}
	if checksums.FilestoreKeysForceReindexCounter != ForceFilestoreKeysReindexCounter {
		reindex = reindex | reindexFileKeys
	}
	if checksums.FilesForceReindexCounter != ForceFilesReindexCounter {
		reindex = reindex | reindexFileObjects
	}
	if checksums.FulltextRebuild != ForceFulltextIndexCounter {
		reindex = reindex | reindexFulltext
	}
	if checksums.BundledTemplates != i.btHash.Hash() {
		reindex = reindex | reindexBundledTemplates
	}
	if checksums.BundledObjects != ForceBundledObjectsReindexCounter {
		reindex = reindex | reindexBundledObjects
	}
	if checksums.IdxRebuildCounter != ForceIdxRebuildCounter {
		reindex = math.MaxUint64
	}
	return i.Reindex(context.WithValue(context.TODO(), metrics.CtxKeyRequest, "reindex_forced"), reindex)
}

func (i *indexer) reindexOutdatedThreads() (toReindex, success int, err error) {
	spc, err := i.anytype.SpaceService().AccountSpace(context.Background())
	if err != nil {
		return
	}

	tids := spc.StoredIds()
	var idsToReindex []string
	for _, tid := range tids {
		logErr := func(err error) {
			log.With("tree", tid).Errorf("reindexOutdatedThreads failed to get tree to reindex: %s", err.Error())
		}

		lastHash, err := i.store.GetLastIndexedHeadsHash(tid)
		if err != nil {
			logErr(err)
			continue
		}
		info, err := spc.Storage().TreeStorage(tid)
		if err != nil {
			logErr(err)
			continue
		}
		heads, err := info.Heads()
		if err != nil {
			logErr(err)
			continue
		}

		hh := headsHash(heads)
		if lastHash != hh {
			log.With("tree", tid).Warnf("not equal indexed heads hash: %s!=%s (%d logs)", lastHash, hh, len(heads))
			idsToReindex = append(idsToReindex, tid)
		}
	}
	if len(idsToReindex) > 0 {
		for _, id := range idsToReindex {
			// TODO: we should reindex it I guess at start
			// if i.anytype.PredefinedBlocks().IsAccount(id) {
			//	continue
			// }

			ctx := context.WithValue(context.Background(), metrics.CtxKeyRequest, "reindexOutdatedThreads")
			d, err := i.doc.GetDocInfo(ctx, id)
			if err != nil {
				continue
			}

			err = i.Index(ctx, d)
			if err == nil {
				success++
			} else {
				log.With("thread", id).Errorf("reindexOutdatedThreads failed to index doc: %s", err.Error())
			}
		}
	}
	return len(idsToReindex), success, nil
}

func (i *indexer) getIdsForTypes(sbt ...smartblock.SmartBlockType) ([]string, error) {
	var ids []string
	for _, t := range sbt {
		st, err := i.source.SourceTypeBySbType(t)
		if err != nil {
			return nil, err
		}
		idsT, err := st.ListIds()
		if err != nil {
			return nil, err
		}
		ids = append(ids, idsT...)
	}
	return ids, nil
}

func (i *indexer) Reindex(ctx context.Context, reindex reindexFlags) (err error) {
	if reindex != 0 {
		log.Infof("start store reindex (eraseIndexes=%v, reindexFileObjects=%v, reindexThreadObjects=%v, reindexBundledRelations=%v, reindexBundledTypes=%v, reindexFulltext=%v, reindexBundledTemplates=%v, reindexBundledObjects=%v, reindexFileKeys=%v)", reindex&eraseIndexes != 0, reindex&reindexFileObjects != 0, reindex&reindexThreadObjects != 0, reindex&reindexBundledRelations != 0, reindex&reindexBundledTypes != 0, reindex&reindexFulltext != 0, reindex&reindexBundledTemplates != 0, reindex&reindexBundledObjects != 0, reindex&reindexFileKeys != 0)
	}

	if reindex&reindexFileKeys != 0 {
		err = i.anytype.FileStore().RemoveEmpty()
		if err != nil {
			log.Errorf("reindex failed to RemoveEmpty filekeys: %v", err.Error())
		} else {
			log.Infof("RemoveEmpty filekeys succeed")
		}
	}

	if reindex&removeAllIndexedObjects != 0 {
		ids, err := i.store.ListIds()
		if err != nil {
			log.Errorf("reindex failed to get all ids(removeAllIndexedObjects): %v", err.Error())
		}
		for _, id := range ids {
			err = i.store.DeleteDetails(id)
			if err != nil {
				log.Errorf("reindex failed to delete details(removeAllIndexedObjects): %v", err.Error())
			}
		}

		defer func() {
			i.relationMigratorMu.Lock()
			defer i.relationMigratorMu.Unlock()
			if i.relationBulkMigration == nil {
				return
			}
			err2 := i.relationBulkMigration.Commit()
			i.relationBulkMigration = nil
			if err2 != nil {
				log.Errorf("reindex relation migration error: %s", err2.Error())
			}
		}()
		i.relationMigratorMu.Lock()
		i.relationBulkMigration = i.relationService.CreateBulkMigration()
		i.relationMigratorMu.Unlock()
	}
	var indexesWereRemoved bool
	if reindex&eraseIndexes != 0 {
		err = i.store.EraseIndexes()
		if err != nil {
			log.Errorf("reindex failed to erase indexes: %v", err.Error())
		} else {
			log.Infof("all store indexes succesfully erased")
			// store this flag because underlying localstore needs to now if it needs to amend indexes based on the prev value
			indexesWereRemoved = true
		}
	}

	// We derive or init predefined blocks here in order to ensure consistency of object store.
	// If we call this method before removing objects from store, we will end up with inconsistent state
	// because indexing of predefined objects will not run again
	err = i.anytype.EnsurePredefinedBlocks(ctx)
	if err != nil {
		return err
	}

	if reindex > 0 {
		d, err := i.doc.GetDocInfo(ctx, i.anytype.PredefinedBlocks().Archive)
		if err != nil {
			log.Errorf("reindex failed to open archive: %s", err.Error())
		} else {
			for _, target := range d.Links {
				i.archivedMap[target] = struct{}{}
			}
		}

		d, err = i.doc.GetDocInfo(ctx, i.anytype.PredefinedBlocks().Home)
		if err != nil {
			log.Errorf("reindex failed to open archive: %s", err.Error())
		} else {
			for _, b := range d.Links {
				i.favoriteMap[b] = struct{}{}
			}
		}
	}

	// for all ids except home and archive setting cache timeout for reindexing
	// ctx = context.WithValue(ctx, ocache.CacheTimeout, cacheTimeout)
	if reindex&reindexThreadObjects != 0 {
		ids, err := i.getIdsForTypes(
			smartblock.SmartBlockTypePage,
			smartblock.SmartBlockTypeSet,
			smartblock.SmartBlockTypeObjectType,
			smartblock.SmartBlockTypeProfilePage,
			smartblock.SmartBlockTypeTemplate,
			smartblock.SmartblockTypeMarketplaceType,
			smartblock.SmartblockTypeMarketplaceTemplate,
			smartblock.SmartblockTypeMarketplaceRelation,
			smartblock.SmartBlockTypeArchive,
			smartblock.SmartBlockTypeHome,
			smartblock.SmartBlockTypeWorkspaceOld,
		)
		if err != nil {
			return err
		}
		start := time.Now()
		successfullyReindexed := i.reindexIdsIgnoreErr(ctx, indexesWereRemoved, ids...)
		if metrics.Enabled {
			metrics.SharedClient.RecordEvent(metrics.ReindexEvent{
				ReindexType:    metrics.ReindexTypeThreads,
				Total:          len(ids),
				Success:        successfullyReindexed,
				SpentMs:        int(time.Since(start).Milliseconds()),
				IndexesRemoved: indexesWereRemoved,
			})
		}
		log.Infof("%d/%d objects have been successfully reindexed", successfullyReindexed, len(ids))
	} else {
		go func() {
			start := time.Now()
			total, success, err := i.reindexOutdatedThreads()
			if err != nil {
				log.Infof("failed to reindex outdated objects: %s", err.Error())
			} else {
				log.Infof("%d/%d outdated objects have been successfully reindexed", success, total)
			}
			if metrics.Enabled && total > 0 {
				metrics.SharedClient.RecordEvent(metrics.ReindexEvent{
					ReindexType:    metrics.ReindexTypeOutdatedHeads,
					Total:          total,
					Success:        success,
					SpentMs:        int(time.Since(start).Milliseconds()),
					IndexesRemoved: indexesWereRemoved,
				})
			}
		}()
	}

	if reindex&reindexFileObjects != 0 {
		ids, err := i.getIdsForTypes(smartblock.SmartBlockTypeFile)
		if err != nil {
			return err
		}
		start := time.Now()
		successfullyReindexed := i.reindexIdsIgnoreErr(ctx, indexesWereRemoved, ids...)
		if metrics.Enabled && len(ids) > 0 {
			metrics.SharedClient.RecordEvent(metrics.ReindexEvent{
				ReindexType:    metrics.ReindexTypeFiles,
				Total:          len(ids),
				Success:        successfullyReindexed,
				SpentMs:        int(time.Since(start).Milliseconds()),
				IndexesRemoved: indexesWereRemoved,
			})
		}
		msg := fmt.Sprintf("%d/%d files have been successfully reindexed", successfullyReindexed, len(ids))
		if len(ids)-successfullyReindexed != 0 {
			log.Error(msg)
		} else {
			log.Info(msg)
		}
	}
	if reindex&reindexBundledRelations != 0 {
		ids, err := i.getIdsForTypes(smartblock.SmartBlockTypeBundledRelation)
		if err != nil {
			return err
		}
		start := time.Now()
		successfullyReindexed := i.reindexIdsIgnoreErr(ctx, indexesWereRemoved, ids...)
		if metrics.Enabled && len(ids) > 0 {
			metrics.SharedClient.RecordEvent(metrics.ReindexEvent{
				ReindexType:    metrics.ReindexTypeBundledRelations,
				Total:          len(ids),
				Success:        successfullyReindexed,
				SpentMs:        int(time.Since(start).Milliseconds()),
				IndexesRemoved: indexesWereRemoved,
			})
		}
		msg := fmt.Sprintf("%d/%d bundled relations have been successfully reindexed", successfullyReindexed, len(ids))
		if len(ids)-successfullyReindexed != 0 {
			log.Error(msg)
		} else {
			log.Info(msg)
		}
	}
	if reindex&reindexBundledTypes != 0 {
		// lets add anytypeProfile here, because it's seems too much to create one more counter especially for it
		ids, err := i.getIdsForTypes(smartblock.SmartBlockTypeBundledObjectType, smartblock.SmartBlockTypeAnytypeProfile)
		if err != nil {
			return err
		}
		start := time.Now()
		successfullyReindexed := i.reindexIdsIgnoreErr(ctx, indexesWereRemoved, ids...)
		if metrics.Enabled && len(ids) > 0 {
			metrics.SharedClient.RecordEvent(metrics.ReindexEvent{
				ReindexType:    metrics.ReindexTypeBundledTypes,
				Total:          len(ids),
				Success:        successfullyReindexed,
				SpentMs:        int(time.Since(start).Milliseconds()),
				IndexesRemoved: indexesWereRemoved,
			})
		}
		msg := fmt.Sprintf("%d/%d bundled types have been successfully reindexed", successfullyReindexed, len(ids))
		if len(ids)-successfullyReindexed != 0 {
			log.Error(msg)
		} else {
			log.Info(msg)
		}

		var ots = make([]string, 0, len(bundle.SystemTypes))
		for _, ot := range bundle.SystemTypes {
			ots = append(ots, ot.BundledURL())
		}

		for _, ot := range bundle.InternalTypes {
			ots = append(ots, ot.BundledURL())
		}

		var rels = make([]*model.Relation, 0, len(bundle.RequiredInternalRelations))
		for _, rel := range bundle.SystemRelations {
			rels = append(rels, bundle.MustGetRelation(rel))
		}
		i.migrateObjectTypes(ots)
		i.migrateRelations(rels)
	}
	if reindex&reindexBundledObjects != 0 {
		// hardcoded for now
		ids := []string{addr.AnytypeProfileId}
		start := time.Now()
		successfullyReindexed := i.reindexIdsIgnoreErr(ctx, indexesWereRemoved, ids...)
		if metrics.Enabled && len(ids) > 0 {
			metrics.SharedClient.RecordEvent(metrics.ReindexEvent{
				ReindexType: metrics.ReindexTypeBundledTemplates,
				Total:       len(ids),
				Success:     successfullyReindexed,
				SpentMs:     int(time.Since(start).Milliseconds()),
			})
		}
		msg := fmt.Sprintf("%d/%d bundled objects have been successfully reindexed", successfullyReindexed, len(ids))
		if len(ids)-successfullyReindexed != 0 {
			log.Error(msg)
		} else {
			log.Info(msg)
		}
	}

	if reindex&reindexBundledTemplates != 0 {
		existsRec, _, err := i.store.QueryObjectInfo(database.Query{}, []smartblock.SmartBlockType{smartblock.SmartBlockTypeBundledTemplate})
		if err != nil {
			return err
		}
		existsIds := make([]string, 0, len(existsRec))
		for _, rec := range existsRec {
			existsIds = append(existsIds, rec.Id)
		}
		ids, err := i.getIdsForTypes(smartblock.SmartBlockTypeBundledTemplate)
		if err != nil {
			return err
		}
		var removed int
		for _, eId := range existsIds {
			if slice.FindPos(ids, eId) == -1 {
				removed++
				i.store.DeleteObject(eId)
			}
		}
		successfullyReindexed := i.reindexIdsIgnoreErr(ctx, indexesWereRemoved, ids...)
		msg := fmt.Sprintf("%d/%d bundled templates have been successfully reindexed; removed: %d", successfullyReindexed, len(ids), removed)
		if len(ids)-successfullyReindexed != 0 {
			log.Error(msg)
		} else {
			log.Info(msg)
		}
	}
	if reindex&reindexFulltext != 0 {
		var ids []string
		ids, err := i.getIdsForTypes(smartblock.SmartBlockTypePage, smartblock.SmartBlockTypeFile, smartblock.SmartBlockTypeBundledRelation, smartblock.SmartBlockTypeBundledObjectType, smartblock.SmartBlockTypeAnytypeProfile)
		if err != nil {
			return err
		}

		var addedToQueue int
		for _, id := range ids {
			if err := i.store.AddToIndexQueue(id); err != nil {
				log.Errorf("failed to add to index queue: %v", err)
			} else {
				addedToQueue++
			}
		}
		msg := fmt.Sprintf("%d/%d objects have been successfully added to the fulltext queue", addedToQueue, len(ids))
		if len(ids)-addedToQueue != 0 {
			log.Error(msg)
		} else {
			log.Info(msg)
		}
	}

	return i.saveLatestChecksums()
}

func extractOldRelationsFromState(s *state.State) []*model.Relation {
	var rels []*model.Relation
	if objRels := s.OldExtraRelations(); len(objRels) > 0 {
		rels = append(rels, s.OldExtraRelations()...)
	}

	if dvBlock := s.Pick(template.DataviewBlockId); dvBlock != nil {
		rels = append(rels, dvBlock.Model().GetDataview().GetRelations()...)
	}

	return rels
}

func (i *indexer) migrateRelations(rels []*model.Relation) {
	if len(rels) == 0 {
		return
	}
	i.relationMigratorMu.Lock()
	defer i.relationMigratorMu.Unlock()

	if i.relationBulkMigration != nil {
		i.relationBulkMigration.AddRelations(rels)
	} else {
		err := i.relationService.MigrateRelations(rels)
		if err != nil {
			log.Errorf("migrateRelations got error: %s", err.Error())
		}
	}
}

func (i *indexer) migrateObjectTypes(ots []string) {
	if len(ots) == 0 {
		return
	}

	var typesModels []*model.ObjectType // do not make
	for _, ot := range ots {
		t, err := bundle.GetTypeByUrl(ot)
		if err != nil {
			continue
		}

		typesModels = append(typesModels, t)
	}

	if len(typesModels) == 0 {
		return
	}

	i.relationMigratorMu.Lock()
	defer i.relationMigratorMu.Unlock()

	if i.relationBulkMigration != nil {
		i.relationBulkMigration.AddObjectTypes(typesModels)
	} else {
		err := i.relationService.MigrateObjectTypes(typesModels)
		if err != nil {
			log.Errorf("migrateObjectTypes got error: %s", err.Error())
		}
	}
}

func (i *indexer) reindexDoc(ctx context.Context, id string, indexesWereRemoved bool) error {
	t, err := i.typeProvider.Type(id)
	if err != nil {
		return fmt.Errorf("incorrect sb type: %v", err)
	}

	d, err := i.doc.GetDocInfo(ctx, id)
	if err != nil {
		log.Errorf("reindexDoc failed to open %s: %s", id, err.Error())
		return fmt.Errorf("failed to open doc: %s", err.Error())
	}

	indexDetails, indexLinks := t.Indexable()
	if indexLinks {
		if err := i.store.UpdateObjectLinks(d.Id, d.Links); err != nil {
			log.With("thread", d.Id).Errorf("failed to save object links: %v", err)
		}
	}

	if !indexDetails {
		i.store.DeleteDetails(d.Id)
		return nil
	}

	details := d.State.CombinedDetails()
	_, isArchived := i.archivedMap[id]
	_, isFavorite := i.favoriteMap[id]

	details.Fields[bundle.RelationKeyIsArchived.String()] = pbtypes.Bool(isArchived)
	details.Fields[bundle.RelationKeyIsFavorite.String()] = pbtypes.Bool(isFavorite)
	details.Fields[bundle.RelationKeyLinks.String()] = pbtypes.StringList(d.Links)

	var curDetails *types.Struct
	curDetailsO, _ := i.store.GetDetails(id)
	if curDetailsO.GetDetails().GetFields() != nil {
		curDetails = curDetailsO.Details
	}
	// compare only real object scoped details
	detailsObjectScope := pbtypes.StructCutKeys(details, bundle.LocalRelationsKeys)
	curDetailsObjectScope := pbtypes.StructCutKeys(curDetails, bundle.LocalRelationsKeys)
	if indexesWereRemoved || curDetailsObjectScope == nil || !detailsObjectScope.Equal(curDetailsObjectScope) {
		if indexesWereRemoved || curDetails.GetFields() == nil {
			if err := i.store.CreateObject(id, details, d.Links, pbtypes.GetString(details, bundle.RelationKeyDescription.String())); err != nil {
				return fmt.Errorf("can't create object in the store: %v", err)
			}
		} else {
			if err := i.store.UpdateObjectDetails(id, details, true); err != nil {
				return fmt.Errorf("can't update object in the store: %v", err)
			}
		}
		if headsHash := headsHash(d.Heads); headsHash != "" {
			err = i.store.SaveLastIndexedHeadsHash(id, headsHash)
			if err != nil {
				log.With("thread", id).Errorf("failed to save indexed heads hash: %v", err)
			}
		}

		var skipFulltext bool
		if i.store.FTSearch() != nil {
			// skip fulltext if we already has the object indexed
			if exists, _ := i.store.FTSearch().Has(id); exists {
				skipFulltext = true
			}
		}

		if !skipFulltext {
			if err = i.store.AddToIndexQueue(id); err != nil {
				log.With("thread", id).Errorf("can't add to index: %v", err)
			}
		}
	}
	return nil
}

func (i *indexer) reindexIdsIgnoreErr(ctx context.Context, indexRemoved bool, ids ...string) (successfullyReindexed int) {
	for _, id := range ids {
		err := i.reindexDoc(ctx, id, indexRemoved)
		if err != nil {
			log.With("thread", id).Errorf("failed to reindex: %v", err)
		} else {
			successfullyReindexed++
		}
	}
	return
}

func (i *indexer) Index(ctx context.Context, info doc.DocInfo) error {
	startTime := time.Now()
	sbType, err := i.typeProvider.Type(info.Id)
	if err != nil {
		sbType = smartblock.SmartBlockTypePage
	}
	saveIndexedHash := func() {
		if headsHash := headsHash(info.Heads); headsHash != "" {
			err = i.store.SaveLastIndexedHeadsHash(info.Id, headsHash)
			if err != nil {
				log.With("thread", info.Id).Errorf("failed to save indexed heads hash: %v", err)
			}
		}
	}

	indexDetails, indexLinks := sbType.Indexable()
	if sbType != smartblock.SmartBlockTypeSubObject && sbType != smartblock.SmartBlockTypeWorkspace && sbType != smartblock.SmartBlockTypeBreadcrumbs {
		// avoid recursions
		log.With("migratedtype", sbType).Warn("migrating types")
		if pbtypes.GetString(info.State.CombinedDetails(), bundle.RelationKeyCreator.String()) != addr.AnytypeProfileId {
			i.migrateRelations(extractOldRelationsFromState(info.State))
			i.migrateObjectTypes(info.State.ObjectTypesToMigrate())
		}
	}
	if !indexDetails && !indexLinks {
		saveIndexedHash()
		return nil
	}

	details := info.State.CombinedDetails()
	details.Fields[bundle.RelationKeyLinks.String()] = pbtypes.StringList(info.Links)
	setCreator := pbtypes.GetString(info.State.LocalDetails(), bundle.RelationKeyCreator.String())
	if setCreator == "" {
		setCreator = i.anytype.ProfileID()
	}
	indexSetTime := time.Now()
	var hasError bool
	if indexLinks {
		if err = i.store.UpdateObjectLinks(info.Id, info.Links); err != nil {
			hasError = true
			log.With("thread", info.Id).Errorf("failed to save object links: %v", err)
		}
	}

	indexLinksTime := time.Now()
	if indexDetails {
		if err := i.store.UpdateObjectDetails(info.Id, details, false); err != nil {
			hasError = true
			log.With("thread", info.Id).Errorf("can't update object store: %v", err)
		}
		if err := i.store.AddToIndexQueue(info.Id); err != nil {
			log.With("thread", info.Id).Errorf("can't add id to index queue: %v", err)
		} else {
			log.With("thread", info.Id).Debugf("to index queue")
		}
	} else {
		_ = i.store.DeleteDetails(info.Id)
	}
	indexDetailsTime := time.Now()
	detailsCount := 0
	if details.GetFields() != nil {
		detailsCount = len(details.GetFields())
	}

	if !hasError {
		saveIndexedHash()
	}

	metrics.SharedClient.RecordEvent(metrics.IndexEvent{
		ObjectId:                info.Id,
		IndexLinksTimeMs:        indexLinksTime.Sub(indexSetTime).Milliseconds(),
		IndexDetailsTimeMs:      indexDetailsTime.Sub(indexLinksTime).Milliseconds(),
		IndexSetRelationsTimeMs: indexSetTime.Sub(startTime).Milliseconds(),
		RelationsCount:          len(info.State.PickRelationLinks()),
		DetailsCount:            detailsCount,
	})

	return nil
}

func (i *indexer) ftLoop() {
	ticker := time.NewTicker(ftIndexInterval)
	i.ftIndex()
	var lastForceIndex time.Time
	i.mu.Lock()
	quit := i.quit
	i.mu.Unlock()
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			i.ftIndex()
		case <-i.forceFt:
			if time.Since(lastForceIndex) > ftIndexForceMinInterval {
				i.ftIndex()
				lastForceIndex = time.Now()
			}
		}
	}
}

func (i *indexer) ftIndex() {
	if err := i.store.IndexForEach(i.ftIndexDoc); err != nil {
		log.Errorf("store.IndexForEach error: %v", err)
	}
}

func (i *indexer) ftIndexDoc(id string, _ time.Time) (err error) {
	st := time.Now()
	// ctx := context.WithValue(context.Background(), ocache.CacheTimeout, cacheTimeout)
	ctx := context.WithValue(context.Background(), metrics.CtxKeyRequest, "index_fulltext")

	info, err := i.doc.GetDocInfo(ctx, id)
	if err != nil {
		return
	}

	sbType, err := i.typeProvider.Type(info.Id)
	if err != nil {
		sbType = smartblock.SmartBlockTypePage
	}
	indexDetails, _ := sbType.Indexable()
	if !indexDetails {
		return nil
	}

	if err = i.store.UpdateObjectSnippet(id, info.State.Snippet()); err != nil {
		return
	}

	if len(info.FileHashes) > 0 {
		// todo: move file indexing to the main indexer as we have  the full state there now
		existingIDs, err := i.store.HasIDs(info.FileHashes...)
		if err != nil {
			log.Errorf("failed to get existing file ids : %s", err.Error())
		}
		newIds := slice.Difference(info.FileHashes, existingIDs)
		for _, hash := range newIds {
			// file's hash is id
			err = i.reindexDoc(ctx, hash, false)
			if err != nil {
				log.With("id", hash).Errorf("failed to reindex file: %s", err.Error())
			}

			err = i.store.AddToIndexQueue(hash)
			if err != nil {
				log.With("id", hash).Error(err.Error())
			}
		}
	}

	if fts := i.store.FTSearch(); fts != nil {
		title := pbtypes.GetString(info.State.Details(), bundle.RelationKeyName.String())
		if info.State.ObjectType() == bundle.TypeKeyNote.String() || title == "" {
			title = info.State.Snippet()
		}
		ftDoc := ftsearch.SearchDoc{
			Id:    id,
			Title: title,
			Text:  info.State.SearchText(),
		}
		if err := fts.Index(ftDoc); err != nil {
			log.Errorf("can't ft index doc: %v", err)
		}
		log.Debugf("ft search indexed with title: '%s'", ftDoc.Title)
	}

	log.With("thread", id).Infof("ft index updated for a %v", time.Since(st))
	return
}

func (i *indexer) ftInit() error {
	if ft := i.store.FTSearch(); ft != nil {
		docCount, err := ft.DocCount()
		if err != nil {
			return err
		}
		if docCount == 0 {
			ids, err := i.store.ListIds()
			if err != nil {
				return err
			}
			for _, id := range ids {
				if err := i.store.AddToIndexQueue(id); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (i *indexer) Close(ctx context.Context) (err error) {
	i.mu.Lock()
	quit := i.quit
	i.mu.Unlock()
	if quit != nil {
		close(quit)
		i.mu.Lock()
		i.quit = nil
		i.mu.Unlock()
	}
	return nil
}

func headsHash(heads []string) string {
	if len(heads) == 0 {
		return ""
	}
	slices.Sort(heads)

	sum := sha256.Sum256([]byte(strings.Join(heads, ",")))
	return fmt.Sprintf("%x", sum)
}
