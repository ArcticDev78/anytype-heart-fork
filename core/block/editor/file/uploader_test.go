package file_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"

	"github.com/anytypeio/go-anytype-middleware/core/block/editor/file"
	"github.com/anytypeio/go-anytype-middleware/core/block/simple"
	file2 "github.com/anytypeio/go-anytype-middleware/core/block/simple/file"
	"github.com/anytypeio/go-anytype-middleware/core/files"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/core"
	"github.com/anytypeio/go-anytype-middleware/pkg/lib/pb/model"
	"github.com/anytypeio/go-anytype-middleware/util/testMock"
	"github.com/anytypeio/go-anytype-middleware/util/testMock/mockFile"
)

func TestUploader_Upload(t *testing.T) {
	ctx := context.Background()
	newBlock := func(tp model.BlockContentFileType) file2.Block {
		return simple.New(&model.Block{Content: &model.BlockContentOfFile{File: &model.BlockContentFile{Type: tp}}}).(file2.Block)
	}
	t.Run("empty source", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.tearDown()
		res := fx.Upload(ctx)
		require.Error(t, res.Err)
	})
	t.Run("image by block type", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.tearDown()
		im := fx.newImage("123")
		fx.anytype.EXPECT().ImageAdd(gomock.Any(), gomock.Any()).Return(im, nil)
		im.EXPECT().GetOriginalFile(gomock.Any())
		b := newBlock(model.BlockContentFile_Image)
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		res := fx.Uploader.SetBlock(b).SetFile("./testdata/unnamed.jpg").Upload(ctx)
		require.NoError(t, res.Err)
		assert.Equal(t, res.Hash, "123")
		assert.Equal(t, res.Name, "unnamed.jpg")
		assert.Equal(t, b.Model().GetFile().Name, "unnamed.jpg")
	})
	t.Run("image type detect", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.tearDown()
		im := fx.newImage("123")
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		fx.anytype.EXPECT().ImageAdd(gomock.Any(), gomock.Any()).Return(im, nil)
		im.EXPECT().GetOriginalFile(gomock.Any())
		res := fx.Uploader.AutoType(true).SetFile("./testdata/unnamed.jpg").Upload(ctx)
		require.NoError(t, res.Err)
	})
	t.Run("image to file failover", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.tearDown()
		meta := &files.FileMeta{
			Media: "text/text",
			Name:  "test.txt",
			Size:  3,
			Added: time.Now(),
		}
		// fx.anytype.EXPECT().ImageAdd(gomock.Any(), gomock.Any()).Return(nil, image.ErrFormat)
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		fx.anytype.EXPECT().FileAdd(gomock.Any(), gomock.Any()).Return(fx.newFile("123", meta), nil)
		b := newBlock(model.BlockContentFile_Image)
		res := fx.Uploader.SetBlock(b).SetFile("./testdata/test.txt").Upload(ctx)
		require.NoError(t, res.Err)
		assert.Equal(t, res.Hash, "123")
		assert.Equal(t, res.Name, "test.txt")
		assert.Equal(t, b.Model().GetFile().Name, "test.txt")
		assert.Equal(t, b.Model().GetFile().Type, model.BlockContentFile_File)
	})
	t.Run("file from url", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/unnamed.jpg")
		})
		serv := httptest.NewServer(mux)
		defer serv.Close()

		fx := newFixture(t)
		defer fx.tearDown()
		im := fx.newImage("123")
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		fx.anytype.EXPECT().ImageAdd(gomock.Any(), gomock.Any()).Return(im, nil)
		im.EXPECT().GetOriginalFile(gomock.Any())
		res := fx.Uploader.AutoType(true).SetUrl(serv.URL + "/unnamed.jpg").Upload(ctx)
		require.NoError(t, res.Err)
		assert.Equal(t, res.Hash, "123")
		assert.Equal(t, res.Name, "unnamed.jpg")
		res.Size = 1
		b := res.ToBlock()
		assert.Equal(t, b.Model().GetFile().Name, "unnamed.jpg")
	})
	t.Run("file from Content-Disposition", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Disposition", "form-data; name=\"fieldName\"; filename=\"filename\"")
			http.ServeFile(w, r, "./testdata/unnamed.jpg")
		})
		serv := httptest.NewServer(mux)
		defer serv.Close()

		fx := newFixture(t)
		defer fx.tearDown()
		im := fx.newImage("123")
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		fx.anytype.EXPECT().ImageAdd(gomock.Any(), gomock.Any()).Return(im, nil)
		im.EXPECT().GetOriginalFile(gomock.Any())
		res := fx.Uploader.AutoType(true).SetUrl(serv.URL + "/unnamed.jpg").Upload(ctx)
		require.NoError(t, res.Err)
		assert.Equal(t, res.Hash, "123")
		assert.Equal(t, res.Name, "filename")
		res.Size = 1
		b := res.ToBlock()
		assert.Equal(t, b.Model().GetFile().Name, "filename")
	})
	t.Run("file without url params", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./testdata/unnamed.jpg")
		})
		serv := httptest.NewServer(mux)
		defer serv.Close()

		fx := newFixture(t)
		defer fx.tearDown()
		im := fx.newImage("123")
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		fx.anytype.EXPECT().ImageAdd(gomock.Any(), gomock.Any()).Return(im, nil)
		im.EXPECT().GetOriginalFile(gomock.Any())
		res := fx.Uploader.AutoType(true).SetUrl(serv.URL + "/unnamed.jpg?text=text").Upload(ctx)
		require.NoError(t, res.Err)
		assert.Equal(t, res.Hash, "123")
		assert.Equal(t, res.Name, "unnamed.jpg")
		res.Size = 1
		b := res.ToBlock()
		assert.Equal(t, b.Model().GetFile().Name, "unnamed.jpg")
	})
	t.Run("bytes", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.tearDown()
		fx.fileService.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil)
		fx.anytype.EXPECT().FileAdd(gomock.Any(), gomock.Any()).Return(fx.newFile("123", &files.FileMeta{}), nil)
		res := fx.Uploader.SetBytes([]byte("my bytes")).SetName("filename").Upload(ctx)
		require.NoError(t, res.Err)
		assert.Equal(t, res.Hash, "123")
		assert.Equal(t, res.Name, "filename")
	})
}

func newFixture(t *testing.T) *uplFixture {
	fx := &uplFixture{
		ctrl: gomock.NewController(t),
	}
	fx.anytype = testMock.NewMockService(fx.ctrl)
	fx.fileService = mockFile.NewMockBlockService(fx.ctrl)

	fx.Uploader = file.NewUploader(fx.fileService, fx.anytype, core.NewTempDirService(nil))
	return fx
}

type uplFixture struct {
	file.Uploader
	fileService *mockFile.MockBlockService
	anytype     *testMock.MockService
	ctrl        *gomock.Controller
}

func (fx *uplFixture) newImage(hash string) *testMock.MockImage {
	im := testMock.NewMockImage(fx.ctrl)
	im.EXPECT().Hash().Return(hash).AnyTimes()
	return im
}

func (fx *uplFixture) newFile(hash string, meta *files.FileMeta) *testMock.MockFile {
	f := testMock.NewMockFile(fx.ctrl)
	f.EXPECT().Hash().Return(hash).AnyTimes()
	f.EXPECT().Meta().Return(meta).AnyTimes()
	return f
}

func (fx *uplFixture) tearDown() {
	fx.ctrl.Finish()
}
