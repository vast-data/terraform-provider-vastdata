HOSTNAME=hashicorp.com
NAMESPACE=edu
NAME=vastdata
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=x86_64
SWAGGER_CODEGEN_FILE=swagger-codegen-cli-3.0.41.jar
SWAGGER_CODEGEN_URL= https://repo1.maven.org/maven2/io/swagger/codegen/v3/swagger-codegen-cli/3.0.41/
USER_AGENT_UPDATE_FILE = ./vast-client/user_agent_version.go

apifiles:=$(wildcard api/*/api.yaml)
getVer = echo "${1}" |cut -d / -f 2 -
showFile = echo "${1}"
BUILD_VERSIONS = 4.6.0 4.7.0 5.0.0 5.1.0 5.2.0
SHELL:=bash
BUILD_DEST=build
BUILD_DIR=build
GINKGO_FLAGS=""
TFPLUGIN_DOCS_OPTIONS = ""
RESOURCES = $(wildcard examples/resources/*)
FORCE_TAG_MATCH = 1
GOFLAGS = ""
TF_GPG_SIG = ""
USER_AGENT_FILE=user_agent_version


.PHONY: $(BUILD_DEST)/terraform-provider-vastdata

document:
	tfplugindocs $(TFPLUGIN_DOCS_OPTIONS) 

clean-releases:
	rm -rf $(BUILD_DEST)/*.tar.gz
	rm -rf $(BUILD_DEST)/terraform-provider-vastdata*SHA256SUMS
clean:
	rm -rf $(BUILD_DEST)/terraform-provider-vastdata
	rm -rf $(BUILD_DEST)/*_*/
	rm -rf $(BUILD_DEST)/{*.zip,*.tar,*.tar.gz,*.json,*.sig}
	rm -rf docs
	rm -rf $(USER_AGENT_UPDATE_FILE)

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

validate_user_agent_version_file:
	[ -e $(USER_AGENT_FILE) ] && true 


validate_user_agent_version_file_release_ready: validate_user_agent_version_file is-tag
	tag=$$(git describe --tags) ;\
	user_agent_version_file=$$(cat $(USER_AGENT_FILE)) ;\
	echo $${user_agent_version_file} ;\
	[ $${tag} == $${user_agent_version_file} ] && true


build-user-agent-file: validate_user_agent_version_file 
	tag=$$(cat $(USER_AGENT_FILE)) ;\
	echo $${tag}}; \
	echo  -e "package vast_client \\n\\n func init() { \\n user_agent_version=\"$${tag}\" \\n }" > $(USER_AGENT_UPDATE_FILE) ;\
	gofmt -w $(USER_AGENT_UPDATE_FILE) 


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
	go build $${GOFLAGS} -o $(BUILD_DEST)/terraform-provider-vastdata

build: $(BUILD_DEST)/terraform-provider-vastdata

test:
	ginkgo $(GINKGO_FLAGS) ./...

build-all: clean build-formatter build-user-agent-file build document

$(BUILD_DEST)/linux_amd64/terraform-provider-vastdata:
	export GOOS="linux" ;\
	export GOARCH="amd64" ;\
	export CGO_ENABLED=0 ;\
	mkdir -p $(BUILD_DEST)/$${GOOS}_$${GOARCH} ;\
	go build $${GOFLAGS} -o $(BUILD_DEST)/$${GOOS}_$${GOARCH}/terraform-provider-vastdata

build-linux-amd64:   $(BUILD_DEST)/linux_amd64/terraform-provider-vastdata

$(BUILD_DEST)/linux_arm64/terraform-provider-vastdata:
	export GOOS="linux" ;\
	export GOARCH="arm64" ;\
	mkdir -p $(BUILD_DEST)/$${GOOS}_$${GOARCH} ;\
	go build $${GOFLAGS} -o $(BUILD_DEST)/$${GOOS}_$${GOARCH}/terraform-provider-vastdata

build-linux-arm64:   $(BUILD_DEST)/linux_arm64/terraform-provider-vastdata


