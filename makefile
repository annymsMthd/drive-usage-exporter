export VERSION=$(shell gogitver)
ARTIFACT_PATH=./artifacts/drive-usage-exporter

.PHONY: clean build ship build-debian-package

build:
	go build -o $(ARTIFACT_PATH) cmd/drive-usage-exporter/main.go

clean:
	rm -Rf ./artifacts

build-debian-package: clean build
	mkdir -p artifacts/debian/DEBIAN
	mkdir -p artifacts/debian/usr/bin/
	mkdir -p artifacts/debian/lib/systemd/system

	cat ./build/debian/control.template | envsubst > ./artifacts/debian/DEBIAN/control
	cp ./build/debian/postinst ./artifacts/debian/DEBIAN
	cp ./build/debian/prerm ./artifacts/debian/DEBIAN
	cp $(ARTIFACT_PATH) artifacts/debian/usr/bin/drive-usage-exporter
	cp build/debian/drive-usage-exporter.service ./artifacts/debian/lib/systemd/system
	dpkg-deb --build ./artifacts/debian ./artifacts/drive-usage-exporter.deb
	rm -R ./artifacts/debian

package: build-debian-package
