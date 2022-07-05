package bundle

import (
	"fmt"
	"strings"

	coresb "github.com/anytypeio/go-anytype-middleware/pkg/lib/core/smartblock"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/localstore/addr"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	types2 "github.com/gogo/protobuf/types"

	"github.com/anytypeio/go-anytype-middleware/util/pbtypes"
)

// RequiredInternalRelations contains internal relations will be added to any new object type.
// Missing ones will be added to object on opening or during reindex
var RequiredInternalRelations = []RelationKey{
	RelationKeyId,
	RelationKeyName,
	RelationKeyDescription,
	RelationKeySnippet,
	RelationKeyIconEmoji,
	RelationKeyIconImage,
	RelationKeyType,
	RelationKeyLayout,
	RelationKeyLayoutAlign,
	RelationKeyCoverId,
	RelationKeyCoverScale,
	RelationKeyCoverType,
	RelationKeyCoverX,
	RelationKeyCoverY,
	RelationKeyCreatedDate,
	RelationKeyCreator,
	RelationKeyLastModifiedDate,
	RelationKeyLastModifiedBy,
	RelationKeyLastOpenedDate,
	RelationKeyFeaturedRelations,
	RelationKeyIsHidden,
	RelationKeyIsArchived,
	RelationKeyIsFavorite,
	RelationKeyWorkspaceId,
	RelationKeyLinks,
	RelationKeyInternalFlags,
}

var FormatFilePossibleTargetObjectTypes = []string{
	TypeKeyFile.URL(),
	TypeKeyImage.URL(),
	TypeKeyVideo.URL(),
	TypeKeyAudio.URL()}

var DefaultObjectTypePerSmartblockType = map[coresb.SmartBlockType]TypeKey{
	coresb.SmartBlockTypePage:        TypeKeyPage,
	coresb.SmartBlockTypeProfilePage: TypeKeyPage,
	coresb.SmartBlockTypeSet:         TypeKeySet,
	coresb.SmartBlockTypeObjectType:  TypeKeyObjectType,
	coresb.SmartBlockTypeHome:        TypeKeyDashboard,
	coresb.SmartBlockTypeTemplate:    TypeKeyTemplate,
}

// filled in init
var LocalRelationsKeys []string   // stored only in localstore
var DerivedRelationsKeys []string // derived

var ErrNotFound = fmt.Errorf("not found")

func init() {
	for _, r := range relations {
		if r.DataSource == model.Relation_account || r.DataSource == model.Relation_local {
			LocalRelationsKeys = append(LocalRelationsKeys, r.Key)
		} else if r.DataSource == model.Relation_derived {
			DerivedRelationsKeys = append(DerivedRelationsKeys, r.Key)
		}
	}
}

func GetTypeByUrl(u string) (*model.ObjectType, error) {
	if !strings.HasPrefix(u, TypePrefix) {
		return nil, fmt.Errorf("invalid url with no bundled type prefix")
	}
	tk := TypeKey(strings.TrimPrefix(u, TypePrefix))
	if v, exists := types[tk]; exists {
		return pbtypes.CopyObjectType(v), nil
	}

	return nil, ErrNotFound
}

// MustGetType returns built-in object type by predefined TypeKey constant
// PANICS IN CASE RELATION KEY IS NOT EXISTS – DO NOT USE WITH ARBITRARY STRING
func MustGetType(tk TypeKey) *model.ObjectType {
	if v, exists := types[tk]; exists {
		return pbtypes.CopyObjectType(v)
	}

	// we can safely panic in case TypeKey is a generated constant
	panic(ErrNotFound)
}

// MustGetRelation returns built-in relation by predefined RelationKey constant
// PANICS IN CASE RELATION KEY IS NOT EXISTS – DO NOT USE WITH ARBITRARY STRING
func MustGetRelation(rk RelationKey) *model.Relation {
	if v, exists := relations[rk]; exists {
		return pbtypes.CopyRelation(v)
	}

	// we can safely panic in case RelationKey is a generated constant
	panic(ErrNotFound)
}

func MustGetRelations(rks []RelationKey) []*model.Relation {
	rels := make([]*model.Relation, 0, len(rks))
	for _, rk := range rks {
		rels = append(rels, MustGetRelation(rk))
	}
	return rels
}

func GetRelation(rk RelationKey) (*model.Relation, error) {
	if v, exists := relations[rk]; exists {
		return pbtypes.CopyRelation(v), nil
	}

	return nil, ErrNotFound
}

