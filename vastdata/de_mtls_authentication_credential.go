// Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"gopkg.in/yaml.v3"
)

var DeMtlsAuthenticationCredentialSchemaRef = is.NewDataEngineSchemaReference(
	http.MethodPost,
	"mtls-authentication-credentials",
	http.MethodGet,
	"mtls-authentication-credentials",
)

type DeMtlsAuthenticationCredential struct {
	tfstate is.ResourceState
}

func (m *DeMtlsAuthenticationCredential) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &DeMtlsAuthenticationCredential{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeMtlsAuthenticationCredentialSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			// List ignores name/vrn/guid; Read* uses List + local filter.
			SearchableFields: []string{"guid", "vrn", "name"},
			// Cert PEM fields are required by OpenAPI but optional in TF when kubeconfig_data is used.
			NotRequiredSchemaFields: []string{
				"certificate_authority_b64",
				"client_certificate_b64",
				"client_key_b64",
			},
			SensitiveFields: []string{
				"certificate_authority_b64",
				"client_certificate_b64",
				"client_key_b64",
				"kubeconfig_data",
			},
			WriteOnlyFields: []string{"kubeconfig_data"},
			CreateOnlyFields: []string{
				"certificate_authority_b64",
				"client_certificate_b64",
				"client_key_b64",
				deTenantNameAttr,
				"context_name",
			},
			PreserveUserValueFields: []string{deTenantNameAttr, "context_name"},
			CommonModifiersMapping: map[string]string{
				"certificate_authority_b64": ModifierForceNew,
				"client_certificate_b64":    ModifierForceNew,
				"client_key_b64":            ModifierForceNew,
				deTenantNameAttr:            ModifierForceNew,
				"context_name":              ModifierForceNew,
			},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
				"kubeconfig_data": rschema.StringAttribute{
					Optional:  true,
					Sensitive: true,
					WriteOnly: true,
					Description: "Raw kubeconfig YAML. Alternative to certificate_authority_b64 / " +
						"client_certificate_b64 / client_key_b64. Parsed at create time only; never stored in state.",
				},
				"context_name": rschema.StringAttribute{
					Optional: true,
					Description: "Kubeconfig context to use with kubeconfig_data. Defaults to current-context. " +
						"Ignored when certificate_*_b64 fields are set directly.",
					PlanModifiers: []planmodifiers.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
		},
	)}
}

func (m *DeMtlsAuthenticationCredential) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &DeMtlsAuthenticationCredential{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            DeMtlsAuthenticationCredentialSchemaRef,
			PathIDField:          "guid",
			ExcludedSchemaFields: []string{"cursor", "limit", "order_by"},
			SearchableFields:     []string{"guid", "vrn", "name"},
			AdditionalSchemaAttributes: map[string]any{
				deTenantNameAttr: deTenantNameSchemaAttribute(),
			},
		},
	)}
}

func (m *DeMtlsAuthenticationCredential) TransformQueryParams(query params, rest *VMSRest) params {
	return deTransformQueryParams(query, m.tfstate, rest)
}

func (m *DeMtlsAuthenticationCredential) TransformResponseRecord(response Record) Record {
	return ensureTenantNameInRecord(response, m.tfstate, nil)
}

// TransformRequestBody drops provider-only tenant_name and kubeconfig fields from the API body.
func (m *DeMtlsAuthenticationCredential) TransformRequestBody(body params) params {
	body = deTransformRequestBody(body)
	if body == nil {
		return nil
	}
	delete(body, "kubeconfig_data")
	delete(body, "context_name")
	return body
}

func (m *DeMtlsAuthenticationCredential) TfState() is.ResourceState {
	return m.tfstate
}

func (m *DeMtlsAuthenticationCredential) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.DataEngine.MtlsAuthenticationCredentials
}

func (m *DeMtlsAuthenticationCredential) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "mtls-authentication-credentials")
}

func (m *DeMtlsAuthenticationCredential) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "mtls-authentication-credentials")
}

func (m *DeMtlsAuthenticationCredential) LookupForCreate(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return findRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "mtls-authentication-credentials")
}

func (m *DeMtlsAuthenticationCredential) DeleteResource(ctx context.Context, rest *VMSRest) error {
	return deleteRecordByLocalFilter(ctx, m.API(rest), m.tfstate, m.TransformQueryParams(nil, rest), "mtls-authentication-credentials")
}

func (m *DeMtlsAuthenticationCredential) ValidateResourceConfig(context.Context) error {
	ts := m.tfstate
	if ts == nil {
		return nil
	}

	const (
		kubeconfigField = "kubeconfig_data"
		caField         = "certificate_authority_b64"
		certField       = "client_certificate_b64"
		keyField        = "client_key_b64"
	)
	pemFields := []string{caField, certField, keyField}

	// Prefer Raw presence: schema HasAttribute is true even when the key is absent from Raw
	// (e.g. unit tests / partial configs), and IsKnownAndNotNull panics on missing keys.
	fieldSet := func(name string) bool {
		v, ok := ts.RawValues()[name]
		return ok && !v.IsNull() && !v.IsUnknown()
	}

	// Unknown values: skip until known (write-only / count=0 / deferred).
	for _, f := range append([]string{kubeconfigField, "context_name"}, pemFields...) {
		if v, ok := ts.RawValues()[f]; ok && v.IsUnknown() {
			return nil
		}
	}

	hasKC := fieldSet(kubeconfigField)
	anyPEM := false
	for _, f := range pemFields {
		if fieldSet(f) {
			anyPEM = true
			break
		}
	}

	// Nothing configured yet (terraform validate without vars, or resource count=0).
	if !hasKC && !anyPEM {
		return nil
	}

	// Exactly one mode: kubeconfig_data XOR any PEM field.
	// When both modes are present, validateOneOf reports the conflict.
	if hasKC && anyPEM {
		return validateOneOf(ts, kubeconfigField, caField, certField, keyField)
	}

	// PEM mode: all three *_b64 fields are required.
	if !hasKC {
		if err := validateAllOf(ts, caField, certField, keyField); err != nil {
			return err
		}
	}

	if fieldSet("context_name") && !hasKC {
		return fmt.Errorf("context_name requires kubeconfig_data")
	}
	return nil
}

