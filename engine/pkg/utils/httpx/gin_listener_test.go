package httpx

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStartBackgroundTLS(t *testing.T) {
	fixture := httptest.NewTLSServer(http.NotFoundHandler())
	t.Cleanup(fixture.Close)
	cert := fixture.TLS.Certificates[0]
	key, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	files := &CAFile{CertFile: filepath.Join(dir, "cert.pem"), KeyFile: filepath.Join(dir, "key.pem")}
	for path, block := range map[string]*pem.Block{
		files.CertFile: {Type: "CERTIFICATE", Bytes: cert.Certificate[0]},
		files.KeyFile:  {Type: "PRIVATE KEY", Bytes: key},
	} {
		if err := os.WriteFile(path, pem.EncodeToMemory(block), 0600); err != nil {
			t.Fatal(err)
		}
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	gs := &GinServer{
		conf: &Conf{CAFile: &CAFile{CertFile: files.CertFile, KeyFile: filepath.Join(dir, "missing.pem")}},
		server: &http.Server{Addr: addr, Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})},
	}
	if err := gs.StartBackground(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing key: got %v, want file not found", err)
	}
	gs.conf.CAFile = files
	if err := gs.StartBackground(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(gs.Stop)
	client := fixture.Client()
	client.Timeout = 3 * time.Second
	resp, err := client.Get("https://" + addr)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent || resp.TLS == nil {
		t.Fatalf("expected HTTPS 204, got status %d with TLS=%t", resp.StatusCode, resp.TLS != nil)
	}
}
