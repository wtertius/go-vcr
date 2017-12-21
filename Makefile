test:
	go test ./...

check_style: hint vet

hint:
	gohint `find -name '*.go'`
vet:
	go vet `find -name '*.go'`

.PHONY: test check_style hint vet
