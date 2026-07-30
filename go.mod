module lesiw.io/boolset

go 1.25.0

require golang.org/x/tools v0.48.0

require (
	github.com/Antonboom/errname v1.1.2 // indirect
	golang.org/x/mod v0.38.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/telemetry v0.0.0-20260708182218-49f421fb7959 // indirect
	lesiw.io/checker v0.16.0 // indirect
	lesiw.io/clerk v0.2.1-0.20260726115654-97532659cfda // indirect
	lesiw.io/errcheck v1.0.0 // indirect
	lesiw.io/linelen v0.4.0 // indirect
	lesiw.io/plscheck v0.20.0 // indirect
	lesiw.io/tidytypes v0.2.0 // indirect
	lesiw.io/tools v0.2.1-0.20260805015709-c336d0b770b5 // indirect
)

tool (
	golang.org/x/tools/cmd/goimports
	lesiw.io/tools/cmd/clerk
	lesiw.io/tools/cmd/vet
)
