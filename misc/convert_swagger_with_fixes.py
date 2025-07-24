#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

"""
Enhanced Swagger to OpenAPI v3 Converter with Debugging and Auto-fixes

This script:
1. Reads Swagger v2 YAML and converts to JSON
2. Validates and identifies common schema issues
3. Automatically fixes known problems
4. Converts to OpenAPI v3 using Go converter
5. Creates tarball with final output

Usage:
    python convert_swagger_with_fixes.py path/to/swagger.yaml
    python convert_swagger_with_fixes.py path/to/swagger.yaml --output-dir /custom/output
    python convert_swagger_with_fixes.py path/to/swagger.yaml --debug --no-auto-fix
"""

import os
import sys
import json
import argparse
import tempfile
import subprocess
from pathlib import Path
from ruamel.yaml import YAML

class SwaggerConverter:
    def __init__(self, input_file, output_dir="/tmp/apiconv", debug=False, auto_fix=True):
        self.input_file = Path(input_file)
        self.output_dir = Path(output_dir)
        self.debug = debug
        self.auto_fix = auto_fix
        self.fixes_applied = []
        self.warnings = []
        
        # Ensure output directory exists
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # File paths
        self.swagger_json = self.output_dir / "swagger.json"
        self.api_json = self.output_dir / "api.json"
        self.api_tarball = self.output_dir / "api.tar.gz"

    def log(self, message, level="info"):
        """Log messages with appropriate emoji and formatting"""
        icons = {
            "info": "ℹ️ ",
            "success": "✅",
            "warning": "⚠️ ",
            "error": "❌",
            "debug": "🔍",
            "fix": "🔧"
        }
        print(f"{icons.get(level, '')} {message}")

    def yaml_to_json(self):
        """Convert YAML to JSON with duplicate key handling"""
        try:
            self.log(f"Converting YAML to JSON: {self.input_file}")
            
            yaml = YAML()
            yaml.allow_duplicate_keys = True
            
            with open(self.input_file, "r") as f:
                data = yaml.load(f)
            
            with open(self.swagger_json, "w") as f:
                json.dump(data, f, indent=2)
            
            self.log(f"Successfully converted to {self.swagger_json}", "success")
            return data
            
        except Exception as e:
            self.log(f"Failed to convert YAML to JSON: {e}", "error")
            raise

    def validate_and_debug_schema(self, data):
        """Comprehensive validation and debugging of swagger schema"""
        self.log("Validating Swagger schema...", "debug")
        
        # Basic structure validation
        if not isinstance(data, dict):
            raise ValueError("Root document must be an object")
        
        if data.get("swagger") != "2.0":
            self.warnings.append(f"Expected swagger: '2.0', got: {data.get('swagger')}")
        
        info = data.get("info", {})
        self.log(f"API: {info.get('title', 'Unknown')}")
        self.log(f"Version: {info.get('version', 'Unknown')}")
        
        # Validate definitions
        definitions = data.get("definitions", {})
        self.log(f"Definitions count: {len(definitions)}")
        
        problems = self._analyze_definitions(definitions)
        
        # Validate paths
        paths = data.get("paths", {})
        self.log(f"Paths count: {len(paths)}")
        problems.extend(self._analyze_paths(paths))
        
        if problems:
            self.log(f"Found {len(problems)} schema issues", "warning")
            if self.debug:
                for i, problem in enumerate(problems[:10], 1):
                    self.log(f"  {i}. {problem}", "warning")
                if len(problems) > 10:
                    self.log(f"  ... and {len(problems) - 10} more", "warning")
        else:
            self.log("No schema issues found", "success")
        
        return problems

    def _analyze_definitions(self, definitions):
        """Analyze definition schemas for common problems"""
        problems = []
        
        for name, schema in definitions.items():
            if not isinstance(schema, dict):
                problems.append(f"Definition '{name}': not an object")
                continue
            
            # Check for missing type
            schema_type = schema.get('type')
            has_composition = any(k in schema for k in ['$ref', 'allOf', 'oneOf', 'anyOf'])
            
            if not schema_type and not has_composition:
                problems.append(f"Definition '{name}': missing type and composition")
            
            # Check properties
            properties = schema.get('properties', {})
            for prop_name, prop_schema in properties.items():
                if prop_schema is None:
                    problems.append(f"Definition '{name}', property '{prop_name}': null value")
                elif not isinstance(prop_schema, dict):
                    problems.append(f"Definition '{name}', property '{prop_name}': not an object")
                elif not prop_schema.get('type') and '$ref' not in prop_schema:
                    problems.append(f"Definition '{name}', property '{prop_name}': missing type and $ref")
            
            # Check array items
            if schema_type == 'array':
                items = schema.get('items')
                if not items:
                    problems.append(f"Definition '{name}': array without items")
                elif not isinstance(items, dict):
                    problems.append(f"Definition '{name}': items not an object")
                elif not items.get('type') and '$ref' not in items:
                    problems.append(f"Definition '{name}': items missing type and $ref")
            
            # Check additionalProperties
            additional_props = schema.get('additionalProperties')
            if additional_props is not None and isinstance(additional_props, dict):
                if not additional_props.get('type') and '$ref' not in additional_props:
                    problems.append(f"Definition '{name}': additionalProperties missing type and $ref")
        
        return problems

    def _analyze_paths(self, paths):
        """Analyze API paths for schema problems"""
        problems = []
        
        for path, methods in paths.items():
            if not isinstance(methods, dict):
                continue
                
            for method, operation in methods.items():
                if not isinstance(operation, dict):
                    continue
                    
                # Check responses
                responses = operation.get('responses', {})
                for status, response in responses.items():
                    if not isinstance(response, dict):
                        continue
                        
                    schema = response.get('schema')
                    if schema is not None and not isinstance(schema, dict):
                        problems.append(f"Path '{path}' {method} response {status}: schema not an object")
                
                # Check parameters
                parameters = operation.get('parameters', [])
                for i, param in enumerate(parameters):
                    if not isinstance(param, dict):
                        problems.append(f"Path '{path}' {method} parameter {i}: not an object")
                    elif param.get('in') == 'body':
                        schema = param.get('schema')
                        if schema is not None and not isinstance(schema, dict):
                            problems.append(f"Path '{path}' {method} body parameter: schema not an object")
        
        return problems

    def apply_auto_fixes(self, data):
        """Apply automatic fixes for common schema problems"""
        if not self.auto_fix:
            return data
        
        self.log("Applying automatic fixes...", "fix")
        
        definitions = data.get("definitions", {})
        
        for name, schema in definitions.items():
            if not isinstance(schema, dict):
                continue
            
            properties = schema.get('properties', {})
            for prop_name, prop_schema in properties.items():
                if prop_schema is None:
                    # Fix known null properties
                    fixed_schema = self._get_fixed_schema(name, prop_name)
                    properties[prop_name] = fixed_schema
                    self.fixes_applied.append(f"Fixed null property '{name}.{prop_name}'")
                    self.log(f"Fixed Definition '{name}', property '{prop_name}': replaced null with {fixed_schema['type']} schema", "fix")
        
        return data

    def _get_fixed_schema(self, definition_name, property_name):
        """Get appropriate schema fix for known problematic properties"""
        
        # Known fixes for specific properties
        known_fixes = {
            's3_block_v2_authentication': {
                "type": "boolean",
                "description": "Manage s3 blocks v2 authentication"
            }
        }
        
        if property_name in known_fixes:
            return known_fixes[property_name]
        
        # Default fixes based on naming patterns
        if any(keyword in property_name.lower() for keyword in ['enabled', 'enable', 'disabled', 'active', 'flag']):
            return {
                "type": "boolean",
                "description": f"Boolean flag for {property_name}"
            }
        
        if any(keyword in property_name.lower() for keyword in ['id', 'count', 'size', 'length', 'port', 'number']):
            return {
                "type": "integer",
                "description": f"Numeric value for {property_name}"
            }
        
        # Default to string
        return {
            "type": "string",
            "description": f"Optional {property_name} parameter"
        }

    def convert_to_openapi_v3(self):
        """Convert Swagger v2 JSON to OpenAPI v3 using Go converter"""
        try:
            self.log("Converting to OpenAPI v3...")
            
            # Use the existing Go converter
            go_script = Path(__file__).parent / "convert_to_v3.go"
            
            if not go_script.exists():
                raise FileNotFoundError(f"Go converter not found: {go_script}")
            
            # Set environment variables for the Go script
            env = os.environ.copy()
            env['INPUT_PATH'] = str(self.swagger_json)
            env['OUTPUT_DIR'] = str(self.output_dir)
            
            result = subprocess.run(
                ["go", "run", str(go_script)],
                capture_output=True,
                text=True,
                env=env
            )
            
            if result.returncode != 0:
                raise RuntimeError(f"Go conversion failed: {result.stderr}")
            
            self.log("Successfully converted to OpenAPI v3", "success")
            self.log(f"Output: {self.api_json}")
            self.log(f"Archive: {self.api_tarball}")
            
        except Exception as e:
            self.log(f"OpenAPI v3 conversion failed: {e}", "error")
            raise

    def copy_outputs(self, dest_dir="."):
        """Copy final outputs to destination directory"""
        dest_dir = Path(dest_dir)
        
        if self.api_json.exists():
            dest_json = dest_dir / "api.json"
            dest_json.write_bytes(self.api_json.read_bytes())
            self.log(f"Copied {self.api_json} -> {dest_json}", "success")
        
        if self.api_tarball.exists():
            dest_tar = dest_dir / "api.tar.gz"
            dest_tar.write_bytes(self.api_tarball.read_bytes())
            self.log(f"Copied {self.api_tarball} -> {dest_tar}", "success")

    def convert(self, dest_dir="."):
        """Main conversion workflow"""
        try:
            # Step 1: Convert YAML to JSON
            data = self.yaml_to_json()
            
            # Step 2: Validate and debug
            problems = self.validate_and_debug_schema(data)
            
            # Step 3: Apply fixes
            if problems and self.auto_fix:
                data = self.apply_auto_fixes(data)
                
                # Save fixed JSON
                with open(self.swagger_json, "w") as f:
                    json.dump(data, f, indent=2)
                
                # Re-validate after fixes
                remaining_problems = self.validate_and_debug_schema(data)
                if len(remaining_problems) < len(problems):
                    self.log(f"Fixed {len(problems) - len(remaining_problems)} issues", "success")
            
            # Step 4: Convert to OpenAPI v3
            self.convert_to_openapi_v3()
            
            # Step 5: Copy outputs
            self.copy_outputs(dest_dir)
            
            # Report summary
            self.log(f"Conversion completed successfully!", "success")
            if self.fixes_applied:
                self.log(f"Applied {len(self.fixes_applied)} automatic fixes:")
                for fix in self.fixes_applied:
                    self.log(f"  - {fix}")
            
            if self.warnings:
                self.log(f"Warnings ({len(self.warnings)}):")
                for warning in self.warnings:
                    self.log(f"  - {warning}", "warning")
            
            return True
            
        except Exception as e:
            self.log(f"Conversion failed: {e}", "error")
            return False


def main():
    parser = argparse.ArgumentParser(
        description="Enhanced Swagger to OpenAPI v3 converter with auto-fixes",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s swagger.yaml
  %(prog)s swagger.yaml --output-dir /tmp/output
  %(prog)s swagger.yaml --debug --no-auto-fix
        """
    )
    
    parser.add_argument("input_file", help="Path to Swagger YAML file")
    parser.add_argument("--output-dir", default="/tmp/apiconv", 
                       help="Output directory for intermediate files (default: /tmp/apiconv)")
    parser.add_argument("--dest-dir", default=".", 
                       help="Destination directory for final outputs (default: current directory)")
    parser.add_argument("--debug", action="store_true", 
                       help="Enable detailed debugging output")
    parser.add_argument("--no-auto-fix", action="store_true", 
                       help="Disable automatic fixes")
    
    args = parser.parse_args()
    
    if not Path(args.input_file).exists():
        print(f"❌ Input file not found: {args.input_file}")
        sys.exit(1)
    
    converter = SwaggerConverter(
        input_file=args.input_file,
        output_dir=args.output_dir,
        debug=args.debug,
        auto_fix=not args.no_auto_fix
    )
    
    success = converter.convert(args.dest_dir)
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main() 