$(BUILD_DEST)/darwin_amd64/terraform-provider-vastdata:
	export GOOS="darwin" ;\
	export GOARCH="amd64" ;\
	export CGO_ENABLED=0 ;\
	mkdir -p $(BUILD_DEST)/$${GOOS}_$${GOARCH} ;\
	go build $${GOFLAGS} -o $(BUILD_DEST)/$${GOOS}_$${GOARCH}/terraform-provider-vastdata

build-darwin-amd64: $(BUILD_DEST)/darwin_amd64/terraform-provider-vastdata


$(BUILD_DEST)/darwin_arm64/terraform-provider-vastdata:
	export GOOS="darwin" ;\
	export GOARCH="arm64" ;\
	export CGO_ENABLED=0 ;\
	mkdir -p $(BUILD_DEST)/$${GOOS}_$${GOARCH} ;\
	go build $${GOFLAGS} -o $(BUILD_DEST)/$${GOOS}_$${GOARCH}/terraform-provider-vastdata

build-darwin-arm64: $(BUILD_DEST)/darwin_arm64/terraform-provider-vastdata


$(BUILD_DEST)/windows_amd64/terraform-provider-vastdata:
	export GOOS="windows" ;\
	export GOARCH="amd64" ;\
	export CGO_ENABLED=0 ;\
	mkdir -p $(BUILD_DEST)/$${GOOS}_$${GOARCH} ;\
	go build -o $(BUILD_DEST)/$${GOOS}_$${GOARCH}/terraform-provider-vastdata

build-windows-amd64: $(BUILD_DEST)/windows_amd64/terraform-provider-vastdata

$(BUILD_DEST)/windows_arm64/terraform-provider-vastdata:
	export GOOS="windows" ;\
	export GOARCH="arm64" ;\
	export CGO_ENABLED=0 ;\
	mkdir -p $(BUILD_DEST)/$${GOOS}_$${GOARCH} ;\
	go build $${GOFLAGS} -o $(BUILD_DEST)/$${GOOS}_$${GOARCH}/terraform-provider-vastdata

build-windows-arm64: $(BUILD_DEST)/windows_arm64/terraform-provider-vastdata

build-archs: build-windows-arm64 build-windows-amd64 build-darwin-arm64 build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64

build-all-archs: clean build-formatter build-archs document

#check if this is a tagged version#
is-tag:
	if [ "$(FORCE_TAG_MATCH)" = "1" ] ; then \
		git describe --exact-match --tags $$(git log -n1 --pretty='%h') ;\
	fi

pack-archs: clean-releases is-tag
	tag=$$(git describe --tags); \
	tf_tag=$$(git describe --tags|sed -r s/'v(.*)'/'\1'/g ); \
	for arch in $(BUILD_DEST)/*_*/terraform-provider-vastdata ; do \
		mv -v $${arch} $${arch}_$${tag} ; \
		arch_build=$$(echo $${arch}|awk -F '$(BUILD_DEST)' '{print $$2}'|tr "/" " "|awk -F " " '{print $$1}'); \
		zip_file=$(BUILD_DEST)/terraform-provider-vastdata_$${tf_tag}_$${arch_build}.zip; \
		echo "Creating Zip File $${zip_file}"; \
		zip -j $${zip_file} $${arch}_$${tag}; \
	done; \
	echo "Generating SHA256SUMS file";\
	pushd $${PWD};\
	cd $(BUILD_DEST);\
	shasum -a 256 *.zip > terraform-provider-vastdata_$${tf_tag}_SHA256SUMS;\
	popd;\
	gpg --detach-sign --local-user ${TF_GPG_SIG} build/terraform-provider-vastdata_$${tf_tag}_SHA256SUMS

pack-all-archs: build-all-archs pack-archs

github-pre-release: validate_user_agent_version_file_release_ready is-tag
	tag=$$(git describe --tags); \
	gh release create $${tag}   ./build/*.zip ./build/*.sig  ./build/*_SHA256SUMS --prerelease --title "Release $${tag}" --generate-notes

github-release: validate_user_agent_version_file_release_ready is-tag
	tag=$$(git describe --tags); \
	gh release create $${tag} ./build/*.zip ./build/*.sig  ./build/*_SHA256SUMS --title "Release $${tag}" --generate-notes
