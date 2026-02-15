# Development Notes

## Project Overview

AozoraEpub3 (Java) の Go 移植プロジェクト。青空文庫形式テキストを EPUB3 に変換する。

## AozoraEpub3 bootstrap

`temp/AozoraEpub3` または `temp/AozoraEpub3-bin` が存在せずエラーが発生した場合は [docs/aozoraepub3-bootstrap.md](docs/aozoraepub3-bootstrap.md) を参照。

## Quality Checks

```bash
go test ./...
go vet ./...
gofmt -l -d .
go tool golangci-lint run ./...
```
