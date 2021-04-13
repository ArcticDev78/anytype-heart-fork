package editor

import (
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/basic"
	_import "github.com/anytypeio/go-anytype-middleware/core/block/editor/import"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/smartblock"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/state"
	"github.com/anytypeio/go-anytype-middleware/core/block/editor/template"
	"github.com/anytypeio/go-anytype-middleware/core/block/meta"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/bundle"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
)

func NewDashboard(m meta.Service, importServices _import.Services) *Dashboard {
	sb := smartblock.New(m)
	return &Dashboard{
		SmartBlock: sb,
		Basic:      basic.NewBasic(sb),
		Import:     _import.NewImport(sb, importServices),
	}
}

type Dashboard struct {
	smartblock.SmartBlock
	basic.Basic
	_import.Import
}

func (p *Dashboard) Init(ctx *smartblock.InitContext) (err error) {
	if err = p.SmartBlock.Init(ctx); err != nil {
		return
	}
	p.DisableLayouts()
	return p.init(ctx.State)
}

func (p *Dashboard) init(s *state.State) (err error) {
	state.CleanupLayouts(s)
	if err = template.ApplyTemplate(p, s,
		template.WithEmpty,
		template.WithDetailName("Home"),
		template.WithDetailIconEmoji("🏠"),
		template.WithRootLink(p.Anytype().PredefinedBlocks().Archive, model.BlockContentLink_Archive),
		template.WithRootLink(p.Anytype().PredefinedBlocks().SetPages, model.BlockContentLink_Dataview),
		template.WithObjectTypesAndLayout([]string{bundle.TypeKeyDashboard.URL()}),
		template.WithRequiredRelations(),
	); err != nil {
		return
	}

	log.Infof("create default structure for dashboard: %v", s.RootId())
	return
}
