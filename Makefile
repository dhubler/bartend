export GOPATH=$(abspath .)

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
