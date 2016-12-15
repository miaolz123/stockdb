#!/bin/sh

gox -os "darwin linux windows" -arch "amd64" -output "{{.Dir}}_{{.OS}}_{{.Arch}}/{{.Dir}}"
