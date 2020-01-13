{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "bk-cmdb.name" -}}
{{- default "bk-cmdb" .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "bk-cmdb.fullname" -}}
{{- $name := default "bk-cmdb" .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Helm required labels */}}
{{- define "bk-cmdb.labels" -}}
heritage: {{ .Release.Service }}
release: {{ .Release.Name }}
chart: {{ .Chart.Name }}
app: "{{ template "bk-cmdb.name" . }}"
{{- end -}}

{{/* matchLabels */}}
{{- define "bk-cmdb.matchLabels" -}}
release: {{ .Release.Name }}
app: "{{ template "bk-cmdb.name" . }}"
{{- end -}}

{{- define "bk-cmdb.autoGenCert" -}}
  {{- if and .Values.expose.tls.enabled (not .Values.expose.tls.secretName) -}}
    {{- printf "true" -}}
  {{- else -}}
    {{- printf "false" -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.autoGenCertForIngress" -}}
  {{- if and (eq (include "bk-cmdb.autoGenCert" .) "true") (eq .Values.expose.type "ingress") -}}
    {{- printf "true" -}}
  {{- else -}}
    {{- printf "false" -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.autoGenCertForNginx" -}}
  {{- if and (eq (include "bk-cmdb.autoGenCert" .) "true") (ne .Values.expose.type "ingress") -}}
    {{- printf "true" -}}
  {{- else -}}
    {{- printf "false" -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.host" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- template "bk-cmdb.database" . }}
  {{- else -}}
    {{- .Values.database.external.host -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.port" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "5432" -}}
  {{- else -}}
    {{- .Values.database.external.port -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.username" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "postgres" -}}
  {{- else -}}
    {{- .Values.database.external.username -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.rawPassword" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- .Values.database.internal.password -}}
  {{- else -}}
    {{- .Values.database.external.password -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.encryptedPassword" -}}
  {{- include "bk-cmdb.database.rawPassword" . | b64enc | quote -}}
{{- end -}}

{{- define "bk-cmdb.database.coreDatabase" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "registry" -}}
  {{- else -}}
    {{- .Values.database.external.coreDatabase -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.clairDatabase" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "postgres" -}}
  {{- else -}}
    {{- .Values.database.external.clairDatabase -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.notaryServerDatabase" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "notaryserver" -}}
  {{- else -}}
    {{- .Values.database.external.notaryServerDatabase -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.notarySignerDatabase" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "notarysigner" -}}
  {{- else -}}
    {{- .Values.database.external.notarySignerDatabase -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.sslmode" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "disable" -}}
  {{- else -}}
    {{- .Values.database.external.sslmode -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.database.clair" -}}
postgres://{{ template "bk-cmdb.database.username" . }}:{{ template "bk-cmdb.database.rawPassword" . }}@{{ template "bk-cmdb.database.host" . }}:{{ template "bk-cmdb.database.port" . }}/{{ template "bk-cmdb.database.clairDatabase" . }}?sslmode={{ template "bk-cmdb.database.sslmode" . }}
{{- end -}}

{{- define "bk-cmdb.database.notaryServer" -}}
postgres://{{ template "bk-cmdb.database.username" . }}:{{ template "bk-cmdb.database.rawPassword" . }}@{{ template "bk-cmdb.database.host" . }}:{{ template "bk-cmdb.database.port" . }}/{{ template "bk-cmdb.database.notaryServerDatabase" . }}?sslmode={{ template "bk-cmdb.database.sslmode" . }}
{{- end -}}

{{- define "bk-cmdb.database.notarySigner" -}}
postgres://{{ template "bk-cmdb.database.username" . }}:{{ template "bk-cmdb.database.rawPassword" . }}@{{ template "bk-cmdb.database.host" . }}:{{ template "bk-cmdb.database.port" . }}/{{ template "bk-cmdb.database.notarySignerDatabase" . }}?sslmode={{ template "bk-cmdb.database.sslmode" . }}
{{- end -}}

{{- define "bk-cmdb.redis.host" -}}
  {{- if eq .Values.redis.type "internal" -}}
    {{- template "bk-cmdb.redis" . -}}
  {{- else -}}
    {{- .Values.redis.external.host -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.redis.port" -}}
  {{- if eq .Values.redis.type "internal" -}}
    {{- printf "%s" "6379" -}}
  {{- else -}}
    {{- .Values.redis.external.port -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.redis.coreDatabaseIndex" -}}
  {{- if eq .Values.redis.type "internal" -}}
    {{- printf "%s" "0" }}
  {{- else -}}
    {{- .Values.redis.external.coreDatabaseIndex -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.redis.jobserviceDatabaseIndex" -}}
  {{- if eq .Values.redis.type "internal" -}}
    {{- printf "%s" "1" }}
  {{- else -}}
    {{- .Values.redis.external.jobserviceDatabaseIndex -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.redis.registryDatabaseIndex" -}}
  {{- if eq .Values.redis.type "internal" -}}
    {{- printf "%s" "2" }}
  {{- else -}}
    {{- .Values.redis.external.registryDatabaseIndex -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.redis.chartmuseumDatabaseIndex" -}}
  {{- if eq .Values.redis.type "internal" -}}
    {{- printf "%s" "3" }}
  {{- else -}}
    {{- .Values.redis.external.chartmuseumDatabaseIndex -}}
  {{- end -}}
{{- end -}}

{{- define "bk-cmdb.redis.rawPassword" -}}
  {{- if and (eq .Values.redis.type "external") .Values.redis.external.password -}}
    {{- .Values.redis.external.password -}}
  {{- end -}}
{{- end -}}

{{/*the username redis is used for a placeholder as no username needed in redis*/}}
{{- define "bk-cmdb.redisForJobservice" -}}
  {{- if (include "bk-cmdb.redis.rawPassword" . ) -}}
    {{- printf "redis://redis:%s@%s:%s/%s" (include "bk-cmdb.redis.rawPassword" . ) (include "bk-cmdb.redis.host" . ) (include "bk-cmdb.redis.port" . ) (include "bk-cmdb.redis.jobserviceDatabaseIndex" . ) }}
  {{- else }}
    {{- template "bk-cmdb.redis.host" . }}:{{ template "bk-cmdb.redis.port" . }}/{{ template "bk-cmdb.redis.jobserviceDatabaseIndex" . }}
  {{- end -}}
{{- end -}}

{{/*the username redis is used for a placeholder as no username needed in redis*/}}
{{- define "bk-cmdb.redisForGC" -}}
  {{- if (include "bk-cmdb.redis.rawPassword" . ) -}}
    {{- printf "redis://redis:%s@%s:%s/%s" (include "bk-cmdb.redis.rawPassword" . ) (include "bk-cmdb.redis.host" . ) (include "bk-cmdb.redis.port" . ) (include "bk-cmdb.redis.registryDatabaseIndex" . ) }}
  {{- else }}
    {{- printf "redis://%s:%s/%s" (include "bk-cmdb.redis.host" . ) (include "bk-cmdb.redis.port" . ) (include "bk-cmdb.redis.registryDatabaseIndex" . ) -}}
  {{- end -}}
{{- end -}}

{{/*
host:port,pool_size,password
100 is the default value of pool size
*/}}
{{- define "bk-cmdb.redisForCore" -}}
  {{- template "bk-cmdb.redis.host" . }}:{{ template "bk-cmdb.redis.port" . }},100,{{ template "bk-cmdb.redis.rawPassword" . }}
{{- end -}}

{{- define "bk-cmdb.portal" -}}
  {{- printf "%s-portal" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.core" -}}
  {{- printf "%s-core" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.redis" -}}
  {{- printf "%s-redis" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.jobservice" -}}
  {{- printf "%s-jobservice" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.registry" -}}
  {{- printf "%s-registry" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.adminserver" -}}
  {{- printf "%s-adminserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.apiserver" -}}
  {{- printf "%s-apiserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.auditcontroller" -}}
  {{- printf "%s-auditcontroller" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.datacollection" -}}
  {{- printf "%s-datacollection" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.eventserver" -}}
  {{- printf "%s-eventserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.hostcontroller" -}}
  {{- printf "%s-hostcontroller" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.hostserver" -}}
  {{- printf "%s-hostserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.objectcontroller" -}}
  {{- printf "%s-objectcontroller" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.proccontroller" -}}
  {{- printf "%s-proccontroller" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.toposerver" -}}
  {{- printf "%s-toposerver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.webserver" -}}
  {{- printf "%s-webserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.procserver" -}}
  {{- printf "%s-procserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.taskserver" -}}
  {{- printf "%s-taskserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.tmserver" -}}
  {{- printf "%s-tmserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.operationserver" -}}
  {{- printf "%s-operationserver" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.coreservice" -}}
  {{- printf "%s-coreservice" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.ingress" -}}
  {{- printf "%s-ingress" (include "bk-cmdb.fullname" .) -}}
{{- end -}}
