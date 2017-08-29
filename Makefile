export GOPATH=$(abspath .)
C2_YANG = ./src/vendor/github.com/c2stack/c2g/yang
export YANGPATH=$(abspath ./etc/yang):$(abspath $(C2_YANG))

C2DOC = ./bin/c2-doc

PKGS = \
    bartend

all: \
	test \
	bartend \
	docs \
	archive

bartend : bartend-host bartend-pi
bartend-pi : GOOS = linux
bartend-pi : GOARCH = arm

bartend-% :
	GOOS=$(GOOS) GOARCH=$(GOARCH) go install ./src/cmd/bartend

Test% :
	$(MAKE) test TEST=$@

TEST='Test*'
test :
	go test $(PKGS) -run $(TEST)

archive : bartend-pi
	! test -d bartend || rm -rf bartend
	mkdir -p bartend/bin bartend/etc/yang
	cp ./bin/linux_arm/bartend ./bartend/bin/bartend
	cp ./etc/bartend.service ./bartend/
	sed -e 's|:8080|:80|' ./etc/bartend.json > ./bartend/etc/bartend.json
	cp ./etc/yang/*.yang ./bartend/etc/yang
	cp $(C2_YANG)/*.yang ./bartend/etc/yang
	rsync -aRL ./web/ ./bartend/
	tar -czf bartend.tgz bartend

clean:
	rm -rf ./bin/* ./pkg/*

.PHONY: doc-tools 
doc-tools :	$(C2DOC)

docs : doc-tools doc-bartend doc-restconf

doc-bartend:
	$(C2DOC) -module bartend -title 'Bartend' -builtin md > docs/api.md
	$(C2DOC) -module bartend --builtin dot > .api.dot
	dot -T svg -o ./docs/api.svg .api.dot

doc-restconf:
	$(C2DOC) -module restconf -title 'RESTCONF' -builtin md > docs/restconf.md
	$(C2DOC) -module restconf -builtin dot > .restconf.dot
	dot -T svg -o ./docs/restconf.svg .restconf.dot

$(C2DOC) :
	go install ./src/vendor/github.com/c2stack/c2g/cmd/c2-doc
