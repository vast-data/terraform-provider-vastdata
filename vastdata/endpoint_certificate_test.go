// Copyright (c) HashiCorp, Inc.
package provider

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

func testGeneratePEMCert(t *testing.T) (certPEM, keyPEM string, info *pemCertificateInfo) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "tf-test.example.com", Organization: []string{"VastData"}},
		NotBefore:    time.Now().Add(-time.Hour).UTC(),
		NotAfter:     time.Now().Add(24 * time.Hour).UTC(),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)

	certPEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	keyPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))

	parsed, err := parsePEMCertificateInfo(certPEM)
	require.NoError(t, err)
	return certPEM, keyPEM, parsed
}

func TestParsePEMCertificateInfo(t *testing.T) {
	certPEM, _, want := testGeneratePEMCert(t)

	got, err := parsePEMCertificateInfo(certPEM)
	require.NoError(t, err)
	assert.Equal(t, want.Subject, got.Subject)
	assert.Equal(t, want.FingerprintSHA256, got.FingerprintSHA256)
	assert.Equal(t, want.NotAfter, got.NotAfter)
	assert.Equal(t, want.MD5, got.MD5)
	assert.Contains(t, got.Subject, "tf-test.example.com")
	assert.Len(t, got.FingerprintSHA256, 64)
	assert.Len(t, got.MD5, 32)
	assert.NotEqual(t, got.MD5, got.FingerprintSHA256)
}

func TestParsePEMCertificateInfo_ChainUsesLeaf(t *testing.T) {
	leaf, _, leafInfo := testGeneratePEMCert(t)
	ca, _, _ := testGeneratePEMCert(t)

	got, err := parsePEMCertificateInfo(leaf + ca)
	require.NoError(t, err)
	assert.Equal(t, leafInfo.FingerprintSHA256, got.FingerprintSHA256)
}

func TestParsePEMCertificateInfo_RejectsNonCert(t *testing.T) {
	_, err := parsePEMCertificateInfo("not a cert")
	require.Error(t, err)
}

func TestLooksLikePEMCertificate(t *testing.T) {
	assert.True(t, looksLikePEMCertificate("-----BEGIN CERTIFICATE-----\nMII\n-----END CERTIFICATE-----"))
	assert.False(t, looksLikePEMCertificate("*****"))
	assert.False(t, looksLikePEMCertificate(""))
}

func TestEndpointCertificateOps_MetadataRecord(t *testing.T) {
	o := &endpointCertificateOps{parentAttr: "vms_id", certAPI: "ssl_certificate", certInfoAPI: "ssl_certificate_info"}
	certPEM, _, info := testGeneratePEMCert(t)

	rec := o.metadataRecord(7, certPEM, nil)
	assert.Equal(t, int64(7), rec["vms_id"])
	_, hasID := rec["id"]
	assert.False(t, hasID)
	assert.Equal(t, info.NotAfter, rec[attrNotAfter])
	assert.Equal(t, info.Subject, rec[attrSubject])
	assert.Equal(t, info.FingerprintSHA256, rec[attrFingerprintSHA256])
	assert.Equal(t, info.MD5, rec[attrCertificateMD5])
	assert.Len(t, info.MD5, 32)

	redacted := o.metadataRecord(7, "*****", nil)
	assert.Equal(t, int64(7), redacted["vms_id"])
	_, hasFP := redacted[attrFingerprintSHA256]
	assert.False(t, hasFP)
	_, hasMD5 := redacted[attrCertificateMD5]
	assert.False(t, hasMD5)

	fromAPI := o.metadataRecord(7, "***", Record{
		"ssl_certificate_info": map[string]any{
			"notAfter": "Oct  6 08:12:43 2027 GMT",
			"subject": []any{
				[]any{[]any{"commonName", "l3118-vms"}},
				[]any{[]any{"organizationName", "VastData Test"}},
				[]any{[]any{"countryName", "US"}},
			},
		},
	})
	assert.Equal(t, "CN=l3118-vms,O=VastData Test,C=US", fromAPI[attrSubject])
	assert.Equal(t, "2027-10-06T08:12:43Z", fromAPI[attrNotAfter])
}

// TERF-310: VMS returns subject RDNs in certificate encoding order (O,OU,CN for
// OpenSSL multi-attribute subjects), while PEM parsing uses pkix.Name.String()
// (CN,OU,O). Read must canonicalize to the PEM form or every plan re-uploads.
func TestFormatAPIDistinguishedName_MatchesPEMSubjectOrder(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       []string{"VAST Data QA"},
			OrganizationalUnit: []string{"Terraform"},
			CommonName:         "tf-s3-cert-a",
		},
		NotBefore: time.Now().Add(-time.Hour).UTC(),
		NotAfter:  time.Now().Add(24 * time.Hour).UTC(),
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	certPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))

	pemInfo, err := parsePEMCertificateInfo(certPEM)
	require.NoError(t, err)
	assert.Equal(t, "CN=tf-s3-cert-a,OU=Terraform,O=VAST Data QA", pemInfo.Subject)

	// API order mirrors OpenSSL encoding order (O, then OU, then CN).
	apiSubject := formatAPIDistinguishedName([]any{
		[]any{[]any{"organizationName", "VAST Data QA"}},
		[]any{[]any{"organizationalUnitName", "Terraform"}},
		[]any{[]any{"commonName", "tf-s3-cert-a"}},
	})
	assert.Equal(t, pemInfo.Subject, apiSubject)
	assert.NotEqual(t, "O=VAST Data QA,OU=Terraform,CN=tf-s3-cert-a", apiSubject)
}

func TestEndpointCertificateOps_ParentID(t *testing.T) {
	schema := rschema.Schema{
		Attributes: map[string]rschema.Attribute{
			"cluster_id": rschema.Int64Attribute{Required: true},
		},
	}
	o := &endpointCertificateOps{parentAttr: "cluster_id"}

	ts := is.NewTFStateMust(map[string]attr.Value{
		"cluster_id": types.Int64Value(3),
	}, schema, nil)
	id, err := o.parentID(ts)
	require.NoError(t, err)
	assert.Equal(t, int64(3), id)

	ts = is.NewTFStateMust(map[string]attr.Value{
		"cluster_id": types.Int64Null(),
	}, schema, nil)
	_, err = o.parentID(ts)
	require.Error(t, err)
}

func TestSnakeCaseEndpointCertificateResources(t *testing.T) {
	assert.Equal(t, "vms_certificate", is.SnakeCaseName(&VmsCertificate{}))
	assert.Equal(t, "s3_certificate", is.SnakeCaseName(&S3Certificate{}))
	assert.Equal(t, "nfs_certificate", is.SnakeCaseName(&NfsCertificate{}))
}
