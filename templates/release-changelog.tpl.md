{{- $release := index .Versions 0 -}}
## [{{ $release.Tag.Name }}] - {{ $release.Tag.Date.Format "2006-01-02" }}

{{- range $release.CommitGroups }}
### {{ .Title }}
{{- range .Commits }}
- {{ .Subject }}
{{- end }}

{{- end }}
