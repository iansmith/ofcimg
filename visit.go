package ofcimg

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ofcimg/gen"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type visit struct {
	ID int64 `param:"id"`
	q  *gen.Queries
}

func (v *visit) createVisit(c echo.Context) error {
	cvp := gen.CreateVisitParams{
		StartTimeUnix: 1698426197,
		LengthSecond:  60 * 15,
	}
	id, err := v.q.CreateVisit(context.Background(), cvp)
	if err != nil {
		return err
	}
	log.Printf("id is %d", id)

	c.Response().Header().Add("created", fmt.Sprint(id))
	c.HTML(http.StatusOK, "ok")

	return nil
}

// listVisit returns all the visits and the client has to
// sort them out because they might be in the past.
func (v *visit) listVisit(c echo.Context) error {
	all, err := v.q.ListVisit(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "not found")
			return nil
		}
		return err
	}
	jsonEncodeResult(c, all)
	return nil
}

// getSingleVisit returns a visit based on the id provided in the url
func (v *visit) getSingleVisit(c echo.Context) error {
	only, err := v.q.GetVisit(context.Background(), v.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "not found")
			return nil
		}
		return err
	}
	jsonEncodeResult(c, only)
	return nil
}

func (v *visit) getSingleVisitImage(c echo.Context) error {
	only, err := v.q.GetVisit(context.Background(), v.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "not found")
			return nil
		}
		return err
	}
	if !only.Filepath.Valid {
		return fmt.Errorf("unable to get a valid path back from db")
	}
	c.Inline(filepath.Join("data", only.Filepath.String), "image")
	return nil
}

// jsonEncodeResult is a utility for stuffing everything from an
// go object to json.
func jsonEncodeResult(c echo.Context, all interface{}) {
	resp := c.Response()
	buf := &bytes.Buffer{}

	enc := json.NewEncoder(buf)
	if err := enc.Encode(all); err != nil {
		panic("problem doing json encode:" + err.Error())
	}

	resp.Status = http.StatusOK
	resp.Writer.Write(buf.Bytes())
	log.Printf("num bytes %d", buf.Len())
	resp.Flush()
	resp.Header().Add("Content-Type", "application/json")

}
