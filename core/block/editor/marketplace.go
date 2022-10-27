package editor

import (
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/smartblock"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/template"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/bundle"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/localstore/addr"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/util/pbtypes"
)

const (
	viewIdLibrary     = "library"
	viewIdMarketplace = "marketplace"
)

func NewMarketplaceType() *MarketplaceType {
	return &MarketplaceType{Set: NewSet()}
}

type MarketplaceType struct {
	*Set
}

func (p *MarketplaceType) Init(ctx *smartblock.InitContext) (err error) {
	err = p.SmartBlock.Init(ctx)
	if err != nil {
		return err
	}

	ot := bundle.TypeKeyObjectType.URL()
	templates := []template.StateTransformer{
		template.WithObjectTypesAndLayout([]string{bundle.TypeKeySet.URL()}),
		template.WithTitle,
		template.WithDefaultFeaturedRelations,
		template.WithDescription,
		template.WithFeaturedRelations,
		template.WithForcedDetail(bundle.RelationKeySetOf, pbtypes.StringList([]string{ot}))}
	dataview := model.BlockContentOfDataview{
		Dataview: &model.BlockContentDataview{
			Source:        []string{ot},
			RelationLinks: bundle.MustGetType(bundle.TypeKeyObjectType).RelationLinks,
			Views: []*model.BlockContentDataviewView{
				{
					Id:    viewIdMarketplace,
					Type:  model.BlockContentDataviewView_Gallery,
					Name:  "Marketplace",
					Sorts: []*model.BlockContentDataviewSort{{RelationKey: bundle.RelationKeyName.String(), Type: model.BlockContentDataviewSort_Asc}},
					Relations: []*model.BlockContentDataviewRelation{
						{Key: bundle.RelationKeyId.String(), IsVisible: false},
						{Key: bundle.RelationKeyName.String(), IsVisible: true},
						{Key: bundle.RelationKeyDescription.String(), IsVisible: true},
						{Key: bundle.RelationKeyIconImage.String(), IsVisible: true},
						{Key: bundle.RelationKeyCreator.String(), IsVisible: true}},
					Filters: []*model.BlockContentDataviewFilter{{
						RelationKey: bundle.RelationKeyIsHidden.String(),
						Condition:   model.BlockContentDataviewFilter_NotEqual,
						Value:       pbtypes.Bool(true),
					}},
				},
				{
					Id:    viewIdLibrary,
					Type:  model.BlockContentDataviewView_Gallery,
					Name:  "Library",
					Sorts: []*model.BlockContentDataviewSort{{RelationKey: bundle.RelationKeyName.String(), Type: model.BlockContentDataviewSort_Asc}},
					Relations: []*model.BlockContentDataviewRelation{
						{Key: bundle.RelationKeyId.String(), IsVisible: false},
						{Key: bundle.RelationKeyName.String(), IsVisible: true},
						{Key: bundle.RelationKeyDescription.String(), IsVisible: true},
						{Key: bundle.RelationKeyIconImage.String(), IsVisible: true},
						{Key: bundle.RelationKeyCreator.String(), IsVisible: true}},
					Filters: []*model.BlockContentDataviewFilter{{
						RelationKey: bundle.RelationKeyIsHidden.String(),
						Condition:   model.BlockContentDataviewFilter_NotEqual,
						Value:       pbtypes.Bool(true),
					}},
				},
			},
		},
	}
	templates = append(templates,
		template.WithDataview(dataview, true),
		template.WithDetailName("Types"),
		template.WithDetailIconEmoji("📒"),
		template.WithRequiredRelations(),
	)

	return smartblock.ObjectApplyTemplate(p, ctx.State, templates...)
}

type MarketplaceRelation struct {
	*Set
}

func NewMarketplaceRelation() *MarketplaceRelation {
	return &MarketplaceRelation{Set: NewSet()}
}

