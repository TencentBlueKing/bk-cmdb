/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package types

import (
	"math/big"
	"time"
)

// MountPropagationMode describes mount propagation.
// +enum
type MountPropagationMode string

// Protocol defines network protocols supported for things like container ports.
// +enum
type Protocol string

// A toleration operator is the set of operators that can be used in a toleration.
// +enum
type TolerationOperator string

// URIScheme identifies the scheme used for connection to a host for Get actions
// +enum
type URIScheme string

// +enum
type TaintEffect string

// PodQOSClass defines the supported qos classes of Pods.
// +enum
type PodQOSClass string

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[ResourceName]Quantity

// PersistentVolumeMode describes how a volume is intended to be consumed, either Block or Filesystem.
type PersistentVolumeMode string

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// +enum
type HostPathType string

// ManagedFieldsOperationType is the type of operation which lead to a ManagedFieldsEntry being created.
type ManagedFieldsOperationType string

// PersistentVolumeAccessMode defines various access modes for PV.
type PersistentVolumeAccessMode string

// These are the valid values for PersistentVolumeAccessMode
const (
	// can be mounted read/write mode to exactly 1 host
	ReadWriteOnce PersistentVolumeAccessMode = "ReadWriteOnce"
	// can be mounted in read-only mode to many hosts
	ReadOnlyMany PersistentVolumeAccessMode = "ReadOnlyMany"
	// can be mounted in read/write mode to many hosts
	ReadWriteMany PersistentVolumeAccessMode = "ReadWriteMany"
	// can be mounted read/write mode to exactly 1 pod
	// cannot be used in combination with other access modes
	ReadWriteOncePod PersistentVolumeAccessMode = "ReadWriteOncePod"
)

// UID is a type that holds unique ID values, including UUIDs.  Because we
// don't ONLY use UUIDs, this is an alias to string.  Being a type captures
// intent and helps make sure that UIDs and names do not get conflated.
type UID string

// StorageMedium defines ways that storage can be allocated to a volume.
type StorageMedium string

// Scale is used for getting and setting the base-10 scaled value.
// Base-2 scales are omitted for mathematical simplicity.
// See Quantity.ScaledValue for more details.
type Scale int32

// Format lists the three possible formattings of a quantity.
type Format string

// +enum
type AzureDataDiskCachingMode string

// +enum
type AzureDataDiskKind string

// IP address information for entries in the (plural) PodIPs field.
// Each entry includes:
//    IP: An IP address allocated to the pod. Routable at least within the cluster.
type PodIP struct {
	// ip is an IP address (IPv4 or IPv6) assigned to the pod
	IP string `json:"ip,omitempty" bson:"ip"`
}

// Volume represents a named volume in a pod that may be accessed by any container in the pod.
type Volume struct {
	// name of the volume.
	// Must be a DNS_LABEL and unique within the pod.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" bson:"name"`
	// volumeSource represents the location and type of the mounted volume.
	// If not specified, the Volume is implied to be an EmptyDir.
	// This implied behavior is deprecated and will be removed in a future version.
	VolumeSource `json:",inline" bson:",inline"`
}

// Represents the source of a volume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	// hostPath represents a pre-existing file or directory on the host
	// machine that is directly exposed to the container. This is generally
	// used for system agents or other privileged things that are allowed
	// to see the host machine. Most containers will NOT need this.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	// ---
	// TODO(jonesdl) We need to restrict who can use host directory mounts and who can/can not
	// mount host directories as read/write.
	// +optional
	HostPath *HostPathVolumeSource `json:"hostPath,omitempty" bson:"hostPath"`
	// emptyDir represents a temporary directory that shares a pod's lifetime.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
	// +optional
	EmptyDir *EmptyDirVolumeSource `json:"emptyDir,omitempty" bson:"emptyDir"`
	// gcePersistentDisk represents a GCE Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	GCEPersistentDisk *GCEPersistentDiskVolumeSource `json:"gcePersistentDisk,omitempty" bson:"gcePersistentDisk"`
	// awsElasticBlockStore represents an AWS Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	AWSElasticBlockStore *AWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty" bson:"awsElasticBlockStore"`
	// gitRepo represents a git repository at a particular revision.
	// DEPRECATED: GitRepo is deprecated. To provision a container with a git repo, mount an
	// EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir
	// into the Pod's container.
	// +optional
	GitRepo *GitRepoVolumeSource `json:"gitRepo,omitempty" bson:"gitRepo"`
	// secret represents a secret that should populate this volume.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	// +optional
	Secret *SecretVolumeSource `json:"secret,omitempty" bson:"secret"`
	// nfs represents an NFS mount on the host that shares a pod's lifetime
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// +optional
	NFS *NFSVolumeSource `json:"nfs,omitempty" bson:"nfs"`
	// iscsi represents an ISCSI Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://examples.k8s.io/volumes/iscsi/README.md
	// +optional
	ISCSI *ISCSIVolumeSource `json:"iscsi,omitempty" bson:"iscsi"`
	// glusterfs represents a Glusterfs mount on the host that shares a pod's lifetime.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md
	// +optional
	Glusterfs *GlusterfsVolumeSource `json:"glusterfs,omitempty" bson:"glusterfs"`
	// persistentVolumeClaimVolumeSource represents a reference to a
	// PersistentVolumeClaim in the same namespace.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty" bson:"persistentVolumeClaim"`
	// rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.
	// More info: https://examples.k8s.io/volumes/rbd/README.md
	// +optional
	RBD *RBDVolumeSource `json:"rbd,omitempty" bson:"rbd"`
	// flexVolume represents a generic volume resource that is
	// provisioned/attached using an exec based plugin.
	// +optional
	FlexVolume *FlexVolumeSource `json:"flexVolume,omitempty" bson:"flexVolume"`
	// cinder represents a cinder volume attached and mounted on kubelets host machine.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	Cinder *CinderVolumeSource `json:"cinder,omitempty" bson:"cinder"`
	// cephFS represents a Ceph FS mount on the host that shares a pod's lifetime
	// +optional
	CephFS *CephFSVolumeSource `json:"cephfs,omitempty" bson:"cephfs"`
	// flocker represents a Flocker volume attached to a kubelet's host machine. This depends on the Flocker control
	// service being running
	// +optional
	Flocker *FlockerVolumeSource `json:"flocker,omitempty" bson:"flocker"`
	// downwardAPI represents downward API about the pod that should populate this volume
	// +optional
	DownwardAPI *DownwardAPIVolumeSource `json:"downwardAPI,omitempty" bson:"downwardAPI"`
	// fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod.
	// +optional
	FC *FCVolumeSource `json:"fc,omitempty" bson:"fc"`
	// azureFile represents an Azure File Service mount on the host and bind mount to the pod.
	// +optional
	AzureFile *AzureFileVolumeSource `json:"azureFile,omitempty" bson:"azureFile"`
	// configMap represents a configMap that should populate this volume
	// +optional
	ConfigMap *ConfigMapVolumeSource `json:"configMap,omitempty" bson:"configMap"`
	// vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	VsphereVolume *VsphereVirtualDiskVolumeSource `json:"vsphereVolume,omitempty" bson:"vsphereVolume"`
	// quobyte represents a Quobyte mount on the host that shares a pod's lifetime
	// +optional
	Quobyte *QuobyteVolumeSource `json:"quobyte,omitempty" bson:"quobyte"`
	// azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
	// +optional
	AzureDisk *AzureDiskVolumeSource `json:"azureDisk,omitempty" bson:"azureDisk"`
	// photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine
	// NOCC:tosa/linelength(忽略长度)
	PhotonPersistentDisk *PhotonPersistentDiskVolumeSource `json:"photonPersistentDisk,omitempty" bson:"photonPersistentDisk"`
	// projected items for all in one resources secrets, configmaps, and downward API
	Projected *ProjectedVolumeSource `json:"projected,omitempty" bson:"projected"`
	// portworxVolume represents a portworx volume attached and mounted on kubelets host machine
	// +optional
	PortworxVolume *PortworxVolumeSource `json:"portworxVolume,omitempty" bson:"portworxVolume"`
	// scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.
	// +optional
	ScaleIO *ScaleIOVolumeSource `json:"scaleIO,omitempty" bson:"scaleIO"`
	// storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes.
	// +optional
	StorageOS *StorageOSVolumeSource `json:"storageos,omitempty" bson:"storageos"`
	// csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI
	// drivers (Beta feature).
	// +optional
	CSI *CSIVolumeSource `json:"csi,omitempty" bson:"csi"`
	// ephemeral represents a volume that is handled by a cluster storage driver.
	// The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,
	// and deleted when the pod is removed.
	//
	// Use this if:
	// a) the volume is only needed while the pod runs,
	// b) features of normal volumes like restoring from snapshot or capacity
	//    tracking are needed,
	// c) the storage driver is specified through a storage class, and
	// d) the storage driver supports dynamic volume provisioning through
	//    a PersistentVolumeClaim (see EphemeralVolumeSource for more
	//    information on the connection between this volume type
	//    and PersistentVolumeClaim).
	//
	// Use PersistentVolumeClaim or one of the vendor-specific
	// APIs for volumes that persist for longer than the lifecycle
	// of an individual pod.
	//
	// Use CSI for light-weight local ephemeral volumes if the CSI driver is meant to
	// be used that way - see the documentation of the driver for
	// more information.
	//
	// A pod can use both types of ephemeral volumes and
	// persistent volumes at the same time.
	//
	// +optional
	Ephemeral *EphemeralVolumeSource `json:"ephemeral,omitempty" bson:"ephemeral"`
}