// MustGetLayout returns built-in layout by predefined Layout constant
// PANICS IN CASE RELATION KEY IS NOT EXISTS – DO NOT USE WITH ARBITRARY STRING
func MustGetLayout(lk model.ObjectTypeLayout) *model.Layout {
	if v, exists := Layouts[lk]; exists {
		return pbtypes.CopyLayout(&v)
	}

	// we can safely panic in case RelationKey is a generated constant
	panic(ErrNotFound)
}

func GetLayout(lk model.ObjectTypeLayout) (*model.Layout, error) {
	if v, exists := Layouts[lk]; exists {
		return pbtypes.CopyLayout(&v), nil
	}

	return nil, ErrNotFound
}

func ListRelations() []*model.Relation {
	var rels []*model.Relation
	for _, rel := range relations {
		rels = append(rels, pbtypes.CopyRelation(rel))
	}

	return rels
}

func ListRelationsKeys() []RelationKey {
	var keys []RelationKey
	for k, _ := range relations {
		keys = append(keys, k)
	}

	return keys
}

func ListRelationsUrls() []string {
	var keys []string
	for k, _ := range relations {
		keys = append(keys, addr.BundledRelationURLPrefix+k.String())
	}

	return keys
}

func HasRelation(key string) bool {
	_, exists := relations[RelationKey(key)]

	return exists
}

func HasObjectType(key string) bool {
	_, exists := types[TypeKey(key)]

	return exists
}

func EqualWithRelation(key string, rel *model.Relation) (equal bool, exists bool) {
	v, exists := relations[RelationKey(key)]
	if !exists {
		return false, false
	}

	return pbtypes.RelationEqualOmitDictionary(v, rel), true
}

func ListTypes() ([]*model.ObjectType, error) {
	var otypes []*model.ObjectType
	for _, ot := range types {
		otypes = append(otypes, ot)
	}

	return otypes, nil
}

func ListTypesKeys() []TypeKey {
	var keys []TypeKey
	for k, _ := range types {
		keys = append(keys, k)
	}

	return keys
}

func GetDetailsForRelation(bundled bool, rel *model.Relation) ([]*model.Relation, *types2.Struct) {
	var prefix string
	if bundled {
		prefix = addr.BundledRelationURLPrefix
	} else {
		prefix = addr.CustomRelationURLPrefix
	}

	d := &types2.Struct{Fields: map[string]*types2.Value{
		RelationKeyName.String():             pbtypes.String(rel.Name),
		RelationKeyDescription.String():      pbtypes.String(rel.Description),
		RelationKeyId.String():               pbtypes.String(prefix + rel.Key),
		RelationKeyType.String():             pbtypes.String(TypeKeyRelation.URL()),
		RelationKeyCreator.String():          pbtypes.String(rel.Creator),
		RelationKeyLayout.String():           pbtypes.Float64(float64(model.ObjectType_relation)),
		RelationKeyRelationFormat.String():   pbtypes.Float64(float64(rel.Format)),
		RelationKeyIsHidden.String():         pbtypes.Bool(rel.Hidden),
		RelationKeyIsReadonly.String():       pbtypes.Bool(rel.ReadOnlyRelation),
		RelationKeyMpAddedToLibrary.String(): pbtypes.Bool(true), // temp
	}}

	var rels []*model.Relation
	for k := range d.Fields {
		rels = append(rels, MustGetRelation(RelationKey(k)))
	}
	return rels, d
}

func MergeRelationsKeys(rels1 []RelationKey, rels2 []RelationKey) []RelationKey {
	if rels1 == nil {
		return rels2
	}
	if rels2 == nil {
		return rels1
	}

	rels := make([]RelationKey, 0, len(rels2)+len(rels1))
	for _, rel := range rels2 {
		rels = append(rels, rel)
	}

	for _, rel := range rels1 {
		if HasRelationKey(rels, rel) {
			continue
		}
		rels = append(rels, rel)
	}

	return rels
}

func GetRelationsKeys(rels []*model.Relation) []RelationKey {
	if rels == nil {
		return nil
	}

	relsKeys := make([]RelationKey, 0, len(rels))
	for _, rel := range rels {
		relsKeys = append(relsKeys, RelationKey(rel.Key))
	}

	return relsKeys
}

func HasRelationKey(rels []RelationKey, rel RelationKey) bool {
	for _, rel1 := range rels {
		if rel1 == rel {
			return true
		}
	}

	return false
}
