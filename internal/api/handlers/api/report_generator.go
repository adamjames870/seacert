package api

import (
	"fmt"
	"time"

	"github.com/adamjames870/seacert/internal/dto"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func GenerateCertificatesReport(certs []dto.Certificate) (core.Document, error) {
	cfg := config.NewBuilder().
		WithDefaultFont(&props.Font{
			Family: fontfamily.Helvetica,
		}).
		WithOrientation(orientation.Horizontal).
		WithLeftMargin(25).
		WithRightMargin(25).
		WithTopMargin(15).
		WithMaxGridSize(24).
		Build()

	m := maroto.New(cfg)

	// Header
	m.AddRows(
		row.New(20).Add(
			col.New(16).Add(
				text.New("Certificate Summary Report", props.Text{
					Size:  18,
					Style: fontstyle.Bold,
					Align: align.Left,
					Top:   3,
				}),
			),
			col.New(8).Add(
				text.New(fmt.Sprintf("Report Created %s", time.Now().Format("02 Jan 2006")), props.Text{
					Size:  10,
					Style: fontstyle.Italic,
					Align: align.Right,
					Top:   6,
				}),
			),
		),
		line.NewRow(1, props.Line{
			Color:         &props.Color{Red: 100, Green: 100, Blue: 100},
			SizePercent:   100,
			OffsetPercent: 50,
		}),
		row.New(5), // Spacer
	)

	// Table Header
	m.AddRows(
		row.New(7).Add(
			col.New(6).Add(text.New("Certificate Name", props.Text{Style: fontstyle.Bold, Size: 10, Top: 1})),
			col.New(2).Add(text.New("Short", props.Text{Style: fontstyle.Bold, Size: 9, Top: 1})),
			col.New(2).Add(text.New("STCW", props.Text{Style: fontstyle.Bold, Size: 9, Top: 1})),
			col.New(4).Add(text.New("Cert Number", props.Text{Style: fontstyle.Bold, Size: 9, Top: 1})),
			col.New(4).Add(text.New("Issuer", props.Text{Style: fontstyle.Bold, Size: 9, Top: 1})),
			col.New(3).Add(text.New("Issued", props.Text{Style: fontstyle.Bold, Size: 10, Top: 1, Align: align.Center})),
			col.New(3).Add(text.New("Expiry", props.Text{Style: fontstyle.Bold, Size: 10, Top: 1, Align: align.Center})),
		),
		line.NewRow(0.5, props.Line{
			Color:         &props.Color{Red: 50, Green: 50, Blue: 50},
			SizePercent:   100,
			OffsetPercent: 50,
		}),
	)

	// Dynamically adjust font size based on number of certificates to fit on one page
	fontSize := 9.0
	if len(certs) > 25 {
		fontSize = 8.0
	}
	if len(certs) > 35 {
		fontSize = 7.0
	}

	for _, cert := range certs {
		if cert.HasSuccessors {
			continue
		}
		expiryStr := "No Expiry"
		expiryColor := &props.Color{Red: 0, Green: 0, Blue: 0}
		if !cert.ExpiryDate.IsZero() {
			expiryStr = cert.ExpiryDate.Format("02 Jan 2006")
			if cert.ExpiryDate.Before(time.Now()) {
				expiryColor = &props.Color{Red: 200, Green: 0, Blue: 0}
			} else if cert.ExpiryDate.Before(time.Now().AddDate(0, 3, 0)) {
				expiryColor = &props.Color{Red: 200, Green: 150, Blue: 0}
			}
		}

		issuer := cert.IssuerName
		if cert.IssuerCountry != "" {
			issuer = fmt.Sprintf("%s (%s)", cert.IssuerName, cert.IssuerCountry)
		}

		m.AddRows(
			row.New(7).Add(
				col.New(6).Add(text.New(cert.CertTypeName, props.Text{Size: fontSize, Top: 1})),
				col.New(2).Add(text.New(cert.CertTypeShortName, props.Text{Size: fontSize - 1, Top: 1})),
				col.New(2).Add(text.New(cert.CertTypeStcwRef, props.Text{Size: fontSize - 1, Top: 1})),
				col.New(4).Add(text.New(cert.CertNumber, props.Text{Size: fontSize - 1, Top: 1})),
				col.New(4).Add(text.New(issuer, props.Text{Size: fontSize - 1, Top: 1})),
				col.New(3).Add(text.New(cert.IssuedDate.Format("02 Jan 2006"), props.Text{Size: fontSize, Top: 1, Align: align.Center})),
				col.New(3).Add(text.New(expiryStr, props.Text{
					Size:  fontSize,
					Top:   1,
					Align: align.Center,
					Color: expiryColor,
					Style: func() fontstyle.Type {
						if expiryColor.Red > 0 {
							return fontstyle.Bold
						}
						return fontstyle.Normal
					}(),
				})),
			),
			line.NewRow(0.2, props.Line{
				Color:         &props.Color{Red: 230, Green: 230, Blue: 230},
				SizePercent:   100,
				OffsetPercent: 50,
			}),
		)
	}

	// Footer
	m.RegisterFooter(
		row.New(15).Add(
			col.New(24).Add(
				text.New("Report Created at seacert.app", props.Text{
					Size:  8,
					Align: align.Center,
					Color: &props.Color{Red: 150, Green: 150, Blue: 150},
					Top:   5,
				}),
			),
		),
	)

	return m.Generate()
}
