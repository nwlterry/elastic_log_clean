#!/usr/bin/env python3
import sys
import tempfile
import unittest
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT))

from elastic_log_clean import Cleaner, default_config, main  # noqa: E402


class CleanerTests(unittest.TestCase):
    def setUp(self) -> None:
        self.cleaner = Cleaner(default_config())

    def test_consistent_ip_and_skip_loopback(self) -> None:
        text = "a 10.20.30.41 b 10.20.30.41 c 127.0.0.1 d 10.20.30.99 version 8.18.4"
        out = self.cleaner.obfuscate_text(text)
        self.assertIn("127.0.0.1", out)
        self.assertIn("8.18.4", out)
        self.assertNotIn("10.20.30.41", out)
        self.assertNotIn("10.20.30.99", out)
        token_a = out.split()[1]
        self.assertTrue(token_a.startswith("x-ipv4-"))
        self.assertEqual(out.split()[3], token_a)

    def test_secrets_and_pem(self) -> None:
        sample = (ROOT / "examples" / "sample-elasticsearch.log").read_text(encoding="utf-8")
        out = self.cleaner.obfuscate_text(sample)
        self.assertNotIn("s3cretKeystore", out)
        self.assertNotIn("Sup3rS3cret", out)
        self.assertNotIn("LdapBind#99", out)
        self.assertNotIn("Changeme123", out)
        self.assertNotIn("aGhpc2lzYWZha2VhcGlrZXlmb3J0ZXN0", out)
        self.assertNotIn("BEGIN RSA PRIVATE KEY", out)
        self.assertIn("x-redacted-secret-x", out)
        self.assertIn("x-redacted-pem-x", out)
        self.assertNotIn("sre@corp.example.com", out)
        self.assertIn("x-email-", out)

    def test_omit_keystore(self) -> None:
        self.assertTrue(self.cleaner.should_omit("nodes/es1/elasticsearch.keystore"))
        self.assertTrue(self.cleaner.should_omit("certs/http.p12"))
        self.assertFalse(self.cleaner.should_omit("logs/prod.log"))

    def test_zip_roundtrip(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            tmp_path = Path(tmp)
            src_zip = tmp_path / "local-diagnostics-2026.zip"
            with zipfile.ZipFile(src_zip, "w") as zf:
                zf.writestr("logs/elasticsearch.log", "node 10.1.2.3 password=abc123\n")
                zf.writestr("certs/http.p12", b"\x00binary")
            out_zip = tmp_path / "out.zip"
            report = tmp_path / "report.yaml"
            rc = main(["-i", str(src_zip), "-o", str(out_zip), "-r", str(report), "-w", "1"])
            self.assertEqual(rc, 0)
            with zipfile.ZipFile(out_zip) as zf:
                names = zf.namelist()
                self.assertTrue(any(n.endswith("elasticsearch.log") for n in names))
                self.assertFalse(any(n.endswith(".p12") for n in names))
                body = zf.read([n for n in names if n.endswith(".log")][0]).decode()
            self.assertNotIn("10.1.2.3", body)
            self.assertNotIn("abc123", body)


if __name__ == "__main__":
    unittest.main()
