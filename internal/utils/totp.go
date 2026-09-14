package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

// GenerateTOTP menghasilkan kode numerik sepanjang `digits` yang berubah otomatis
// setiap `period`, dihitung dari HMAC-SHA256(secret, counter-waktu). Deterministik:
// dipanggil kapan saja dalam window waktu yang sama akan selalu menghasilkan kode yang
// sama, tanpa perlu scheduler/cron untuk "merotasi" — sehingga valid dipakai bersama
// oleh banyak jadwal ujian yang berjalan bersamaan (token bersifat global, tidak terikat
// satu jadwal tertentu).
func GenerateTOTP(secret string, t time.Time, period time.Duration, digits int) string {
	counter := uint64(t.Unix() / int64(period.Seconds()))
	return computeTOTP(secret, counter, digits)
}

// ValidateTOTP mencocokkan token terhadap kode yang berlaku saat ini, dengan toleransi
// mundur sebanyak `graceSteps` window sebelumnya — mengatasi kasus token baru saja
// berganti tepat saat peserta submit.
func ValidateTOTP(secret, token string, t time.Time, period time.Duration, digits int, graceSteps int) bool {
	counter := uint64(t.Unix() / int64(period.Seconds()))
	for i := 0; i <= graceSteps && uint64(i) <= counter; i++ {
		if computeTOTP(secret, counter-uint64(i), digits) == token {
			return true
		}
	}
	return false
}

// RemainingSeconds menghitung sisa detik sebelum token saat ini berganti ke window berikutnya.
func RemainingSeconds(t time.Time, period time.Duration) int {
	periodSec := int64(period.Seconds())
	elapsed := t.Unix() % periodSec
	return int(periodSec - elapsed)
}

func computeTOTP(secret string, counter uint64, digits int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(buf)
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	code := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff

	mod := uint32(math.Pow10(digits))
	return fmt.Sprintf("%0*d", digits, code%mod)
}