// Represents an ephemeral volume that is handled by a normal storage driver.
type EphemeralVolumeSource struct {
	// Will be used to create a stand-alone PVC to provision the volume.
	// The pod in which this EphemeralVolumeSource is embedded will be the
	// owner of the PVC, i.e. the PVC will be deleted together with the
	// pod.  The name of the PVC will be `<pod name>-<volume name>` where
	// `<volume name>` is the name from the `PodSpec.Volumes` array
	// entry. Pod validation will reject the pod if the concatenated name
	// is not valid for a PVC (for example, too long).
	//
	// An existing PVC with that name that is not owned by the pod
	// will *not* be used for the pod to avoid using an unrelated
	// volume by mistake. Starting the pod is then blocked until
	// the unrelated PVC is removed. If such a pre-created PVC is
	// meant to be used by the pod, the PVC has to updated with an
	// owner reference to the pod once the pod exists. Normally
	// this should not be necessary, but it may be useful when
	// manually reconstructing a broken cluster.
	//
	// This field is read-only and no changes will be made by Kubernetes
	// to the PVC after it has been created.
	// Required, must not be nil.
	// NOCC:tosa/linelength(忽略长度)
	VolumeClaimTemplate *PersistentVolumeClaimTemplate `json:"volumeClaimTemplate,omitempty" bson:"volumeClaimTemplate"`

	// ReadOnly is tombstoned to show why 2 is a reserved protobuf tag.
	// ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
}

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// users must create.
type ObjectMeta struct {
	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	// +optional
	Name string `json:"name,omitempty" bson:"name"`

	// GenerateName is an optional prefix, used by the server, to generate a unique
	// name ONLY IF the Name field has not been provided.
	// If this field is used, the name returned to the client will be different
	// than the name passed. This value will also be combined with a unique suffix.
	// The provided value has the same validation rules as the Name field,
	// and may be truncated by the length of the suffix required to make the value
	// unique on the server.
	//
	// If this field is specified and the generated name exists, the server will return a 409.
	//
	// Applied only if Name is not specified.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
	// +optional
	GenerateName string `json:"generateName,omitempty" bson:"generateName"`

	// Namespace defines the space within which each name must be unique. An empty namespace is
	// equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for
	// those objects will be empty.
	//
	// Must be a DNS_LABEL.
	// Cannot be updated.
	// More info: http://kubernetes.io/docs/user-guide/namespaces
	// +optional
	Namespace string `json:"namespace,omitempty" bson:"namespace"`

	// Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
	// +optional
	SelfLink string `json:"selfLink,omitempty" bson:"selfLink"`

	// UID is the unique in time and space value for this object. It is typically generated by
	// the server on successful creation of a resource and is not allowed to change on PUT
	// operations.
	//
	// Populated by the system.
	// Read-only.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#uids
	// +optional
	UID UID `json:"uid,omitempty" bson:"uid"`

	// An opaque value that represents the internal version of this object that can
	// be used by clients to determine when objects have changed. May be used for optimistic
	// concurrency, change detection, and the watch operation on a resource or set of resources.
	// Clients must treat these values as opaque and passed unmodified back to the server.
	// They may only be valid for a particular resource or set of resources.
	//
	// Populated by the system.
	// Read-only.
	// Value must be treated as opaque by clients and .
	// NOCC:tosa/linelength(忽略长度)
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" bson:"resourceVersion"`

	// A sequence number representing a specific generation of the desired state.
	// Populated by the system. Read-only.
	// +optional
	Generation int64 `json:"generation,omitempty" bson:"generation"`

	// CreationTimestamp is a timestamp representing the server time when this object was
	// created. It is not guaranteed to be set in happens-before order across separate operations.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	CreationTimestamp Time `json:"creationTimestamp,omitempty" bson:"creationTimestamp"`

	// DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted. This
	// field is set by the server when a graceful deletion is requested by the user, and is not
	// directly settable by a client. The resource is expected to be deleted (no longer visible
	// from resource lists, and not reachable by name) after the time in this field, once the
	// finalizers list is empty. As long as the finalizers list contains items, deletion is blocked.
	// Once the deletionTimestamp is set, this value may not be unset or be set further into the
	// future, although it may be shortened or the resource may be deleted prior to this time.
	// For example, a user may request that a pod is deleted in 30 seconds. The Kubelet will react
	// by sending a graceful termination signal to the containers in the pod. After that 30 seconds,
	// the Kubelet will send a hard termination signal (SIGKILL) to the container and after cleanup,
	// remove the pod from the API. In the presence of network partitions, this object may still
	// exist after this timestamp, until an administrator or automated process can determine the
	// resource is fully terminated.
	// If not set, graceful deletion of the object has not been requested.
	//
	// Populated by the system when a graceful deletion is requested.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	DeletionTimestamp *Time `json:"deletionTimestamp,omitempty" bson:"deletionTimestamp"`

	// Number of seconds allowed for this object to gracefully terminate before
	// it will be removed from the system. Only set when deletionTimestamp is also set.
	// May only be shortened.
	// Read-only.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	DeletionGracePeriodSeconds *int64 `json:"deletionGracePeriodSeconds,omitempty" bson:"deletionGracePeriodSeconds"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty" bson:"labels"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" bson:"annotations"`

	// List of objects depended by this object. If ALL objects in the list have
	// been deleted, this object will be garbage collected. If this object is managed by a controller,
	// then an entry in this list will point to this controller, with the controller field set to true.
	// There cannot be more than one managing controller.
	// +optional
	// +patchMergeKey=uid
	// +patchStrategy=merge
	OwnerReferences []OwnerReference `json:"ownerReferences,omitempty" bson:"ownerReferences"`

	// Must be empty before the object is deleted from the registry. Each entry
	// is an identifier for the responsible component that will remove the entry
	// from the list. If the deletionTimestamp of the object is non-nil, entries
	// in this list can only be removed.
	// Finalizers may be processed and removed in any order.  Order is NOT enforced
	// because it introduces significant risk of stuck finalizers.
	// finalizers is a shared field, any actor with permission can reorder it.
	// If the finalizer list is processed in order, then this can lead to a situation
	// in which the component responsible for the first finalizer in the list is
	// waiting for a signal (field value, external system, or other) produced by a
	// component responsible for a finalizer later in the list, resulting in a deadlock.
	// Without enforced ordering finalizers are free to order amongst themselves and
	// are not vulnerable to ordering changes in the list.
	// +optional
	// +patchStrategy=merge
	Finalizers []string `json:"finalizers,omitempty" bson:"finalizers"`

	// Deprecated: ClusterName is a legacy field that was always cleared by
	// the system and never used; it will be removed completely in 1.25.
	//
	// The name in the go struct is changed to help clients detect
	// accidental use.
	//
	// +optional
	ZZZ_DeprecatedClusterName string `json:"clusterName,omitempty" bson:"clusterName"`

	// ManagedFields maps workflow-id and version to the set of fields
	// that are managed by that workflow. This is mostly for internal
	// housekeeping, and users typically shouldn't need to set or
	// understand this field. A workflow can be the user's name, a
	// controller's name, or the name of a specific apply path like
	// "ci-cd". The set of fields is always in the version that the
	// workflow used when modifying the object.
	//
	// +optional
	ManagedFields []ManagedFieldsEntry `json:"managedFields,omitempty" bson:"managedFields"`
}

