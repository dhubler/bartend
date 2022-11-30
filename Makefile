export YANGPATH=$(abspath ./etc/yang)

PKGS = \
    bartend

all: \
	test \
	bartend \
	docs \
	web-deps \
	archive

bartend : bartend-host bartend-pi
bartend-pi : GOOS = linux
bartend-pi : GOARCH = arm

bartend-% :
	GOOS=$(GOOS) GOARCH=$(GOARCH) go install ./src/bartend/cmd/bartend

test :
	go test .

.PHONY: web
web :
	cd web; \
		npx parcel build --public-url /web/ ./index.html

archive : bartend-pi
	! test -d bartend || rm -rf bartend
	mkdir -p bartend/bin bartend/etc/yang
	cp ./bin/linux_arm/bartend ./bartend/bin/bartend
	cp ./etc/bartend.service ./bartend/
	sed -e 's|:8080|:80|' ./etc/bartend.json > ./bartend/etc/bartend.json
	cp ./etc/yang/*.yang ./bartend/etc/yang
	rsync -aRL ./web/ ./bartend/
	tar -czf bartend.tgz bartend

run:
	go run ./cmd/bartend/main.go -config ./etc/bartend.json

fc-yang:
	go run github.com/freeconf/yang/cmd/fc-yang get -dir ./etc/yang

clean:
	rm -rf ./bin/*

docs : docs/api.md

docs/api.md : ./etc/yang/acc.yang
	go run github.com/freeconf/yang/cmd/fc-yang doc -module bartend -title 'Bartend' -f md > $@

