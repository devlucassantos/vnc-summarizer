package datetime

import (
	"fmt"
	"time"
)

func GetCurrentDateTimeInBrazil() (*time.Time, error) {
	saoPauloLocation, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Error loading Brazil time zone (Based on São Paulo time zone): ", err)
		return nil, err
	}

	currentDateTimeInBrazil := time.Now().In(saoPauloLocation)
	return &currentDateTimeInBrazil, err
}
