// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

const (
	attrCertificate       = "certificate"
	attrPrivateKey        = "private_key"
	attrCertificateMD5    = "certificate_md5"
	attrNotAfter          = "not_after"
	attrSubject           = "subject"
	attrFingerprintSHA256 = "fingerprint_sha256"
)

// pemCertificateInfo is the non-secret metadata derived from a PEM certificate.
type pemCertificateInfo struct {
	NotAfter          string
	Subject           string
	FingerprintSHA256 string
	MD5               string
}

func parsePEMCertificateInfo(pemData string) (*pemCertificateInfo, error) {
	rest := []byte(pemData)
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			return nil, fmt.Errorf("no PEM CERTIFICATE block found")
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse certificate: %w", err)
		}
		sum := sha256.Sum256(cert.Raw)
		md5sum := md5.Sum(cert.Raw) //nolint:gosec // stored hash for Terraform plan equality, not a security control
		return &pemCertificateInfo{
			NotAfter:          cert.NotAfter.UTC().Format(time.RFC3339),
			Subject:           cert.Subject.String(),
			FingerprintSHA256: hex.EncodeToString(sum[:]),
			MD5:               hex.EncodeToString(md5sum[:]),
		}, nil
	}
}

func looksLikePEMCertificate(s string) bool {
	return strings.Contains(s, "BEGIN CERTIFICATE")
}

func endpointCertificateHints(parentAttr, parentDesc, description, markdown string) *is.TFStateHints {
	return &is.TFStateHints{
		Importable:      &notImportable,
		WriteOnlyFields: []string{attrCertificate, attrPrivateKey},
		SensitiveFields: []string{attrCertificate, attrPrivateKey},
		TFStateHintsForCustom: &is.TFStateHintsForCustom{
			Description:         description,
			MarkdownDescription: markdown,
			SchemaAttributes: map[string]any{
				parentAttr: rschema.Int64Attribute{
					Required:    true,
					Description: parentDesc,
					PlanModifiers: []planmodifiers.Int64{
						int64planmodifier.RequiresReplace(),
					},
				},
				attrCertificate: rschema.StringAttribute{
					Required:    true,
					Sensitive:   true,
					WriteOnly:   true,
					Description: "PEM-encoded certificate (leaf plus optional chain). Write-only: never stored in Terraform state.",
				},
				attrPrivateKey: rschema.StringAttribute{
					Required:    true,
					Sensitive:   true,
					WriteOnly:   true,
					Description: "PEM-encoded private key. Write-only: never stored in Terraform state.",
				},
				attrCertificateMD5: rschema.StringAttribute{
					Computed:    true,
					Description: "MD5 of the leaf certificate DER (lowercase hex). Computed from certificate at plan time so a new certificate is applied automatically. Not settable.",
				},
				attrNotAfter: rschema.StringAttribute{
					Computed:    true,
					Description: "Not-after (expiry) of the leaf certificate, RFC3339 UTC.",
				},
				attrSubject: rschema.StringAttribute{
					Computed:    true,
					Description: "Subject of the leaf certificate.",
				},
				attrFingerprintSHA256: rschema.StringAttribute{
					Computed:    true,
					Description: "SHA-256 fingerprint of the leaf certificate DER, lowercase hex.",
				},
			},
		},
	}
}

type endpointCertificateOps struct {
	parentAttr  string
	certAPI     string
	certInfoAPI string
	apply       func(ctx context.Context, rest *VMSRest, parentID int64, cert, key string) (Record, error)
	get         func(ctx context.Context, rest *VMSRest, parentID int64) (Record, error)
}

func (o *endpointCertificateOps) parentID(ts *is.TFState) (int64, error) {
	if ts.IsKnownAndNotNull(o.parentAttr) {
		if id := ts.Int64(o.parentAttr); id != 0 {
			return id, nil
		}
	}
	return 0, fmt.Errorf("%s is required", o.parentAttr)
}

func (o *endpointCertificateOps) metadataRecord(parentID int64, pemData string, parent Record) Record {
	out := Record{
		o.parentAttr: parentID,
	}
	if looksLikePEMCertificate(pemData) {
		if info, err := parsePEMCertificateInfo(pemData); err == nil {
			out[attrNotAfter] = info.NotAfter
			out[attrSubject] = info.Subject
			out[attrFingerprintSHA256] = info.FingerprintSHA256
			out[attrCertificateMD5] = info.MD5
			return out
		}
	}
	// GET redacts PEM material; VMS still returns ssl_certificate_info /
	// s3_certificate_info / nfs4_certificate_info for drift.
	if parent != nil && o.certInfoAPI != "" {
		if info, err := parseAPICertificateInfo(parent[o.certInfoAPI]); err == nil {
			if info.NotAfter != "" {
				out[attrNotAfter] = info.NotAfter
			}
			if info.Subject != "" {
				out[attrSubject] = info.Subject
			}
			if info.FingerprintSHA256 != "" {
				out[attrFingerprintSHA256] = info.FingerprintSHA256
			}
		}
	}
	return out
}

func parseAPICertificateInfo(raw any) (*pemCertificateInfo, error) {
	obj, ok := raw.(map[string]any)
	if !ok || len(obj) == 0 {
		return nil, fmt.Errorf("certificate_info is not an object")
	}
	info := &pemCertificateInfo{
		Subject: formatAPIDistinguishedName(obj["subject"]),
	}
	if v, ok := obj["notAfter"].(string); ok && v != "" {
		if t, err := parseAPICertTime(v); err == nil {
			info.NotAfter = t.UTC().Format(time.RFC3339)
		}
	}
	if v, ok := obj["fingerprint_sha256"].(string); ok && v != "" {
		info.FingerprintSHA256 = strings.ToLower(strings.ReplaceAll(v, ":", ""))
	}
	if info.Subject == "" && info.NotAfter == "" {
		return nil, fmt.Errorf("certificate_info has no subject or notAfter")
	}
	return info, nil
}

func parseAPICertTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	layouts := []string{
		"Jan _2 15:04:05 2006 MST",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}
	var last error
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
		last = err
	}
	return time.Time{}, last
}

// formatAPIDistinguishedName converts VMS certificate_info subject RDNs into the
// same string form as x509.Certificate.Subject.String() (via pkix.Name).
// The API returns RDNs in certificate encoding order (often O,OU,CN for OpenSSL
// subjects), while Go's String() uses a canonical reverse RDN order (CN,OU,O).
// Matching that form avoids perpetual plan diffs / re-uploads after Read.
func formatAPIDistinguishedName(raw any) string {
	rdns, ok := raw.([]any)
	if !ok {
		return ""
	}
	var name pkix.Name
	for _, rdn := range rdns {
		avas, ok := rdn.([]any)
		if !ok {
			continue
		}
		for _, ava := range avas {
			pair, ok := ava.([]any)
			if !ok || len(pair) < 2 {
				continue
			}
			typ := fmt.Sprint(pair[0])
			val := fmt.Sprint(pair[1])
			switch typ {
			case "commonName":
				name.CommonName = val
			case "organizationName":
				name.Organization = append(name.Organization, val)
			case "organizationalUnitName":
				name.OrganizationalUnit = append(name.OrganizationalUnit, val)
			case "countryName":
				name.Country = append(name.Country, val)
			case "localityName":
				name.Locality = append(name.Locality, val)
			case "stateOrProvinceName":
				name.Province = append(name.Province, val)
			case "serialNumber":
				name.SerialNumber = val
			case "streetAddress":
				name.StreetAddress = append(name.StreetAddress, val)
			case "postalCode":
				name.PostalCode = append(name.PostalCode, val)
			case "emailAddress":
				name.ExtraNames = append(name.ExtraNames, pkix.AttributeTypeAndValue{
					Type:  []int{1, 2, 840, 113549, 1, 9, 1}, // emailAddress
					Value: val,
				})
			}
		}
	}
	return name.String()
}

func (o *endpointCertificateOps) Read(ctx context.Context, rest *VMSRest, ts *is.TFState) (DisplayableRecord, error) {
	parentID, err := o.parentID(ts)
	if err != nil {
		return nil, err
	}
	parent, err := o.get(ctx, rest, parentID)
	if err != nil {
		return nil, err
	}
	pemData := ""
	if parent != nil {
		if s, ok := parent[o.certAPI].(string); ok {
			pemData = s
		}
	}
	return o.metadataRecord(parentID, pemData, parent), nil
}

func (o *endpointCertificateOps) Apply(ctx context.Context, rest *VMSRest, ts *is.TFState) (DisplayableRecord, error) {
	parentID, err := o.parentID(ts)
	if err != nil {
		return nil, err
	}
	cert := ts.String(attrCertificate)
	key := ts.String(attrPrivateKey)
	if cert == "" || key == "" {
		return nil, fmt.Errorf("%s and %s are required", attrCertificate, attrPrivateKey)
	}
	if _, err := parsePEMCertificateInfo(cert); err != nil {
		return nil, fmt.Errorf("%s: %w", attrCertificate, err)
	}
	record, err := o.apply(ctx, rest, parentID, cert, key)
	if err != nil {
		return nil, err
	}
	if record != nil && !record.Empty() {
		if err := handleMaybeAsyncTask(ctx, rest, record, nil); err != nil {
			return nil, err
		}
	}
	parent, err := o.get(ctx, rest, parentID)
	if err != nil {
		// Certificate was applied; still expose metadata from the PEM we sent.
		return o.metadataRecord(parentID, cert, nil), nil
	}
	pemData := cert
	if parent != nil {
		if s, ok := parent[o.certAPI].(string); ok && looksLikePEMCertificate(s) {
			pemData = s
		}
	}
	return o.metadataRecord(parentID, pemData, parent), nil
}

func (o *endpointCertificateOps) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	if req.Plan.Raw.IsNull() {
		return nil
	}
	set := func(name string, v attr.Value) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root(name), v)...)
	}
	var cert types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(attrCertificate), &cert)...)
	if resp.Diagnostics.HasError() {
		return nil
	}
	if cert.IsUnknown() {
		set(attrCertificateMD5, types.StringUnknown())
		set(attrNotAfter, types.StringUnknown())
		set(attrSubject, types.StringUnknown())
		set(attrFingerprintSHA256, types.StringUnknown())
		return nil
	}
	if cert.IsNull() || cert.ValueString() == "" {
		return nil
	}
	info, err := parsePEMCertificateInfo(cert.ValueString())
	if err != nil {
		return err
	}
	set(attrCertificateMD5, types.StringValue(info.MD5))
	set(attrNotAfter, types.StringValue(info.NotAfter))
	set(attrSubject, types.StringValue(info.Subject))
	set(attrFingerprintSHA256, types.StringValue(info.FingerprintSHA256))
	return nil
}

func newResourceTFState(raw map[string]attr.Value, schema any, hints *is.TFStateHints) *is.TFState {
	return is.NewTFStateMust(raw, schema, hints)
}
