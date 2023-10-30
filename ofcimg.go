package ofcimg

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"
	"ofcimg/gen"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	el "github.com/labstack/gommon/log"

	_ "time/tzdata"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

//go:embed static
var embeddedFiles embed.FS

var nyc *time.Location

func Main() {
	ctx := context.Background()
	e := echo.New()
	e.Logger.SetLevel(el.DEBUG)

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 12, // 4 KB
		LogLevel:  el.ERROR,
	}))

	var err error
	nyc, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatalf("Unable to load New York location:%v", err)
		return
	}

	db, err := openDB(ctx)
	if err != nil {
		log.Fatalf("unable to open db: %v", err)
	}
	query := gen.New(db)
	vptr := &visit{q: query, ID: 0}

	initRouteVisit(vptr, e)
	initStatic(e)

	// this was not specified in assignment, so I just did
	// something simple with html page
	e.POST("/upload", func(c echo.Context) error {
		return upload(c, query)
	})

	e.Start(":9000")
}

// handle static files at / but also do a redir on /
func initStatic(e *echo.Echo) {
	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	log.Printf("checking for static files...are we in live mode? %v", useOS)
	assetHandler := http.FileServer(getFileSystem(useOS))
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/static/index.html")
	})
	e.GET("/*", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

}

// initRoute sets up the mapping from url's to functions to call.
func initRouteVisit(visit *visit, e *echo.Echo) {
	g := e.Group("/api")
	g.POST("/visit", func(c echo.Context) error {
		return visit.createVisit(c)
	})
	g.GET("/visit", func(c echo.Context) error {
		return visit.listVisit(c)
	})
	g.GET("/visit/:id", func(c echo.Context) error {
		if err := c.Bind(visit); err != nil {
			return err
		}
		log.Printf("visit is %d", visit.ID)
		return visit.getSingleVisit(c)
	})
	g.GET("/visit/:id/image", func(c echo.Context) error {
		if err := c.Bind(visit); err != nil {
			return err
		}
		return visit.getSingleVisitImage(c)
	})
}

// open sqlite as a file, insure tables are set up
func openDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}
	return db, nil
}
