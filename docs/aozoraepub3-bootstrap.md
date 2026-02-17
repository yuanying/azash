# AozoraEpub3 Bootstrap

## Source bootstrap

`temp/AozoraEpub3` に上流の Java ソースを配置する。

```bash
mkdir -p temp
git clone --depth 1 https://github.com/kyukyunyorituryo/AozoraEpub3.git temp/AozoraEpub3
```

## Binary bootstrap

互換比較ハーネスは `temp/AozoraEpub3-bin` 配下のバイナリ配布物を使用する。

```bash
mkdir -p temp/AozoraEpub3-bin
gh release download \
  --repo kyukyunyorituryo/AozoraEpub3 \
  --pattern 'AozoraEpub3-*.zip' \
  --dir temp
unzip -oq temp/AozoraEpub3-*.zip -d temp/AozoraEpub3-bin
rm -f temp/AozoraEpub3-*.zip
```
