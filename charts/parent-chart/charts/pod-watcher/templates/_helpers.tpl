{{/*
Expand the name of the chart.
*/}}
{{- define "pod-watcher.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "pod-watcher.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "pod-watcher.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end }}
