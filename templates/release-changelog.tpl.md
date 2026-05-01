# Changelog

## Version History

{{ range .Versions -}}
* [{{ .Tag.Name }}](#{{ .Tag.Name }})
{{ end }}
## Changes
{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
### {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }}

> {{ datetime "2006-01-02" .Tag.Date }}

{{ range .CommitGroups -}}
{{ if eq .Title "Features" }}#### 🚀 Enhancements{{ else if eq .Title "Performance" }}#### 🔥 Performance{{ else if eq .Title "Bug Fixes" }}#### 🩹 Fixes{{ else if eq .Title "Refactoring" }}#### 💅 Refactors{{ else if eq .Title "Documentation" }}#### 📖 Documentation{{ else if eq .Title "Build" }}#### 📦 Build{{ else if eq .Title "Chores" }}#### 🏡 Chore{{ else if eq .Title "Tests" }}#### ✅ Tests{{ else if eq .Title "Styles" }}#### 🎨 Styles{{ end }}

{{ range .Commits -}}
* {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
#### Pull Requests

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
#### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}