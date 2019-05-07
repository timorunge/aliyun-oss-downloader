# aliyun-oss-downloader

[![](https://goreportcard.com/badge/github.com/timorunge/aliyun-oss-downloader)](https://goreportcard.com/report/github.com/timorunge/aliyun-oss-downloader)

`aliyun-oss-downloader` is a simple CLI tool written in Go for downloading
all objects from an [Aliyun OSS](https://www.alibabacloud.com/product/oss)
bucket.

## Install

The installation of `aliyun-oss-downloader` is simple. Just run `go get`
to get the latest version.

```sh
go get -u github.com/timorunge/aliyun-oss-downloader
```

## Configuration

Since `aliyun-oss-downloader` is based on
[Cobra](https://github.com/spf13/cobra) it's providing an easy way to store all
flags in a `yaml` config file.

You should use this functionality at least for sensitive flags like
`accessKeyID` and `accessKeySecret`.

If not otherwiese specified `aliyun-oss-downloader` will try load its default
config file which is located in `$HOME/.aliyun-oss-downloader.yaml`. You can
override this with the `--config` flag.

### Example config:

```yaml
accessKeyID: Aecahl7ieghie6rae
accessKeySecret: Aigi2amaiyohRia5aithe7OivaiM6Da
bucket: myBucket
destinationDir: /mnt/aliyun/oss/myBucket
region: eu-central-1
threads: 20
```

## Usage

If `$GOPATH/bin` is not in your `$PATH` call `aliyun-oss-downloader`
directly via `$GOPATH/bin/aliyun-oss-downloader`.

```sh
Usage:
  aliyun-oss-downloader [flags]
  aliyun-oss-downloader [command]

Available Commands:
  help        Help about any command
  version     Version

Flags:
      --accessKeyID string       Your access key ID
      --accessKeySecret string   Your access key secret
  -b, --bucket string            The name of the OSS bucket which should be downloaded
      --config string            Config file (default "$HOME/.aliyun-oss-downloader.yaml")
      --createDestinationDir     Create the (local) destination directory if not existing
      --destinationDir string    The (local) destination directory
      --disableTLS               Use a HTTP connection instead of a HTTP over TLS (HTTPS) connection
      --exclude string           Regex for excluding items from the download
  -h, --help                     help for aliyun-oss-downloader
      --marker string            The marker to start the download
      --maxKeys int              The amount of objects which are fetched in a single request (default 250)
      --prefix string            The prefix to filter objects in the bucket
  -r, --region string            The name of the OSS region in which you have stored your bucket (default "eu-central-1")
      --threads int              The amount of threads to use (default 5)

Use "aliyun-oss-downloader [command] --help" for more information about a command.
```

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

## Author Information

- Timo Runge
