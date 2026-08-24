package gateway

import "testing"

func TestIsLoopbackHost(t *testing.T) {
	for _, h := range []string{"127.0.0.1", "localhost", "::1", ""} {
		if !IsLoopbackHost(h) {
			t.Errorf("expected loopback: %q", h)
		}
	}
	for _, h := range []string{"0.0.0.0", "10.0.0.1", "192.168.1.5"} {
		if IsLoopbackHost(h) {
			t.Errorf("did not expect loopback: %q", h)
		}
	}
}

func TestGenerateGatewayToken(t *testing.T) {
	a, err := GenerateGatewayToken()
	if err != nil {
		t.Fatal(err)
	}
	b, err := GenerateGatewayToken()
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("tokens should be unique")
	}
	if len(a) < 20 || a[:6] != "icxgw_" {
		t.Fatalf("unexpected token shape: %s", a)
	}
}