// ManagedFieldsEntry is a workflow-id, a FieldSet and the group version of the resource
// that the fieldset applies to.
type ManagedFieldsEntry struct {
	// Manager is an identifier of the workflow managing these fields.
	Manager string `json:"manager,omitempty" bson:"manager"`
	// Operation is the type of operation which lead to this ManagedFieldsEntry being created.
	// The only valid values for this field are 'Apply' and 'Update'.
	// NOCC:tosa/linelength(忽略长度)
	Operation ManagedFieldsOperationType `json:"operation,omitempty" bson:"operation"`
	// APIVersion defines the version of this resource that this field set
	// applies to. The format is "group/version" just like the top-level
	// APIVersion field. It is necessary to track the version of a field
	// set because it cannot be automatically converted.
	APIVersion string `json:"apiVersion,omitempty" bson:"apiVersion"`
	// Time is the timestamp of when the ManagedFields entry was added. The
	// timestamp will also be updated if a field is added, the manager
	// changes any of the owned fields value or removes a field. The
	// timestamp does not update when a field is removed from the entry
	// because another manager took it over.
	// +optional
	Time *Time `json:"time,omitempty" bson:"time"`

	// Fields is tombstoned to show why 5 is a reserved protobuf tag.
	//Fields *Fields `json:"fields,omitempty" protobuf:"bytes,5,opt,name=fields,casttype=Fields"`

	// FieldsType is the discriminator for the different fields format and version.
	// There is currently only one possible value: "FieldsV1"
	FieldsType string `json:"fieldsType,omitempty" bson:"fieldsType"`
	// FieldsV1 holds the first JSON version format as described in the "FieldsV1" type.
	// +optional
	FieldsV1 *FieldsV1 `json:"fieldsV1,omitempty" bson:"fieldsV1"`

	// Subresource is the name of the subresource used to update that object, or
	// empty string if the object was updated through the main resource. The
	// value of this field is used to distinguish between managers, even if they
	// share the same name. For example, a status update will be distinct from a
	// regular update using the same manager name.
	// Note that the APIVersion field is not related to the Subresource field and
	// it always corresponds to the version of the main resource.
	Subresource string `json:"subresource,omitempty" bson:"subresource"`
}

// FieldsV1 stores a set of fields in a data structure like a Trie, in JSON format.
//
// Each key is either a '.' representing the field itself, and will always map to an empty set,
// or a string representing a sub-field or item. The string will follow one of these four formats:
// 'f:<name>', where <name> is the name of a field in a struct, or key in a map
// 'v:<value>', where <value> is the exact json formatted value of a list item
// 'i:<index>', where <index> is position of a item in a list
// 'k:<keys>', where <keys> is a map of  a list item's key fields to their unique values
// If a key maps to an empty Fields value, the field that key represents is part of the set.
//
// The exact format is defined in sigs.k8s.io/structured-merge-diff
// +protobuf.options.(gogoproto.goproto_stringer)=false
type FieldsV1 struct {
	// Raw is the underlying serialization of this object.
	Raw []byte `json:"-"`
}

// OwnerReference contains enough information to let you identify an owning
// object. An owning object must be in the same namespace as the dependent, or
// be cluster-scoped, so there is no namespace field.
// +structType=atomic
type OwnerReference struct {
	// API version of the referent.
	APIVersion string `json:"apiVersion" bson:"apiVersion"`
	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind" bson:"kind"`
	// Name of the referent.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name" bson:"name"`
	// UID of the referent.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#uids
	UID UID `json:"uid" bson:"uid"`
	// If true, this reference points to the managing controller.
	// +optional
	Controller *bool `json:"controller,omitempty" bson:"controller"`
	// If true, AND if the owner has the "foregroundDeletion" finalizer, then
	// the owner cannot be deleted from the key-value store until this
	// reference is removed.
	// See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion
	// for how the garbage collector interacts with this field and enforces the foreground deletion.
	// Defaults to false.
	// To set this field, a user needs "delete" permission of the owner,
	// otherwise 422 (Unprocessable Entity) will be returned.
	// +optional
	BlockOwnerDeletion *bool `json:"blockOwnerDeletion,omitempty" bson:"blockOwnerDeletion"`
}

// Time is a wrapper around time.Time which supports correct
// marshaling to YAML and JSON.  Wrappers are provided for many
// of the factory methods that the time package offers.
//
// +protobuf.options.marshal=false
// +protobuf.as=Timestamp
// +protobuf.options.(gogoproto.goproto_stringer)=false
type Time struct {
	time.Time
}

// PersistentVolumeClaimTemplate is used to produce
// PersistentVolumeClaim objects as part of an EphemeralVolumeSource.
type PersistentVolumeClaimTemplate struct {
	// May contain labels and annotations that will be copied into the PVC
	// when creating it. No other fields are allowed and will be rejected during
	// validation.
	//
	// +optional
	ObjectMeta `json:"metadata,omitempty" bson:"metadata"`

	// The specification for the PersistentVolumeClaim. The entire content is
	// copied unchanged into the PVC that gets created from this
	// template. The same fields as in a PersistentVolumeClaim
	// are also valid here.
	Spec PersistentVolumeClaimSpec `json:"spec" bson:"spec"`
}

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// +optional
	Limits ResourceList `json:"limits,omitempty" bson:"limits"`
	// Requests describes the minimum amount of compute resources required.
	// If Request is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value
	// +optional
	Requests ResourceList `json:"requests,omitempty" bson:"requests"`
}

// PersistentVolumeClaimSpec describes the common attributes of storage devices
// and allows a Source for provider-specific attributes
type PersistentVolumeClaimSpec struct {
	// Contains the types of access modes required
	// +optional
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" bson:"accessModes"`
	// A label query over volumes to consider for binding. This selector is
	// ignored when VolumeName is set
	// +optional
	Selector *LabelSelector `json:"selector,omitempty" bson:"selector"`
	// Resources represents the minimum resources required
	// If RecoverVolumeExpansionFailure feature is enabled users are allowed to specify resource requirements
	// that are lower than previous value but must still be higher than capacity recorded in the
	// status field of the claim.
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty" bson:"resources"`
	// VolumeName is the binding reference to the PersistentVolume backing this
	// claim. When set to non-empty value Selector is not evaluated
	// +optional
	VolumeName string `json:"volumeName,omitempty" bson:"volumeName"`
	// Name of the StorageClass required by the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#class-1
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty" bson:"storageClassName"`
	// volumeMode defines what type of volume is required by the claim.
	// Value of Filesystem is implied when not included in claim spec.
	// +optional
	VolumeMode *PersistentVolumeMode `json:"volumeMode,omitempty" bson:"volumeMode"`
	// This field can be used to specify either:
	// * An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot)
	// * An existing PVC (PersistentVolumeClaim)
	// If the provisioner or an external controller can support the specified data source,
	// it will create a new volume based on the contents of the specified data source.
	// If the AnyVolumeDataSource feature gate is enabled, this field will always have
	// the same contents as the DataSourceRef field.
	// +optional
	DataSource *TypedLocalObjectReference `json:"dataSource,omitempty" bson:"dataSource"`
	// Specifies the object from which to populate the volume with data, if a non-empty
	// volume is desired. This may be any local object from a non-empty API group (non
	// core object) or a PersistentVolumeClaim object.
	// When this field is specified, volume binding will only succeed if the type of
	// the specified object matches some installed volume populator or dynamic
	// provisioner.
	// This field will replace the functionality of the DataSource field and as such
	// if both fields are non-empty, they must have the same value. For backwards
	// compatibility, both fields (DataSource and DataSourceRef) will be set to the same
	// value automatically if one of them is empty and the other is non-empty.
	// There are two important differences between DataSource and DataSourceRef:
	// * While DataSource only allows two specific types of objects, DataSourceRef
	//   allows any non-core object, as well as PersistentVolumeClaim objects.
	// * While DataSource ignores disallowed values (dropping them), DataSourceRef
	//   preserves all values, and generates an error if a disallowed value is
	//   specified.
	// (Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled.
	// +optional
	DataSourceRef *TypedLocalObjectReference `json:"dataSourceRef,omitempty" bson:"dataSourceRef"`
}

