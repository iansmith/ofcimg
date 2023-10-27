package ofcimg

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"net/http"
	"ofcimg/gen"
	"os"

	"github.com/labstack/echo/v4"
	el "github.com/labstack/gommon/log"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

//go:embed static
var embeddedFiles embed.FS

func Main() {
	ctx := context.Background()
	e := echo.New()
	e.Logger.SetLevel(el.DEBUG)

	db, err := openDB(ctx)
	if err != nil {
		log.Fatalf("unable to open db: %v", err)
	}
	query := gen.New(db)
	vptr := &visit{q: query, ID: 0}

	initRouteVisit(vptr, e)
	initStatic(e)
	// e.GET("/bleah", func(c echo.Context) error {
	// 	return vptr.createVisit(c)
	// })
	e.Start("localhost:9000")
}

// handle static files at / but also do a redir on /
func initStatic(e *echo.Echo) {
	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOS))
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})
	e.GET("/*", echo.WrapHandler(assetHandler))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))

	// this was not specified in assignment, so I just did
	// something simple with html page
	e.POST("/upload", upload)

}

// initRoute sets up the mapping from url's to functions to call.
func initRouteVisit(visit *visit, e *echo.Echo) {
	g := e.Group("/api/visit")
	g.POST("/", func(c echo.Context) error {
		return visit.createVisit(c)
	})
	g.GET("/", func(c echo.Context) error {
		return visit.listVisit(c)
	})
	g.GET("/:id", func(c echo.Context) error {
		if err := c.Bind(visit); err != nil {
			return err
		}
		log.Printf("visit is %d", visit.ID)
		return visit.getSingleVisit(c)
	})
	g.GET("/:id/image", func(c echo.Context) error {
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
