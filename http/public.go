package http

import (
	"net/http"
	"strings"

	"github.com/filebrowser/filebrowser/v2/files"
)

var withHashFile = func(fn handleFunc, trim bool) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		id, path := ifPathWithName(r, trim)
			link, err := d.store.Share.GetByHash(id)
			if err != nil {
				return errToStatus(err), err
			}

		user, err := d.store.Users.Get(d.server.Root, link.UserID)
		if err != nil {
			return errToStatus(err), err
		}

		d.user = user

		file, err := files.NewFileInfo(files.FileOptions{
			Fs:      d.user.Fs,
			Path:    link.Path + path,
			Modify:  d.user.Perm.Modify,
			Expand:  true,
			Checker: d,
		})
		if err != nil {
			return errToStatus(err), err
		}

		d.raw = file
		return fn(w, r, d)
	}
}

// ref to https://github.com/filebrowser/filebrowser/pull/727
// `/api/public/dl/MEEuZK-v/file-name.txt` for old browsers to save file with correct name
func ifPathWithName(r *http.Request, trim bool) (string, string) {
	pathElements := strings.Split(r.URL.Path, "/")
	// prevent maliciously constructed parameters like `/api/public/dl/XZzCDnK2_not_exists_hash_name`
	// len(pathElements) will be 1, and golang will panic `runtime error: index out of range`
	if len(pathElements) < 2 { //nolint: mnd
		return r.URL.Path, ""
	}
	id := pathElements[0]
	if trim {
		return id, strings.Join(pathElements[1:len(pathElements)-1], "/")
	}
	return id, strings.Join(pathElements[1:], "/")
}

var publicShareHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)

	if file.IsDir {
		file.Listing.Sorting = files.Sorting{By: "name", Asc: false}
		file.Listing.ApplySort()
		return renderJSON(w, r, file)
	}

	return renderJSON(w, r, file)
}, false)

var publicDlHandler = withHashFile(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	file := d.raw.(*files.FileInfo)
	if !file.IsDir {
		return rawFileHandler(w, r, file)
	}

	return rawDirHandler(w, r, d, file)
}, true)
