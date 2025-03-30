package handler

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strconv"

	token "github.com/Festivals-App/festivals-identity-server/auth"
	servertools "github.com/Festivals-App/festivals-server-tools"
	"github.com/rs/zerolog/log"
)

func GetValidationKey(auth *token.AuthService, db *sql.DB, w http.ResponseWriter, r *http.Request) {

	crt, err := os.Open(auth.ValidationKeyFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open public validation certificate.")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	respondFile(w, crt)
}

// respondFile makes the response with payload as json format
func respondFile(w http.ResponseWriter, file *os.File) {

	// calculate content size
	fileInfo, err := file.Stat()
	if err != nil || fileInfo == nil {
		log.Error().Err(err).Msg("Failed to read file stats for file: '" + file.Name() + "'")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	size := fileInfo.Size()

	// calculate content type dynamically
	contentType, err := getFileContentType(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve content type for file: '" + file.Name() + "'")
		servertools.RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send write file to response")
	}
}

func getFileContentType(file *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the read pointer if necessary.
	_, _ = file.Seek(0, 0)

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
