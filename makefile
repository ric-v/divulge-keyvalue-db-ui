build:
	echo 'Compiling...'
	rm -rf bin/*
	rm -rf ui/build/*
	cd ui/ && npm i
	cd ui/ && npm run build
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -X main.Version=${shell cat VERSION}" -o bin/win64/divulge-viewer-${shell cat VERSION}-amd64.exe cmd/main.go
	GOOS=linux   GOARCH=amd64 go build -ldflags "-s -X main.Version=${shell cat VERSION}" -o bin/linux64/divulge-viewer-${shell cat VERSION}-amd64   cmd/main.go
	GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -X main.Version=${shell cat VERSION}" -o bin/darwin64/divulge-viewer-${shell cat VERSION}-amd64   cmd/main.go
	GOOS=windows GOARCH=arm64 go build -ldflags "-s -X main.Version=${shell cat VERSION}" -o bin/win64/divulge-viewer-${shell cat VERSION}-arm64.exe cmd/main.go
	GOOS=linux   GOARCH=arm64 go build -ldflags "-s -X main.Version=${shell cat VERSION}" -o bin/linux64/divulge-viewer-${shell cat VERSION}-arm64   cmd/main.go
	GOOS=darwin  GOARCH=arm64 go build -ldflags "-s -X main.Version=${shell cat VERSION}" -o bin/darwin64/divulge-viewer-${shell cat VERSION}-arm64   cmd/main.go
	echo 'Compiling for windows and linux done.'

package:

	echo 'Packaging...'
	rm -rf pkg/
	mkdir -p pkg

	mkdir -p divulge-viewer-linux-amd64-${shell cat VERSION}/ui/
	mkdir -p divulge-viewer-darwin-amd64-${shell cat VERSION}/ui/
	mkdir -p divulge-viewer-windows-amd64-${shell cat VERSION}/ui/
	mkdir -p divulge-viewer-linux-arm64-${shell cat VERSION}/ui/
	mkdir -p divulge-viewer-darwin-arm64-${shell cat VERSION}/ui/
	mkdir -p divulge-viewer-windows-arm64-${shell cat VERSION}/ui/

	cp bin/darwin64/divulge-viewer-${shell cat VERSION}-amd64   divulge-viewer-darwin-amd64-${shell cat VERSION}/
	cp bin/linux64/divulge-viewer-${shell cat VERSION}-amd64   divulge-viewer-linux-amd64-${shell cat VERSION}/
	cp bin/win64/divulge-viewer-${shell cat VERSION}-amd64.exe divulge-viewer-windows-amd64-${shell cat VERSION}/
	cp bin/darwin64/divulge-viewer-${shell cat VERSION}-arm64   divulge-viewer-darwin-arm64-${shell cat VERSION}/
	cp bin/linux64/divulge-viewer-${shell cat VERSION}-arm64   divulge-viewer-linux-arm64-${shell cat VERSION}/
	cp bin/win64/divulge-viewer-${shell cat VERSION}-arm64.exe divulge-viewer-windows-arm64-${shell cat VERSION}/

	cp -Rf ui/build/ divulge-viewer-linux-amd64-${shell cat VERSION}/ui/
	cp -Rf ui/build/ divulge-viewer-darwin-amd64-${shell cat VERSION}/ui/
	cp -Rf ui/build/ divulge-viewer-windows-amd64-${shell cat VERSION}/ui/
	cp -Rf ui/build/ divulge-viewer-linux-arm64-${shell cat VERSION}/ui/
	cp -Rf ui/build/ divulge-viewer-darwin-arm64-${shell cat VERSION}/ui/
	cp -Rf ui/build/ divulge-viewer-windows-arm64-${shell cat VERSION}/ui/

	zip -r divulge-viewer-darwin-${shell cat VERSION}-amd64.zip  divulge-viewer-darwin-amd64-${shell cat VERSION}
	zip -r divulge-viewer-linux-${shell cat VERSION}-amd64.zip  divulge-viewer-linux-amd64-${shell cat VERSION}
	zip -r divulge-viewer-window-${shell cat VERSION}-amd64.zip divulge-viewer-windows-amd64-${shell cat VERSION}
	zip -r divulge-viewer-darwin-${shell cat VERSION}-arm64.zip  divulge-viewer-darwin-arm64-${shell cat VERSION}
	zip -r divulge-viewer-linux-${shell cat VERSION}-arm64.zip  divulge-viewer-linux-arm64-${shell cat VERSION}
	zip -r divulge-viewer-window-${shell cat VERSION}-arm64.zip divulge-viewer-windows-arm64-${shell cat VERSION}

	tar -czf divulge-viewer-linux-${shell cat VERSION}-amd64.tar.gz  divulge-viewer-linux-amd64-${shell cat VERSION}
	tar -czf divulge-viewer-darwin-${shell cat VERSION}-amd64.tar.gz  divulge-viewer-darwin-amd64-${shell cat VERSION}
	tar -czf divulge-viewer-window-${shell cat VERSION}-amd64.tar.gz divulge-viewer-windows-amd64-${shell cat VERSION}
	tar -czf divulge-viewer-linux-${shell cat VERSION}-arm64.tar.gz  divulge-viewer-linux-arm64-${shell cat VERSION}
	tar -czf divulge-viewer-darwin-${shell cat VERSION}-arm64.tar.gz  divulge-viewer-darwin-arm64-${shell cat VERSION}
	tar -czf divulge-viewer-window-${shell cat VERSION}-arm64.tar.gz divulge-viewer-windows-arm64-${shell cat VERSION}

	mv *.zip pkg/
	mv *.tar.gz pkg/

	rm -rf public-dist/
	rm -rf divulge-viewer-darwin-amd64-${shell cat VERSION}
	rm -rf divulge-viewer-linux-amd64-${shell cat VERSION}
	rm -rf divulge-viewer-windows-amd64-${shell cat VERSION}
	rm -rf divulge-viewer-darwin-arm64-${shell cat VERSION}
	rm -rf divulge-viewer-linux-arm64-${shell cat VERSION}
	rm -rf divulge-viewer-windows-arm64-${shell cat VERSION}


docker_build:
	docker build -t divulge-viewer:latest .

docker_run:
	docker run --name divulge-viewer --rm -it -p 8080:8080 divulge-viewer:latest
