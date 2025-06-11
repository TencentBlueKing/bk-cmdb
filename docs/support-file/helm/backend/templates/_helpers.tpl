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
  {{- $base := printf "mongodb://%s:%s@%s/cmdb"
      (include "cmdb.mongodb.usr" . | trim)
      (include "cmdb.mongodb.pwd" . | trim)
      (include "cmdb.mongodb.addr" . | trim)
  -}}
  {{- /* Check if CA certificate is provided, indicating TLS is enabled */ -}}
  {{- if .Values.mongodb.tls.caFile -}}
    {{- $tlsParams := printf "?tls=true&tlsInsecure=%v&tlsCAFile=%s/mongodb/%s"
        .Values.mongodb.tls.insecureSkipVerify
        .Values.certPath
        .Values.mongodb.tls.caFile
    -}}
    {{- /* Check if both client certificate and key are provided for mutual TLS */ -}}
    {{- if and .Values.mongodb.tls.pemFile -}}
      {{- $tlsParams = printf "%s&tlsCertificateKeyFile=%s/mongodb/%s"
          $tlsParams
          .Values.certPath
          .Values.mongodb.tls.pemFile
      -}}
    {{- end -}}
    {{- /* Append TLS parameters to the base URL */ -}}
    {{- printf "%s%s" $base $tlsParams -}}
  {{- else -}}
    {{- /* No CA provided, use non-TLS base URL */ -}}
    {{- $base -}}
  {{- end -}}
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

{{- define "cmdb.redis.certVolumeMount" -}}
{{- if or .Values.redisCert.redis.ca .Values.redisCert.redis.key .Values.redisCert.redis.cert .Values.redisCert.redis.tlsSecretName }}
- name: redis-certs
  mountPath: {{ .Values.certPath }}/redis
{{- end }}
{{- end -}}

{{- define "cmdb.redis.certVolume" -}}
{{- if or .Values.redisCert.redis.ca .Values.redisCert.redis.key .Values.redisCert.redis.cert }}
- name: redis-certs
  configMap:
    name: {{ template "bk-cmdb.fullname" . }}-redis-certs
{{- else if .Values.redisCert.redis.tlsSecretName }}
- name: redis-certs
  secret:
    secretName: {{ .Values.redisCert.redis.tlsSecretName }}
{{- end }}
{{- end -}}

{{- define "cmdb.redis.snapshotCertVolumeMount" -}}
{{- if or .Values.redisCert.snapshotRedis.ca .Values.redisCert.snapshotRedis.key .Values.redisCert.snapshotRedis.cert .Values.redisCert.snapshotRedis.tlsSecretName }}
- name: snapshot-redis-certs
  mountPath: {{ .Values.certPath }}/snapshot-redis
{{- end }}
{{- end -}}

{{- define "cmdb.redis.snapshotCertVolume" -}}
{{- if or .Values.redisCert.snapshotRedis.ca .Values.redisCert.snapshotRedis.key .Values.redisCert.snapshotRedis.cert }}
- name: snapshot-redis-certs
  configMap:
    name: {{ template "bk-cmdb.fullname" . }}-snapshot-redis-certs
{{- else if .Values.redisCert.snapshotRedis.tlsSecretName }}
- name: snapshot-redis-certs
  secret:
    secretName: {{ .Values.redisCert.snapshotRedis.tlsSecretName }}
{{- end }}
{{- end -}}

{{- define "cmdb.imagePullSecrets" -}}
{{- if .Values.image.pullSecretName }}
imagePullSecrets:
- name: {{ .Values.image.pullSecretName }}
{{- end }}
{{- end -}}

{{- define "cmdb.mongodb.certVolumeMount" -}}
{{- if or .Values.mongodbCert.mongodb.cert .Values.mongodbCert.mongodb.key .Values.mongodbCert.mongodb.ca .Values.mongodbCert.mongodb.tlsSecretName }}
- name: mongodb-certs
  mountPath: {{ .Values.certPath }}/mongodb
{{- end }}
{{- end -}}

