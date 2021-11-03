version := 99.99.0
provider_macos_path = registry.terraform.io/Kekenika/typesense/$(version)/darwin_amd64/

.PHONY: doc
doc:
	go generate ./...

.PHONY: build
build: 
	@go build

.PHONY: install_macos
install_macos: build
	@mkdir -p ~/Library/Application\ Support/io.terraform/plugins/$(provider_macos_path)
	@mv terraform-provider-typesense ~/Library/Application\ Support/io.terraform/plugins/$(provider_macos_path)
