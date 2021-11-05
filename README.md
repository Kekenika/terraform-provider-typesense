# Terraform Provider for Typesense

This is a [Terraform](https://www.terraform.io/) provider for [Typesense](https://typesense.org/).

## Maintainers

This provider plugin is maintained by:

* [@KeisukeYamashita](https;//github.com/KeisukeYamashita)

## Support

- Supports v0.21.0 version of Typesense.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= v0.12.0 (v0.11.x may work but not supported actively)

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/Kekenika/terraform-provider-typesense`

```console
$ mkdir -p $GOPATH/src/github.com/Kekenika; cd $GOPATH/src/github.com/Kekenika
$ git clone git@github.com:Kekenika/terraform-provider-typesense
Enter the provider directory and build the provider

$ cd $GOPATH/src/github.com/Kekenika/terraform-provider-typesense
$ make build
```
