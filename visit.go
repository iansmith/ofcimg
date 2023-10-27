package ofcimg

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ofcimg/gen"

	"github.com/labstack/echo/v4"
)

var ok = []byte("OK")

type visit struct {
	q *gen.Queries
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

	resp := c.Response()
	resp.Status = http.StatusOK
	resp.Writer.Write(ok)

	return nil
}

type OutputVisit struct {
	ID           int32
	StartUnix    int32
	LengthSecond int32
}

// getVisit returns all the visits and the client has to
// sort them out because they might be in the past.
func (v *visit) getVisit(c echo.Context) error {
	log.Printf("got to get visit")
	all, err := v.q.ListVisit(context.Background())
	if err != nil {
		log.Printf("err was xxx %+v", err)
		return err
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(all); err != nil {
		log.Fatalf("problem doing json encode: %v", err)
	}
	resp := c.Response()
	resp.Status = http.StatusOK
	resp.Writer.Write(buf.Bytes())
	log.Printf("num bytes %d", buf.Len())
	resp.Flush()
	resp.Header().Add("Content-Type", "application/json")
	return nil
}
