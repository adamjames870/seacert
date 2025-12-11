package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
)

type DummyCert struct {
	Name       string `json:"name"`
	CertNumber string `json:"cert-number"`
	Issuer     string `json:"issuer"`
	IssuedDate string `json:"issued-date"`
}

const FileName = "dummy_certs.json"
const PostUrl = "http://localhost:8080/api/certificates"

func LoadDummyCerts() error {

	data, errData := os.ReadFile(FileName)
	if errData != nil {
		return fmt.Errorf("unable to read file: %w", errData)
	}

	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var certs []DummyCert
	if errUnmarshal := json.Unmarshal(data, &certs); errUnmarshal != nil {
		return fmt.Errorf("unable to unmarshal JSON: %w", errUnmarshal)
	}

	for _, cert := range certs {

		certType := models.ParamsAddCertificateType{
			Name:                 cert.Name,
			ShortName:            "xx",
			StcwReference:        "xx",
			NormalValidityMonths: 0,
		}

		certTypeBody, errCertTypeBody := json.Marshal(certType)
		if errCertTypeBody != nil {
			return fmt.Errorf("unable to marshal cert type: %w", errCertTypeBody)
		}

		fmt.Println(string(certTypeBody))

		respCertType, errPostCertType := http.Post("http://localhost:8080/api/certificate-types", "application/json", bytes.NewReader(certTypeBody))
		if errPostCertType != nil {
			return fmt.Errorf("unable to post cert type: %w", errPostCertType)
		}

		decoder := json.NewDecoder(respCertType.Body)
		certTypeReturned := database.CertificateType{}
		errDecode := decoder.Decode(&certTypeReturned)
		if errDecode != nil {
			return fmt.Errorf("unable to decode cert type: %w", errDecode)
		}

		newCert := models.ParamsAddCertificate{
			CertTypeId:      certTypeReturned.ID.String(),
			CertNumber:      cert.CertNumber,
			Issuer:          cert.Issuer,
			IssuedDate:      cert.IssuedDate,
			AlternativeName: "",
			Remarks:         "",
		}

		certBody, errCertBody := json.Marshal(newCert)
		if errCertBody != nil {
			return fmt.Errorf("unable to marshal cert: %w", errCertBody)
		}

		resp, errPost := http.Post(PostUrl, "application/json", bytes.NewReader(certBody))
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
