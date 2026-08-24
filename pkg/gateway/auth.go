package gateway

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	gatewayTokenFile    = ".icx-gateway-token"
	defaultMaxBodyBytes = int64(1 << 20) // 1 MiB
	defaultMaxTools     = 256
)

// IsLoopbackHost reports whether host is a loopback bind target.
func IsLoopbackHost(host string) bool {
	h := strings.TrimSpace(host)
	if h == "" || h == "127.0.0.1" || h == "localhost" || h == "::1" {
		return true
	}
	ip := net.ParseIP(h)
	return ip != nil && ip.IsLoopback()
}

// GenerateGatewayToken returns a new icxgw_ token.
func GenerateGatewayToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "icxgw_" + hex.EncodeToString(buf), nil
}

// ResolveGatewayToken returns a token from env, flag, existing file, or a newly generated file.
func ResolveGatewayToken(explicit string) (token string, source string, err error) {
	if env := strings.TrimSpace(os.Getenv("ICX_GATEWAY_TOKEN")); env != "" {
		return env, "ICX_GATEWAY_TOKEN", nil
	}
	if strings.TrimSpace(explicit) != "" {
		return strings.TrimSpace(explicit), "flag", nil
	}
	if data, readErr := os.ReadFile(gatewayTokenFile); readErr == nil {
		tok := strings.TrimSpace(string(data))
		if tok != "" {
			return tok, gatewayTokenFile, nil
		}
	}
	tok, genErr := GenerateGatewayToken()
	if genErr != nil {
		return "", "", genErr
	}
	if writeErr := os.WriteFile(gatewayTokenFile, []byte(tok+"\n"), 0600); writeErr != nil {
		return "", "", fmt.Errorf("failed to persist gateway token: %w", writeErr)
	}
	return tok, "generated:" + gatewayTokenFile, nil
}

func tokenEqual(got, want string) bool {
	if want == "" {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1 {
		return true
	}
	return false
}

func extractBearer(r *http.Request) string {
	if tok := strings.TrimSpace(r.Header.Get("X-ICX-Gateway-Token")); tok != "" {
		return tok
	}
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return ""
}

func isLocalhostOrigin(origin string) bool {
	o := strings.ToLower(strings.TrimSpace(origin))
	return strings.HasPrefix(o, "http://127.0.0.1") || strings.HasPrefix(o, "http://localhost") ||
		strings.HasPrefix(o, "https://127.0.0.1") || strings.HasPrefix(o, "https://localhost")
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error": {"message": %q, "type": "invalid_request_error"}}`, msg)))
}
