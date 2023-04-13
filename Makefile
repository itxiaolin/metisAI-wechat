BINARY      := openai-wechat
LDFLAGS     := -ldflags "-w -s"
ENV         :=

build:
	[ -n "${ENV}" ] && export ${ENV}; \
	go build -o ${BINARY} ${LDFLAGS} ./; \

linux-env:
	$(eval ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64)
linux:	linux-env build


windows-env:
	$(eval ENV := CGO_ENABLED=0 GOOS=windows GOARCH=arm64)

windows : windows-env build

clean:
	rm -f "${BINARY}"

.PHONY: clean
