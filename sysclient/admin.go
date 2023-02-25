package sysclient

import (
	"crypto/rc4"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ricoschulte/go-myapps/encryption"
)

// Sysclient message header
const (
	MessageTypeAdmin byte = 0x01
)

// Admin message header
var AdminSendIdentify []byte = []byte{0x00, 0x00}             // Send identify
var AdminReceiveSysclientPassword []byte = []byte{0x00, 0x01} // Receive sysclient password
var AdminReceiveChallenge []byte = []byte{0x00, 0x02}         // Receive challenge
var AdminReceiveNewAdminPassword []byte = []byte{0x00, 0x03}  // Receive new administrative password
var AdminProvision []byte = []byte{0x00, 0x04}                // Provision
var AdminProvisionResult []byte = []byte{0x00, 0x05}          // Provision result
var AdminConfiguration []byte = []byte{0x00, 0x06}            // Configuration

// SMT | Byte offset
//     |  0      1      2      3      4      5      6      7      8      9     ...
//     +------+------+------+------+------+------+------+------+------+------+------+
//  1  | Msg- |   Adm-Msg   | Data.........

type AdminMessage struct {
	Type    byte
	Command []byte // 2 bytes
	Data    []byte // the rest of the data

}

func NewAdminMessage(message []byte) (*AdminMessage, error) {
	if len(message) < 3 { // 3 + data
		return nil, errors.New("invalid length of AdminMessage")
	}

	m := &AdminMessage{
		Type:    message[0],
		Command: message[1:3],
		Data:    message[3:],
	}

	if m.Type != MessageTypeAdmin {
		return nil, errors.New("message is not a AdminMessage")
	}

	return m, nil
}
func (am *AdminMessage) AsBytes() []byte {
	message_as_bytes := []byte{am.Type}
	message_as_bytes = append(message_as_bytes, am.Command...)
	message_as_bytes = append(message_as_bytes, am.Data...)
	return message_as_bytes
}

/*
receives the sysclient password and stores it.
no answer is required
*/
func (am *AdminMessage) HandleAdminReceiveSysclientPassword(secretKey []byte, fileSysclientpassword string) error {
	password, err_from_json := NewPassword(am.Data)
	if err_from_json != nil {
		return fmt.Errorf("parsing password from JSON failed: %v", err_from_json)
	}
	if len(password.Password) != 32 {
		return fmt.Errorf("the password parsed from JSON has a invalid length of: %v", len(password.Password))
	}

	//err := os.WriteFile(fileSysclientpassword, []byte(password.Password), 0644)
	err := encryption.EncryptFileSha256AES256(secretKey, []byte(password.Password), fileSysclientpassword, 0644)
	if err != nil {
		return fmt.Errorf("error while writing adminpassword to file '%s': %v", fileSysclientpassword, err)
	}
	return nil
}

func (am *AdminMessage) HandleAdminReceiveChallenge(secretKey []byte, deviceInfo *Identity, fileSysclientpassword string) (*AdminMessage, error) {
	challenge, err_from_json := NewChallenge(am.Data)
	if err_from_json != nil {
		return nil, fmt.Errorf("parsing challenge from json failed: %v", err_from_json)
	}

	if challenge.Challenge == "" {
		return nil, fmt.Errorf("couldn't parse Challenge Message")
	}
	fileinfo, err_read_file := os.Stat(fileSysclientpassword)
	if err_read_file != nil {
		return nil, fmt.Errorf("could not read sysclientpassword from file '%s': %v", fileSysclientpassword, err_read_file)
	}

	if fileinfo.IsDir() {
		return nil, fmt.Errorf("path is a directory: %s", fileinfo.Name())
	}
	//password, err := os.ReadFile(fileSysclientpassword)
	password, err := encryption.DecryptFileSha256AES256(secretKey, fileSysclientpassword)
	if err != nil {
		fmt.Printf("error while reading password file: %v\n", err)
		return nil, err
	}
	digest, err_digest := am.GetLoginDigest(deviceInfo.Id, deviceInfo.Product, deviceInfo.Version, challenge.Challenge, string(password))
	if err_digest != nil {
		return nil, err_digest
	}

	deviceInfo.Digest = digest
	idbytes, err_to_bytes := deviceInfo.ToBytes()
	if err_to_bytes != nil {
		return nil, err_to_bytes
	}

	message_as_bytes := []byte{am.Type}
	message_as_bytes = append(message_as_bytes, AdminSendIdentify...)
	message_as_bytes = append(message_as_bytes, idbytes...)
	return NewAdminMessage(message_as_bytes)

}

