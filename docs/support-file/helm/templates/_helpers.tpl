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
{{- printf "%s" $name | trunc 63 | trimSuffix "-" -}}
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

{{- define "bk-cmdb.adminserver" -}}
  {{- printf "%s-admin" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.apiserver" -}}
  {{- printf "%s-api" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.authserver" -}}
  {{- printf "%s-auth" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.cacheservice" -}}
  {{- printf "%s-cache" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.cloudserver" -}}
  {{- printf "%s-cloud" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.datacollection" -}}
  {{- printf "%s-datacollection" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.eventserver" -}}
  {{- printf "%s-event" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.hostserver" -}}
  {{- printf "%s-host" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.toposerver" -}}
  {{- printf "%s-topo" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.webserver" -}}
  {{- printf "%s-web" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.procserver" -}}
  {{- printf "%s-proc" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.taskserver" -}}
  {{- printf "%s-task" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.operationserver" -}}
  {{- printf "%s-operation" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.coreservice" -}}
  {{- printf "%s-core" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.synchronizeserver" -}}
  {{- printf "%s-synchronize" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{- define "bk-cmdb.ingress" -}}
  {{- printf "%s-ingress" (include "bk-cmdb.fullname" .) -}}
{{- end -}}

{{/*
Create a default fully qualified redis name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "cmdb.redis.fullname" -}}
{{- $name := default "redis" .Values.redis.nameOverride -}}
{{- printf "%s-%s-master" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cmdb.redis.host" -}}
  {{- if eq .Values.redis.enabled true -}}
    {{- template "cmdb.redis.fullname" . -}}:{{- printf "%s" "6379" -}}
  {{- else -}}
    {{- .Values.redis.redis.host -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.redis.pwd" -}}
{{- if .Values.redis.enabled -}}
      {{- required "redis.auth.password is required" .Values.redis.auth.password  -}}
{{- else }}
     {{- .Values.redis.redis.pwd -}}
{{- end -}}
{{- end -}}

{{- define "cmdb.redis.snap.host" -}}
  {{- if eq .Values.redis.enabled true -}}
    {{- template "cmdb.redis.fullname" . -}}:{{- printf "%s" "6379" -}}
  {{- else -}}
    {{- .Values.redis.snapshotRedis.host -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.redis.snap.pwd" -}}
{{- if .Values.redis.enabled -}}
     {{- .Values.redis.auth.password -}}
{{- else }}
     {{- .Values.redis.snapshotRedis.pwd -}}
{{- end -}}
{{- end -}}

{{- define "cmdb.redis.discover.host" -}}
  {{- if eq .Values.redis.enabled true -}}
    {{- template "cmdb.redis.fullname" . -}}:{{- printf "%s" "6379" -}}
  {{- else -}}
    {{- .Values.redis.discoverRedis.host -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.redis.discover.pwd" -}}
{{- if .Values.redis.enabled -}}
     {{- .Values.redis.auth.password -}}
{{- else }}
     {{- .Values.redis.discoverRedis.pwd -}}
{{- end -}}
{{- end -}}

{{- define "cmdb.redis.netcollect.host" -}}
  {{- if eq .Values.redis.enabled true -}}
    {{- template "cmdb.redis.fullname" . -}}:{{- printf "%s" "6379" -}}
  {{- else -}}
    {{- .Values.redis.netCollectRedis.host -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.redis.netcollect.pwd" -}}
{{- if .Values.redis.enabled -}}
     {{- .Values.redis.auth.password -}}
{{- else }}
     {{- .Values.redis.netCollectRedis.pwd -}}
{{- end -}}
{{- end -}}

{{- define "cmdb.zookeeper.fullname" -}}
{{- $name := default "zookeeper" .Values.zookeeper.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cmdb.configAndServiceCenter.addr" -}}
  {{- if eq .Values.zookeeper.enabled true -}}
    {{- template "cmdb.zookeeper.fullname" . -}}:{{- printf "%s" "2181" -}}
  {{- else -}}
    {{- .Values.configAndServiceCenter.addr -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.addr" -}}
  {{- if eq .Values.mongodb.enabled true -}}
    {{ .Release.Name }}-mongodb-0.{{ .Release.Name }}-{{- .Values.mongodb.host -}}:{{- printf "%s" "27017" -}}
  {{- else -}}
    {{- .Values.mongodb.externalMongodb.host -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.usr" -}}
  {{- if eq .Values.mongodb.enabled true -}}
    {{- .Values.mongodb.auth.username -}}
  {{- else -}}
    {{- .Values.mongodb.externalMongodb.usr -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.pwd" -}}
  {{- if eq .Values.mongodb.enabled true -}}
    {{- required "mongodb.auth.password is required" .Values.mongodb.auth.password -}}
  {{- else -}}
    {{- .Values.mongodb.externalMongodb.pwd -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.watch.addr" -}}
  {{- if eq .Values.mongodb.enabled true -}}
    {{ .Release.Name }}-mongodb-0.{{ .Release.Name }}-{{- .Values.mongodb.host -}}:{{- printf "%s" "27017" -}}
  {{- else -}}
    {{- .Values.mongodb.watch.host -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.watch.usr" -}}
  {{- if eq .Values.mongodb.enabled true -}}
    {{- .Values.mongodb.auth.username -}}
  {{- else -}}
    {{- .Values.mongodb.watch.usr -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.watch.pwd" -}}
  {{- if eq .Values.mongodb.enabled true -}}
    {{- .Values.mongodb.auth.password -}}
  {{- else -}}
    {{- .Values.mongodb.watch.pwd -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.mongodb.mongo-url" -}}
    mongodb://{{ include "cmdb.mongodb.usr" . }}:{{ include "cmdb.mongodb.pwd" . }}@{{- template "cmdb.mongodb.addr" . -}}/cmdb
{{- end -}}

{{- define "cmdb.elasticsearch.urlAndPort" -}}
    {{- if eq .Values.elasticsearch.enabled true -}}
      {{- $name := default "elasticsearch" .Values.elasticsearch.nameOverride -}}
      {{- printf "http://%s-%s-coordinating-only" .Release.Name $name | trunc 63 | trimSuffix "-" -}}:{{- printf "%s" "9200" -}}
    {{- else -}}
      {{- .Values.common.es.url -}}
    {{- end -}}
{{- end -}}

{{- define "cmdb.basicImagesAddress" -}}
    {{ .Values.image.registry }}/{{ .Values.migrate.image.repository }}:v{{ default .Chart.AppVersion .Values.migrate.image.tag }}
{{- end -}}

{{- define "cmdb.webserver.bkLoginUrl" -}}
  {{- if eq .Values.common.webServer.login.version "opensource" -}}
    {{- printf "%s" "" -}}
  {{- else -}}
    {{- printf "%s" .Values.bkPaasUrl -}}/login/?app_id=%s&c_url=%s
  {{- end -}}
{{- end -}}

{{- define "cmdb.webserver.bkHttpsLoginUrl" -}}
  {{- if eq .Values.common.webServer.login.version "opensource" -}}
    {{- printf "%s" "" -}}
  {{- else -}}
    {{- printf "%s" .Values.bkPaasUrl -}}/login/?app_id=%s&c_url=%s
  {{- end -}}
{{- end -}}

{{- define "cmdb.webserver.bk_desktop_url" -}}
  {{- if eq .Values.common.webServer.login.version "opensource" -}}
    {{- printf "%s" "" -}}
  {{- else -}}
    {{- printf "%s" .Values.bkPaasUrl -}}
  {{- end -}}
{{- end -}}

{{- define "cmdb.webserver.paas_domain_url" -}}
  {{- if eq .Values.common.webServer.login.version "opensource" -}}
    {{- printf "%s" "" -}}
  {{- else -}}
    {{- printf "%s" .Values.bkComponentApiUrl -}}
  {{- end -}}
{{- end -}}