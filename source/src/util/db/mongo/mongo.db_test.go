package db

import (
	"fmt"

	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load("../../../.env"); err != nil {
		fmt.Println(err)
	}
}
