package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type MyAppsUtils struct{}

func (mu *MyAppsUtils) GetDigestHashForUserLogingToAppService(app string, domain string, sip string, guid string, dn string, info string, challenge string, password string) string {
	str := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s", app, domain, sip, guid, dn, info, challenge, password)
	bytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(bytes[:])
}

func (mu *MyAppsUtils) GetDigestForAppLoginFromJson(message_json, password, challenge string) (string, error) {
	var msgin AppLogin
	if err := json.Unmarshal([]byte(message_json), &msgin); err != nil {
		return "", err
	}
	var calculated_digest string

	if msgin.PbxObj != "" {
		// user/admin login
		log.Trace("GetDigestForAppLoginFromJson user/admin login")
		info, _ := msgin.InfoAsUserDigestString()
		calculated_digest = mu.GetDigestHashForUserLogingToAppService(msgin.App, msgin.Domain, msgin.Sip, msgin.Guid, msgin.Dn, info, challenge, password)
	} else {
		// pbxobj login
		log.Trace("GetDigestForAppLoginFromJson pbxobj login")
		info, _ := msgin.InfoAsPbxobjectDigestString()
		calculated_digest = mu.GetDigestHashForUserLogingToAppService(msgin.App, msgin.Domain, msgin.Sip, msgin.Guid, msgin.Dn, info, challenge, password)
	}

	return calculated_digest, nil
}

func (mu *MyAppsUtils) GetRandomHexString(n int) string {
	charPool := "abcdef0123456789"
	rand.Seed(time.Now().UnixNano()) // does
	b := strings.Builder{}
	for i := 0; i < n; i++ {
		b.WriteByte(charPool[rand.Intn(len(charPool))])
	}
	return b.String()
}

func CheckAppPasswordForMaximumLength(password string) error {
	if password == "" {
		return errors.Errorf("Error: App password cant be empty")
	}
	if len(password) > 15 {
		return errors.Errorf("Error: App password invalid. The app password cant be longer than 15 chars. This is a limitation of passwords stored in pbx-objects. the currently set password has a length of %v", len(password))
	}
	return nil
}
