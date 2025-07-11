package pdfcontentextractor

type VncPdfContentExtractor interface {
	MakeRequest(pdfUrl string) (string, error)
}