// TypedLocalObjectReference contains enough information to
// let you locate the typed referenced object inside the same namespace.
type TypedLocalObjectReference struct {
	// APIGroup is the group for the resource being referenced.
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	// +optional
	APIGroup *string `json:"apiGroup" bson:"apiGroup"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind" bson:"kind"`
	// Name is the name of resource being referenced
	Name string `json:"name" bson:"name"`
}

// Represents a source location of a volume to mount, managed by an external CSI driver
type CSIVolumeSource struct {
	// driver is the name of the CSI driver that handles this volume.
	// Consult with your admin for the correct name as registered in the cluster.
	Driver string `json:"driver" bson:"driver"`

	// readOnly specifies a read-only configuration for the volume.
	// Defaults to false (read/write).
	// +optional
	ReadOnly *bool `json:"readOnly,omitempty" bson:"readOnly"`

	// fsType to mount. Ex. "ext4", "xfs", "ntfs".
	// If not provided, the empty value is passed to the associated CSI driver
	// which will determine the default filesystem to apply.
	// +optional
	FSType *string `json:"fsType,omitempty" bson:"fsType"`

	// volumeAttributes stores driver-specific properties that are passed to the CSI
	// driver. Consult your driver's documentation for supported values.
	// +optional
	VolumeAttributes map[string]string `json:"volumeAttributes,omitempty" bson:"volumeAttributes"`

	// nodePublishSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// NodePublishVolume and NodeUnpublishVolume calls.
	// This field is optional, and  may be empty if no secret is required. If the
	// secret object contains more than one secret, all secret references are passed.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	NodePublishSecretRef *LocalObjectReference `json:"nodePublishSecretRef,omitempty" bson:"nodePublishSecretRef"`
}

// Represents a StorageOS persistent volume resource.
type StorageOSVolumeSource struct {
	// volumeName is the human-readable name of the StorageOS volume.  Volume
	// names are only unique within a namespace.
	VolumeName string `json:"volumeName,omitempty" bson:"volumeName"`
	// volumeNamespace specifies the scope of the volume within StorageOS.  If no
	// namespace is specified then the Pod's namespace will be used.  This allows the
	// Kubernetes name scoping to be mirrored within StorageOS for tighter integration.
	// Set VolumeName to any name to override the default behaviour.
	// Set to "default" if you are not using namespaces within StorageOS.
	// Namespaces that do not pre-exist within StorageOS will be created.
	// +optional
	VolumeNamespace string `json:"volumeNamespace,omitempty" bson:"volumeNamespace"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
	// secretRef specifies the secret to use for obtaining the StorageOS API
	// credentials.  If not specified, default values will be attempted.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" bson:"secretRef"`
}

// ScaleIOVolumeSource represents a persistent ScaleIO volume
type ScaleIOVolumeSource struct {
	// gateway is the host address of the ScaleIO API Gateway.
	Gateway string `json:"gateway" bson:"gateway"`
	// system is the name of the storage system as configured in ScaleIO.
	System string `json:"system" bson:"system"`
	// secretRef references to the secret for ScaleIO user and other
	// sensitive information. If this is not provided, Login operation will fail.
	SecretRef *LocalObjectReference `json:"secretRef" bson:"secretRef"`
	// sslEnabled Flag enable/disable SSL communication with Gateway, default false
	// +optional
	SSLEnabled bool `json:"sslEnabled,omitempty" bson:"sslEnabled"`
	// protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.
	// +optional
	ProtectionDomain string `json:"protectionDomain,omitempty" bson:"protectionDomain"`
	// storagePool is the ScaleIO Storage Pool associated with the protection domain.
	// +optional
	StoragePool string `json:"storagePool,omitempty" bson:"storagePool"`
	// storageMode indicates whether the storage for a volume should be ThickProvisioned or ThinProvisioned.
	// Default is ThinProvisioned.
	// +optional
	StorageMode string `json:"storageMode,omitempty" bson:"storageMode"`
	// volumeName is the name of a volume already created in the ScaleIO system
	// that is associated with this volume source.
	VolumeName string `json:"volumeName,omitempty" bson:"volumeName"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs".
	// Default is "xfs".
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// PortworxVolumeSource represents a Portworx volume resource.
type PortworxVolumeSource struct {
	// volumeID uniquely identifies a Portworx volume
	VolumeID string `json:"volumeID" bson:"volumeID"`
	// fSType represents the filesystem type to mount
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs". Implicitly inferred to be "ext4" if unspecified.
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Projection that may be projected along with other supported volume types
type VolumeProjection struct {
	// all types below are the supported types for projection into the same volume

	// secret information about the secret data to project
	// +optional
	Secret *SecretProjection `json:"secret,omitempty" bson:"secret"`
	// downwardAPI information about the downwardAPI data to project
	// +optional
	DownwardAPI *DownwardAPIProjection `json:"downwardAPI,omitempty" bson:"downwardAPI"`
	// configMap information about the configMap data to project
	// +optional
	ConfigMap *ConfigMapProjection `json:"configMap,omitempty" bson:"configMap"`
	// serviceAccountToken is information about the serviceAccountToken data to project
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	ServiceAccountToken *ServiceAccountTokenProjection `json:"serviceAccountToken,omitempty" bson:"serviceAccountToken"`
}

// ServiceAccountTokenProjection represents a projected service account token
// volume. This projection can be used to insert a service account token into
// the pods runtime filesystem for use against APIs (Kubernetes API Server or
// otherwise).
type ServiceAccountTokenProjection struct {
	// audience is the intended audience of the token. A recipient of a token
	// must identify itself with an identifier specified in the audience of the
	// token, and otherwise should reject the token. The audience defaults to the
	// identifier of the apiserver.
	//+optional
	Audience string `json:"audience,omitempty" bson:"audience"`
	// expirationSeconds is the requested duration of validity of the service
	// account token. As the token approaches expiration, the kubelet volume
	// plugin will proactively rotate the service account token. The kubelet will
	// start trying to rotate the token if the token is older than 80 percent of
	// its time to live or if the token is older than 24 hours.Defaults to 1 hour
	// and must be at least 10 minutes.
	//+optional
	ExpirationSeconds *int64 `json:"expirationSeconds,omitempty" bson:"expirationSeconds"`
	// path is the path relative to the mount point of the file to project the
	// token into.
	Path string `json:"path" bson:"path"`
}

type ConfigMapProjection struct {
	LocalObjectReference `json:",inline" bson:",inline"`
	// items if unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KeyToPath `json:"items,omitempty" bson:"items"`
	// optional specify whether the ConfigMap or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" bson:"optional"`
}

// Represents downward API info for projecting into a projected volume.
// Note that this is identical to a downwardAPI volume source without the default
// mode.
type DownwardAPIProjection struct {
	// Items is a list of DownwardAPIVolume file
	// +optional
	Items []DownwardAPIVolumeFile `json:"items,omitempty" bson:"items"`
}

