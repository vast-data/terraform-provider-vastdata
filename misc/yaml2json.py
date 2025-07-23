#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.

"""
Swagger YAML to JSON Converter CLI

This script:
1. Reads a Swagger/OpenAPI YAML file.
2. Converts it to JSON.
3. Outputs a prettified JSON version (2-space indentation).

Usage:
    python yaml2json.py path/to/swagger.yaml
    python yaml2json.py path/to/swagger.yaml -o path/to/output.json

Arguments:
    source            Required. Path to the source YAML file.

Options:
    -o, --output      Optional. Destination path for the JSON file.
                     Default: /tmp/apiconv/swagger.json

Dependencies:
    - ruamel.yaml (Python): pip install ruamel.yaml
"""

import os
import argparse
import json
from ruamel.yaml import YAML

def ensure_dir_exists(path: str):
    dir_path = os.path.dirname(path)
    if dir_path:
        os.makedirs(dir_path, exist_ok=True)

def yaml_to_json(source_path: str, dest_path: str):
    yaml = YAML()
    yaml.allow_duplicate_keys = True
    with open(source_path, "r") as f:
        data = yaml.load(f)

    with open(dest_path, "w") as out:
        json.dump(data, out, indent=2)  # Prettified with 2-space indent

    print(f"✅ Converted {source_path} to {dest_path}")

def main():
    parser = argparse.ArgumentParser(description="Convert Swagger/OpenAPI YAML to JSON.")
    parser.add_argument("source", help="Path to the source YAML file")
    parser.add_argument("-o", "--output", help="Output JSON file path (default: /tmp/apiconv/swagger.json)")

    args = parser.parse_args()
    output_path = args.output or "/tmp/apiconv/swagger.json"

    ensure_dir_exists(output_path)
    yaml_to_json(args.source, output_path)

if __name__ == "__main__":
    main()
