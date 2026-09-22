package clean

import (
	"strings"
	"testing"
)

func TestConsistentIPAndLoopback(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	out := c.ObfuscateText("a 10.20.30.41 b 10.20.30.41 c 127.0.0.1 d 10.20.30.99 version 8.18.4")
	if !strings.Contains(out, "127.0.0.1") || !strings.Contains(out, "8.18.4") {
		t.Fatalf("loopback/version rewritten: %s", out)
	}
	if strings.Contains(out, "10.20.30.41") || strings.Contains(out, "10.20.30.99") {
		t.Fatalf("ip left in place: %s", out)
	}
	parts := strings.Fields(out)
	if parts[1] != parts[3] || !strings.HasPrefix(parts[1], "x-ipv4-") {
		t.Fatalf("inconsistent tokens: %s", out)
	}
}

func TestSecrets(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	out := c.ObfuscateText("password=s3cretKeystore ApiKey abcdefghijklmnop cloud.auth=elastic:Sup3rS3cret")
	if strings.Contains(out, "s3cretKeystore") || strings.Contains(out, "abcdefghijklmnop") || strings.Contains(out, "Sup3rS3cret") {
		t.Fatalf("leaked: %s", out)
	}
}

func TestOmit(t *testing.T) {
	c := NewCleaner(DefaultConfig(), nil)
	if !c.ShouldOmit("nodes/es1/elasticsearch.keystore") || !c.ShouldOmit("certs/http.p12") || c.ShouldOmit("logs/prod.log") {
		t.Fatal("omit rules")
	}
}
