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
	"strconv"
	"time"

	//_ "time/tzdata"

	"github.com/labstack/echo/v4"
)

type visit struct {
	ID int64 `param:"id"`
	q  *gen.Queries
}

type VisitFormData struct {
	Day   int `param:"day"`
	Month int `param:"month"`
	Year  int `param:"year"`
	Hour  int `param:"hour"`
	Min   int `param:"min"`
	Len   int `param:"len"`
}

func (v *visit) createVisit(c echo.Context) error {
	// XXX why doesn't this bind work?
	// err := c.Bind(fd)
	// if err != nil {
	// 	return err
	// }
	vfd := convertFormData(c)
	log.Printf("all fields %+v", vfd)
	unixTime, len := formDataToTimeAndLen(vfd)
	cvp := gen.CreateVisitParams{
		StartTimeUnix: unixTime,
		LengthSecond:  len,
	}

	id, err := v.q.CreateVisit(context.Background(), cvp)
	if err != nil {
		return err
	}

	c.Response().Header().Add("created", fmt.Sprint(id))
	//c.Redirect(http.StatusTemporaryRedirect, "/index.html")
	c.HTML(http.StatusOK, fmt.Sprintf("id is %d", id))

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
	resp.Flush()
	resp.Header().Add("Content-Type", "application/json")
}

func mustIntConvert(s string) int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic("unable to convert integer:" + err.Error())
	}
	return int(i)
}

// xxx why is this needed? Why doesn't using Bind() on the context
// xxx do this for us?
func convertFormData(c echo.Context) *VisitFormData {
	result := &VisitFormData{}
	result.Day = mustIntConvert(c.FormValue("day"))
	result.Month = mustIntConvert(c.FormValue("month"))
	result.Year = mustIntConvert(c.FormValue("year"))
	if result.Year < 100 {
		result.Year += 2000
	}
	result.Hour = mustIntConvert(c.FormValue("hour"))
	result.Min = mustIntConvert(c.FormValue("min"))
	result.Len = mustIntConvert(c.FormValue("len"))
	return result
}

const tz = "America/NewYork"

func formDataToTimeAndLen(vfd *VisitFormData) (int64, int64) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatalf("unable to find timezone %s: %v", tz, err)
	}
	if vfd.Year < 100 {
		vfd.Year += 2000
	}
	mon := time.Month(vfd.Month)
	t := time.Date(vfd.Year, mon, vfd.Day, vfd.Hour, vfd.Min, 0, 0, loc)
	return int64(t.Unix()), int64(vfd.Len * 60)
}
