package utils

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	UserPrefix        = "USR"
	UserAddressPrefix = "USRADR"
)

func GeneratePrefixCode(prefix string) string {
	var (
		code string
		now  = time.Now().In(time.UTC)
	)

	rand.Seed(now.UnixNano())
	code, _ = Generate(`[A-Z]{5}`)
	return fmt.Sprintf("%s%s%s", prefix, now.Format("20060102"), code)
}
