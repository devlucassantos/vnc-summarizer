package dicontainer

import (
	"vnc-summarizer/adapters/apis/dalle"
	interfaces "vnc-summarizer/core/interfaces/dalle"
)

func GetDallEApi() interfaces.DallE {
	return dalle.NewDallEApi()
}