{{- define "cmdb.mongodb.certVolume" -}}
{{- if or .Values.mongodbCert.mongodb.cert .Values.mongodbCert.mongodb.key .Values.mongodbCert.mongodb.ca }}
- name: mongodb-certs
  configMap:
    name: {{ template "bk-cmdb.fullname" . }}-mongodb-certs
{{- else if .Values.mongodbCert.mongodb.tlsSecretName }}
- name: mongodb-certs
  secret:
    secretName: {{ .Values.mongodbCert.mongodb.tlsSecretName }}
{{- end }}
{{- end -}}

{{- define "cmdb.mongodb.watch.certVolumeMount" -}}
{{- if or .Values.mongodbCert.watch.cert .Values.mongodbCert.watch.key .Values.mongodbCert.watch.ca .Values.mongodbCert.watch.tlsSecretName }}
- name: mongodb-watch-certs
  mountPath: {{ .Values.certPath }}/mongodb-watch
{{- end }}
{{- end -}}

{{- define "cmdb.mongodb.watch.certVolume" -}}
{{- if or .Values.mongodbCert.watch.cert .Values.mongodbCert.watch.key .Values.mongodbCert.watch.ca }}
- name: mongodb-watch-certs
  configMap:
    name: {{ template "bk-cmdb.fullname" . }}-mongodb-watch-certs
{{- else if .Values.mongodbCert.watch.tlsSecretName }}
- name: mongodb-watch-certs
  secret:
    secretName: {{ .Values.mongodbCert.watch.tlsSecretName }}
{{- end }}
{{- end -}}

{{- define  "cmdb.configAndServiceCenter.certVolumeMount" -}}
{{- if or .Values.zookeeperCert.cert .Values.zookeeperCert.key .Values.zookeeperCert.ca .Values.zookeeperCert.tlsSecretName }}
- name: zookeeper-certs
  mountPath: {{ .Values.certPath }}/zookeeper
{{- end }}
{{- end -}}

{{- define  "cmdb.configAndServiceCenter.certVolume" -}}
{{- if or .Values.zookeeperCert.cert .Values.zookeeperCert.key .Values.zookeeperCert.ca }}
- name: zookeeper-certs
  configMap:
    name: {{ template "bk-cmdb.fullname" $ }}-zookeeper-certs
{{- else if .Values.zookeeperCert.tlsSecretName }}
- name: zookeeper-certs
  secret:
    secretName: {{ .Values.zookeeperCert.tlsSecretName }}
{{- end }}
{{- end }}

{{- define "cmdb.configAndServiceCenter.certCommand" -}}
{{- if .Values.configAndServiceCenter.tls.caFile }}
- --regdiscv-cafile={{ .Values.certPath }}/zookeeper/{{ .Values.configAndServiceCenter.tls.caFile }}
- --regdiscv-skipverify={{ .Values.configAndServiceCenter.tls.insecureSkipVerify }}
{{- end }}
{{- if and .Values.configAndServiceCenter.tls.certFile .Values.configAndServiceCenter.tls.keyFile }}
- --regdiscv-certfile={{ .Values.certPath }}/zookeeper/{{ .Values.configAndServiceCenter.tls.certFile }}
- --regdiscv-keyfile={{ .Values.certPath }}/zookeeper/{{ .Values.configAndServiceCenter.tls.keyFile }}
{{- end }}
{{- if .Values.configAndServiceCenter.tls.password }}
- --regdiscv-certpassword={{ .Values.certPath }}/{{ .Values.configAndServiceCenter.tls.password }}
{{- end }}
{{- end -}}

{{- define "cmdb.es.certVolumeMount" -}}
{{- if or .Values.esCert.ca .Values.esCert.cert .Values.esCert.key .Values.esCert.tlsSecretName }}
- name: es-certs
  mountPath: {{ .Values.certPath }}/elasticsearch
{{- end }}
{{- end -}}

{{- define "cmdb.es.certVolume" -}}
{{- if or .Values.esCert.ca .Values.esCert.cert .Values.esCert.key }}
- name: es-certs
  configMap:
    name: {{ template "bk-cmdb.fullname" . }}-elasticsearch-certs
{{- else if .Values.esCert.tlsSecretName }}
- name: es-certs
  secret:
    secretName: {{ .Values.esCert.tlsSecretName }}
{{- end }}
{{- end -}}

{{- define "bk-cmdb.tenantID" -}}
    {{- if eq .Values.tenant.enableMultiTenantMode true -}}
      system
    {{- else -}}
      default
    {{- end -}}
{{- end -}}
