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
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var ddl string

//go:embed static
var embeddedFiles embed.FS

func Main() {
	ctx := context.Background()
	e := echo.New()

	db, err := openDB(ctx)
	if err != nil {
		log.Fatalf("unable to open db: %v", err)
	}
	query := gen.New(db)
	visit := &visit{query}

	initRouteVisit(visit, e)
	initStatic(e)
	e.GET("/bleah", func(c echo.Context) error {
		log.Printf("xxx %+v", visit)
		return visit.createVisit(c)
	})
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
}

// initRoute sets up the mapping from url's to functions to call.
func initRouteVisit(visit *visit, e *echo.Echo) {
	g := e.Group("/api/visit")
	g.POST("/", func(c echo.Context) error {
		return visit.createVisit(c)
	})
	g.GET("/", func(c echo.Context) error {
		return visit.getVisit(c)
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
