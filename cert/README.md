# cert

Reads a PEM-encoded X.509 certificate from disk and prints its main fields, using only `encoding/pem` and `crypto/x509` from the standard library.

## What the code does

`main.go` reads `test.crt`, decodes the first PEM block, and parses it with `x509.ParseCertificate`. It then prints these fields:

- Issuer
- Subject
- Subject Common Name
- NotBefore and NotAfter
- Serial number

It panics if the file is missing, has no PEM block, or fails to parse.

The bundled `test.crt` is self-signed (Issuer == Subject: `CN=fenxx.com,O=FenXX,L=Taipei,ST=Some-State,C=TW,1.2.840.113549.1.9.1=xxxxx@xxxxx.com`; Go prints the emailAddress attribute as its raw OID) and valid from 2026-09-15 to 2026-10-15 UTC.

## Concepts you can learn

- **PEM vs DER**: PEM is a base64 text wrapper (`-----BEGIN CERTIFICATE-----`) around DER bytes. `pem.Decode` removes the wrapper, and `x509` parses the DER bytes.
- **Certificate identity fields**: Issuer and Subject are distinguished names. When they're equal, the certificate is self-signed.
- **Validity window and serial number**: these are the fields to check first when debugging TLS errors.

## How to run

```bash
cd cert
go run .
```

Run it from inside `cert/`, because `test.crt` is opened with a relative path.

## Notes

- The code only parses the certificate. It doesn't verify the signature, chain, or validity dates, so it prints the details of an expired certificate (like the bundled one, which expired 2026-10-15) without complaint.
- Only the first PEM block is read. Any later certificates in a bundle are ignored.
