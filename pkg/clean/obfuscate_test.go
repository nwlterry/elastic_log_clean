package clean

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConsistentIPAndLoopback(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	out := c.ObfuscateText("a 10.20.30.41 b 10.20.30.41 c 127.0.0.1 d 10.20.30.99 version 8.18.4")
	if !strings.Contains(out, "127.0.0.1") {
		t.Fatalf("loopback rewritten: %s", out)
	}
	if !strings.Contains(out, "8.18.4") {
		t.Fatalf("version rewritten: %s", out)
	}
	if strings.Contains(out, "10.20.30.41") || strings.Contains(out, "10.20.30.99") {
		t.Fatalf("ip left in place: %s", out)
	}
	parts := strings.Fields(out)
	if !strings.HasPrefix(parts[1], "x-ipv4-") {
		t.Fatalf("token: %s", parts[1])
	}
	if parts[1] != parts[3] {
		t.Fatalf("inconsistent tokens %s vs %s", parts[1], parts[3])
	}
}

func TestSecretsAndPEM(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	samplePath := filepath.Join("..", "..", "examples", "sample-elasticsearch.log")
	b, err := os.ReadFile(samplePath)
	if err != nil {
		t.Fatal(err)
	}
	out := c.ObfuscateText(string(b))
	for _, bad := range []string{"s3cretKeystore", "Sup3rS3cret", "LdapBind#99", "Changeme123", "aGhpc2lzYWZha2VhcGlrZXlmb3J0ZXN0", "BEGIN RSA PRIVATE KEY", "sre@corp.example.com"} {
		if strings.Contains(out, bad) {
			t.Fatalf("leaked %q in:\n%s", bad, out)
		}
	}
}

func TestOmitKeystore(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	if !c.ShouldOmit("nodes/es1/elasticsearch.keystore") {
		t.Fatal("keystore not omitted")
	}
	if !c.ShouldOmit("certs/http.p12") {
		t.Fatal("p12 not omitted")
	}
	if c.ShouldOmit("logs/prod.log") {
		t.Fatal("log omitted")
	}
}

func TestModuleCoordNotEmailAndClockNotIPv6(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	in := `[2026-09-24T07:29:19,123][INFO] loaded io.netty.transport@4.1.135.Final org.elasticsearch.server@8.18.4 "07:29:19" user sre@corp.example.com ipv6 2001:db8::1`
	out := c.ObfuscateText(in)
	for _, keep := range []string{
		"io.netty.transport@4.1.135.Final",
		"org.elasticsearch.server@8.18.4",
		"07:29:19",
		"2026-09-24T07:29:19,123",
	} {
		if !strings.Contains(out, keep) {
			t.Fatalf("false positive, lost %q in:\n%s", keep, out)
		}
	}
	if strings.Contains(out, "sre@corp.example.com") {
		t.Fatalf("real email not redacted:\n%s", out)
	}
	if strings.Contains(out, "2001:db8::1") {
		t.Fatalf("real ipv6 not redacted:\n%s", out)
	}
	if !strings.Contains(out, "x-email-") {
		t.Fatalf("missing email token:\n%s", out)
	}
	if !strings.Contains(out, "x-ipv6-") && !strings.Contains(out, "x:x:x") {
		t.Fatalf("missing ipv6 token:\n%s", out)
	}
}
