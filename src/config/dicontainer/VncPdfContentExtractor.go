package dicontainer

import (
	"vnc-summarizer/adapters/apis/pdfcontentextractor"
	interfaces "vnc-summarizer/core/interfaces/pdfcontentextractor"
)

func GetVncPdfContentExtractorApi() interfaces.VncPdfContentExtractor {
	return pdfcontentextractor.NewVncPdfContentExtractorApi()
}
