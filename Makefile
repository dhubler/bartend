export GOPATH=$(abspath .)

C2DOC = ./bin/c2-doc

PKGS = \
    bartend

all: \
	test \
	bartend \
	build \
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

archive : bartend-pi ./web/build.html
	! test -d bartend || rm -rf bartend
	mkdir -p bartend/web/images
	rsync -av ./bin/linux_arm/ ./bartend/bin/
	rsync -avR ./etc/ ./bartend/
	rsync -av ./web/images ./bartend/web
	cp \
	  ./docs/api.html \
	  ./docs/bartend-model.svg \
	  ./web/*.js \
	  ./bartend/web
	cp ./web/build.html ./bartend/web/index.html
	tar -cvzf bartend.tgz bartend

build : 
	cd web ; \
		vulcanize index.html > build.html

clean:
	rm -rf ./bin/* ./pkg/*

.PHONY: doc-tools 
doc-tools :	$(C2DOC)

docs : doc-tools
	YANGPATH=./etc/yang $(C2DOC) -module bartend -title 'Bartend'  > docs/bartend-api.html
	YANGPATH=./etc/yang $(C2DOC) -module bartend -tmpl dot > docs/bartend-model.dot
	dot -T svg -o ./docs/bartend-model.svg docs/bartend-model.dot

$(C2DOC) :
	go install ./src/vendor/github.com/c2stack/c2g/cmd/c2-doc
