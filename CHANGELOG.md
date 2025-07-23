## 2.0.0

NOTES:

* The provider has been **fully rewritten** using the new [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), replacing the legacy SDKv2 implementation.

* The provider now supports **dynamic schema generation** by parsing the VAST OpenAPI specification at runtime. This enables automatic exposure of all supported VAST API resources as Terraform resources and data sources, significantly improving coverage and maintainability.

BREAKING CHANGES:

* Legacy resources defined using SDKv2 have been removed or migrated to the new Plugin Framework format. Users may need to review their configurations for compatibility, especially if relying on custom behaviors or non-standard attributes.

FEATURES:

* Dynamic schema generation based on the VAST OpenAPI schema.
* Automatic creation of Terraform resources and data sources aligned with VAST API endpoints.

ENHANCEMENTS:

* Improved maintainability and extensibility via the Plugin Framework.
* Consistent handling of resource plans, state, and diagnostics.
* Enhanced API errors handling and validation using [go-vast-client](https://github.com/vast-data/go-vast-client) library.
