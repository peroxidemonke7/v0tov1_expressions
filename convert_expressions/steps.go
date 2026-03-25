package convertexpressions

var StepsConversionRules = []ConversionRule{
	{"identifier", "id"},
}

// inside step.spec
var StepSpecConversionRules = map[string][]ConversionRule{
	StepTypeRun: {
		{"command", "spec.script"},
		{"envVariables", "spec.env"},
		{"connectorRef", "spec.container.connector"},
		{"image", "spec.container.image"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"shell", "spec.shell"},
		{"outputVariables", "spec.outputs"},
		{"reports", "spec.report"},
		{"resources.limits.cpu", "spec.container.cpu"},
		{"resources.limits.memory", "spec.container.memory"},
		{"runAsUser", "spec.container.user"},
	},

	StepTypeBackground: {
		{"command", "spec.script"},
		{"image", "spec.container.image"},
		{"connectorRef", "spec.container.connector"},
		{"imagePullPolicy", "spec.container.pull"},
		{"privileged", "spec.container.privileged"},
		{"envVariables", "spec.env"},
		{"shell", "spec.shell"},
		{"entrypoint", "spec.entrypoint"},
		{"portBindings", "spec.ports"},
		{"reports", "spec.report"},
	},

	StepTypeHTTP: {
		{"url", "spec.env.PLUGIN_URL"},
		{"method", "spec.env.PLUGIN_METHOD"},
		{"headers", "spec.env.PLUGIN_HEADERS"},
		{"requestBody", "spec.env.PLUGIN_BODY"},
		{"assertion", "spec.env.PLUGIN_ASSERTION"},
		{"inputVariables", "spec.env.PLUGIN_INPUT_VARIABLES"},
		{"outputVariables", "spec.env.PLUGIN_OUTPUT_VARIABLES"},
		{"certificate", "spec.env.PLUGIN_CLIENT_CERT"},
		{"certificateKey", "spec.env.PLUGIN_CLIENT_KEY"},
	},

	StepTypeGCSUpload: {
		{"sourcePath", "spec.sourcePath"},
		{"bucket", "spec.bucket"},
		{"target", "spec.target"},
	},

	StepTypeS3Upload: {
		{"step.spec.bucket", "step.steps.s3Upload.spec.with.BUCKET"},
		{"step.spec.region", "step.steps.s3Upload.spec.with.REGION"},
		{"step.spec.endpoint", "step.steps.s3Upload.spec.with.ENDPOINT"},
		{"step.spec.sourcePath", "step.steps.s3Upload.spec.with.SOURCE"},
		{"step.spec.target", "step.steps.s3Upload.spec.with.TARGET"},
		// {"step.spec.connectorRef", ""},
	},

	StepTypeRestoreCacheS3: {
		{"bucket", "steps.restoreCacheS3.spec.with.BUCKET"},
		{"key", "steps.restoreCacheS3.spec.with.CACHE_KEY"},
		{"region", "steps.restoreCacheS3.spec.with.REGION"},
		{"endpoint", "steps.restoreCacheS3.spec.with.ENDPOINT"},
		{"archiveFormat", "steps.restoreCacheS3.spec.with.ARCHIVE_FORMAT"},
		{"pathStyle", "steps.restoreCacheS3.spec.with.PATH_STYLE"},
		{"failIfKeyNotFound", "steps.restoreCacheS3.spec.with.FAIL_RESTORE_IF_KEY_NOT_PRESENT"},
		// {"connectorRef", ""},
	},

	StepTypeSaveCacheGCS: {
		{"bucket", "steps.saveCacheGCS.spec.with.BUCKET"},
		{"key", "steps.saveCacheGCS.spec.with.CACHE_KEY"},
		{"sourcePaths", "steps.saveCacheGCS.spec.with.MOUNT"},
		{"archiveFormat", "steps.saveCacheGCS.spec.with.ARCHIVE_FORMAT"},
		{"override", "steps.saveCacheGCS.spec.with.OVERRIDE"},
		// {"spec.connectorRef", ""},

	},

	StepTypeSaveCacheS3: {
		{"bucket", "spec.with.BUCKET"},
		{"key", "spec.with.CACHE_KEY"},
		{"sourcePaths", "spec.with.MOUNT"},
		{"region", "spec.with.REGION"},
		{"endpoint", "spec.with.ENDPOINT"},
		{"archiveFormat", "spec.with.ARCHIVE_FORMAT"},
		{"pathStyle", "spec.with.PATH_STYLE"},
		{"override", "spec.with.OVERRIDE"},
		// {"spec.connectorRef", ""},

	},

	StepTypeShellScript: {
		{"shell", "spec.shell"},
		{"source.spec.script", "spec.script"},
		{"environmentVariables", "spec.env"},
		{"outputVariables", "spec.outputs"},
	},

	StepTypeCustomApproval: {
		{"shellType", "runStepInfo.shell.yamlName"},
		{"source.spec.script", "runStepInfo.script"},
		{"outputVariables.*", "runStepInfo.outputs.*"}, // TODO: verify
	},

	StepTypeJiraApproval: {
		{"connectorRef", "runStepInfo.env.PLUGIN_HARNESS_CONNECTOR"},
		{"issueKey", "runStepInfo.env.PLUGIN_ISSUE_KEY"},
		// {"issueType", "runStepInfo.env.PLUGIN_ISSUE_TYPE"}, // TODO: verify
		// {"projectKey", "runStepInfo.env.PLUGIN_PROJECT_KEY"}, // TODO: verify
	},

	StepTypeServiceNowApproval: {
		{"connectorRef", "runStepInfo.env.PLUGIN_HARNESS_CONNECTOR"},
		{"ticketType", "runStepInfo.env.PLUGIN_TICKET_TYPE"},
		{"ticketNumber", "runStepInfo.env.PLUGIN_TICKET_NUMBER"},
	},

	// StepRestoreCacheFromGCS
	StepTypeRestoreCacheGCS: {
		{"bucket", "steps.restoreCacheGCS.spec.with.BUCKET"},
		{"key", "steps.restoreCacheGCS.spec.with.CACHE_KEY"},
		{"archiveFormat", "steps.restoreCacheGCS.spec.with.ARCHIVE_FORMAT"},
		{"failIfKeyNotFound", "steps.restoreCacheGCS.spec.with.FAIL_RESTORE_IF_KEY_NOT_PRESENT"},
		// {"connectorRef", ""},
	},

	StepTypeArtifactoryUpload: {
		{"sourcePath", "steps.jfrogArtifactory.spec.sourcePath"},
		{"target", "steps.jfrogArtifactory.spec.target"},
		// {"connectorRef", ""},
	},

	StepTypeGitClone: {
		// {"repoName", "spec.env.DRONE_REMOTE_URL"}, // Note: v0 has "repoName", v1 expects full URL
		{"cloneDirectory", "spec.env.DRONE_WORKSPACE"},
		{"depth", "spec.with.DEPTH"},
	},

	// ============================================================
	// K8S STEPS
	// ============================================================

	// K8sRollingDeploy (v1: k8sRollingDeployStep)
	// Step group: k8sRollingPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction
	StepTypeK8sRollingDeploy: {
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		{"pruningEnabled", "steps.k8sRollingPrepareAction.spec.env.PLUGIN_PRUNING_ENABLED"},
		// pruningEnabled also maps to: steps.k8sApplyAction.spec.env.PLUGIN_RELEASE_PRUNING_ENABLED
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sRollingRollback (v1: k8sRollingRollbackStep)
	// Step group: k8sRollingRollbackAction
	StepTypeK8sRollingRollback: {
		{"pruningEnabled", "steps.k8sRollingRollbackAction.spec.env.PLUGIN_PRUNING"},
		{"flags", "steps.k8sRollingRollbackAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sApply (v1: k8sApplyStep)
	// Step group: k8sApplyAction, k8sSteadyStateCheckAction
	StepTypeK8sApply: {
		{"filePaths", "steps.k8sApplyAction.spec.env.PLUGIN_MANIFEST_PATH"},
		// COMPLEX: each path prefixed with "<+runtime.manifestPath>/"
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		// skipSteadyStateCheck -> NO_ENV_MAPPING (controls conditional execution)
		// skipRendering -> NO_TEMPLATE_MAPPING
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sBGSwapServices (v1: k8sBlueGreenSwapServicesSelectorsStep)
	// Step group: k8sBlueGreenSwapServicesSelectorsAction
	// All fields are COMPLEX (populated from exported variables, not from v0 spec)
	StepTypeK8sBGSwapServices: {
		{"stable_service", "steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STABLE_SERVICE"},
		{"stage_service", "steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STAGE_SERVICE"},
		{"is_openshift", "steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.HARNESS_IS_OPENSHIFT"},
	},

	// K8sBlueGreenStageScaleDown (v1: k8sBlueGreenStageScaleDownStep)
	// Step group: k8sBlueGreenStageScaleDownAction, (conditional) k8sDeleteAction
	// deleteResources -> NO_ENV_MAPPING (controls conditional execution)
	StepTypeK8sBlueGreenStageScaleDown: {},

	// K8sCanaryDelete (v1: k8sCanaryDeleteStep)
	// Step group: k8sCanaryDeleteAction (uses k8sDeleteAction)
	// All fields are COMPLEX (populated from exported variables)
	StepTypeK8sCanaryDelete: {
		{"resources", "steps.k8sCanaryDeleteAction.spec.env.PLUGIN_RESOURCES"},
		{"is_openshift", "steps.k8sCanaryDeleteAction.spec.env.HARNESS_IS_OPENSHIFT"},
		// select_delete_resources (Rollback) -> NO_ENV_MAPPING
	},

	// K8sDiff (v1: k8sDiffStep)
	// Step group: k8sDiffAction
	// No field mappings (v0 spec and v1 with are both empty)
	StepTypeK8sDiff: {},

	// K8sRollout (v1: k8sRolloutStep)
	// Step group: k8sRolloutAction (id: k8sRolloutStep in step group)
	StepTypeK8sRollout: {
		{"command", "steps.k8sRolloutStep.spec.env.PLUGIN_COMMAND"},
		// resources.type -> NO_ENV_MAPPING (controls UI visibility)
		{"resources.spec.resourceNames", "steps.k8sRolloutStep.spec.env.PLUGIN_RESOURCES"},
		// only when type=ResourceName
		{"resources.spec.manifestPaths", "steps.k8sRolloutStep.spec.env.PLUGIN_MANIFESTS"},
		// only when type=ManifestPath
		{"flags", "steps.k8sRolloutStep.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sScale (v1: k8sScaleStep)
	// Step group: k8sScaleAction, (conditional) k8sSteadyStateCheckAction
	StepTypeK8sScale: {
		{"instanceSelection.type", "steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE"},
		// COMPLEX: "Count" -> "count", "Percentage" -> "percentage"
		{"instanceSelection.spec.count", "steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Count
		{"instanceSelection.spec.percentage", "steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Percentage
		{"workload", "steps.k8sScaleAction.spec.env.PLUGIN_WORKLOAD"},
		// skipSteadyStateCheck -> NO_ENV_MAPPING (controls conditional execution)
	},

	// K8sDryRun (v1: k8sDryRunStep)
	// Step group: k8sDryRunAction (uses k8sApplyAction with only_dry_run=true)
	StepTypeK8sDryRun: {
		{"encryptYamlOutput", "steps.k8sDryRunAction.spec.env.PLUGIN_ENCRYPT_YAML_OUTPUT"},
	},

	// K8sDelete (v1: k8sDeleteStep)
	// Step group: k8sDeleteAction
	StepTypeK8sDelete: {
		// deleteResources.type -> NO_ENV_MAPPING (controls UI visibility)
		{"deleteResources.spec.resourceNames", "steps.k8sDeleteAction.spec.env.PLUGIN_RESOURCES"},
		// only when type=ResourceName
		{"deleteResources.spec.manifestPaths", "steps.k8sDeleteAction.spec.env.PLUGIN_MANIFESTS"},
		// COMPLEX: each path prefixed with "<+runtime.manifestPath>/"; only when type=ManifestPath
		{"deleteResources.spec.deleteNamespace", "steps.k8sDeleteAction.spec.env.PLUGIN_INCLUDE_NAMESPACES"},
		// only when type=ReleaseName
		{"flags", "steps.k8sDeleteAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sTrafficRouting (v1: k8sTrafficRoutingStep)
	// Step group: k8sTrafficShiftAction
	StepTypeK8sTrafficRouting: {
		{"trafficRouting.provider", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_PROVIDER"},
		{"trafficRouting.spec.name", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_RESOURCE_NAME"},
		{"trafficRouting.spec.hosts", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_HOSTNAMES"},
		{"trafficRouting.spec.gateways", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_GATEWAYS"},
		{"trafficRouting.spec.routes", "steps.k8sTrafficShiftAction.spec.env.PLUGIN_ROUTES"},
		// COMPLEX: converted to JSON string
		// trafficRouting.spec.rootService -> NO_TEMPLATE_MAPPING
		// type -> NO_TEMPLATE_MAPPING (v1 always sets config="new")
	},

	// K8sCanaryDeploy (v1: k8sCanaryDeployStep)
	// Step group: k8sCanaryPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction, (conditional) k8sTrafficRoutingAction
	StepTypeK8sCanaryDeploy: {
		{"instanceSelection.type", "steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE"},
		// COMPLEX: "Count" -> "count", "Percentage" -> "percentage"
		{"instanceSelection.spec.count", "steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Count
		{"instanceSelection.spec.percentage", "steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Percentage
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		{"trafficRouting.provider", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_PROVIDER"},
		{"trafficRouting.spec.name", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_RESOURCE_NAME"},
		{"trafficRouting.spec.hosts", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_HOSTNAMES"},
		{"trafficRouting.spec.gateways", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_GATEWAYS"},
		{"trafficRouting.spec.routes", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_ROUTES"},
		// COMPLEX: converted to JSON string
		// trafficRouting.spec.rootService -> NO_TEMPLATE_MAPPING
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sBlueGreenDeploy (v1: k8sBlueGreenDeployStep)
	// Step group: k8sBlueGreenPrepareAction, k8sApplyAction, k8sSteadyStateCheckAction, (conditional) k8sTrafficRoutingAction
	StepTypeK8sBlueGreenDeploy: {
		{"skipDryRun", "steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN"},
		{"pruningEnabled", "steps.k8sApplyAction.spec.env.PLUGIN_RELEASE_PRUNING_ENABLED"},
		{"skipUnchangedManifest", "steps.k8sBlueGreenPrepareAction.spec.env.PLUGIN_SKIP_UNCHANGED_MANIFEST"},
		{"trafficRouting.provider", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_PROVIDER"},
		{"trafficRouting.spec.name", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_RESOURCE_NAME"},
		{"trafficRouting.spec.hosts", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_HOSTNAMES"},
		{"trafficRouting.spec.gateways", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_GATEWAYS"},
		{"trafficRouting.spec.routes", "steps.k8sTrafficRoutingAction.spec.env.PLUGIN_ROUTES"},
		// COMPLEX: converted to JSON string
		// trafficRouting.spec.rootService -> NO_TEMPLATE_MAPPING
		{"flags", "steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON"},
	},

	// K8sPatch (v1: k8sPatchStep)
	// Step group: k8sPatchAction, (conditional) k8sSteadyStateCheckAction
	StepTypeK8sPatch: {
		{"workload", "steps.k8sPatchAction.spec.env.PLUGIN_WORKLOAD"},
		// skipSteadyStateCheck -> NO_ENV_MAPPING (controls conditional execution)
		{"mergeStrategyType", "steps.k8sPatchAction.spec.env.PLUGIN_MERGE_STRATEGY"},
		// COMPLEX: value is lowercased: "Json"->"json", "Strategic"->"strategic", "Merge"->"merge"
		// source.type -> NO_TEMPLATE_MAPPING (used to determine how to extract content)
		{"source.spec.content", "steps.k8sPatchAction.spec.env.PLUGIN_CONTENT"},
		// only when source.type=Inline
		// source.spec (Git/GitLab/etc) -> NO_TEMPLATE_MAPPING (remote source not supported)
		// recordChangeCause -> NO_TEMPLATE_MAPPING
	},

	// ============================================================
	// HELM STEPS
	// ============================================================

	// HelmBGDeploy (v1: helmDeployBluegreenStep)
	// Step group: helmBluegreenDeployAction (uses helmDeployAction with strategy=blue-green), (conditional) helmTestAction
	StepTypeHelmBGDeploy: {
		{"ignoreReleaseHistFailStatus", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE"},
		{"skipSteadyStateCheck", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// useUpgradeInstall -> NO_ENV_MAPPING
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		{"flags", "steps.helmBluegreenDeployAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmBlueGreenSwapStep (v1: helmBluegreenSwapStep)
	// Step group: helmBluegreenSwapAction
	StepTypeHelmBlueGreenSwapStep: {
		{"flags", "steps.helmBluegreenSwapAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmCanaryDeploy (v1: helmDeployCanaryStep)
	// Step group: helmBasicDeployAction (uses helmDeployAction with strategy=canary), (conditional) helmTestAction
	StepTypeHelmCanaryDeploy: {
		{"ignoreReleaseHistFailStatus", "steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE"},
		{"skipSteadyStateCheck", "steps.helmBasicDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// useUpgradeInstall -> NO_ENV_MAPPING
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmBasicDeployAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		{"instanceSelection.type", "steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE"},
		// COMPLEX: "Count" -> "count", "Percentage" -> "percentage"
		{"instanceSelection.spec.count", "steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Count
		{"instanceSelection.spec.percentage", "steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES"},
		// only when type=Percentage
		{"flags", "steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmCanaryDelete (v1: helmCanaryDeleteStep)
	// Step group: helmUninstallAction (uses helmDeployAction — no with params from v0)
	// No field mappings (v0 spec is empty)
	StepTypeHelmCanaryDelete: {},

	// HelmDelete (v1: helmDeleteStep)
	// Step group: helmUninstallAction
	StepTypeHelmDelete: {
		{"releaseName", "steps.helmUninstallAction.spec.env.PLUGIN_RELEASE_NAME"},
		{"dryRun", "steps.helmUninstallAction.spec.env.PLUGIN_DRY_RUN"},
		{"commandFlags", "steps.helmUninstallAction.spec.env.PLUGIN_FLAGS"},
		// environmentVariables -> NO_ENV_MAPPING (helmUninstallAction does not have PLUGIN_ENV_VARS)
	},

	// HelmDeploy (v1: helmDeployBasicStep)
	// Step group: helmBasicDeployAction (uses helmDeployAction with strategy=basic), (conditional) helmTestAction
	StepTypeHelmDeploy: {
		{"ignoreReleaseHistFailStatus", "steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE"},
		{"skipSteadyStateCheck", "steps.helmBasicDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// useUpgradeInstall -> NO_ENV_MAPPING
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmBasicDeployAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		// skipDryRun -> NO_TEMPLATE_MAPPING (intentionally omitted in v1)
		// skipCleanup -> NO_TEMPLATE_MAPPING (intentionally omitted in v1)
		{"flags", "steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS"},
	},

	// HelmRollback (v1: helmRollbackStep)
	// Step group: helmRollbackAction, (conditional) helmTestAction
	StepTypeHelmRollback: {
		{"skipSteadyStateCheck", "steps.helmRollbackAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK"},
		// runChartTests -> NO_ENV_MAPPING (controls conditional execution of helmTestAction)
		{"environmentVariables", "steps.helmRollbackAction.spec.env.PLUGIN_ENV_VARS"},
		// COMPLEX: map[string]string converted to []map with key/value entries
		// skipDryRun -> NO_TEMPLATE_MAPPING (not used in conversion)
		{"flags", "steps.helmRollbackAction.spec.env.PLUGIN_FLAGS"},
	},

	// ============================================================
	// BUILD & PUSH STEPS
	// ============================================================
	// Note: Build & Push steps use conditional plugin selection (buildx/docker/kaniko).
	// Env var names are identical across all plugin variants.
	// "pushWithBuildx" is used as the representative action ID.

	// BuildAndPushDockerRegistry (v1: buildAndPushToDocker)
	// Step group: pushWithBuildx / pushWithDocker / pushWithKaniko (conditional)
	StepTypeBuildAndPushDockerRegistry: {
		// connectorRef -> COMPLEX (resolves to USERNAME, PASSWORD, REGISTRY via connector)
		{"repo", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// baseImageConnectorRefs -> COMPLEX (resolves to BASE_IMAGE_REGISTRY, BASE_IMAGE_USERNAME, BASE_IMAGE_PASSWORD via connector)
		// envVariables -> NO_ENV_MAPPING (merged into plugin env, not a single env var)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		// optimize -> NO_TEMPLATE_MAPPING
		// privileged -> NO_TEMPLATE_MAPPING
		// remoteCacheRepo -> NO_TEMPLATE_MAPPING
		// reports -> NO_TEMPLATE_MAPPING
		// resources -> NO_TEMPLATE_MAPPING
		// runAsUser -> NO_TEMPLATE_MAPPING
	},

	// BuildAndPushECR (v1: buildAndPushToECR)
	// Step group: pushWithBuildx / pushWithECR / pushWithKaniko (conditional)
	StepTypeBuildAndPushECR: {
		// connectorRef -> COMPLEX (resolves to ACCESS_KEY, SECRET_KEY, ASSUME_ROLE, EXTERNAL_ID, OIDC_TOKEN_ID via connector)
		{"region", "steps.pushWithBuildx.spec.with.REGION"},
		// account + region -> COMPLEX (constructed as "<account>.dkr.ecr.<region>.amazonaws.com" for REGISTRY)
		{"imageName", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// envVariables -> NO_ENV_MAPPING
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		// baseImageConnectorRefs -> COMPLEX (resolves via connector)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		// account -> NO_TEMPLATE_MAPPING (consumed to build registry URL)
		// runAsUser -> NO_TEMPLATE_MAPPING
	},

	// BuildAndPushGAR (v1: buildAndPushToGAR)
	// Step group: pushWithBuildx / pushWithGAR / pushWithKaniko (conditional)
	StepTypeBuildAndPushGAR: {
		// connectorRef -> COMPLEX (resolves to JSON_KEY, OIDC_TOKEN_ID, PROJECT_NUMBER, POOL_ID, PROVIDER_ID, SERVICE_ACCOUNT_EMAIL via connector)
		// host + projectID -> COMPLEX (constructed as "<host>/<projectID>" for REGISTRY)
		{"imageName", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// envVariables -> NO_ENV_MAPPING
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		// baseImageConnectorRefs -> COMPLEX (resolves via connector)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		// host -> NO_TEMPLATE_MAPPING (consumed to build registry URL)
		// projectID -> NO_TEMPLATE_MAPPING (consumed to build registry URL)
		// runAsUser -> NO_TEMPLATE_MAPPING
	},

	// BuildAndPushGCR (v1: N/A)
	// No conversion implemented
	StepTypeBuildAndPushGCR: {},

	// BuildAndPushACR (v1: buildAndPushToACR)
	// Step group: pushWithBuildx / pushWithACR / pushWithKaniko (conditional)
	StepTypeBuildAndPushACR: {
		// connectorRef -> COMPLEX (resolves to CLIENT_ID, TENANT_ID, CLIENT_SECRET, CLIENT_CERTIFICATE, OIDC_TOKEN_ID via connector)
		{"registry", "steps.pushWithBuildx.spec.with.REGISTRY"},
		{"imageName", "steps.pushWithBuildx.spec.with.REPO"},
		{"tags", "steps.pushWithBuildx.spec.with.TAGS"},
		// caching -> NO_ENV_MAPPING (controls which plugin variant is selected)
		// envVariables -> NO_ENV_MAPPING
		{"labels", "steps.pushWithBuildx.spec.with.CUSTOM_LABELS"},
		{"buildArgs", "steps.pushWithBuildx.spec.with.BUILD_ARGS"},
		// baseImageConnectorRefs -> COMPLEX (resolves via connector)
		{"dockerfile", "steps.pushWithBuildx.spec.with.DOCKERFILE"},
		{"context", "steps.pushWithBuildx.spec.with.CONTEXT"},
		{"target", "steps.pushWithBuildx.spec.with.TARGET"},
		{"subscriptionId", "steps.pushWithBuildx.spec.with.SUBSCRIPTION_ID"},
		// runAsUser -> NO_TEMPLATE_MAPPING
	},
}

// inside step.output
var StepOutputConversionRules = map[string][]ConversionRule{
	StepTypeHTTP: {
		{"httpUrl", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_URL"},
		{"httpMethod", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_METHOD"},
		{"httpResponseCode", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE"},
		{"httpResponseBody", "steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_BODY_BYTES"},
		{"status", "steps.httpStep.output.outputVariables.PLUGIN_EXECUTION_STATUS"},
		{"responseHeaders", "steps.httpStep.output.outputVariables.PLUGIN_RESPONSE_HEADERS"},
	},
}