func (am *AdminMessage) GetLoginDigest(id, product, version, challenge, password string) (string, error) {
	if id == "" {
		return "", errors.New("id cant be empty")
	}
	if product == "" {
		return "", errors.New("product cant be empty")
	}
	if version == "" {
		return "", errors.New("version cant be empty")
	}
	if challenge == "" {
		return "", errors.New("challenge cant be empty")
	}
	if password == "" {
		return "", errors.New("password cant be empty")
	}

	bytes := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%s:%s", id, product, version, challenge, password)))
	digest := hex.EncodeToString(bytes[:])

	return digest, nil
}

/*
receives a admin password and stores it, no answer is required
*/
func (am *AdminMessage) handleAdminReceiveNewAdministrativePassword(secretKey []byte, fileSysclientpassword string, fileAdministrativePassword string) error {
	if fileSysclientpassword == "" {
		return errors.New("fileSysclientpassword cant be empty")
	}

	payload, err := NewAdministrativePassword(am.Data)
	if err != nil {
		return fmt.Errorf("couldn't parse AdministrativePassword Message: %v", err)
	}

	//passwordBytes, err := os.ReadFile(fileSysclientpassword)
	passwordBytes, err := encryption.DecryptFileSha256AES256(secretKey, fileSysclientpassword)

	if err != nil {
		return err
	}
	decryped_adminpassword, err_decrypt := am.DecryptAdminPassword(string(passwordBytes), payload.Seed, payload.Pwd)
	if err_decrypt != nil {
		return err_decrypt
	}

	//err = os.WriteFile(fileAdministrativePassword, decryped_adminpassword, 0644)
	err_write_adminpassword := encryption.EncryptFileSha256AES256(secretKey, decryped_adminpassword, fileAdministrativePassword, 0644)
	if err_write_adminpassword != nil {
		return err_write_adminpassword
	}
	return nil //resp, nil
}

func (am *AdminMessage) DecryptAdminPassword(sysclientPassword string, seed string, pwd string) ([]byte, error) {
	if sysclientPassword == "" {
		return nil, errors.New("sysclientPassword cant be empty")
	}
	if len(sysclientPassword) != 32 {
		return nil, fmt.Errorf("sysclientPassword has a invalid length of '%v'. should have a length of 32", len(sysclientPassword))
	}
	if seed == "" {
		return nil, errors.New("seed cant be empty")
	}
	if len(seed) != 16 {
		return nil, fmt.Errorf("seed has a invalid length of '%v'. should have a length of 16", len(seed))
	}

	if pwd == "" {
		return nil, errors.New("pwd cant be empty")
	}

	seedPassword := seed + ":" + sysclientPassword
	key := []byte(seedPassword)
	block, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(pwd))
	block.XORKeyStream(decrypted, []byte(pwd))

	return decrypted, nil
}

type Password struct {
	Password string `json:"password"`
}

func NewPassword(json_bytes []byte) (*Password, error) {
	var payload Password
	err_from_json := json.Unmarshal(json_bytes, &payload)
	if err_from_json != nil {
		return nil, err_from_json
	} else {
		return &payload, nil
	}
}

type Challenge struct {
	Challenge string `json:"challenge"`
}

func NewChallenge(json_bytes []byte) (*Challenge, error) {
	var payload Challenge
	err_from_json := json.Unmarshal(json_bytes, &payload)
	if err_from_json != nil {
		return nil, err_from_json
	} else {
		return &payload, nil
	}
}

type AdministrativePassword struct {
	Seed string `json:"seed"`
	Pwd  string `json:"pwd"`
}

func NewAdministrativePassword(json_bytes []byte) (*AdministrativePassword, error) {
	var payload AdministrativePassword
	err_from_json := json.Unmarshal(json_bytes, &payload)
	if err_from_json != nil {
		return nil, err_from_json
	} else {
		return &payload, nil
	}
}
