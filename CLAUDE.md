# Development Notes

## AozoraEpub3 source bootstrap

This repository expects the upstream Java source under `temp/AozoraEpub3`.

If `temp/AozoraEpub3` does not exist, clone it with:

```bash
mkdir -p temp
git clone --depth 1 https://github.com/kyukyunyorituryo/AozoraEpub3.git temp/AozoraEpub3
```

## Quality Checks

```bash
go test ./...
go vet ./...
gofmt -l -d .
go tool golangci-lint run ./...
```

