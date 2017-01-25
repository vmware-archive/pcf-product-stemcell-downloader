# PCF Product Stemcell Downloader

This is a tool for unpacking a PCF product tile, figuring out which stemcell it
requires, and downloading that stemcell.

It is used by [PCF Pipelines](https://github.com/c0-ops/pcf-pipelines).

## Installing

Download the [latest release](https://github.com/c0-ops/pcf-product-stemcell-downloader/releases/latest).

### Install from source

Requirements:

* [glide](https://github.com/masterminds/glide)
* [go](https://golang.org)

```
mkdir -p $GOPATH/src/github.com/c0-ops/pcf-product-stemcell-downloader
git clone git@github.com:c0-ops/pcf-product-stemcell-downloader.git $GOPATH/src/github.com/c0-ops/pcf-product-stemcell-downloader
cd $GOPATH/src/github.com/c0-ops/pcf-product-stemcell-downloader
glide install
GOARCH=amd64 GOOS=linux go install github.com/c0-ops/pcf-product-stemcell-downloader
```

## Usage

`pcf-product-stemcell-downloader [OPTIONS]`

### Options

All options are required.

* `download-dir`: Directory to place downloaded stemcell in
* `iaas-type`: Stemcell for this IaaS will be downloaded. Valid options are `aws|openstack|vcloud|vsphere|azure|gcp`
* `product-file`: Path to .pivotal product file to extract stemcell requirements from.
* `product-name`: Name of product. This will be used to find the correct metadata in the product file.

## Developing

Note that the tests currently download a stemcell from bosh.io, and thus require Internet access.

To run all of the tests in a Docker container:

`./testrunner`

To continually run the tests during development:

* `ginkgo watch -r -p`
