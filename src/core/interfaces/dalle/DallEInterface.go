package dalle

type DallE interface {
	MakeRequest(prompt, purpose string) (string, error)
}
