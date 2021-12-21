GO ?= 1.17.2
BRANCH ?= main
VERSION ?= 0.0.0

help:
	@echo "cli|dev|css|platforms|release"

devcli:
	@go run cmd/pngsource.go

cli:
	@go build -o dist/cli/pngsource cmd/pngsource.go 

devweb:
	@go run pngsource/webview.go

css:
	@yarn css

# Assuming Linux... yup.
buildlinux:
	@go build -o dist/linux/pngsourceapp pngsource/webview.go

linux: buildlinux

buildwindows:
	@xgo --branch=$(BRANCH) --go=$(GO) --dest dist/windows --pkg pngsource --targets=windows/amd64 github.com/fusion/pngsource

packagewindows:
	@cp -r packaging/windows/* dist/windows/ \
	&& cd dist/windows \
	&& cat pngsource.nsi.tmpl | sed "s/{{VERSION}}/$(VERSION)/g"  > pngsource.nsi \
	&& makensis pngsource.nsi

windows: buildwindows packagewindows

buildmacos:
	@xgo --branch=$(BRANCH) --go=$(GO) --dest dist/macos --pkg pngsource --targets=darwin/arm64 github.com/fusion/pngsource

packagemacos:
	@rm -rf dist/macos/pngsource.app  \
	&& cp -r packaging/macos/* dist/macos/  \
	&& cat dist/macos/pngsource.app/Contents/Info.plist.tmpl | sed  "s/{{VERSION}}/$(VERSION)/g"  \
		> dist/macos/pngsource.app/Contents/Info.plist \
	&& mkdir -p dist/macos/pngsource.app/Contents/MacOS  \
	&& cp dist/macos/pngsource-darwin-10.12-arm64 dist/macos/pngsource.app/Contents/MacOS/  \
	&&  dd if=/dev/zero of=dist/macos/PNGSource.dmg bs=1M count=6 status=progress  \
	&& mkfs.hfsplus -v PNGSource dist/macos/PNGSource.dmg  \
	&& sudo mkdir -pv /mnt/dmgwork  \
	&& sudo mount -o loop dist/macos/PNGSource.dmg /mnt/dmgwork  \
	&& sudo cp -arv dist/macos/pngsource.app /mnt/dmgwork/  \
	&& sudo umount /mnt/dmgwork

macos: buildmacos packagemacos

platforms: linux windows macos

release: cli css platforms