// Adapts a secret into a projected volume.
//
// The contents of the target Secret's Data field will be presented in a
// projected volume as files using the keys in the Data field as the file names.
// Note that this is identical to a secret volume source without the default
// mode.
type SecretProjection struct {
	LocalObjectReference `json:",inline" bson:",inline"`
	// items if unspecified, each key-value pair in the Data field of the referenced
	// Secret will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the Secret,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KeyToPath `json:"items,omitempty" bson:"items"`
	// optional field specify whether the Secret or its key must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" bson:"optional"`
}

// Represents a projected volume source
type ProjectedVolumeSource struct {
	// sources is the list of volume projections
	// +optional
	Sources []VolumeProjection `json:"sources" bson:"sources"`
	// defaultMode are the mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" bson:"defaultMode"`
}

// Represents a Photon Controller persistent disk resource.
type PhotonPersistentDiskVolumeSource struct {
	// pdID is the ID that identifies Photon Controller persistent disk
	PdID string `json:"pdID" bson:"pdID"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	FSType string `json:"fsType,omitempty" bson:"fsType"`
}

// AzureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
type AzureDiskVolumeSource struct {
	// diskName is the Name of the data disk in the blob storage
	DiskName string `json:"diskName" bson:"diskName"`
	// diskURI is the URI of data disk in the blob storage
	DataDiskURI string `json:"diskURI" bson:"diskURI"`
	// cachingMode is the Host Caching mode: None, Read Only, Read Write.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	CachingMode *AzureDataDiskCachingMode `json:"cachingMode,omitempty" bson:"cachingMode"`
	// fsType is Filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType *string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly *bool `json:"readOnly,omitempty" bson:"readOnly"`
	// kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob disk
	// per storage account  Managed: azure managed data disk (only in managed availability set). defaults to shared
	Kind *AzureDataDiskKind `json:"kind,omitempty" bson:"kind"`
}

// Represents a Quobyte mount that lasts the lifetime of a pod.
// Quobyte volumes do not support ownership management or SELinux relabeling.
type QuobyteVolumeSource struct {
	// registry represents a single or multiple Quobyte Registry services
	// specified as a string as host:port pair (multiple entries are separated with commas)
	// which acts as the central registry for volumes
	Registry string `json:"registry" bson:"registry"`

	// volume is a string that references an already created Quobyte volume by name.
	Volume string `json:"volume" bson:"volume"`

	// readOnly here will force the Quobyte volume to be mounted with read-only permissions.
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`

	// user to map volume access to
	// Defaults to serivceaccount user
	// +optional
	User string `json:"user,omitempty" bson:"user"`

	// group to map volume access to
	// Default is no group
	// +optional
	Group string `json:"group,omitempty" bson:"group"`

	// tenant owning the given Quobyte volume in the Backend
	// Used with dynamically provisioned Quobyte volumes, value is set by the plugin
	// +optional
	Tenant string `json:"tenant,omitempty" bson:"tenant"`
}

// HTTPGetAction describes an action based on HTTP Get requests.
type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	Path string `json:"path,omitempty" bson:"path"`
	// Name or number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	Port IntOrString `json:"port" bson:"port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional
	Host string `json:"host,omitempty" bson:"host"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +optional
	Scheme URIScheme `json:"scheme,omitempty" bson:"scheme"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty" bson:"httpHeaders"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
type HTTPHeader struct {
	// The header field name
	Name string `json:"name" bson:"name"`
	// The header field value
	Value string `json:"value" bson:"value"`
}

// ExecAction describes a "run in container" action.
type ExecAction struct {
	// Command is the command line to execute inside the container, the working directory for the
	// command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
	// not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
	// a shell, you need to explicitly call out to that shell.
	// Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
	// +optional
	Command []string `json:"command,omitempty" bson:"command"`
}

// TCPSocketAction describes an action based on opening a socket
type TCPSocketAction struct {
	// Number or name of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	Port IntOrString `json:"port" bson:"port"`
	// Optional: Host name to connect to, defaults to the pod IP.
	// +optional
	Host string `json:"host,omitempty" bson:"host"`
}

// ProbeHandler defines a specific action that should be taken in a probe.
// One and only one of the fields must be specified.
type ProbeHandler struct {
	// Exec specifies the action to take.
	// +optional
	Exec *ExecAction `json:"exec,omitempty" bson:"exec"`
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty" bson:"httpGet"`
	// TCPSocket specifies an action involving a TCP port.
	// +optional
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty" bson:"tcpSocket"`

	// GRPC specifies an action involving a GRPC port.
	// This is a beta field and requires enabling GRPCContainerProbe feature gate.
	// +featureGate=GRPCContainerProbe
	// +optional
	GRPC *GRPCAction `json:"grpc,omitempty" bson:"grpc"`
}

// EnvVar represents an environment variable present in a Container.
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	Name string `json:"name" bson:"name"`

	// Optional: no more than one of the following may be specified.

	// Variable references $(VAR_NAME) are expanded
	// using the previously defined environment variables in the container and
	// any service environment variables. If a variable cannot be resolved,
	// the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e.
	// "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)".
	// Escaped references will never be expanded, regardless of whether the variable
	// exists or not.
	// Defaults to "".
	// +optional
	Value string `json:"value,omitempty" bson:"value"`
	// Source for the environment variable's value. Cannot be used if value is not empty.
	// +optional
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty" bson:"valueFrom"`
}

// EnvVarSource represents a source for the value of an EnvVar.
type EnvVarSource struct {
	// Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`,
	// `metadata.annotations['<KEY>']`, spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP,
	// status.podIPs.
	// +optional
	FieldRef *ObjectFieldSelector `json:"fieldRef,omitempty" bson:"fieldRef"`
	// Selects a resource of the container: only resources limits and requests
	// (limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and
	// requests.ephemeral-storage) are currently supported.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	ResourceFieldRef *ResourceFieldSelector `json:"resourceFieldRef,omitempty" bson:"resourceFieldRef"`
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty" bson:"configMapKeyRef"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty" bson:"secretKeyRef"`
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	// This must match the Name of a Volume.
	Name string `json:"name" bson:"name"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath" bson:"mountPath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty" bson:"subPath"`
	// mountPropagation determines how mounts are propagated from the host
	// to container and the other way around.
	// When not set, MountPropagationNone is used.
	// This field is beta in 1.10.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	MountPropagation *MountPropagationMode `json:"mountPropagation,omitempty" bson:"mountPropagation"`
	// Expanded path within the volume from which the container's volume should be mounted.
	// Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the
	// container's environment. Defaults to "" (volume's root). SubPathExpr and SubPath are mutually exclusive.
	// +optional
	SubPathExpr string `json:"subPathExpr,omitempty" bson:"subPathExpr"`
}

// SecretKeySelector selects a key of a Secret.
// +structType=atomic
type SecretKeySelector struct {
	// The name of the secret in the pod's namespace to select from.
	LocalObjectReference `json:",inline" bson:",inline"`
	// The key of the secret to select from.  Must be a valid secret key.
	Key string `json:"key" bson:"key"`
	// Specify whether the Secret or its key must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" bson:"optional"`
}

// Selects a key from a ConfigMap.
// +structType=atomic
type ConfigMapKeySelector struct {
	// The ConfigMap to select from.
	LocalObjectReference `json:",inline" bson:",inline"`
	// The key to select.
	Key string `json:"key" bson:"key"`
	// Specify whether the ConfigMap or its key must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" bson:"optional"`
}

