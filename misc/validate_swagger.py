#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

"""
Swagger/OpenAPI Schema Validator

This script provides comprehensive validation and diagnostic reporting for Swagger v2 and OpenAPI schemas.
It identifies schema issues, validates structure, and provides detailed reports with exact locations
and actionable recommendations.

Usage:
    python validate_swagger.py path/to/swagger.yaml
    python validate_swagger.py path/to/swagger.yaml --debug
    python validate_swagger.py path/to/swagger.yaml --json-output report.json
"""

import os
import sys
import json
import argparse
from pathlib import Path
from ruamel.yaml import YAML

class SwaggerValidator:
    def __init__(self, input_file, debug=False):
        self.input_file = Path(input_file)
        self.debug = debug
        self.detailed_issues = []
        self.warnings = []
        
    def log(self, message, level="info"):
        """Log messages with appropriate emoji and formatting"""
        icons = {
            "info": "ℹ️ ",
            "success": "✅",
            "warning": "⚠️ ",
            "error": "❌",
            "debug": "🔍",
            "issue": "🚨"
        }
        print(f"{icons.get(level, '')} {message}")

    def report_issue(self, severity, location, issue_type, description, current_value=None, expected_value=None, context=None, blocks_compilation=False):
        """Report a detailed issue with the schema"""
        issue = {
            "severity": severity,  # "error", "warning", "info"
            "location": location,  # JSON path like "definitions.User.properties.name"
            "issue_type": issue_type,  # Type of issue like "invalid_required_field"
            "description": description,  # Human readable description
            "current_value": current_value,  # What we found
            "expected_value": expected_value,  # What should be there
            "context": context,  # Additional context
            "blocks_compilation": blocks_compilation  # Whether this prevents OpenAPI conversion
        }
        self.detailed_issues.append(issue)
        
        # Also log immediately if debug mode
        if self.debug:
            icon = "🚨" if severity == "error" else "⚠️" if severity == "warning" else "ℹ️"
            print(f"{icon} [{severity.upper()}] {location}: {description}")
            if current_value is not None:
                print(f"    Current: {current_value}")
            if expected_value is not None:
                print(f"    Expected: {expected_value}")
            if context:
                print(f"    Context: {context}")

    def load_schema(self):
        """Load and parse the schema file"""
        try:
            self.log(f"Loading schema file: {self.input_file}")
            
            yaml = YAML()
            yaml.allow_duplicate_keys = True
            
            with open(self.input_file, "r") as f:
                data = yaml.load(f)
            
            self.log(f"Successfully loaded schema", "success")
            return data
            
        except Exception as e:
            self.log(f"Failed to load schema: {e}", "error")
            raise

    def validate_schema(self, data):
        """Main validation entry point"""
        if not isinstance(data, dict):
            self.report_issue(
                "error", "root", "invalid_root_type",
                "Root document must be an object",
                current_value=type(data).__name__,
                expected_value="object/dict",
                context="OpenAPI/Swagger documents must be JSON objects"
            )
            return False

        # Check basic OpenAPI/Swagger structure
        self._validate_basic_structure(data)
        
        # Validate definitions/components
        if 'definitions' in data:
            self._validate_definitions(data['definitions'])
        elif 'components' in data and 'schemas' in data['components']:
            self._validate_definitions(data['components']['schemas'], is_openapi3=True)
        
        # Validate paths
        if 'paths' in data:
            self._validate_paths(data['paths'])
        
        return True

    def _validate_basic_structure(self, data):
        """Validate basic document structure"""
        # Check for required fields
        if 'swagger' in data:
            # Swagger 2.0
            if data.get('swagger') != '2.0':
                self.report_issue(
                    "warning", "swagger", "invalid_swagger_version",
                    f"Unexpected Swagger version: {data.get('swagger')}",
                    current_value=data.get('swagger'),
                    expected_value="2.0",
                    context="This validator is designed for Swagger 2.0"
                )
        elif 'openapi' in data:
            # OpenAPI 3.x
            version = data.get('openapi', '')
            if not version.startswith('3.'):
                self.report_issue(
                    "warning", "openapi", "unsupported_openapi_version",
                    f"Unsupported OpenAPI version: {version}",
                    current_value=version,
                    expected_value="3.x.x",
                    context="This validator primarily supports OpenAPI 3.x"
                )
        else:
            self.report_issue(
                "error", "root", "missing_version_field",
                "Document missing 'swagger' or 'openapi' version field",
                current_value="missing",
                expected_value="swagger: '2.0' or openapi: '3.x.x'",
                context="OpenAPI documents must specify their version"
            )

        # Check for info object
        if 'info' not in data:
            self.report_issue(
                "error", "info", "missing_info_object",
                "Document missing required 'info' object",
                current_value="missing",
                expected_value="info object with title and version",
                context="All OpenAPI documents must have an info section"
            )
        else:
            self._validate_info_object(data['info'])

    def _validate_info_object(self, info):
        """Validate the info object"""
        if not isinstance(info, dict):
            self.report_issue(
                "error", "info", "invalid_info_type",
                "Info field must be an object",
                current_value=type(info).__name__,
                expected_value="object",
                context="The info field contains API metadata"
            )
            return

        # Check required info fields
        if 'title' not in info:
            self.report_issue(
                "error", "info.title", "missing_title",
                "Info object missing required 'title' field",
                current_value="missing",
                expected_value="API title string",
                context="Every API must have a title"
            )

        if 'version' not in info:
            self.report_issue(
                "error", "info.version", "missing_version",
                "Info object missing required 'version' field",
                current_value="missing",
                expected_value="version string",
                context="Every API must specify its version"
            )

    def _validate_definitions(self, definitions, is_openapi3=False):
        """Validate schema definitions"""
        location_prefix = "components.schemas" if is_openapi3 else "definitions"
        
        if not isinstance(definitions, dict):
            self.report_issue(
                "error", location_prefix, "invalid_definitions_type",
                "Definitions must be an object",
                current_value=type(definitions).__name__,
                expected_value="object",
                context="Schema definitions must be a map of name to schema"
            )
            return

        for name, schema in definitions.items():
            location = f"{location_prefix}.{name}"
            self._validate_schema_object(schema, location)

    def _validate_schema_object(self, schema, location):
        """Detailed validation of a schema object"""
        if not isinstance(schema, dict):
            self.report_issue(
                "error", location, "invalid_schema_type",
                f"Schema must be an object",
                current_value=f"{type(schema).__name__}: {schema}",
                expected_value="schema object",
                context="All schemas must be objects with type, properties, etc."
            )
            return

        schema_type = schema.get('type')
        has_composition = any(k in schema for k in ['$ref', 'allOf', 'oneOf', 'anyOf'])

        # Check for type or composition
        if not schema_type and not has_composition:
            self.report_issue(
                "warning", location, "missing_type_or_composition",
                f"Schema has no type or composition keywords",
                current_value="missing",
                expected_value="type field or $ref/allOf/oneOf/anyOf",
                context="Schemas should have either a type or use composition"
            )

        # Validate properties
        if 'properties' in schema:
            self._validate_properties(schema['properties'], location)

        # Check array-specific issues
        if schema_type == 'array':
            self._validate_array_schema(schema, location)

        # Check object-specific issues
        if schema_type == 'object' or 'properties' in schema:
            self._validate_object_schema(schema, location)

        # Recursively validate nested schemas
        self._validate_nested_schemas(schema, location)

    def _validate_properties(self, properties, location):
        """Validate schema properties"""
        if not isinstance(properties, dict):
            self.report_issue(
                "error", f"{location}.properties", "invalid_properties_type",
                "Properties must be an object",
                current_value=f"{type(properties).__name__}: {properties}",
                expected_value="object mapping property names to schemas",
                context="Properties define the fields of an object schema"
            )
            return

        for prop_name, prop_schema in properties.items():
            prop_location = f"{location}.properties.{prop_name}"

            if prop_schema is None:
                self.report_issue(
                    "error", prop_location, "null_property_schema",
                    f"Property '{prop_name}' has null schema",
                    current_value="null",
                    expected_value="schema object with type, description, etc.",
                    context="Properties cannot be null, they need proper schema definitions",
                    blocks_compilation=True
                )
                continue

            if not isinstance(prop_schema, dict):
                self.report_issue(
                    "error", prop_location, "invalid_property_schema",
                    f"Property '{prop_name}' schema is not an object",
                    current_value=f"{type(prop_schema).__name__}: {prop_schema}",
                    expected_value="schema object",
                    context="Property schemas must be objects with type, description, etc.",
                    blocks_compilation=True
                )
                continue

            # Check for invalid required fields in properties
            if 'required' in prop_schema:
                required_value = prop_schema['required']
                prop_type = prop_schema.get('type')
                
                # Only flag as error if this is not an object schema with a valid required array
                # Object schemas can legitimately have required arrays listing their required properties
                if prop_type != 'object' or not isinstance(required_value, list):
                    self.report_issue(
                        "error", prop_location, "invalid_required_field",
                        f"Property '{prop_name}' has invalid 'required' field",
                        current_value=required_value,
                        expected_value="For object types: array of required property names. For other types: remove required field",
                        context="Only object schemas can have 'required' arrays. Non-object properties should not have 'required' fields",
                        blocks_compilation=True
                    )

            # Check for nested type issues
            if 'type' in prop_schema:
                type_value = prop_schema['type']
                if isinstance(type_value, dict):
                    self.report_issue(
                        "error", f"{prop_location}.type", "nested_type_field",
                        f"Property '{prop_name}' has nested type object instead of string",
                        current_value=type_value,
                        expected_value="Simple string like 'string', 'integer', 'boolean'",
                        context="Type fields should be simple strings, not nested objects"
                    )

            # Check for nested description issues
            if 'description' in prop_schema:
                desc_value = prop_schema['description']
                if isinstance(desc_value, dict):
                    self.report_issue(
                        "error", f"{prop_location}.description", "nested_description_field",
                        f"Property '{prop_name}' has nested description object instead of string",
                        current_value=desc_value,
                        expected_value="Simple description string",
                        context="Description fields should be simple strings"
                    )

            # Check if property has proper type or $ref
            if (not prop_schema.get('type') and '$ref' not in prop_schema and 
                'allOf' not in prop_schema and 'oneOf' not in prop_schema and 'anyOf' not in prop_schema):
                self.report_issue(
                    "warning", prop_location, "missing_property_type",
                    f"Property '{prop_name}' missing type or reference",
                    current_value="missing",
                    expected_value="type field or $ref",
                    context="Properties should specify their data type or reference another schema"
                )

            # Recursively validate the property schema
            self._validate_schema_object(prop_schema, prop_location)

    def _validate_array_schema(self, schema, location):
        """Validate array-specific schema issues"""
        # Check for invalid properties field on arrays
        if 'properties' in schema:
            self.report_issue(
                "error", f"{location}.properties", "array_properties",
                "Array schema has 'properties' field - arrays should use 'items'",
                current_value=schema['properties'],
                expected_value="Remove 'properties', use 'items' instead",
                context="Arrays define their element schema with 'items', not 'properties'",
                blocks_compilation=True
            )

        # Check for missing or invalid items
        items = schema.get('items')
        if not items:
            self.report_issue(
                "error", f"{location}.items", "missing_array_items",
                "Array schema missing 'items' definition",
                current_value="missing",
                expected_value="items schema object",
                context="Arrays must define what type of elements they contain"
            )
        elif not isinstance(items, dict):
            self.report_issue(
                "error", f"{location}.items", "invalid_array_items",
                "Array 'items' is not an object",
                current_value=f"{type(items).__name__}: {items}",
                expected_value="schema object",
                context="Array items must be a schema object"
            )
        else:
            # Recursively validate items schema
            self._validate_schema_object(items, f"{location}.items")

        # Check for invalid required field on arrays
        if 'required' in schema and schema['required'] is not None:
            required_value = schema['required']
            self.report_issue(
                "error", f"{location}.required", "invalid_array_required",
                "Array schema has 'required' field - arrays cannot have required properties",
                current_value=required_value,
                expected_value="Remove 'required' field or move to items schema if applicable",
                context="Only object schemas can have required fields, not arrays",
                blocks_compilation=True
            )

    def _validate_object_schema(self, schema, location):
        """Validate object-specific schema issues"""
        # Validate required field
        if 'required' in schema:
            required_value = schema['required']
            if not isinstance(required_value, list):
                self.report_issue(
                    "error", f"{location}.required", "invalid_required_type",
                    "Required field must be an array of strings",
                    current_value=f"{type(required_value).__name__}: {required_value}",
                    expected_value="array of property names",
                    context="Required field should list property names that are mandatory"
                )
            else:
                # Check that required properties exist
                properties = schema.get('properties', {})
                for req_prop in required_value:
                    if not isinstance(req_prop, str):
                        self.report_issue(
                            "error", f"{location}.required", "invalid_required_item",
                            f"Required array contains non-string item: {req_prop}",
                            current_value=f"{type(req_prop).__name__}: {req_prop}",
                            expected_value="string property name",
                            context="Required array should only contain property names as strings"
                        )
                    elif req_prop not in properties:
                        self.report_issue(
                            "warning", f"{location}.required", "required_property_missing",
                            f"Required property '{req_prop}' not defined in properties",
                            current_value=f"'{req_prop}' in required but not in properties",
                            expected_value="property defined in properties section",
                            context="All required properties should be defined in the properties section"
                        )

    def _validate_nested_schemas(self, schema, location):
        """Recursively validate nested schemas"""
        for key, value in schema.items():
            if key in ['additionalProperties'] and isinstance(value, dict):
                self._validate_schema_object(value, f"{location}.{key}")
            elif key in ['allOf', 'oneOf', 'anyOf'] and isinstance(value, list):
                for i, sub_schema in enumerate(value):
                    if isinstance(sub_schema, dict):
                        self._validate_schema_object(sub_schema, f"{location}.{key}[{i}]")

    def _validate_paths(self, paths):
        """Validate API paths"""
        if not isinstance(paths, dict):
            self.report_issue(
                "error", "paths", "invalid_paths_type",
                "Paths must be an object",
                current_value=type(paths).__name__,
                expected_value="object mapping path patterns to operations",
                context="Paths define the API endpoints"
            )
            return

        for path, methods in paths.items():
            if not isinstance(methods, dict):
                continue

            for method, operation in methods.items():
                if not isinstance(operation, dict):
                    continue

                operation_location = f"paths.{path}.{method}"
                self._validate_operation(operation, operation_location)

    def _validate_operation(self, operation, location):
        """Validate an API operation"""
        # Check parameters
        if 'parameters' in operation:
            parameters = operation['parameters']
            if isinstance(parameters, list):
                for i, param in enumerate(parameters):
                    param_location = f"{location}.parameters[{i}]"
                    self._validate_parameter(param, param_location)

        # Check responses
        if 'responses' in operation:
            responses = operation['responses']
            if isinstance(responses, dict):
                for status, response in responses.items():
                    response_location = f"{location}.responses.{status}"
                    self._validate_response(response, response_location)

    def _validate_parameter(self, param, location):
        """Validate a parameter object"""
        if not isinstance(param, dict):
            self.report_issue(
                "error", location, "invalid_parameter",
                "Parameter must be an object",
                current_value=f"{type(param).__name__}: {param}",
                expected_value="parameter object",
                context="Parameters must be objects with name, in, type, etc."
            )
            return

        # Check required parameter fields
        if 'name' not in param:
            self.report_issue(
                "error", f"{location}.name", "missing_parameter_name",
                "Parameter missing required 'name' field",
                current_value="missing",
                expected_value="parameter name string",
                context="All parameters must have a name"
            )

        if 'in' not in param:
            self.report_issue(
                "error", f"{location}.in", "missing_parameter_location",
                "Parameter missing required 'in' field",
                current_value="missing",
                expected_value="query, header, path, formData, or body",
                context="All parameters must specify where they are located"
            )

        # Check parameter required field (should be boolean)
        if 'required' in param:
            required_value = param['required']
            if not isinstance(required_value, bool):
                self.report_issue(
                    "error", f"{location}.required", "parameter_required",
                    "Parameter 'required' field must be a boolean",
                    current_value=f"{type(required_value).__name__}: {required_value}",
                    expected_value="true or false",
                    context="Parameter required field should be boolean, not array or other type",
                    blocks_compilation=True
                )

        # Check body parameter schema
        if param.get('in') == 'body' and 'schema' in param:
            schema = param['schema']
            if isinstance(schema, dict):
                self._validate_schema_object(schema, f"{location}.schema")

        # Check parameter type (for non-body parameters)
        if param.get('in') != 'body' and 'type' not in param and '$ref' not in param:
            self.report_issue(
                "warning", f"{location}.type", "missing_parameter_type",
                "Parameter missing type definition",
                current_value="missing",
                expected_value="string, integer, boolean, array, etc.",
                context="Non-body parameters should specify their type"
            )

    def _validate_response(self, response, location):
        """Validate a response object"""
        if not isinstance(response, dict):
            return

        if 'schema' in response:
            schema = response['schema']
            if isinstance(schema, dict):
                self._validate_schema_object(schema, f"{location}.schema")

    def print_detailed_report(self):
        """Print a comprehensive diagnostic report"""
        if not self.detailed_issues:
            self.log("✅ No schema issues detected - schema is valid!", "success")
            return True

        print("\n" + "="*80)
        print("📋 SWAGGER/OPENAPI SCHEMA VALIDATION REPORT")
        print("="*80)
        
        # Group issues by severity and compilation-blocking status
        errors = [i for i in self.detailed_issues if i["severity"] == "error"]
        warnings = [i for i in self.detailed_issues if i["severity"] == "warning"]
        info = [i for i in self.detailed_issues if i["severity"] == "info"]
        
        # Critical issues that prevent OpenAPI conversion/compilation
        critical_issues = [i for i in self.detailed_issues if i.get("blocks_compilation", False)]
        
        print(f"\n📊 SUMMARY:")
        print(f"   🚨 Errors: {len(errors)}")
        print(f"   ⚠️  Warnings: {len(warnings)}")
        print(f"   ℹ️  Info: {len(info)}")
        print(f"   💥 CRITICAL (Blocks Compilation): {len(critical_issues)}")
        print(f"   📝 Total Issues: {len(self.detailed_issues)}")
        
        # Show critical issues first if any exist
        if critical_issues:
            print(f"\n💥 CRITICAL COMPILATION-BLOCKING ISSUES ({len(critical_issues)}):")
            print("="*80)
            print("❗ These issues MUST be fixed before OpenAPI v3 conversion will work!")
            print("-" * 80)
            
            # Group critical issues by type
            critical_by_type = {}
            for issue in critical_issues:
                issue_type = issue["issue_type"]
                if issue_type not in critical_by_type:
                    critical_by_type[issue_type] = []
                critical_by_type[issue_type].append(issue)
            
            for issue_type, type_issues in critical_by_type.items():
                print(f"\n  💥 {issue_type.replace('_', ' ').title()} ({len(type_issues)} issues):")
                
                # Show ALL critical issues (don't collapse them)
                for i, issue in enumerate(type_issues, 1):
                    print(f"\n    {i}. Location: {issue['location']}")
                    print(f"       Description: {issue['description']}")
                    
                    if issue['current_value'] is not None:
                        current_str = str(issue['current_value'])
                        if len(current_str) > 100:
                            current_str = current_str[:97] + "..."
                        print(f"       Current Value: {current_str}")
                    
                    if issue['expected_value'] is not None:
                        print(f"       Expected: {issue['expected_value']}")
            
            print("\n" + "="*80)
        
        for severity, issues, icon in [("ERROR", errors, "🚨"), ("WARNING", warnings, "⚠️"), ("INFO", info, "ℹ️")]:
            if not issues:
                continue
                
            print(f"\n{icon} {severity}S ({len(issues)}):")
            print("-" * 50)
            
            # Group by issue type
            by_type = {}
            for issue in issues:
                issue_type = issue["issue_type"]
                if issue_type not in by_type:
                    by_type[issue_type] = []
                by_type[issue_type].append(issue)
            
            for issue_type, type_issues in by_type.items():
                print(f"\n  📂 {issue_type.replace('_', ' ').title()} ({len(type_issues)} issues):")
                
                # Show first 10 issues of each type
                for i, issue in enumerate(type_issues[:10], 1):
                    print(f"\n    {i}. Location: {issue['location']}")
                    print(f"       Description: {issue['description']}")
                    
                    if issue['current_value'] is not None:
                        current_str = str(issue['current_value'])
                        if len(current_str) > 100:
                            current_str = current_str[:97] + "..."
                        print(f"       Current Value: {current_str}")
                    
                    if issue['expected_value'] is not None:
                        print(f"       Expected: {issue['expected_value']}")
                    
                    if issue['context']:
                        print(f"       Context: {issue['context']}")
                
                if len(type_issues) > 10:
                    print(f"\n    ... and {len(type_issues) - 10} more {issue_type.replace('_', ' ')} issues")
        
        print("\n" + "="*80)
        print("💡 RECOMMENDATIONS:")
        print("="*80)
        
        # Show critical compilation-blocking recommendations first
        critical_types = set(issue["issue_type"] for issue in critical_issues)
        if critical_types:
            print("\n💥 CRITICAL FIXES REQUIRED FOR COMPILATION:")
            print("-" * 50)
            
            if "invalid_required_field" in critical_types:
                print("\n🚨 Fix Required Field Issues (CRITICAL):")
                print("   - Remove 'required: true/false' from individual properties")
                print("   - Add 'required: [\"prop1\", \"prop2\"]' at the schema level instead")
                print("   - Example: Move 'required: true' from property to schema-level array")
                
            if "array_properties" in critical_types:
                print("\n🚨 Fix Array Schema Issues (CRITICAL):")
                print("   - Remove 'properties' field from array schemas")
                print("   - Use 'items' to define the schema of array elements")
                print("   - Arrays cannot have 'properties', only 'items'")
                
            if "parameter_required" in critical_types:
                print("\n🚨 Fix Parameter Required Fields (CRITICAL):")
                print("   - Change parameter 'required' from array to boolean")
                print("   - Example: Change 'required: [\"username\", \"password\"]' to 'required: true'")
                
            if "invalid_array_required" in critical_types:
                print("\n🚨 Fix Array Required Fields (CRITICAL):")
                print("   - Remove 'required' field from array schemas")
                print("   - Move required fields to the items schema if needed")
                print("   - Only object schemas can have required fields")
                
            if "null_property_schema" in critical_types:
                print("\n🚨 Fix Null Property Schemas (CRITICAL):")
                print("   - Replace null property values with proper schema objects")
                print("   - Add 'type', 'description' and other schema fields")
                
            if "invalid_property_schema" in critical_types:
                print("\n🚨 Fix Invalid Property Schemas (CRITICAL):")
                print("   - Convert string property values to proper schema objects")
                print("   - Property schemas must be objects with 'type', 'description', etc.")
        
        # Provide general recommendations based on all issue types
        issue_types = set(issue["issue_type"] for issue in self.detailed_issues)
        
        print(f"\n📚 GENERAL RECOMMENDATIONS:")
        print("-" * 50)
        
        if "invalid_required_field" in issue_types:
            print("\n🔧 Required Field Issues:")
            print("   - In OpenAPI/Swagger, 'required' should be an array of strings at the schema level")
            print("   - Remove 'required: true/false' from individual properties")
            print("   - Add 'required: [\"prop1\", \"prop2\"]' at the schema level instead")
            
        if "array_properties" in issue_types:
            print("\n🔧 Array Schema Issues:")
            print("   - Array schemas should have 'items' property, not 'properties'")
            print("   - Use 'items' to define the schema of array elements")
            
        if "parameter_required" in issue_types or "missing_parameter_location" in issue_types:
            print("\n🔧 Parameter Issues:")
            print("   - Parameter 'required' field should be a boolean (true/false)")
            print("   - All parameters must have 'name' and 'in' fields")
            print("   - 'in' should be one of: query, header, path, formData, body")

        if "nested_type_field" in issue_types or "nested_description_field" in issue_types:
            print("\n🔧 Nested Field Issues:")
            print("   - Type fields should be simple strings like 'string', 'integer', 'boolean'")
            print("   - Description fields should be simple strings")
            print("   - Avoid nested objects in these fields")

        print("\n" + "="*80)
        
        return len(errors) == 0  # Return True if no errors (only warnings/info)

    def export_json_report(self, output_file):
        """Export detailed report as JSON"""
        critical_issues = [i for i in self.detailed_issues if i.get("blocks_compilation", False)]
        
        report = {
            "file": str(self.input_file),
            "summary": {
                "total_issues": len(self.detailed_issues),
                "errors": len([i for i in self.detailed_issues if i["severity"] == "error"]),
                "warnings": len([i for i in self.detailed_issues if i["severity"] == "warning"]),
                "info": len([i for i in self.detailed_issues if i["severity"] == "info"]),
                "critical_compilation_blocking": len(critical_issues)
            },
            "critical_issues": critical_issues,
            "issues": self.detailed_issues
        }
        
        with open(output_file, 'w') as f:
            json.dump(report, f, indent=2, default=str)
        
        self.log(f"JSON report exported to: {output_file}", "success")

    def validate(self, json_output=None):
        """Main validation workflow"""
        try:
            # Load and parse schema
            data = self.load_schema()
            
            # Validate schema
            self.validate_schema(data)
            
            # Print detailed report
            is_valid = self.print_detailed_report()
            
            # Export JSON report if requested
            if json_output:
                self.export_json_report(json_output)
            
            return is_valid
            
        except Exception as e:
            self.log(f"Validation failed: {e}", "error")
            if self.debug:
                import traceback
                traceback.print_exc()
            return False


