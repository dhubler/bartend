export YANGPATH=$(abspath ./etc/yang)

test :
	go test .

# run on host in emulation mode
run:
	go run ./cmd/bartend/main.go -config ./etc/bartend.json

# build web source into package
.PHONY: web
web : ./web/dist/index.html
./web/dist/index.html : $(wildcard web/src/*.js) web/index.html $(wildcard web/images/*.*)
	cd web; \
		npx parcel build --public-url /web/ ./index.html
	touch ./web/dist/index.html

# build raspberry pi package
pi : clean
	$(MAKE) web
	mkdir -p bartend/bin bartend/etc/yang
	GOOS=linux GOARCH=arm go build -o ./bartend/bin/bartend ./cmd/bartend
	cp ./etc/bartend.service ./bartend/
	sed -e 's|:8080|:80|' ./etc/bartend.json > ./bartend/etc/bartend.json
	cp ./etc/yang/*.yang ./bartend/etc/yang
	rsync -aRL ./web/dist ./bartend/web
	tar -czf bartend.tgz bartend

clean:
	-rm -rf ./bartend ./web/dist docs/api.md

# REST API docs that host in github
docs : docs/api.md
docs/api.md : ./etc/yang/bartend.yang
	go run github.com/freeconf/yang/cmd/fc-yang \
		doc -module bartend -title 'Bartend' -f md > $@

# update the local copy of the yang model files from freeconf. Only has to
# be run when freeconf dep is updated
update-yang :
	go run github.com/freeconf/yang/cmd/fc-yang get -dir ./etc/yang