// GRPCAction grpc service
type GRPCAction struct {
	// Port number of the gRPC service. Number must be in the range 1 to 65535.
	Port int32 `json:"port" bson:"port"`

	// Service is the name of the service to place in the gRPC HealthCheckRequest
	// (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
	//
	// If this is not specified, the default behavior is defined by gRPC.
	// +optional
	// +default=""
	Service *string `json:"service" bson:"service"`
}

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
type Probe struct {
	// The action taken to determine the health of a container
	ProbeHandler `json:",inline" bson:",inline"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty" bson:"initialDelaySeconds"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty" bson:"timeoutSeconds"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty" bson:"periodSeconds"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty" bson:"successThreshold"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty" bson:"failureThreshold"`
	// Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
	// value overrides the value provided by the pod spec.
	// Value must be non-negative integer. The value zero indicates stop immediately via
	// the kill signal (no opportunity to shut down).
	// This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
	// Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty" bson:"terminationGracePeriodSeconds"`
}

// Represents a vSphere volume resource.
type VsphereVirtualDiskVolumeSource struct {
	// volumePath is the path that identifies vSphere volume vmdk
	VolumePath string `json:"volumePath" bson:"volumePath"`
	// fsType is filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// storagePolicyName is the storage Policy Based Management (SPBM) profile name.
	// +optional
	StoragePolicyName string `json:"storagePolicyName,omitempty" bson:"storagePolicyName"`
	// storagePolicyID is the storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName.
	// +optional
	StoragePolicyID string `json:"storagePolicyID,omitempty" bson:"storagePolicyID"`
}

// Adapts a ConfigMap into a volume.
//
// The contents of the target ConfigMap's Data field will be presented in a
// volume as files using the keys in the Data field as the file names, unless
// the items element is populated with specific mappings of keys to paths.
// ConfigMap volumes support ownership management and SELinux relabeling.
type ConfigMapVolumeSource struct {
	LocalObjectReference `json:",inline" bson:",inline"`
	// items if unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KeyToPath `json:"items,omitempty" bson:"items"`
	// defaultMode is optional: mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" bson:"defaultMode"`
	// optional specify whether the ConfigMap or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" bson:"optional"`
}

// AzureFile represents an Azure File Service mount on the host and bind mount to the pod.
type AzureFileVolumeSource struct {
	// secretName is the  name of secret that contains Azure Storage Account Name and Key
	SecretName string `json:"secretName" bson:"secretName"`
	// shareName is the azure share Name
	ShareName string `json:"shareName" bson:"shareName"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Represents a Fibre Channel volume.
// Fibre Channel volumes can only be mounted as read/write once.
// Fibre Channel volumes support ownership management and SELinux relabeling.
type FCVolumeSource struct {
	// targetWWNs is Optional: FC target worldwide names (WWNs)
	// +optional
	TargetWWNs []string `json:"targetWWNs,omitempty" bson:"targetWWNs"`
	// lun is Optional: FC target lun number
	// +optional
	Lun *int32 `json:"lun,omitempty" bson:"lun"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
	// wwids Optional: FC volume world wide identifiers (wwids)
	// Either wwids or combination of targetWWNs and lun must be set, but not both simultaneously.
	// +optional
	WWIDs []string `json:"wwids,omitempty" bson:"wwids"`
}

// ObjectFieldSelector selects an APIVersioned field of an object.
// +structType=atomic
type ObjectFieldSelector struct {
	// Version of the schema the FieldPath is written in terms of, defaults to "v1".
	// +optional
	APIVersion string `json:"apiVersion,omitempty" bson:"apiVersion"`
	// Path of the field to select in the specified API version.
	FieldPath string `json:"fieldPath" bson:"fieldPath"`
}

// ResourceFieldSelector represents container resources (cpu, memory) and their output format
// +structType=atomic
type ResourceFieldSelector struct {
	// Container name: required for volumes, optional for env vars
	// +optional
	ContainerName string `json:"containerName,omitempty" bson:"containerName"`
	// Required: resource to select
	Resource string `json:"resource" bson:"resource"`
	// Specifies the output format of the exposed resources, defaults to "1"
	// +optional
	Divisor Quantity `json:"divisor,omitempty" bson:"divisor"`
}

// DownwardAPIVolumeFile represents information to create the file containing the pod field
type DownwardAPIVolumeFile struct {
	// Required: Path is  the relative path name of the file to be created. Must not be absolute or contain the '..'
	// path. Must be utf-8 encoded. The first item of the relative path must not start with '..'
	Path string `json:"path" bson:"path"`
	// Required: Selects a field of the pod: only annotations, labels, name and namespace are supported.
	// +optional
	FieldRef *ObjectFieldSelector `json:"fieldRef,omitempty" bson:"fieldRef"`
	// Selects a resource of the container: only resources limits and requests
	// (limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	ResourceFieldRef *ResourceFieldSelector `json:"resourceFieldRef,omitempty" bson:"resourceFieldRef"`
	// Optional: mode bits used to set permissions on this file, must be an octal value
	// between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	Mode *int32 `json:"mode,omitempty" bson:"mode"`
}

// DownwardAPIVolumeSource represents a volume containing downward API info.
// Downward API volumes support ownership management and SELinux relabeling.
type DownwardAPIVolumeSource struct {
	// Items is a list of downward API volume file
	// +optional
	Items []DownwardAPIVolumeFile `json:"items,omitempty" bson:"items"`
	// Optional: mode bits to use on created files by default. Must be a
	// Optional: mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" bson:"defaultMode"`
}

// Represents a Flocker volume mounted by the Flocker agent.
// One and only one of datasetName and datasetUUID should be set.
// Flocker volumes do not support ownership management or SELinux relabeling.
type FlockerVolumeSource struct {
	// datasetName is Name of the dataset stored as metadata -> name on the dataset for Flocker
	// should be considered as deprecated
	// +optional
	DatasetName string `json:"datasetName,omitempty" bson:"datasetName"`
	// datasetUUID is the UUID of the dataset. This is unique identifier of a Flocker dataset
	// +optional
	DatasetUUID string `json:"datasetUUID,omitempty" bson:"datasetUUID"`
}

// Represents a host path mapped into a pod.
// Host path volumes do not support ownership management or SELinux relabeling.
type HostPathVolumeSource struct {
	// path of the directory on the host.
	// If the path is a symlink, it will follow the link to the real path.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	Path string `json:"path" bson:"path"`
	// type for HostPath Volume
	// Defaults to ""
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	// +optional
	Type *HostPathType `json:"type,omitempty" bson:"type"`
}

// Represents an empty directory for a pod.
// Empty directory volumes support ownership management and SELinux relabeling.
type EmptyDirVolumeSource struct {
	// medium represents what type of storage medium should back this directory.
	// The default is "" which means to use the node's default medium.
	// Must be an empty string (default) or Memory.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
	// +optional
	Medium StorageMedium `json:"medium,omitempty" bson:"medium"`
	// sizeLimit is the total amount of local storage required for this EmptyDir volume.
	// The size limit is also applicable for memory medium.
	// The maximum usage on memory medium EmptyDir would be the minimum value between
	// the SizeLimit specified here and the sum of memory limits of all containers in a pod.
	// The default is nil which means that the limit is undefined.
	// More info: http://kubernetes.io/docs/user-guide/volumes#emptydir
	// +optional
	SizeLimit *Quantity `json:"sizeLimit,omitempty" bson:"sizeLimit"`
}

