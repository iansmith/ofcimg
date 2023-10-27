package ofcimg

import (
	"io/fs"
	"log"
	"net/http"
	"os"
)

// copied from https://echo.labstack.com/docs/cookbook/embed-resources
// allows switch from the local filesystem (for dev) and the embedded
// filesystem above
func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")
		return http.FS(os.DirFS("static"))
	}
	log.Print("using embed mode")
	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}
