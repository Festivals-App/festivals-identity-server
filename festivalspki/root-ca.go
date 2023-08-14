package festivalspki

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func DownloadRootCERTIfNeeded(url string, cachePath string) error {
	if _, err := os.Stat(cachePath); errors.Is(err, os.ErrNotExist) {

		err := downloadFile(cachePath, url)
		if err != nil {
			return err
		}
		log.Info().Msg("Server did download the FestivalsApp Root CA public certificate.")
	}
	return nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
