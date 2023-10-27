package ofcimg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

// modified a bit but started with https://echo.labstack.com/docs/cookbook/file-upload
func upload(c echo.Context) error {
	// Read form fields
	raw := c.FormValue("id")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filepath.Join("data", file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with id=%d.</p>", file.Filename, id))
}
