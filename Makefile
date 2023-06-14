HOSTNAME=hashicorp.com
NAMESPACE=edu
NAME=vastdata
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=x86_64
SWAGGER_CODEGEN_FILE=swagger-codegen-cli-3.0.41.jar
SWAGGER_CODEGEN_URL= https://repo1.maven.org/maven2/io/swagger/codegen/v3/swagger-codegen-cli/3.0.41/

apifiles:=$(wildcard api/*/api.yaml)
getVer = echo "${1}" |cut -d / -f 2 -
showFile = echo "${1}"
BUILD_VERSIONS = 4.6.0 4.7.0
SHELL:=bash
BUILD_DEST=build
BUILD_DIR=build
GINKGO_FLAGS=""
TFPLUGIN_DOCS_OPTIONS = ""
RESOURCES = $(wildcard examples/resources/*)

document_import:
	for r in $(RESOURCES) ; do\
	     echo $${r} ;\
	     n=$$(echo $$(basename $${r})) ;\
	     t=$$(echo $${n} |awk -F 'vastdata_' '{print $$2}') ;\
	     echo terrafrom import $${n}.$${t} "<guid>" >  examples/resources/$${n}/import.sh;\
	done ;

document: document_import
	tfplugindocs $(TFPLUGIN_DOCS_OPTIONS) 

clean:
	rm -rf $(BUILD_DEST)/terraform-provider-vastdata
	rm -rf docs

$(BUILD_DIR)/swagger-codegen-cli.jar:
	(! test -e $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE)  && wget $(SWAGGER_CODEGEN_URL)$(SWAGGER_CODEGEN_FILE) -O $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE)) || ( test -e $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE))


flush_codegen: $(BUILD_DIR)/swagger-codegen-cli.jar
	rm -rf codegen/*
#	(! test -e $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE)  && wget $(SWAGGER_CODEGEN_URL)$(SWAGGER_CODEGEN_FILE) -O $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE)) || ( test -e $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE))

codegen: flush_codegen
	for i in $(BUILD_VERSIONS); do \
		echo "=================Building Structs For Version $${i}===================" ; \
		mkdir -p  codegen/$${i}/ ;\
		java -Dmodels -jar $(BUILD_DIR)/$(SWAGGER_CODEGEN_FILE) generate -l go -i api/$${i}/api.yaml -o codegen/$${i}/ ; \
	done 

sort-versions: codegen 
	rm -f /tmp/versions.txt
	export BUILD_VERSIONS="$(BUILD_VERSIONS)"; \
	for i in $(BUILD_VERSIONS); do \
		echo $${i}>>/tmp/versions.txt; \
	done ; \
	cp -av codegen/$$(cat /tmp/versions.txt |sort -V|tail -1) codegen/latest
	cp -av api/$$(cat /tmp/versions.txt |sort -V|tail -1)/api.yaml codegen/latest

build-provider: sort-versions
	export BUILD_VERSIONS="$(BUILD_VERSIONS)"; \
	cd codegen_tools; \
	go run *.go

build-formatter: build-provider
	echo "################Formatting datasources code################"; \
	go fmt ./datasources/; \
	echo "################Formatting resources code################" ; \
	go fmt ./resources/ 
	echo "################Formatting Vast Versions code################" ; \
	go fmt ./vast_versions/ 

$(BUILD_DEST)/terraform-provider-vastdata:
	go build -o $(BUILD_DEST)/terraform-provider-vastdata

build: $(BUILD_DEST)/terraform-provider-vastdata

test:
	ginkgo $(GINKGO_FLAGS) ./...

build-all: clean build-formatter build document
