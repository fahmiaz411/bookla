package encrypt

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/fahmiaz411/bookla/config"
)

func Token(v any) (token string, err error) {
	var bytes []byte
	bytes, err = json.Marshal(v)
	if err != nil {
		return
	}

	token = encrypt(string(bytes), hex.EncodeToString([]byte(config.Env.EncKey)))

	return
}

func DeToken(v string, out any) (err error) {
	str, err := decrypt(v, hex.EncodeToString([]byte(config.Env.EncKey)))
	if err != nil {
		return errors.New("Invalid Token")
	}
	return json.Unmarshal([]byte(str), out)
}

func BearerToken(s string) string {
	return strings.Replace(s, "Bearer ", "", 1)
}
