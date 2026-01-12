::go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
set GOFLAGS=-buildvcs=false
tfplugindocs generate --disable-vcs