func (p *MarketplaceRelation) Init(ctx *smartblock.InitContext) (err error) {
	err = p.SmartBlock.Init(ctx)
	if err != nil {
		return err
	}

	ot := bundle.TypeKeyRelation.URL()
	templates := []template.StateTransformer{
		template.WithTitle,
		template.WithForcedDetail(bundle.RelationKeySetOf, pbtypes.StringList([]string{ot})),
		template.WithObjectTypesAndLayout([]string{bundle.TypeKeySet.URL()})}
	dataview := model.BlockContentOfDataview{
		Dataview: &model.BlockContentDataview{
			Source:        []string{ot},
			RelationLinks: bundle.MustGetType(bundle.TypeKeyRelation).RelationLinks,
			Views: []*model.BlockContentDataviewView{
				{
					Id:    viewIdMarketplace,
					Type:  model.BlockContentDataviewView_Gallery,
					Name:  "Marketplace",
					Sorts: []*model.BlockContentDataviewSort{{RelationKey: bundle.RelationKeyName.String(), Type: model.BlockContentDataviewSort_Asc}},
					Relations: []*model.BlockContentDataviewRelation{
						{Key: bundle.RelationKeyId.String(), IsVisible: false},
						{Key: bundle.RelationKeyDescription.String(), IsVisible: true},
						{Key: bundle.RelationKeyIconImage.String(), IsVisible: true},
						{Key: bundle.RelationKeyCreator.String(), IsVisible: true}},
					Filters: []*model.BlockContentDataviewFilter{{
						RelationKey: bundle.RelationKeyWorkspaceId.String(),
						Condition:   model.BlockContentDataviewFilter_Equal,
						Value:       pbtypes.String(addr.AnytypeMarketplaceWorkspace),
					}},
				},
				{
					Id:    viewIdLibrary,
					Type:  model.BlockContentDataviewView_Gallery,
					Name:  "Library",
					Sorts: []*model.BlockContentDataviewSort{{RelationKey: bundle.RelationKeyName.String(), Type: model.BlockContentDataviewSort_Asc}},
					Relations: []*model.BlockContentDataviewRelation{
						{Key: bundle.RelationKeyId.String(), IsVisible: false},
						{Key: bundle.RelationKeyDescription.String(), IsVisible: true},
						{Key: bundle.RelationKeyIconImage.String(), IsVisible: true},
						{Key: bundle.RelationKeyCreator.String(), IsVisible: true}},
					Filters: []*model.BlockContentDataviewFilter{{
						RelationKey: bundle.RelationKeyWorkspaceId.String(),
						Condition:   model.BlockContentDataviewFilter_NotEqual,
						Value:       pbtypes.String(addr.AnytypeMarketplaceWorkspace),
					}},
				},
			},
		},
	}
	templates = append(templates, template.WithDataview(dataview, true), template.WithDetailName("Relations"), template.WithDetailIconEmoji("📒"), template.WithRequiredRelations())

	return smartblock.ObjectApplyTemplate(p, ctx.State, templates...)
}

type MarketplaceTemplate struct {
	*Set
}

func NewMarketplaceTemplate() *MarketplaceTemplate {
	return &MarketplaceTemplate{Set: NewSet()}
}

func (p *MarketplaceTemplate) Init(ctx *smartblock.InitContext) (err error) {
	err = p.SmartBlock.Init(ctx)
	if err != nil {
		return err
	}

	ot := bundle.TypeKeyTemplate.URL()
	templates := []template.StateTransformer{
		template.WithTitle,
		template.WithForcedDetail(bundle.RelationKeySetOf, pbtypes.StringList([]string{ot})),
		template.WithObjectTypesAndLayout([]string{bundle.TypeKeySet.URL()})}
	dataview := model.BlockContentOfDataview{
		Dataview: &model.BlockContentDataview{
			Source:        []string{ot},
			RelationLinks: bundle.MustGetType(bundle.TypeKeyTemplate).RelationLinks,
			Views: []*model.BlockContentDataviewView{
				{
					Id:    viewIdMarketplace,
					Type:  model.BlockContentDataviewView_Gallery,
					Name:  "Marketplace",
					Sorts: []*model.BlockContentDataviewSort{{RelationKey: bundle.RelationKeyName.String(), Type: model.BlockContentDataviewSort_Asc}},
					Relations: []*model.BlockContentDataviewRelation{
						{Key: bundle.RelationKeyId.String(), IsVisible: false},
						{Key: bundle.RelationKeyDescription.String(), IsVisible: true},
						{Key: bundle.RelationKeyIconImage.String(), IsVisible: true},
						{Key: bundle.RelationKeyCreator.String(), IsVisible: true}},
					Filters: []*model.BlockContentDataviewFilter{{
						RelationKey: bundle.RelationKeyIsHidden.String(),
						Condition:   model.BlockContentDataviewFilter_NotEqual,
						Value:       pbtypes.Bool(true),
					}},
				},
				{
					Id:    viewIdLibrary,
					Type:  model.BlockContentDataviewView_Gallery,
					Name:  "Library",
					Sorts: []*model.BlockContentDataviewSort{{RelationKey: bundle.RelationKeyName.String(), Type: model.BlockContentDataviewSort_Asc}},
					Relations: []*model.BlockContentDataviewRelation{
						{Key: bundle.RelationKeyId.String(), IsVisible: false},
						{Key: bundle.RelationKeyDescription.String(), IsVisible: true},
						{Key: bundle.RelationKeyIconImage.String(), IsVisible: true},
						{Key: bundle.RelationKeyCreator.String(), IsVisible: true}},
					Filters: []*model.BlockContentDataviewFilter{{
						RelationKey: bundle.RelationKeyIsHidden.String(),
						Condition:   model.BlockContentDataviewFilter_NotEqual,
						Value:       pbtypes.Bool(true),
					}},
				},
			},
		},
	}
	templates = append(templates,
		template.WithDataview(dataview, true),
		template.WithDetailName("Relations"),
		template.WithDetailIconEmoji("📒"),
	)

	return smartblock.ObjectApplyTemplate(p, ctx.State, templates...)
}