def main():
    parser = argparse.ArgumentParser(
        description="Swagger/OpenAPI Schema Validator",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s swagger.yaml                    # Basic validation
  %(prog)s swagger.yaml --debug            # Detailed debugging output  
  %(prog)s swagger.yaml --json-output report.json  # Export JSON report
        """
    )
    
    parser.add_argument("input_file", help="Path to Swagger/OpenAPI YAML or JSON file")
    parser.add_argument("--debug", action="store_true", 
                       help="Enable detailed debugging output")
    parser.add_argument("--json-output", 
                       help="Export detailed report as JSON to specified file")
    
    args = parser.parse_args()
    
    if not Path(args.input_file).exists():
        print(f"❌ Input file not found: {args.input_file}")
        sys.exit(1)
    
    print("🔍 Swagger/OpenAPI Schema Validation")
    print(f"📁 Input file: {args.input_file}")
    if args.json_output:
        print(f"📄 JSON report: {args.json_output}")
    print()
    
    validator = SwaggerValidator(
        input_file=args.input_file,
        debug=args.debug
    )
    
    is_valid = validator.validate(json_output=args.json_output)
    
    print(f"\n📋 Validation complete for: {args.input_file}")
    if is_valid:
        print("✅ Schema is valid!")
    else:
        print("❌ Schema has issues that need to be addressed")
    
    sys.exit(0 if is_valid else 1)


if __name__ == "__main__":
    main()