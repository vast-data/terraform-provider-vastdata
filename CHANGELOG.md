## 2.0.0

NOTES:

* The provider has been **fully rewritten** using the new [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), replacing the legacy SDKv2 implementation.

* The provider now supports **dynamic schema generation** by parsing the VAST OpenAPI specification at runtime. This enables easier exposure of VAST API endpoints as Terraform resources and data sources, significantly improving coverage and maintainability.

BREAKING CHANGES:

* Legacy resources defined using SDKv2 have been removed or migrated to the new Plugin Framework format. You may need to review their configurations for compatibility, especially if relying on custom behaviors or non-standard attributes.
* We've prepared a migration script that automatically converts your existing configuration files into the format supported by the new Terraform provider. You can find the usage instructions [here](migration/README.md)

ENHANCEMENTS:

* Improved maintainability and extensibility via the Plugin Framework.
* Consistent handling of resource plans, state, and diagnostics.
* Enhanced API errors handling and validation using [go-vast-client](https://github.com/vast-data/go-vast-client) library.