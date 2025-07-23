# Copyright (c) HashiCorp, Inc.

FROM golang:1.24 AS builder

# Install system dependencies
RUN apt-get update && apt-get install -y \
    zip \
    git \
    curl \
    make \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

ENV UPX_VERSION=4.2.1
RUN curl -L -o upx.tar.xz https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz && \
    tar -xf upx.tar.xz && \
    mv upx-${UPX_VERSION}-amd64_linux/upx /usr/local/bin/upx && \
    chmod +x /usr/local/bin/upx && \
    rm -rf upx.tar.xz upx-${UPX_VERSION}-amd64_linux

# Install latest Terraform from releases.hashicorp.com
RUN set -eux; \
    TERRAFORM_URL=$(curl -s https://checkpoint-api.hashicorp.com/v1/check/terraform | grep -oP '"current_version": *"\K[^"]+') && \
    curl -L -o terraform.zip "https://releases.hashicorp.com/terraform/${TERRAFORM_URL}/terraform_${TERRAFORM_URL}_linux_amd64.zip" && \
    unzip terraform.zip -d /usr/local/bin && \
    chmod +x /usr/local/bin/terraform && \
    rm terraform.zip

WORKDIR /app
