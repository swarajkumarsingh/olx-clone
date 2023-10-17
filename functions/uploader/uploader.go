package uploader

import (
	"mime/multipart"
	"net/http"
	"olx-clone/functions/general"
	"strings"
)

// CheckSupportedFormats checks for supported formats and returns file object, extension, error http code and error message if any
func CheckSupportedFormats(r *http.Request, extensions, contentTypes []string) (multipart.File, string, int, string) {
	// check for file types using file name
	file, header, _ := r.FormFile("file")
	if header == nil {
		// logger.WithRequest(r).Println("header is null")
		return nil, "", http.StatusBadRequest, "couldn't determine the file headers, please try again"
	}
	fileNameArr := strings.Split(header.Filename, ".")
	if len(fileNameArr) < 2 {
		// no extension case
		return nil, "", http.StatusBadRequest, "no file extension found"
	}
	extension := strings.ToLower(fileNameArr[len(fileNameArr)-1]) // change case to lower for easy checks
	if !general.InArrStr(extension, extensions) {                 // []string{"pdf", "jpeg", "png", "jpg"}
		return nil, "", http.StatusBadRequest, "file type not supported"
	}

	// now check using content type
	buff := make([]byte, 512) // why 512 bytes, see http://golang.org/pkg/net/http/#DetectContentType
	_, err := file.Read(buff)
	if err != nil {
		// logger.WithRequest(r).Errorln(err)
		return nil, "", http.StatusBadRequest, "couldn't read the file, please try again"
	}
	contentType := strings.ToLower(http.DetectContentType(buff))
	if !general.InArrStr(contentType, contentTypes) { // []string{"image/jpeg", "image/jpg", "image/png", "application/pdf"}
		return nil, "", http.StatusBadRequest, "file type not supported"
	}

	// reset the read pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		// logger.WithRequest(r).Errorln(err)
	}

	return file, extension, 0, ""
}

func GetImageFromMultipartRequest(r *http.Request) (multipart.File, string) {
	extensions := []string{"jpeg", "png", "jpg"}
	contentTypes := []string{"image/jpeg", "image/png", "image/jpg"}
	file, extension, _, _ := CheckSupportedFormats(r, extensions, contentTypes)
	return file, extension
}
