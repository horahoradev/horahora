package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"

	"github.com/horahoradev/horahora/user_service/internal/model"
	userproto "github.com/horahoradev/horahora/user_service/protocol"
	"golang.org/x/crypto/bcrypt"
	jose "gopkg.in/square/go-jose.v2"
)

const hashCost = 5

func Login(username, password string, privateKey *rsa.PrivateKey, u *model.UserModel) (string, error) {

	uid, err := u.GetUserWithUsername(username)
	if err != nil {
		return "", err
	}

	passHash, err := u.GetPassHash(uid)
	if err != nil {
		return "", err
	}

	inpPassHash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", err
	}

	isValid, err := compareHashedPassword([]byte(passHash), inpPassHash)
	if err != nil {
		return "", err
	}

	if !isValid {
		return "", errors.New("invalid password")
	}

	// Password is valid
	payload := JWTPayload{UID: uid}
	return CreateJWT(payload, privateKey)
}

func Register(username, email, password string, u *model.UserModel, privateKey *rsa.PrivateKey, foreignUser bool,
	foreignUserID string, foreignWebsite userproto.Site) (string, error) {
	pwBytes := []byte(password)

	var passHash []byte
	var err error
	// TODO: salt + pepper?
	if !foreignUser {
		passHash, err = bcrypt.GenerateFromPassword(pwBytes, hashCost)
		if err != nil {
			return "", err
		}
	}

	uid, err := u.NewUser(username, email, passHash, foreignUser, foreignUserID, foreignWebsite)
	if err != nil {
		return "", err
	}

	payload := JWTPayload{UID: uid}
	return CreateJWT(payload, privateKey)
}

type JWTPayload struct {
	UID int64 `json:"uid"`
}

func CreateJWT(payload JWTPayload, privateKey *rsa.PrivateKey) (string, error) {
	// Instantiate a signer using RSASSA-PSS (SHA512) with the given private key.
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateKey}, nil)
	if err != nil {
		return "", err
	}

	// Sign a sample payload. Calling the signer returns a protected JWS object,
	// which can then be serialized for output afterwards. An error would
	// indicate a problem in an underlying cryptographic primitive.

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	object, err := signer.Sign(payloadBytes)
	if err != nil {
		return "", err
	}

	// Serialize the encrypted object using the full serialization format.
	// Alternatively you can also use the compact format here by calling
	// object.CompactSerialize() instead.
	serialized := object.FullSerialize()

	return serialized, nil
	// Parse the serialized, protected JWS object. An error would indicate that
	// the given input did not represent a valid message.
	//object, err = ParseSigned(serialized)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Now we can verify the signature on the payload. An error here would
	//// indicate the the message failed to verify, e.g. because the signature was
	//// broken or the message was tampered with.
	//output, err := object.Verify(&privateKey.PublicKey)
	//if err != nil {
	//	panic(err)
	//}
}

func ValidateJWT(jwt string, privateKey rsa.PrivateKey) (int64, error) {
	object, err := jose.ParseSigned(jwt)
	if err != nil {
		return 0, err
	}

	// Now we can verify the signature on the payload. An error here would
	// indicate the the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	output, err := object.Verify(&privateKey.PublicKey)
	if err != nil {
		return 0, err
	}

	payload := JWTPayload{}

	err = json.Unmarshal(output, &payload)
	if err != nil {
		return 0, err
	}

	return payload.UID, nil
}

// Must be pem encoded
func ParsePrivateKey(keypair string) (*rsa.PrivateKey, error) {
	p, _ := pem.Decode([]byte(keypair))
	if p == nil {
		return nil, errors.New("no PEM data found in ParsePrivateKey")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Returns true if equal
func compareHashedPassword(password, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	switch {
	case err == bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}