// Quantity is a fixed-point representation of a number.
// It provides convenient marshaling/unmarshaling in JSON and YAML,
// in addition to String() and AsInt64() accessors.
//
// The serialization format is:
//
// <quantity>        ::= <signedNumber><suffix>
//   (Note that <suffix> may be empty, from the "" case in <decimalSI>.)
// <digit>           ::= 0 | 1 | ... | 9
// <digits>          ::= <digit> | <digit><digits>
// <number>          ::= <digits> | <digits>.<digits> | <digits>. | .<digits>
// <sign>            ::= "+" | "-"
// <signedNumber>    ::= <number> | <sign><number>
// <suffix>          ::= <binarySI> | <decimalExponent> | <decimalSI>
// <binarySI>        ::= Ki | Mi | Gi | Ti | Pi | Ei
//   (International System of units; See: http://physics.nist.gov/cuu/Units/binary.html)
// <decimalSI>       ::= m | "" | k | M | G | T | P | E
//   (Note that 1024 = 1Ki but 1000 = 1k; I didn't choose the capitalization.)
// <decimalExponent> ::= "e" <signedNumber> | "E" <signedNumber>
//
// No matter which of the three exponent forms is used, no quantity may represent
// a number greater than 2^63-1 in magnitude, nor may it have more than 3 decimal
// places. Numbers larger or more precise will be capped or rounded up.
// (E.g.: 0.1m will rounded up to 1m.)
// This may be extended in the future if we require larger or smaller quantities.
//
// When a Quantity is parsed from a string, it will remember the type of suffix
// it had, and will use the same type again when it is serialized.
//
// Before serializing, Quantity will be put in "canonical form".
// This means that Exponent/suffix will be adjusted up or down (with a
// corresponding increase or decrease in Mantissa) such that:
//   a. No precision is lost
//   b. No fractional digits will be emitted
//   c. The exponent (or suffix) is as large as possible.
// The sign will be omitted unless the number is negative.
//
// Examples:
//   1.5 will be serialized as "1500m"
//   1.5Gi will be serialized as "1536Mi"
//
// Note that the quantity will NEVER be internally represented by a
// floating point number. That is the whole point of this exercise.
//
// Non-canonical values will still parse as long as they are well formed,
// but will be re-emitted in their canonical form. (So always use canonical
// form, or don't diff.)
//
// This format is intended to make it difficult to use these numbers without
// writing some sort of special handling code in the hopes that that will
// cause implementors to also use a fixed point implementation.
//
// +protobuf=true
// +protobuf.embed=string
// +protobuf.options.marshal=false
// +protobuf.options.(gogoproto.goproto_stringer)=false
// +k8s:deepcopy-gen=true
// +k8s:openapi-gen=true
type Quantity struct {
	// i is the quantity in int64 scaled form, if d.Dec == nil
	i int64Amount
	// d is the quantity in inf.Dec form if d.Dec != nil
	d infDecAmount
	// s is the generated value of this quantity to avoid recalculation
	s string

	// Change Format at will. See the comment for Canonicalize for
	// more details.
	Format
}

// int64Amount represents a fixed precision numerator and arbitrary scale exponent. It is faster
// than operations on inf.Dec for values that can be represented as int64.
// +k8s:openapi-gen=true
type int64Amount struct {
	value int64
	scale Scale
}

// infDecAmount implements common operations over an inf.Dec that are specific to the quantity
// representation.
type infDecAmount struct {
	*Dec
}

// A Dec represents a signed arbitrary-precision decimal.
// It is a combination of a sign, an arbitrary-precision integer coefficient
// value, and a signed fixed-precision exponent value.
// The sign and the coefficient value are handled together as a signed value
// and referred to as the unscaled value.
// (Positive and negative zero values are not distinguished.)
// Since the exponent is most commonly non-positive, it is handled in negated
// form and referred to as scale.
//
// The mathematical value of a Dec equals:
//
//  unscaled * 10**(-scale)
//
// Note that different Dec representations may have equal mathematical values.
//
//  unscaled  scale  String()
//  -------------------------
//         0      0    "0"
//         0      2    "0.00"
//         0     -2    "0"
//         1      0    "1"
//       100      2    "1.00"
//        10      0   "10"
//         1     -1   "10"
//
// The zero value for a Dec represents the value 0 with scale 0.
//
// Operations are typically performed through the *Dec type.
// The semantics of the assignment operation "=" for "bare" Dec values is
// undefined and should not be relied on.
//
// Methods are typically of the form:
//
//	func (z *Dec) Op(x, y *Dec) *Dec
//
// and implement operations z = x Op y with the result as receiver; if it
// is one of the operands it may be overwritten (and its memory reused).
// To enable chaining of operations, the result is also returned. Methods
// returning a result other than *Dec take one of the operands as the receiver.
//
// A "bare" Quo method (quotient / division operation) is not provided, as the
// result is not always a finite decimal and thus in general cannot be
// represented as a Dec.
// Instead, in the common case when rounding is (potentially) necessary,
// QuoRound should be used with a Scale and a Rounder.
// QuoExact or QuoRound with RoundExact can be used in the special cases when it
// is known that the result is always a finite decimal.
//
type Dec struct {
	unscaled big.Int
	scale    Scale
}

