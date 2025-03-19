package utils

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"
)

var CurrentVersion string

func IfElse[T any](cond bool, t T, f T) T {
	if cond {
		return t
	}
	return f
}

// Ptr returns a pointer to the value passed in.
// The pointer is to a shallow copy, not the original value.
func Ptr[T any](v T) *T {
	var val = v
	return &val
}

func GenerateCode() string {
	return fmt.Sprintf("%d", (1+rand.Intn(10))*10000+rand.Intn(10000))
}

func HashPassword(password string, salt string) (string, error) {
	preHashed := fmt.Sprintf("%s-%s", password, salt)
	h := sha256.New()
	_, err := h.Write([]byte(preHashed))
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)
	hashed := fmt.Sprintf("%x", bs)
	return hashed, nil
}

func GenerateHashAndSalt(password string) (string, string, error) {
	salt := fmt.Sprintf("salt-%d", rand.Int31())
	hashedPassword, err := HashPassword(password, salt)
	if err != nil {
		return "", "", err
	}
	return salt, hashedPassword, nil
}

func JSONConvert[T any, U any](v T, u *U) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, "JSONConvert failed")
	}
	return json.Unmarshal(raw, u)
}

func TryMarshal(o any) string {
	raw, err := json.Marshal(o)
	if err != nil {
		return "<unmarshalable>"
	}
	return string(raw)
}

func UnwrapOrDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

func TestTCPConnection(ctx context.Context, host string, port int32, timeout time.Duration) error {
	var d net.Dialer
	d.Timeout = timeout

	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func TruncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
