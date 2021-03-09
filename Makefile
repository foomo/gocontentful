spaceid = ${SPACEID}
cmakey = ${CMAKEY}
package = ${PACKAGE}
contenttypes = ${CONTENTTYPES}
all: generate killemptylines gofmt
generate:
	go run cmd/contentfulerm.go -spaceid=$(spaceid) -cmakey=$(cmakey) -package=$(package) -contenttypes=$(contenttypes)
killemptylines:
	find ./generated/$(package) -type f -name '*.go' | xargs sed -i .bak "/^[[:space:]]*$$/d"
	find ./generated/$(package) -type f -name '*.go.bak' | xargs rm
gofmt:
	gofmt -w generated/$(package)/*.go