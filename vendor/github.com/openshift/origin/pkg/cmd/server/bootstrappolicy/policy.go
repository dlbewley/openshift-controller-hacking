package bootstrappolicy

import (
	"fmt"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/util/sets"

	"github.com/openshift/origin/pkg/api"
	authorizationapi "github.com/openshift/origin/pkg/authorization/api"
	imageapi "github.com/openshift/origin/pkg/image/api"
)

func GetBootstrapOpenshiftRoles(openshiftNamespace string) []authorizationapi.Role {
	roles := []authorizationapi.Role{
		{
			ObjectMeta: kapi.ObjectMeta{
				Name:      OpenshiftSharedResourceViewRoleName,
				Namespace: openshiftNamespace,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "list"),
					Resources: sets.NewString("templates", authorizationapi.ImageGroupName),
				},
				{
					// so anyone can pull from openshift/* image streams
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("imagestreams/layers"),
				},
			},
		},
	}

	// we don't want to expose the resourcegroups externally because it makes it very difficult for customers to learn from
	// our default roles and hard for them to reason about what power they are granting their users
	for i := range roles {
		for j := range roles[i].Rules {
			roles[i].Rules[j].Resources = authorizationapi.NormalizeResources(roles[i].Rules[j].Resources)
		}
	}

	return roles

}
func GetBootstrapClusterRoles() []authorizationapi.ClusterRole {
	roles := []authorizationapi.ClusterRole{
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ClusterAdminRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{authorizationapi.APIGroupAll},
					Verbs:     sets.NewString(authorizationapi.VerbAll),
					Resources: sets.NewString(authorizationapi.ResourceAll),
				},
				{
					Verbs:           sets.NewString(authorizationapi.VerbAll),
					NonResourceURLs: sets.NewString(authorizationapi.NonResourceAll),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ClusterReaderRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString(authorizationapi.NonEscalatingResourcesGroupName),
				},
				{
					APIGroups: []string{autoscaling.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("horizontalpodautoscalers"),
				},
				{
					APIGroups: []string{batch.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("jobs"),
				},
				{
					APIGroups: []string{extensions.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("daemonsets", "jobs", "horizontalpodautoscalers", "replicationcontrollers/scale"),
				},
				{ // permissions to check access.  These creates are non-mutating
					Verbs:     sets.NewString("create"),
					Resources: sets.NewString("resourceaccessreviews", "subjectaccessreviews"),
				},
				// Allow read access to node metrics
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString(authorizationapi.NodeMetricsResource),
				},
				// Allow read access to stats
				// Node stats requests are submitted as POSTs.  These creates are non-mutating
				{
					Verbs:     sets.NewString("get", "create"),
					Resources: sets.NewString(authorizationapi.NodeStatsResource),
				},
				{
					Verbs:           sets.NewString("get"),
					NonResourceURLs: sets.NewString(authorizationapi.NonResourceAll),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: BuildStrategyDockerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{api.GroupName},
					Verbs:     sets.NewString("create"),
					Resources: sets.NewString(authorizationapi.DockerBuildResource),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: BuildStrategyCustomRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{api.GroupName},
					Verbs:     sets.NewString("create"),
					Resources: sets.NewString(authorizationapi.CustomBuildResource),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: BuildStrategySourceRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{api.GroupName},
					Verbs:     sets.NewString("create"),
					Resources: sets.NewString(authorizationapi.SourceBuildResource),
				},
			},
		},

		{
			ObjectMeta: kapi.ObjectMeta{
				Name: AdminRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{kapi.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString(
						authorizationapi.KubeExposedGroupName,
						"secrets",
						"pods/attach", "pods/proxy", "pods/exec", "pods/portforward",
						"services/proxy",
						"replicationcontrollers/scale",
					),
				},
				{
					APIGroups: []string{api.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString(
						authorizationapi.OpenshiftExposedGroupName,
						authorizationapi.PermissionGrantingGroupName,
						"projects",
						"deploymentconfigs/scale",
						"imagestreams/secrets",
					),
				},
				{
					APIGroups: []string{autoscaling.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString("horizontalpodautoscalers"),
				},
				{
					APIGroups: []string{batch.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString("jobs"),
				},
				{
					APIGroups: []string{extensions.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString("jobs", "horizontalpodautoscalers", "replicationcontrollers/scale"),
				},
				{
					APIGroups: []string{extensions.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("daemonsets"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString(authorizationapi.PolicyOwnerGroupName, authorizationapi.KubeAllGroupName, authorizationapi.OpenshiftStatusGroupName, authorizationapi.KubeStatusGroupName),
				},
				{
					Verbs: sets.NewString("get", "update"),
					// this is used by verifyImageStreamAccess in pkg/dockerregistry/server/auth.go
					Resources: sets.NewString("imagestreams/layers"),
				},
				// an admin can run routers that write back conditions to the route
				{
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("routes/status"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: EditRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{kapi.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString(
						authorizationapi.KubeExposedGroupName,
						"secrets",
						"pods/attach", "pods/proxy", "pods/exec", "pods/portforward",
						"services/proxy",
						"replicationcontrollers/scale",
					),
				},
				{
					APIGroups: []string{api.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString(
						authorizationapi.OpenshiftExposedGroupName,
						"deploymentconfigs/scale",
						"imagestreams/secrets",
					),
				},
				{
					APIGroups: []string{autoscaling.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString("horizontalpodautoscalers"),
				},
				{
					APIGroups: []string{batch.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString("jobs"),
				},
				{
					APIGroups: []string{extensions.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "create", "update", "patch", "delete", "deletecollection"),
					Resources: sets.NewString("jobs", "horizontalpodautoscalers", "replicationcontrollers/scale"),
				},
				{
					APIGroups: []string{extensions.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("daemonsets"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString(authorizationapi.KubeAllGroupName, authorizationapi.OpenshiftStatusGroupName, authorizationapi.KubeStatusGroupName, "projects"),
				},
				{
					Verbs: sets.NewString("get", "update"),
					// this is used by verifyImageStreamAccess in pkg/dockerregistry/server/auth.go
					Resources: sets.NewString("imagestreams/layers"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ViewRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString(authorizationapi.OpenshiftExposedGroupName, authorizationapi.KubeAllGroupName, authorizationapi.OpenshiftStatusGroupName, authorizationapi.KubeStatusGroupName, "projects"),
				},
				{
					APIGroups: []string{autoscaling.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("horizontalpodautoscalers"),
				},
				{
					APIGroups: []string{batch.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("jobs"),
				},
				{
					APIGroups: []string{extensions.GroupName},
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("daemonsets", "jobs", "horizontalpodautoscalers"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: BasicUserRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{Verbs: sets.NewString("get"), Resources: sets.NewString("users"), ResourceNames: sets.NewString("~")},
				{Verbs: sets.NewString("list"), Resources: sets.NewString("projectrequests")},
				{Verbs: sets.NewString("list", "get"), Resources: sets.NewString("clusterroles")},
				{Verbs: sets.NewString("list"), Resources: sets.NewString("projects")},
				{Verbs: sets.NewString("create"), Resources: sets.NewString("subjectaccessreviews", "localsubjectaccessreviews"), AttributeRestrictions: &authorizationapi.IsPersonalSubjectAccessReview{}},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: SelfProvisionerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{Verbs: sets.NewString("create"), Resources: sets.NewString("projectrequests")},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: StatusCheckerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs: sets.NewString("get"),
					NonResourceURLs: sets.NewString(
						// Health
						"/healthz", "/healthz/*",

						// Server version checking
						"/version",

						// API discovery/negotiation
						"/api", "/api/*",
						"/apis", "/apis/*",
						"/oapi", "/oapi/*",
						"/osapi", "/osapi/", // these cannot be removed until we can drop support for pre 3.1 clients
					),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ImageAuditorRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{imageapi.GroupName},
					Verbs:     sets.NewString("get", "list", "watch", "patch", "update"),
					Resources: sets.NewString("images"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ImagePullerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs: sets.NewString("get"),
					// this is used by verifyImageStreamAccess in pkg/dockerregistry/server/auth.go
					Resources: sets.NewString("imagestreams/layers"),
				},
			},
		},
		{
			// This role looks like a duplicate of ImageBuilderRole, but the ImageBuilder role is specifically for our builder service accounts
			// if we found another permission needed by them, we'd add it there so the intent is different if you used the ImageBuilderRole
			// you could end up accidentally granting more permissions than you intended.  This is intended to only grant enough powers to
			// push an image to our registry
			ObjectMeta: kapi.ObjectMeta{
				Name: ImagePusherRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs: sets.NewString("get", "update"),
					// this is used by verifyImageStreamAccess in pkg/dockerregistry/server/auth.go
					Resources: sets.NewString("imagestreams/layers"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ImageBuilderRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs: sets.NewString("get", "update"),
					// this is used by verifyImageStreamAccess in pkg/dockerregistry/server/auth.go
					Resources: sets.NewString("imagestreams/layers"),
				},
				{
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("builds/details"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ImagePrunerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("delete"),
					Resources: sets.NewString("images"),
				},
				{
					Verbs:     sets.NewString("get", "list"),
					Resources: sets.NewString("images", "imagestreams", "pods", "replicationcontrollers", "buildconfigs", "builds", "deploymentconfigs"),
				},
				{
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("imagestreams/status"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: DeployerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					// replicationControllerGetter
					Verbs:     sets.NewString("get", "list"),
					Resources: sets.NewString("replicationcontrollers"),
				},
				{
					// RecreateDeploymentStrategy.replicationControllerClient
					// RollingDeploymentStrategy.updaterClient
					Verbs:     sets.NewString("get", "update"),
					Resources: sets.NewString("replicationcontrollers"),
				},
				{
					// RecreateDeploymentStrategy.hookExecutor
					// RollingDeploymentStrategy.hookExecutor
					Verbs:     sets.NewString("get", "list", "watch", "create"),
					Resources: sets.NewString("pods"),
				},
				{
					// RecreateDeploymentStrategy.hookExecutor
					// RollingDeploymentStrategy.hookExecutor
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("pods/log"),
				},
				{
					// Deployer.After.TagImages
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("imagestreamtags"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: MasterRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					APIGroups: []string{authorizationapi.APIGroupAll},
					Verbs:     sets.NewString(authorizationapi.VerbAll),
					Resources: sets.NewString(authorizationapi.ResourceAll),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: OAuthTokenDeleterRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("delete"),
					Resources: sets.NewString("oauthaccesstokens", "oauthauthorizetokens"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RouterRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("list", "watch"),
					Resources: sets.NewString("routes", "endpoints"),
				},
				// routers write back conditions to the route
				{
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("routes/status"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RegistryRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "delete"),
					Resources: sets.NewString("images"),
				},
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("imagestreamimages", "imagestreamtags", "imagestreams", "imagestreams/secrets"),
				},
				{
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("imagestreams"),
				},
				{
					Verbs:     sets.NewString("create"),
					Resources: sets.NewString("imagestreammappings"),
				},
				{
					Verbs:     sets.NewString("list"),
					Resources: sets.NewString("resourcequotas"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeProxierRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					// Used to build serviceLister
					Verbs:     sets.NewString("list", "watch"),
					Resources: sets.NewString("services", "endpoints"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeAdminRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				// Allow read-only access to the API objects
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("nodes"),
				},
				// Allow all API calls to the nodes
				{
					Verbs:     sets.NewString("proxy"),
					Resources: sets.NewString("nodes"),
				},
				{
					Verbs:     sets.NewString(authorizationapi.VerbAll),
					Resources: sets.NewString("nodes/proxy", authorizationapi.NodeMetricsResource, authorizationapi.NodeStatsResource, authorizationapi.NodeLogResource),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeReaderRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				// Allow read-only access to the API objects
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("nodes"),
				},
				// Allow read access to node metrics
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString(authorizationapi.NodeMetricsResource),
				},
				// Allow read access to stats
				// Node stats requests are submitted as POSTs.  These creates are non-mutating
				{
					Verbs:     sets.NewString("get", "create"),
					Resources: sets.NewString(authorizationapi.NodeStatsResource),
				},
				// TODO: expose other things like /healthz on the node once we figure out non-resource URL policy across systems
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					// Needed to check API access.  These creates are non-mutating
					Verbs:     sets.NewString("create"),
					Resources: sets.NewString("subjectaccessreviews", "localsubjectaccessreviews"),
				},
				{
					// Needed to build serviceLister, to populate env vars for services
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("services"),
				},
				{
					// Nodes can register themselves
					// TODO: restrict to creating a node with the same name they announce
					Verbs:     sets.NewString("create", "get", "list", "watch"),
					Resources: sets.NewString("nodes"),
				},
				{
					// TODO: restrict to the bound node once supported
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("nodes/status"),
				},

				{
					// TODO: restrict to the bound node as creator once supported
					Verbs:     sets.NewString("create", "update", "patch"),
					Resources: sets.NewString("events"),
				},

				{
					// TODO: restrict to pods scheduled on the bound node once supported
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("pods"),
				},
				{
					// TODO: remove once mirror pods are removed
					// TODO: restrict deletion to mirror pods created by the bound node once supported
					// Needed for the node to create/delete mirror pods
					Verbs:     sets.NewString("get", "create", "delete"),
					Resources: sets.NewString("pods"),
				},
				{
					// TODO: restrict to pods scheduled on the bound node once supported
					Verbs:     sets.NewString("update"),
					Resources: sets.NewString("pods/status"),
				},

				{
					// TODO: restrict to secrets and configmaps used by pods scheduled on bound node once supported
					// Needed for imagepullsecrets, rbd/ceph and secret volumes, and secrets in envs
					// Needed for configmap volume and envs
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("secrets", "configmaps"),
				},
				{
					// TODO: restrict to claims/volumes used by pods scheduled on bound node once supported
					// Needed for persistent volumes
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("persistentvolumeclaims", "persistentvolumes"),
				},
				{
					// TODO: restrict to namespaces of pods scheduled on bound node once supported
					// TODO: change glusterfs to use DNS lookup so this isn't needed?
					// Needed for glusterfs volumes
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("endpoints"),
				},
			},
		},

		{
			ObjectMeta: kapi.ObjectMeta{
				Name: SDNReaderRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("hostsubnets"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("netnamespaces"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("nodes"),
				},
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("clusternetworks"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("namespaces"),
				},
			},
		},

		{
			ObjectMeta: kapi.ObjectMeta{
				Name: SDNManagerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "list", "watch", "create", "delete"),
					Resources: sets.NewString("hostsubnets"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch", "create", "delete"),
					Resources: sets.NewString("netnamespaces"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("nodes"),
				},
				{
					Verbs:     sets.NewString("get", "create"),
					Resources: sets.NewString("clusternetworks"),
				},
			},
		},

		{
			ObjectMeta: kapi.ObjectMeta{
				Name: WebHooksRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "create"),
					Resources: sets.NewString("buildconfigs/webhooks"),
				},
			},
		},

		{
			ObjectMeta: kapi.ObjectMeta{
				Name: DiscoveryRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs: sets.NewString("get"),
					NonResourceURLs: sets.NewString(
						// Server version checking
						"/version",

						// API discovery/negotiation
						"/api", "/api/*",
						"/apis", "/apis/*",
						"/oapi", "/oapi/*",
						"/osapi", "/osapi/", // these cannot be removed until we can drop support for pre 3.1 clients
					),
				},
			},
		},

		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RegistryAdminRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"),
					Resources: sets.NewString("imagestreamimages", "imagestreamimports", "imagestreammappings", "imagestreams", "imagestreams/secrets", "imagestreamtags"),
				},
				{
					Verbs:     sets.NewString("create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"),
					Resources: sets.NewString("localresourceaccessreviews", "localsubjectaccessreviews", "resourceaccessreviews", "rolebindings", "roles", "subjectaccessreviews"),
				},
				{
					Verbs:     sets.NewString("get", "update"),
					Resources: sets.NewString("imagestreams/layers"),
				},
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("policies", "policybindings"),
				},
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("namespaces", "projects"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RegistryViewerRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get", "list", "watch"),
					Resources: sets.NewString("imagestreamimages", "imagestreamimports", "imagestreammappings", "imagestreams", "imagestreamtags"),
				},
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("imagestreams/layers", "namespaces", "projects"),
				},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RegistryEditorRoleName,
			},
			Rules: []authorizationapi.PolicyRule{
				{
					Verbs:     sets.NewString("get"),
					Resources: sets.NewString("namespaces", "projects"),
				},
				{
					Verbs:     sets.NewString("create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"),
					Resources: sets.NewString("imagestreamimages", "imagestreamimports", "imagestreammappings", "imagestreams", "imagestreams/secrets", "imagestreamtags"),
				},
				{
					Verbs:     sets.NewString("get", "update"),
					Resources: sets.NewString("imagestreams/layers"),
				},
			},
		},
	}

	saRoles := InfraSAs.AllRoles()
	for _, saRole := range saRoles {
		for _, existingRole := range roles {
			if existingRole.Name == saRole.Name {
				panic(fmt.Sprintf("clusterrole/%s is already registered", existingRole.Name))
			}
		}
	}

	roles = append(roles, saRoles...)

	// we don't want to expose the resourcegroups externally because it makes it very difficult for customers to learn from
	// our default roles and hard for them to reason about what power they are granting their users
	for i := range roles {
		for j := range roles[i].Rules {
			roles[i].Rules[j].Resources = authorizationapi.NormalizeResources(roles[i].Rules[j].Resources)
		}
	}

	return roles
}

func GetBootstrapOpenshiftRoleBindings(openshiftNamespace string) []authorizationapi.RoleBinding {
	return []authorizationapi.RoleBinding{
		{
			ObjectMeta: kapi.ObjectMeta{
				Name:      OpenshiftSharedResourceViewRoleBindingName,
				Namespace: openshiftNamespace,
			},
			RoleRef: kapi.ObjectReference{
				Name:      OpenshiftSharedResourceViewRoleName,
				Namespace: openshiftNamespace,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}},
		},
	}
}

func GetBootstrapClusterRoleBindings() []authorizationapi.ClusterRoleBinding {
	return []authorizationapi.ClusterRoleBinding{
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: MasterRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: MasterRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: MastersGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeAdminRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: NodeAdminRoleName,
			},
			Subjects: []kapi.ObjectReference{
				// ensure the legacy username in the master's kubelet-client certificate is allowed
				{Kind: authorizationapi.SystemUserKind, Name: LegacyMasterKubeletAdminClientUsername},
				{Kind: authorizationapi.SystemGroupKind, Name: NodeAdminsGroup},
			},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ClusterAdminRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: ClusterAdminRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: ClusterAdminGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: ClusterReaderRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: ClusterReaderRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: ClusterReaderGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: BasicUserRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: BasicUserRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: SelfProvisionerRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: SelfProvisionerRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedOAuthGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: OAuthTokenDeleterRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: OAuthTokenDeleterRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}, {Kind: authorizationapi.SystemGroupKind, Name: UnauthenticatedGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: StatusCheckerRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: StatusCheckerRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}, {Kind: authorizationapi.SystemGroupKind, Name: UnauthenticatedGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RouterRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: RouterRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: RouterGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: RegistryRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: RegistryRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: RegistryGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: NodeRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: NodesGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: NodeProxierRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: NodeProxierRoleName,
			},
			// Allow node identities to run node proxies
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: NodesGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: SDNReaderRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: SDNReaderRoleName,
			},
			// Allow node identities to run SDN plugins
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: NodesGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: WebHooksRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: WebHooksRoleName,
			},
			Subjects: []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}, {Kind: authorizationapi.SystemGroupKind, Name: UnauthenticatedGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{
				Name: DiscoveryRoleBindingName,
			},
			RoleRef: kapi.ObjectReference{
				Name: DiscoveryRoleName,
			},
			Subjects: []kapi.ObjectReference{
				{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup},
				{Kind: authorizationapi.SystemGroupKind, Name: UnauthenticatedGroup},
			},
		},

		// Allow all build strategies by default.
		// Cluster admins can remove these role bindings, and the reconcile-cluster-role-bindings command
		// run during an upgrade won't re-add the "system:authenticated" group
		{
			ObjectMeta: kapi.ObjectMeta{Name: BuildStrategyDockerRoleBindingName},
			RoleRef:    kapi.ObjectReference{Name: BuildStrategyDockerRoleName},
			Subjects:   []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{Name: BuildStrategyCustomRoleBindingName},
			RoleRef:    kapi.ObjectReference{Name: BuildStrategyCustomRoleName},
			Subjects:   []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}},
		},
		{
			ObjectMeta: kapi.ObjectMeta{Name: BuildStrategySourceRoleBindingName},
			RoleRef:    kapi.ObjectReference{Name: BuildStrategySourceRoleName},
			Subjects:   []kapi.ObjectReference{{Kind: authorizationapi.SystemGroupKind, Name: AuthenticatedGroup}},
		},
	}
}