// PrepareCreateResource derives certificate_*_b64 from kubeconfig_data when set.
func (m *DeMtlsAuthenticationCredential) PrepareCreateResource(context.Context, *VMSRest) error {
	ts := m.tfstate
	if ts == nil {
		return fmt.Errorf("mtls authentication credential: empty state")
	}

	hasKC := ts.IsKnownAndNotNull("kubeconfig_data") && strings.TrimSpace(ts.String("kubeconfig_data")) != ""
	caSet := ts.IsKnownAndNotNull("certificate_authority_b64") && strings.TrimSpace(ts.String("certificate_authority_b64")) != ""
	certSet := ts.IsKnownAndNotNull("client_certificate_b64") && strings.TrimSpace(ts.String("client_certificate_b64")) != ""
	keySet := ts.IsKnownAndNotNull("client_key_b64") && strings.TrimSpace(ts.String("client_key_b64")) != ""

	if hasKC {
		ctxName := ""
		if ts.IsKnownAndNotNull("context_name") {
			ctxName = strings.TrimSpace(ts.String("context_name"))
		}
		ca, cert, key, err := mtlsMaterialFromKubeconfig(strings.TrimSpace(ts.String("kubeconfig_data")), ctxName)
		if err != nil {
			return err
		}
		ts.Set("certificate_authority_b64", ca)
		ts.Set("client_certificate_b64", cert)
		ts.Set("client_key_b64", key)
		return nil
	}

	if caSet && certSet && keySet {
		return nil
	}
	return fmt.Errorf(
		"provide kubeconfig_data (optional context_name) or all three of " +
			"certificate_authority_b64, client_certificate_b64, and client_key_b64",
	)
}

// --- kubeconfig parsing -------------------------------------------------------

type kubeconfigFile struct {
	CurrentContext string `yaml:"current-context"`
	Contexts       []struct {
		Name    string `yaml:"name"`
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
	} `yaml:"contexts"`
	Clusters []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			CertificateAuthority     string `yaml:"certificate-authority"`
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificate     string `yaml:"client-certificate"`
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKey             string `yaml:"client-key"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

// mtlsMaterialFromKubeconfig returns base64 PEM values for CA, client cert, and key.
// contextName empty → use current-context.
func mtlsMaterialFromKubeconfig(kubeconfigYAML, contextName string) (caB64, certB64, keyB64 string, err error) {
	var cfg kubeconfigFile
	if err := yaml.Unmarshal([]byte(kubeconfigYAML), &cfg); err != nil {
		return "", "", "", fmt.Errorf("parse kubeconfig_data: %w", err)
	}
	ctxName := strings.TrimSpace(contextName)
	if ctxName == "" {
		ctxName = strings.TrimSpace(cfg.CurrentContext)
	}
	if ctxName == "" {
		return "", "", "", fmt.Errorf("kubeconfig_data: no context_name and no current-context")
	}

	var clusterName, userName string
	for _, c := range cfg.Contexts {
		if c.Name == ctxName {
			clusterName, userName = c.Context.Cluster, c.Context.User
			break
		}
	}
	if clusterName == "" || userName == "" {
		return "", "", "", fmt.Errorf("kubeconfig_data: context %q not found", ctxName)
	}

	for _, c := range cfg.Clusters {
		if c.Name != clusterName {
			continue
		}
		caB64, err = kubeconfigPEMData(c.Cluster.CertificateAuthorityData, c.Cluster.CertificateAuthority)
		if err != nil {
			return "", "", "", fmt.Errorf("kubeconfig cluster %q CA: %w", clusterName, err)
		}
		break
	}
	if caB64 == "" {
		return "", "", "", fmt.Errorf("kubeconfig_data: cluster %q not found or missing CA", clusterName)
	}

	for _, u := range cfg.Users {
		if u.Name != userName {
			continue
		}
		certB64, err = kubeconfigPEMData(u.User.ClientCertificateData, u.User.ClientCertificate)
		if err != nil {
			return "", "", "", fmt.Errorf("kubeconfig user %q client certificate: %w", userName, err)
		}
		keyB64, err = kubeconfigPEMData(u.User.ClientKeyData, u.User.ClientKey)
		if err != nil {
			return "", "", "", fmt.Errorf("kubeconfig user %q client key: %w", userName, err)
		}
		break
	}
	if certB64 == "" || keyB64 == "" {
		return "", "", "", fmt.Errorf("kubeconfig_data: user %q not found or missing client cert/key", userName)
	}
	return caB64, certB64, keyB64, nil
}

// kubeconfigPEMData prefers inline *-data (already base64); otherwise base64-encodes a local file's PEM.
func kubeconfigPEMData(dataField, filePath string) (string, error) {
	if s := strings.TrimSpace(dataField); s != "" {
		// Accept standard or raw base64; validate it decodes.
		if _, err := base64.StdEncoding.DecodeString(s); err != nil {
			if _, err2 := base64.RawStdEncoding.DecodeString(s); err2 != nil {
				return "", fmt.Errorf("invalid base64 in *-data field: %w", err)
			}
		}
		return s, nil
	}
	path := strings.TrimSpace(filePath)
	if path == "" {
		return "", fmt.Errorf("neither *-data nor file path set")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", path, err)
	}
	return base64.StdEncoding.EncodeToString(raw), nil
}
