all:
	CGO_ENABLED=0 GO_EXTLINK_ENABLED=0 gox -osarch="darwin/amd64 linux/amd64 windows/amd64" -output "dist/{{.OS}}_{{.Arch}}_{{.Dir}}"