// Represents a Persistent Disk resource in Google Compute Engine.
//
// A GCE PD must exist before mounting to a container. The disk must
// also be in the same GCE project and zone as the kubelet. A GCE PD
// can only be mounted as read/write once or read-only many times. GCE
// PDs support ownership management and SELinux relabeling.
type GCEPersistentDiskVolumeSource struct {
	// pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	PDName string `json:"pdName" bson:"pdName"`
	// fsType is filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// partition is the partition in the volume that you want to mount.
	// If omitted, the default is to mount by volume name.
	// Examples: For volume /dev/sda1, you specify the partition as "1".
	// Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	Partition int32 `json:"partition,omitempty" bson:"partition"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Represents a Persistent Disk resource in AWS.
//
// An AWS EBS disk must exist before mounting to a container. The disk
// must also be in the same AWS zone as the kubelet. An AWS EBS disk
// can only be mounted as read/write once. AWS EBS volumes support
// ownership management and SELinux relabeling.
type AWSElasticBlockStoreVolumeSource struct {
	// volumeID is unique ID of the persistent disk resource in AWS (Amazon EBS volume).
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	VolumeID string `json:"volumeID" bson:"volumeID"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// partition is the partition in the volume that you want to mount.
	// If omitted, the default is to mount by volume name.
	// Examples: For volume /dev/sda1, you specify the partition as "1".
	// Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
	// +optional
	Partition int32 `json:"partition,omitempty" bson:"partition"`
	// readOnly value true will force the readOnly setting in VolumeMounts.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Represents a volume that is populated with the contents of a git repository.
// Git repo volumes do not support ownership management.
// Git repo volumes support SELinux relabeling.
//
// DEPRECATED: GitRepo is deprecated. To provision a container with a git repo, mount an
// EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir
// into the Pod's container.
type GitRepoVolumeSource struct {
	// repository is the URL
	Repository string `json:"repository" bson:"repository"`
	// revision is the commit hash for the specified revision.
	// +optional
	Revision string `json:"revision,omitempty" bson:"revision"`
	// directory is the target directory name.
	// Must not contain or start with '..'.  If '.' is supplied, the volume directory will be the
	// git repository.  Otherwise, if specified, the volume will contain the git repository in
	// the subdirectory with the given name.
	// +optional
	Directory string `json:"directory,omitempty" bson:"directory"`
}

// Adapts a Secret into a volume.
//
// The contents of the target Secret's Data field will be presented in a volume
// as files using the keys in the Data field as the file names.
// Secret volumes support ownership management and SELinux relabeling.
type SecretVolumeSource struct {
	// secretName is the name of the secret in the pod's namespace to use.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	// +optional
	SecretName string `json:"secretName,omitempty" bson:"secretName"`
	// items If unspecified, each key-value pair in the Data field of the referenced
	// Secret will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the Secret,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KeyToPath `json:"items,omitempty" bson:"items"`
	// defaultMode is Optional: mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values
	// for mode bits. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" bson:"defaultMode"`
	// optional field specify whether the Secret or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" bson:"optional"`
}

// Maps a string key to a path within a volume.
type KeyToPath struct {
	// key is the key to project.
	Key string `json:"key" bson:"key"`

	// path is the relative path of the file to map the key to.
	// May not be an absolute path.
	// May not contain the path element '..'.
	// May not start with the string '..'.
	Path string `json:"path" bson:"path"`
	// mode is Optional: mode bits used to set permissions on this file.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	Mode *int32 `json:"mode,omitempty" bson:"mode"`
}

// Represents an NFS mount that lasts the lifetime of a pod.
// NFS volumes do not support ownership management or SELinux relabeling.
type NFSVolumeSource struct {
	// server is the hostname or IP address of the NFS server.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	Server string `json:"server" bson:"server"`

	// path that is exported by the NFS server.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	Path string `json:"path" bson:"path"`

	// readOnly here will force the NFS export to be mounted with read-only permissions.
	// Defaults to false.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Represents an ISCSI disk.
// ISCSI volumes can only be mounted as read/write once.
// ISCSI volumes support ownership management and SELinux relabeling.
type ISCSIVolumeSource struct {
	// targetPortal is iSCSI Target Portal. The Portal is either an IP or ip_addr:port if the port
	// is other than default (typically TCP ports 860 and 3260).
	TargetPortal string `json:"targetPortal" bson:"targetPortal"`
	// iqn is the target iSCSI Qualified Name.
	IQN string `json:"iqn" bson:"iqn"`
	// lun represents iSCSI Target Lun number.
	Lun int32 `json:"lun" bson:"lun"`
	// iscsiInterface is the interface Name that uses an iSCSI transport.
	// Defaults to 'default' (tcp).
	// +optional
	ISCSIInterface string `json:"iscsiInterface,omitempty" bson:"iscsiInterface"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
	// portals is the iSCSI Target Portal List. The portal is either an IP or ip_addr:port if the port
	// is other than default (typically TCP ports 860 and 3260).
	// +optional
	Portals []string `json:"portals,omitempty" bson:"portals"`
	// chapAuthDiscovery defines whether support iSCSI Discovery CHAP authentication
	// +optional
	DiscoveryCHAPAuth bool `json:"chapAuthDiscovery,omitempty" bson:"chapAuthDiscovery"`
	// chapAuthSession defines whether support iSCSI Session CHAP authentication
	// +optional
	SessionCHAPAuth bool `json:"chapAuthSession,omitempty" bson:"chapAuthSession"`
	// secretRef is the CHAP Secret for iSCSI target and initiator authentication
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" bson:"secretRef"`
	// initiatorName is the custom iSCSI Initiator Name.
	// If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface
	// <target portal>:<volume name> will be created for the connection.
	// +optional
	InitiatorName *string `json:"initiatorName,omitempty" bson:"initiatorName"`
}

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
// +structType=atomic
type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// TODO: Add other useful fields. apiVersion, kind, uid?
	// +optional
	Name string `json:"name,omitempty" bson:"name"`
}

// Represents a Glusterfs mount that lasts the lifetime of a pod.
// Glusterfs volumes do not support ownership management or SELinux relabeling.
type GlusterfsVolumeSource struct {
	// endpoints is the endpoint name that details Glusterfs topology.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	EndpointsName string `json:"endpoints" bson:"endpoints"`

	// path is the Glusterfs volume path.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	Path string `json:"path" bson:"path"`

	// readOnly here will force the Glusterfs volume to be mounted with read-only permissions.
	// Defaults to false.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// PersistentVolumeClaimVolumeSource references the user's PVC in the same namespace.
// This volume finds the bound PV and mounts that volume for the pod. A
// PersistentVolumeClaimVolumeSource is, essentially, a wrapper around another
// type of volume that is owned by someone else (the system).
type PersistentVolumeClaimVolumeSource struct {
	// claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	ClaimName string `json:"claimName" bson:"claimName"`
	// readOnly Will force the ReadOnly setting in VolumeMounts.
	// Default false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Represents a Rados Block Device mount that lasts the lifetime of a pod.
// RBD volumes support ownership management and SELinux relabeling.
type RBDVolumeSource struct {
	// monitors is a collection of Ceph monitors.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	CephMonitors []string `json:"monitors" bson:"monitors"`
	// image is the rados image name.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	RBDImage string `json:"image" bson:"image"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// pool is the rados pool name.
	// Default is rbd.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	RBDPool string `json:"pool,omitempty" bson:"pool"`
	// user is the rados user name.
	// Default is admin.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	RadosUser string `json:"user,omitempty" bson:"user"`
	// keyring is the path to key ring for RBDUser.
	// Default is /etc/ceph/keyring.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	Keyring string `json:"keyring,omitempty" bson:"keyring"`
	// secretRef is name of the authentication secret for RBDUser. If provided
	// overrides keyring.
	// Default is nil.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" bson:"secretRef"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// FlexVolume represents a generic volume resource that is
// provisioned/attached using an exec based plugin.
type FlexVolumeSource struct {
	// driver is the name of the driver to use for this volume.
	Driver string `json:"driver" bson:"driver"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// secretRef is Optional: secretRef is reference to the secret object containing
	// sensitive information to pass to the plugin scripts. This may be
	// empty if no secret object is specified. If the secret object
	// contains more than one secret, all secrets are passed to the plugin
	// scripts.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" bson:"secretRef"`
	// readOnly is Optional: defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
	// options is Optional: this field holds extra command options if any.
	// +optional
	Options map[string]string `json:"options,omitempty" bson:"options"`
}

// Represents a cinder volume resource in Openstack.
// A Cinder volume must exist before mounting to a container.
// The volume must also be in the same region as the kubelet.
// Cinder volumes support ownership management and SELinux relabeling.
type CinderVolumeSource struct {
	// volumeID used to identify the volume in cinder.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	VolumeID string `json:"volumeID" bson:"volumeID"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	FSType string `json:"fsType,omitempty" bson:"fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
	// secretRef is optional: points to a secret object containing parameters used to connect
	// to OpenStack.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" bson:"secretRef"`
}

// Represents a Ceph Filesystem mount that lasts the lifetime of a pod
// Cephfs volumes do not support ownership management or SELinux relabeling.
type CephFSVolumeSource struct {
	// monitors is Required: Monitors is a collection of Ceph monitors
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	Monitors []string `json:"monitors" bson:"monitors"`
	// path is Optional: Used as the mounted root, rather than the full Ceph tree, default is /
	// +optional
	Path string `json:"path,omitempty" bson:"path"`
	// user is optional: User is the rados user name, default is admin
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	User string `json:"user,omitempty" bson:"user"`
	// secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretFile string `json:"secretFile,omitempty" bson:"secretFile"`
	// secretRef is Optional: SecretRef is reference to the authentication secret for User, default is empty.
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" bson:"secretRef"`
	// readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" bson:"readOnly"`
}

// Toleration The pod this Toleration is attached to tolerates any taint that matches
// the triple <key,value,effect> using the matching operator <operator>.
type Toleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	// +optional
	Key string `json:"key,omitempty" bson:"key"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exists is equivalent to wildcard for value, so that a pod can
	// tolerate all taints of a particular category.
	// +optional
	// NOCC:tosa/linelength(忽略长度)
	Operator TolerationOperator `json:"operator,omitempty" bson:"operator"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value should be empty, otherwise just a regular string.
	// +optional
	Value string `json:"value,omitempty" bson:"value"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.
	// +optional
	Effect TaintEffect `json:"effect,omitempty" bson:"effect"`
	// TolerationSeconds represents the period of time the toleration (which must be
	// of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default,
	// it is not set, which means tolerate the taint forever (do not evict). Zero and
	// negative values will be treated as 0 (evict immediately) by the system.
	// +optional
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty" bson:"tolerationSeconds"`
}

// ContainerPort represents a network port in a single container.
type ContainerPort struct {
	// If specified, this must be an IANA_SVC_NAME and unique within the pod. Each
	// named port in a pod must have a unique name. Name for the port that can be
	// referred to by services.
	// +optional
	Name string `json:"name,omitempty" bson:"name"`
	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// Most containers do not need this.
	// +optional
	HostPort int32 `json:"hostPort,omitempty" bson:"hostPort"`
	// Number of port to expose on the pod's IP address.
	// This must be a valid port number, 0 < x < 65536.
	ContainerPort int32 `json:"containerPort" bson:"containerPort"`
	// Protocol for port. Must be UDP, TCP, or SCTP.
	// Defaults to "TCP".
	// +optional
	// +default="TCP"
	Protocol Protocol `json:"protocol,omitempty" bson:"protocol"`
	// What host IP to bind the external port to.
	// +optional
	HostIP string `json:"hostIP,omitempty" bson:"hostIP"`
}
