package unsplash

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/any-sync/app/ocache"
	"github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/hbagdi/go-unsplash/unsplash"
	"golang.org/x/oauth2"

	"github.com/anyproto/anytype-heart/core/anytype/config/loadenv"
	"github.com/anyproto/anytype-heart/pkg/lib/core"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/util/uri"
)

var log = logging.Logger("unsplash")

var DefaultToken = ""

const (
	CName         = "unsplash"
	cacheTTL      = time.Minute * 10
	cacheGCPeriod = time.Minute * 5
	anytypeURL    = "https://unsplash.anytype.io/"
)

type Unsplash interface {
	Search(ctx context.Context, query string, max int) ([]Result, error)
	Download(ctx context.Context, id string) (imgPath string, err error)

	app.ComponentRunnable
}

type unsplashService struct {
	mu              sync.Mutex
	cache           ocache.OCache
	client          *unsplash.Unsplash
	limit           int
	tempDirProvider core.TempDirProvider
}

func New() Unsplash {
	return &unsplashService{}
}

func (l *unsplashService) Init(a *app.App) (err error) {
	l.cache = ocache.New(l.search, ocache.WithTTL(cacheTTL), ocache.WithGCPeriod(cacheGCPeriod))
	l.tempDirProvider = app.MustComponent[core.TempDirProvider](a)
	return
}

func (l *unsplashService) Name() (name string) {
	return CName
}

func (l *unsplashService) Run(_ context.Context) error {
	return nil
}

func (l *unsplashService) Close(_ context.Context) error {
	return l.cache.Close()
}

type Result struct {
	ID              string
	Description     string
	PictureThumbUrl string
	PictureSmallUrl string
	PictureFullUrl  string
	PictureHDUrl    string
	Artist          string
	ArtistURL       string
}

type results struct {
	results []Result
}

func (r results) TryClose(objectTTL time.Duration) (bool, error) {
	return true, r.Close()
}

func (r results) Close() error {
	return nil
}

func newFromPhoto(v unsplash.Photo) (Result, error) {
	if v.ID == nil || v.Urls == nil {
		return Result{}, fmt.Errorf("nil input from unsplash")
	}
	res := Result{ID: *v.ID}
	if v.Urls.Thumb != nil {
		res.PictureThumbUrl = v.Urls.Thumb.String()
	}
	if v.Description != nil && *v.Description != "" {
		res.Description = *v.Description
	} else if v.AltDescription != nil {
		res.Description = *v.AltDescription
	}
	if v.Urls.Small != nil {
		res.PictureSmallUrl = v.Urls.Small.String()
	}
	if v.Urls.Regular != nil {
		fUrl := v.Urls.Regular.String()
		// hack to have full hd instead of 1080w,
		// in case unsplash will change the URL format it will not break things
		u, err := uri.ParseURI(fUrl)
		if err == nil {
			if q := u.Query(); q.Get("w") != "" {
				q.Set("w", "1920")
				u.RawQuery = q.Encode()
			}
		}
		res.PictureHDUrl = u.String()
	}
	if v.Urls.Full != nil {
		res.PictureFullUrl = v.Urls.Full.String()
	}
	if v.Photographer == nil {
		return res, nil
	}
	if v.Photographer.Name != nil {
		res.Artist = *v.Photographer.Name
	}
	if v.Photographer.Links != nil && v.Photographer.Links.HTML != nil {
		res.ArtistURL = v.Photographer.Links.HTML.String()
	}
	return res, nil
}

func (l *unsplashService) lazyInitClient() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.client != nil {
		return
	}
	token := DefaultToken
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	l.client = unsplash.New(oauth2.NewClient(context.Background(), ts))
}

func (l *unsplashService) Search(ctx context.Context, query string, limit int) ([]Result, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	l.limit = limit
	v, err := l.cache.Get(ctx, query)
	if err != nil {
		return nil, err
	}

	if r, ok := v.(results); ok {
		return r.results, nil
	} else {
		panic("invalid cache value")
	}
}

func (l *unsplashService) search(ctx context.Context, query string) (ocache.Object, error) {
	l.lazyInitClient()
	query = strings.ToLower(strings.TrimSpace(query))

	var opt unsplash.RandomPhotoOpt

	opt.Count = l.limit
	opt.SearchQuery = query

	res, _, err := l.client.Photos.Random(&opt)
	if err != nil {
		if strings.Contains("404", err.Error()) {
			return nil, nil
		}
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	var photos = make([]Result, 0, len(*res))
	for _, v := range *res {
		res, err := newFromPhoto(v)
		if err != nil {
			continue
		}

		photos = append(photos, res)
	}

	return results{results: photos}, nil
}

func (l *unsplashService) Download(ctx context.Context, id string) (imgPath string, err error) {
	l.lazyInitClient()
	var picture Result
	l.cache.ForEach(func(v ocache.Object) (isContinue bool) {
		// todo: it will be better to save the last result, but we need another lock for this
		if r, ok := v.(results); ok {
			for _, res := range r.results {
				if res.ID == id {
					picture = res
					break
				}
			}
		}
		return picture.ID == ""
	})

	if picture.ID == "" {
		res, _, err := l.client.Photos.Photo(id, nil)
		if err != nil {
			return "", err
		}
		picture, err = newFromPhoto(*res)
		if err != nil {
			return "", err
		}
	}

	req, err := http.NewRequest("GET", picture.PictureHDUrl, nil)
	if err != nil {
		return "", err
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download file from unsplash: %s", err)
	}
	defer resp.Body.Close()
	tmpfile, err := ioutil.TempFile(l.tempDirProvider.TempDir(), picture.ID)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %s", err)
	}
	_, _ = io.Copy(tmpfile, resp.Body)
	tmpfile.Close()

	err = injectIntoExif(tmpfile.Name(), picture.Artist, picture.ArtistURL, picture.Description)
	if err != nil {
		return "", fmt.Errorf("failed to inject exif: %s", err)
	}
	p, err := filepath.Abs(tmpfile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to inject exif: %s", err)
	}

	go func(cl *unsplash.Unsplash) {
		// we must call download endpoint according to the API guidelines
		// but we can do it in a separate goroutine to make sure we will download the picture as fast as possible
		_, _, err = cl.Photos.DownloadLink(id)
		if err != nil {
			log.Errorf("failed to call unsplash download endpoint: %s", err)
		}
	}(l.client)
	return p, nil
}

func PackArtistNameAndURL(name, url string) string {
	return fmt.Sprintf("%s; %s", name, url)
}

func injectIntoExif(filePath, artistName, artistUrl, description string) error {
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file to read exif: %s", err)
	}
	sl := intfc.(*jpegstructure.SegmentList)
	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		return err
	}
	ifdPath := "IFD0"
	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
	if err != nil {
		return err
	}
	// Artist key in decimal is 315 https://www.exiv2.org/tags.html
	err = ifdIb.SetStandard(315, PackArtistNameAndURL(artistName, artistUrl))
	err = ifdIb.SetStandard(270, description)
	err = sl.SetExif(rootIb)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0755)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("failed to open file to write exif: %s", err)
	}
	err = sl.Write(f)
	if err != nil {
		return fmt.Errorf("failed to write exif: %s", err)
	}
	return nil
}

func init() {
	if DefaultToken == "" {
		DefaultToken = loadenv.Get("UNSPLASH_KEY")
	}

	setAnytypeURL()
}

func setAnytypeURL() {
	unsplash.SetupBaseUrl(anytypeURL)
}
