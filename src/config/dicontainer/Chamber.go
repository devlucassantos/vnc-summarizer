package dicontainer

import (
	"vnc-summarizer/adapters/apis/chamber"
	interfaces "vnc-summarizer/core/interfaces/chamber"
)

func GetChamberApi() interfaces.Chamber {
	return chamber.NewChamberApi()
}
