export GOPATH=$(abspath .)

C2DOC = ./bin/c2-doc

PKGS = \
    bartend

bartend : bartend-host bartend-pi
bartend-pi : GOOS = linux
bartend-pi : GOARCH = arm

bartend-% :
	GOOS=$(GOOS) GOARCH=$(GOARCH) go install ./src/cmd/bartend

Test% :
	$(MAKE) test TEST=$@

TEST='Test*'
test :
	go test -v $(PKGS) -run $(TEST)

clean:
	rm -rf ./bin/* ./pkg/*

docs : $(C2DOC)
	YANGPATH=./etc/yang $(C2DOC) -module bartend -title 'Bartend'  > docs/bartend-api.html
	YANGPATH=./etc/yang $(C2DOC) -module bartend -tmpl dot > docs/bartend-model.dot
	dot -T svg -o ./docs/bartend-model.svg docs/bartend-model.dot

$(C2DOC) :
	go install ./src/vendor/github.com/c2stack/c2g/cmd/c2-doc
