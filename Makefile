GO ?= 1.17.2
BRANCH ?= main
VERSION ?= 0.0.0
TARGET ?= webview

help:
	@echo "cli|dev|css|platforms|release"

devcli:
	@go run cmd/pngsource.go

linuxcli:
	@go build -o dist/linux/cli/pngsource cmd/pngsource.go 

windowscli:
	@xgo --branch=$(BRANCH) --go=$(GO) --dest dist/windows/cli --ldflags="-H windowsgui" --pkg cmd --targets=windows/amd64 github.com/fusion/pngsource \
	&& sudo chown -R $$(id -u) dist \
	&& mv dist/windows/cli/cmd-windows-4.0-amd64.exe dist/windows/cli/pngsource.exe

macoscli:
	@xgo --branch=$(BRANCH) --go=$(GO) --dest dist/macos/cli --pkg cmd --targets=darwin/arm64 github.com/fusion/pngsource \
	&& sudo chown -R $$(id -u) dist \
	&& mv dist/macos/cli/cmd-darwin-10.12-arm64 dist/macos/cli/pngsource

cli: linuxcli windowscli macoscli

devweb:
	@go run pngsource/gui.go

css:
	@yarn css

# Assuming Linux... yup.
buildlinuxapp:
	@go build --tags $(TARGET) -o dist/linux/pngsourceapp pngsource/gui.go

linuxapp: buildlinuxapp

buildwindowsapp:
	@xgo --tags="$(TARGET)" --branch=$(BRANCH) --go=$(GO) --dest dist/windows --ldflags="-H windowsgui" --pkg pngsource --targets=windows/amd64 github.com/fusion/pngsource

packagewindowsapp:
	@cp -r packaging/windows/* dist/windows/ \
	&& cd dist/windows \
	&& cat pngsource.nsi.tmpl | sed "s/{{VERSION}}/$(VERSION)/g"  > pngsource.nsi \
	&& makensis pngsource.nsi

windowsapp: buildwindowsapp packagewindowsapp

buildmacosapp:
	@xgo --tags="$(TARGET)" --branch=$(BRANCH) --go=$(GO) --dest dist/macos --pkg pngsource --targets=darwin/arm64 github.com/fusion/pngsource

packagemacosapp:
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

macosapp: buildmacosapp packagemacosapp

app: linuxapp windowsapp macosapp

clean:
	@rm -rf dist/*

collect:
	@mkdir -p dist/release \
	&& zip dist/release/pngsource-cli-linux.zip dist/linux/cli/pngsource \
	&& zip dist/release/pngsource-cli-windows.zip dist/windows/cli/pngsource.exe \
	&& zip dist/release/pngsource-cli-macos.zip dist/macos/cli/pngsource \
	&& cp dist/linux/pngsourceapp dist/release/ \
	&& cp dist/windows/pngsource_installer.exe dist/release/ \
	&& cp dist/macos/PNGSource.dmg dist/release/

release: cli css app collect
