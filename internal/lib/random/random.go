package random

import (
	"math/rand"
	"time"
)

func GenerateAlias(aliasLength int) string {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	alias := make([]rune, aliasLength)

	chars := []rune("WERTYUIOPASDFGHJKLZXCVBNM" +
		"qwertyuioplasdfghjkzxcvbnm" + "1234567890")

	for i := range alias {
		alias[i] = chars[rnd.Intn(len(chars))]
	}

	return string(alias)
}
