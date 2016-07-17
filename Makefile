export GOPATH=$(abspath .)

PKGS = \
    pi

isaac : isaac-host isaac-pi
isaac-pi : GOOS = linux
isaac-pi : GOARCH = arm

isaac-% :
	GOOS=$(GOOS) GOARCH=$(GOARCH) go install ./src/cmd/isaac

Test% :
	$(MAKE) test TEST=$@

TEST='Test*'
test :
	go test -v $(PKGS) -run $(TEST)

clean:
	rm -rf ./bin/* ./pkg/*
