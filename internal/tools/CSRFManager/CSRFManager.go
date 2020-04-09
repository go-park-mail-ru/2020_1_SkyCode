package CSRFManager

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/2020_1_Skycode/internal/tools"
	"github.com/sirupsen/logrus"
	"hash"
	"strconv"
	"strings"
	"time"
)

type CSRFManager struct {
	HashM hash.Hash
}

func NewCSRFManager() *CSRFManager {
 	return &CSRFManager{
		HashM: sha256.New(),
	}
}

func (m *CSRFManager) GenerateCSRF(userID uint64, cookie string) (string, error) {
	data := fmt.Sprintf("%d:%s", userID, cookie)

	_, err := m.HashM.Write([]byte(data))
	defer m.HashM.Reset()

	if err != nil {
		logrus.Info(err)
		return  "", err
	}

	token := fmt.Sprintf("%s:%d", hex.EncodeToString(m.HashM.Sum(nil)), time.Now().Unix())

	return token, nil
}

func (m *CSRFManager) ValidateCSRF(gotToken string, userID uint64, cookie string) error {
	expectedToken := fmt.Sprintf("%d:%s", userID, cookie)

	_, err := m.HashM.Write([]byte(expectedToken))
	defer m.HashM.Reset()

	if err != nil {
		logrus.Info(err)
		return  tools.HashingError
	}

	expectedToken = hex.EncodeToString(m.HashM.Sum(nil))

	tokenValues := strings.Split(gotToken, ":")

	if len(tokenValues) != 2 {
		logrus.Info("Wrong CSRF token")
		return tools.WrongCSRFtoken
	}

	token := tokenValues[0]
	timestamp, err := strconv.ParseInt(tokenValues[1], 10, 64)

	if token != expectedToken || err != nil{
		logrus.Info(tools.WrongCSRFtoken.Error())
		return tools.WrongCSRFtoken
	}

	if time.Now().Unix() - timestamp > time.Now().AddDate(0, 0, 1).Unix() {
		return tools.ExpiredCSRFError
	}

	return nil
}

