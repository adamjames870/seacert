package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/models"
)

const FileName = "dummy_certs.json"
const PostUrl = "http://localhost:8080/api/certificates"

func LoadDummyCerts() error {

	data, errData := os.ReadFile(FileName)
	if errData != nil {
		return fmt.Errorf("unable to read file: %w", errData)
	}

	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var certs []models.ParamsAddCertificate
	if errUnmarshal := json.Unmarshal(data, &certs); errUnmarshal != nil {
		return fmt.Errorf("unable to unmarshal JSON: %w", errUnmarshal)
	}

	for _, cert := range certs {
		body, errBody := json.Marshal(cert)
		if errBody != nil {
			return fmt.Errorf("unable to marshal cert: %w", errBody)
		}

		resp, errPost := http.Post(PostUrl, "application/json", bytes.NewReader(body))
		if errPost != nil {
			return fmt.Errorf("unable to post cert: %w", errPost)
		}

		if resp.StatusCode >= 300 {
			// Read response body for debugging
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return fmt.Errorf("server error for cert-number %s: %s", cert.CertNumber, string(b))
		}
		fmt.Println("Saved cert-number: " + cert.CertNumber)
		resp.Body.Close()

	}

	return nil

}
