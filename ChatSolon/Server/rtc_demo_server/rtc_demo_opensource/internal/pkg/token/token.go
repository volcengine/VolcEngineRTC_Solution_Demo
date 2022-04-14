package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"time"
)

const (
	Version       = "001"
	VersionLength = 3

	AppIDLength = 24
)

type Privilege uint16

const (
	PrivPublishStream Privilege = iota

	// not exported, do not use directly
	privPublishAudioStream
	privPublishVideoStream
	privPublishDataStream

	PrivSubscribeStream
)

type Token struct {
	AppID      string
	AppKey     string
	RoomID     string
	UserID     string
	IssuedAt   uint32
	ExpireAt   uint32
	Nonce      uint32
	Privileges map[uint16]uint32
	Signature  string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

// New initializes token struct by required parameters.
func New(appID, appKey, roomID, userID string) *Token {
	return &Token{
		AppID:      appID,
		AppKey:     appKey,
		RoomID:     roomID,
		UserID:     userID,
		IssuedAt:   uint32(time.Now().Unix()),
		Nonce:      uint32(random(1, 99999999)),
		Privileges: make(map[uint16]uint32),
	}
}

// Parse retrieves token information from raw string
func Parse(raw string) (*Token, error) {
	// length check
	if len(raw) <= VersionLength+AppIDLength {
		return nil, fmt.Errorf("invalid token length: %d", len(raw))
	}
	if raw[:VersionLength] != Version {
		return nil, fmt.Errorf("expect version: %s, got %s", Version, raw[:VersionLength])
	}

	t := new(Token)

	// app id
	t.AppID = raw[VersionLength : VersionLength+AppIDLength]

	// parse signature + msg from content
	contentEncoded := raw[VersionLength+AppIDLength:]
	content, err := base64.StdEncoding.DecodeString(contentEncoded)
	if err != nil {
		return nil, errors.New("failed to decode token by content: " + contentEncoded)
	}
	msg, signature, err := unPackContent(content)
	if err != nil {
		return nil, errors.New("failed to unpack content")
	}
	t.Signature = signature

	// parse from msg
	in := bytes.NewReader([]byte(msg))
	t.Privileges = make(map[uint16]uint32)
	t.Nonce, err = unPackUint32(in)
	if err != nil {
		return nil, errors.New("failed to unpack nonce")
	}
	t.IssuedAt, err = unPackUint32(in)
	if err != nil {
		return nil, errors.New("failed to unpack issuedAt")
	}
	t.ExpireAt, err = unPackUint32(in)
	if err != nil {
		return nil, errors.New("failed to unpack expireAt")
	}
	t.RoomID, err = unPackString(in)
	if err != nil {
		return nil, errors.New("failed to unpack room id")
	}
	t.UserID, err = unPackString(in)
	if err != nil {
		return nil, errors.New("failed to unpack user id")
	}
	keyLength, err := unPackUint16(in)
	if err != nil {
		return nil, errors.New("failed to unpack key length")
	}
	for i := uint16(0); i < keyLength; i++ {
		key, err := unPackUint16(in)
		if err != nil {
			return nil, errors.New("failed to unpack privilege key")
		}
		value, err := unPackUint32(in)
		if err != nil {
			return nil, errors.New("failed to unpack privilege value")
		}
		t.Privileges[key] = value
	}

	return t, nil
}

// Verify checks if this token valid, called by server side.
func (t *Token) Verify(key string) bool {
	if t.ExpireAt > 0 && uint32(time.Now().Unix()) > t.ExpireAt {
		return false
	}

	t.AppKey = key
	_, sign, err := t.pack()
	if err != nil {
		return false
	}
	return sign == t.Signature
}

// AddPrivilege adds permission for token with an expiration.
func (t *Token) AddPrivilege(p Privilege, expireAt time.Time) {
	if t.Privileges == nil {
		t.Privileges = make(map[uint16]uint32)
	}
	expire := uint32(expireAt.Unix())
	if expireAt.IsZero() {
		expire = 0
	}
	t.Privileges[uint16(p)] = expire
	// add separated publish privileges for now
	if p == PrivPublishStream {
		t.Privileges[uint16(privPublishVideoStream)] = expire
		t.Privileges[uint16(privPublishAudioStream)] = expire
		t.Privileges[uint16(privPublishDataStream)] = expire
	}
}

// ExpireTime sets token expire time, won't expire by default.
// The token will be invalid after expireTime no matter what privilege's expireTime is.
func (t *Token) ExpireTime(et time.Time) {
	if !et.IsZero() {
		t.ExpireAt = uint32(et.Unix())
	}
}

func (t *Token) pack() (string, string, error) {
	bufM := new(bytes.Buffer)
	if err := packUint32(bufM, t.Nonce); err != nil {
		return "", "", err
	}
	if err := packUint32(bufM, t.IssuedAt); err != nil {
		return "", "", err
	}
	if err := packUint32(bufM, t.ExpireAt); err != nil {
		return "", "", err
	}
	if err := packString(bufM, t.RoomID); err != nil {
		return "", "", err
	}
	if err := packString(bufM, t.UserID); err != nil {
		return "", "", err
	}
	if err := packMapUint32(bufM, t.Privileges); err != nil {
		return "", "", err
	}
	bytesM := bufM.Bytes()

	bufSign := hmac.New(sha256.New, []byte(t.AppKey))
	bufSign.Write(bytesM)
	bytesSign := bufSign.Sum(nil)

	return string(bytesM[:]), string(bytesSign[:]), nil
}

// Serialize generates the token string
func (t *Token) Serialize() (string, error) {
	msg, sign, err := t.pack()
	if err != nil {
		return "", err
	}

	bufContent := new(bytes.Buffer)
	if err := packString(bufContent, msg); err != nil {
		return "", err
	}
	if err := packString(bufContent, sign); err != nil {
		return "", err
	}
	bytesContent := bufContent.Bytes()

	return Version + t.AppID + base64.StdEncoding.EncodeToString(bytesContent), nil
}

func packUint16(w io.Writer, n uint16) error {
	return binary.Write(w, binary.LittleEndian, n)
}

func packUint32(w io.Writer, n uint32) error {
	return binary.Write(w, binary.LittleEndian, n)
}

func packString(w io.Writer, s string) error {
	err := packUint16(w, uint16(len(s)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(s))
	return err
}

func packMapUint32(w io.Writer, extra map[uint16]uint32) error {
	var keys []int
	if err := packUint16(w, uint16(len(extra))); err != nil {
		return err
	}
	for k := range extra {
		keys = append(keys, int(k))
	}
	//should sorted keys
	sort.Ints(keys)

	for _, k := range keys {
		v := extra[uint16(k)]
		if err := packUint16(w, uint16(k)); err != nil {
			return err
		}
		if err := packUint32(w, v); err != nil {
			return err
		}
	}
	return nil
}

func unPackUint16(r io.Reader) (uint16, error) {
	var n uint16
	err := binary.Read(r, binary.LittleEndian, &n)
	return n, err
}

func unPackUint32(r io.Reader) (uint32, error) {
	var n uint32
	err := binary.Read(r, binary.LittleEndian, &n)
	return n, err
}

func unPackString(r io.Reader) (string, error) {
	n, err := unPackUint16(r)
	if err != nil {
		return "", err
	}

	buf := make([]byte, n)
	r.Read(buf)
	s := string(buf[:])
	return s, err
}

func unPackContent(buff []byte) (string, string, error) {
	in := bytes.NewReader(buff)
	msg, err := unPackString(in)
	if err != nil {
		return "", "", err
	}

	sig, err := unPackString(in)
	if err != nil {
		return "", "", err
	}
	return msg, sig, nil
}

type GenerateParam struct {
	AppID        string
	AppKey       string
	RoomID       string
	UserID       string
	ExpireAt     int64
	CanPublish   bool
	CanSubscribe bool
}

func GenerateToken(param *GenerateParam) (string, error) {
	expiredAt := param.ExpireAt
	if expiredAt < 1609500000 {
		expiredAt = time.Now().Unix() + param.ExpireAt
	}

	token := New(param.AppID, param.AppKey, param.RoomID, param.UserID)

	token.ExpireTime(time.Unix(expiredAt, 0))

	if param.CanPublish {
		token.AddPrivilege(PrivPublishStream, time.Unix(expiredAt, 0))
	}

	if param.CanSubscribe {
		token.AddPrivilege(PrivSubscribeStream, time.Unix(expiredAt, 0))
	}

	return token.Serialize()
}
