package convertexpressions

import (
	"testing"
)

func TestTrieRules_PipelineLevel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		{
			name:     "stage variables",
			path:     "stage.variables.var1",
			context:  nil,
			expected: "stage.variables.var1",
		},
		{
			name:     "stage identifier FQN",
			path:     "pipeline.stages.build.identifier",
			context:  nil,
			expected: "pipeline.stages.build.id",
		},
		{
			name:     "stage env identifier FQN",
			path:     "pipeline.stages.build.spec.env.identifier",
			context:  nil,
			expected: "pipeline.stages.build.env.id",
		},
		{
			name:     "stage env group identifier FQN",
			path:     "pipeline.stages.build.spec.env.envGroupRef",
			context:  nil,
			expected: "pipeline.stages.build.env.group.id",
		},
		{
			name:     "stage env group name FQN",
			path:     "pipeline.stages.build.spec.env.envGroupName",
			context:  nil,
			expected: "pipeline.stages.build.env.group.name",
		},
		{
			name:     "stage env identifier relative",
			path:     "stage.spec.env.identifier",
			context:  nil,
			expected: "stage.env.id",
		},
		{
			name:     "stage env group identifier relative",
			path:     "spec.env.envGroupRef",
			context:  nil,
			expected: "env.group.id",
		},
		{
			name:     "stage env group name relative",
			path:     "spec.env.envGroupName",
			context:  nil,
			expected: "env.group.name",
		},
		{
			name:     "stage env identifier direct",
			path:     "env.identifier",
			context:  nil,
			expected: "env.id",
		},
		{
			name:     "stage env group identifier direct",
			path:     "env.envGroupRef",
			context:  nil,
			expected: "env.group.id",
		},
		{
			name:     "stage env group name direct",
			path:     "env.envGroupName",
			context:  nil,
			expected: "env.group.name",
		},
		{
			name:     "spec.execution.steps removal FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1",
			context:  nil,
			expected: "pipeline.stages.build.steps.step1.output.outputVariables.var1",
		},
		{
			name:     "is not resolved",
			path:     "expression.isResolved(<+pipeline.variables.var1>)",
			context:  nil,
			expected: "expression.isResolved(<+pipeline.variables.var1>)",
		},
		{
			name:     "nested and function",
			path:     "<+pipeline.variables.var1>.some_func(\"param\")",
			context:  nil,
			expected: "<+pipeline.variables.var1>.some_func(\"param\")",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_StepLevel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// General step rules (no context needed)
		{
			name:     "step identifier FQN - no context",
			path:     "pipeline.stages.build.spec.execution.steps.step1.identifier",
			context:  nil,
			expected: "pipeline.stages.build.steps.step1.id",
		},
		{
			name:     "step identifier FQN - with context",
			path:     "pipeline.stages.build.spec.execution.steps.step1.identifier",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.id",
		},
		{
			name:     "output variables FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.output.outputVariables.var1",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.output.outputVariables.var1",
		},
		{
			name:     "output variables relative",
			path:     "spec.execution.steps.step1.output.outputVariables.var1",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "steps.step1.output.outputVariables.var1",
		},
		{
			name:     "output variables relative with function",
			path:     "spec.execution.steps.step1.output.outputVariables.var1.some_func(\"param\")",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "steps.step1.output.outputVariables.var1.some_func(\"param\")",
		},
		// Failure strategies
		{
			name:     "failure strategy in step - errors",
			path:     "pipeline.stages.build.spec.execution.steps.step1.failureStrategies[0].onFailure.errors",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.onFailure[0].errors",
		},
		{
			name:     "failure strategy in step - retry count",
			path:     "pipeline.stages.build.spec.execution.steps.step1.failureStrategies[0].onFailure.action.specConfig.retryCount",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.onFailure[0].action.retry.attempts",
		},
		// Step-specific rules with context - SaveCacheGCS
		{
			name:     "save cache to gcs FQN",
			path:     "pipeline.stages.build.spec.execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: "pipeline.stages.build.steps.STEPID.steps.saveCacheGCS.spec.with.BUCKET",
		},
		{
			name:     "save cache to gcs relative",
			path:     "execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeSaveCacheGCS},
			expected: "steps.STEPID.steps.saveCacheGCS.spec.with.BUCKET",
		},
		// Step-specific rules with context - RestoreCacheS3
		{
			name:     "restore cache from s3 FQN",
			path:     "pipeline.stages.build.spec.execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "pipeline.stages.build.steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		{
			name:     "restore cache from s3 relative",
			path:     "execution.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		{
			name:     "restore cache from s3 inside step group",
			path:     "pipeline.stages.build.spec.execution.steps.stepGroupID.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "pipeline.stages.build.steps.stepGroupID.steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		{
			name:     "restore cache from s3 relative stepgroup",
			path:     "stepGroup.steps.STEPID.spec.bucket",
			context:  &ConversionContext{StepType: StepTypeRestoreCacheS3},
			expected: "group.steps.STEPID.steps.restoreCacheS3.spec.with.BUCKET",
		},
		// Step-specific rules with context - Run
		{
			name:     "run step command FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.command",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.spec.script",
		},
		{
			name:     "run step image FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.image",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.spec.container.image",
		},
		{
			name:     "run step envVariables FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.envVariables",
			context:  &ConversionContext{StepType: StepTypeRun},
			expected: "pipeline.stages.build.steps.step1.spec.env",
		},
		// Step-specific rules with context - HTTP
		{
			name:     "http step url FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.spec.url",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.spec.env.PLUGIN_URL",
		},
		{
			name:     "http step output httpResponseCode FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.httpResponseCode",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_HTTP_RESPONSE_CODE",
		},
		{
			name:     "http step output status FQN",
			path:     "pipeline.stages.build.spec.execution.steps.http1.output.status",
			context:  &ConversionContext{StepType: StepTypeHTTP},
			expected: "pipeline.stages.build.steps.http1.steps.httpStep.output.outputVariables.PLUGIN_EXECUTION_STATUS",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_K8sSteps(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// K8sRollingDeploy
		{
			name:     "K8sRollingDeploy skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sRollingDeploy pruningEnabled FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.pruningEnabled",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingPrepareAction.spec.env.PLUGIN_PRUNING_ENABLED",
		},
		{
			name:     "K8sRollingDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sRollingRollback
		{
			name:     "K8sRollingRollback pruningEnabled FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.pruningEnabled",
			context:  &ConversionContext{StepType: StepTypeK8sRollingRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingRollbackAction.spec.env.PLUGIN_PRUNING",
		},
		{
			name:     "K8sRollingRollback flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sRollingRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRollingRollbackAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sApply
		{
			name:     "K8sApply filePaths FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.filePaths",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_MANIFEST_PATH",
		},
		{
			name:     "K8sApply skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sApply flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sApply},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sBGSwapServices
		{
			name:     "K8sBGSwapServices stable_service FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.stable_service",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STABLE_SERVICE",
		},
		{
			name:     "K8sBGSwapServices stage_service FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.stage_service",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.PLUGIN_STAGE_SERVICE",
		},
		{
			name:     "K8sBGSwapServices is_openshift FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.is_openshift",
			context:  &ConversionContext{StepType: StepTypeK8sBGSwapServices},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenSwapServicesSelectorsAction.spec.env.HARNESS_IS_OPENSHIFT",
		},
		// K8sCanaryDelete
		{
			name:     "K8sCanaryDelete resources FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.resources",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryDeleteAction.spec.env.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sCanaryDelete is_openshift FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.is_openshift",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryDeleteAction.spec.env.HARNESS_IS_OPENSHIFT",
		},
		// K8sRollout
		{
			name:     "K8sRollout command FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.command",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_COMMAND",
		},
		{
			name:     "K8sRollout resources.spec.resourceNames FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.resources.spec.resourceNames",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sRollout resources.spec.manifestPaths FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.resources.spec.manifestPaths",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_MANIFESTS",
		},
		{
			name:     "K8sRollout flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sRollout},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sRolloutStep.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sScale
		{
			name:     "K8sScale instanceSelection.type FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.type",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE",
		},
		{
			name:     "K8sScale instanceSelection.spec.count FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.spec.count",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.spec.env.PLUGIN_INSTANCES",
		},
		{
			name:     "K8sScale workload FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.workload",
			context:  &ConversionContext{StepType: StepTypeK8sScale},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sScaleAction.spec.env.PLUGIN_WORKLOAD",
		},
		// K8sDryRun
		{
			name:     "K8sDryRun encryptYamlOutput FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.encryptYamlOutput",
			context:  &ConversionContext{StepType: StepTypeK8sDryRun},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDryRunAction.spec.env.PLUGIN_ENCRYPT_YAML_OUTPUT",
		},
		// K8sDelete
		{
			name:     "K8sDelete deleteResources.spec.resourceNames FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.deleteResources.spec.resourceNames",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_RESOURCES",
		},
		{
			name:     "K8sDelete deleteResources.spec.manifestPaths FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.deleteResources.spec.manifestPaths",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_MANIFESTS",
		},
		{
			name:     "K8sDelete deleteResources.spec.deleteNamespace FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.deleteResources.spec.deleteNamespace",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_INCLUDE_NAMESPACES",
		},
		{
			name:     "K8sDelete flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sDeleteAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sTrafficRouting
		{
			name:     "K8sTrafficRouting trafficRouting.provider FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.provider",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_PROVIDER",
		},
		{
			name:     "K8sTrafficRouting trafficRouting.spec.name FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.name",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_RESOURCE_NAME",
		},
		{
			name:     "K8sTrafficRouting trafficRouting.spec.hosts FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.hosts",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_HOSTNAMES",
		},
		{
			name:     "K8sTrafficRouting trafficRouting.spec.routes FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.routes",
			context:  &ConversionContext{StepType: StepTypeK8sTrafficRouting},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficShiftAction.spec.env.PLUGIN_ROUTES",
		},
		// K8sCanaryDeploy
		{
			name:     "K8sCanaryDeploy instanceSelection.type FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.type",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sCanaryPrepareAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE",
		},
		{
			name:     "K8sCanaryDeploy skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sCanaryDeploy trafficRouting.provider FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.provider",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficRoutingAction.spec.env.PLUGIN_PROVIDER",
		},
		{
			name:     "K8sCanaryDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sBlueGreenDeploy
		{
			name:     "K8sBlueGreenDeploy skipDryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
		{
			name:     "K8sBlueGreenDeploy pruningEnabled FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.pruningEnabled",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_RELEASE_PRUNING_ENABLED",
		},
		{
			name:     "K8sBlueGreenDeploy skipUnchangedManifest FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipUnchangedManifest",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sBlueGreenPrepareAction.spec.env.PLUGIN_SKIP_UNCHANGED_MANIFEST",
		},
		{
			name:     "K8sBlueGreenDeploy trafficRouting.spec.gateways FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.trafficRouting.spec.gateways",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sTrafficRoutingAction.spec.env.PLUGIN_GATEWAYS",
		},
		{
			name:     "K8sBlueGreenDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeK8sBlueGreenDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_FLAGS_JSON",
		},
		// K8sPatch
		{
			name:     "K8sPatch workload FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.workload",
			context:  &ConversionContext{StepType: StepTypeK8sPatch},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sPatchAction.spec.env.PLUGIN_WORKLOAD",
		},
		{
			name:     "K8sPatch mergeStrategyType FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.mergeStrategyType",
			context:  &ConversionContext{StepType: StepTypeK8sPatch},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sPatchAction.spec.env.PLUGIN_MERGE_STRATEGY",
		},
		{
			name:     "K8sPatch source.spec.content FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.source.spec.content",
			context:  &ConversionContext{StepType: StepTypeK8sPatch},
			expected: "pipeline.stages.deploy.steps.step1.steps.k8sPatchAction.spec.env.PLUGIN_CONTENT",
		},
		// K8sRollingDeploy relative path
		{
			name:     "K8sRollingDeploy skipDryRun relative",
			path:     "execution.steps.step1.spec.skipDryRun",
			context:  &ConversionContext{StepType: StepTypeK8sRollingDeploy},
			expected: "steps.step1.steps.k8sApplyAction.spec.env.PLUGIN_SKIP_DRY_RUN",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_HelmSteps(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// HelmBGDeploy
		{
			name:     "HelmBGDeploy ignoreReleaseHistFailStatus FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
		{
			name:     "HelmBGDeploy skipSteadyStateCheck FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipSteadyStateCheck",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK",
		},
		{
			name:     "HelmBGDeploy environmentVariables FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.environmentVariables",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_ENV_VARS",
		},
		{
			name:     "HelmBGDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmBlueGreenSwapStep
		{
			name:     "HelmBlueGreenSwapStep flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmBlueGreenSwapStep},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBluegreenSwapAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmCanaryDeploy
		{
			name:     "HelmCanaryDeploy ignoreReleaseHistFailStatus FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
		{
			name:     "HelmCanaryDeploy instanceSelection.type FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.type",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES_UNIT_TYPE",
		},
		{
			name:     "HelmCanaryDeploy instanceSelection.spec.count FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.instanceSelection.spec.count",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_INSTANCES",
		},
		{
			name:     "HelmCanaryDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmCanaryDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmDelete
		{
			name:     "HelmDelete releaseName FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.releaseName",
			context:  &ConversionContext{StepType: StepTypeHelmDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmUninstallAction.spec.env.PLUGIN_RELEASE_NAME",
		},
		{
			name:     "HelmDelete dryRun FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.dryRun",
			context:  &ConversionContext{StepType: StepTypeHelmDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmUninstallAction.spec.env.PLUGIN_DRY_RUN",
		},
		{
			name:     "HelmDelete commandFlags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.commandFlags",
			context:  &ConversionContext{StepType: StepTypeHelmDelete},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmUninstallAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmDeploy
		{
			name:     "HelmDeploy ignoreReleaseHistFailStatus FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
		{
			name:     "HelmDeploy environmentVariables FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.environmentVariables",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_ENV_VARS",
		},
		{
			name:     "HelmDeploy flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmDeploy},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmBasicDeployAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmRollback
		{
			name:     "HelmRollback skipSteadyStateCheck FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.skipSteadyStateCheck",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.spec.env.PLUGIN_SKIP_STEADY_STATE_CHECK",
		},
		{
			name:     "HelmRollback environmentVariables FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.environmentVariables",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.spec.env.PLUGIN_ENV_VARS",
		},
		{
			name:     "HelmRollback flags FQN",
			path:     "pipeline.stages.deploy.spec.execution.steps.step1.spec.flags",
			context:  &ConversionContext{StepType: StepTypeHelmRollback},
			expected: "pipeline.stages.deploy.steps.step1.steps.helmRollbackAction.spec.env.PLUGIN_FLAGS",
		},
		// HelmBGDeploy relative path
		{
			name:     "HelmBGDeploy ignoreReleaseHistFailStatus relative",
			path:     "execution.steps.step1.spec.ignoreReleaseHistFailStatus",
			context:  &ConversionContext{StepType: StepTypeHelmBGDeploy},
			expected: "steps.step1.steps.helmBluegreenDeployAction.spec.env.PLUGIN_IGNORE_HISTORY_FAILURE",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrieRules_BuildAndPushSteps(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		context  *ConversionContext
		expected string
	}{
		// BuildAndPushDockerRegistry
		{
			name:     "BuildAndPushDockerRegistry repo FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.repo",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushDockerRegistry tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushDockerRegistry dockerfile FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.dockerfile",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.DOCKERFILE",
		},
		{
			name:     "BuildAndPushDockerRegistry context FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.context",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.CONTEXT",
		},
		{
			name:     "BuildAndPushDockerRegistry labels FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.labels",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.CUSTOM_LABELS",
		},
		{
			name:     "BuildAndPushDockerRegistry buildArgs FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.buildArgs",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.BUILD_ARGS",
		},
		{
			name:     "BuildAndPushDockerRegistry target FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.target",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TARGET",
		},
		// BuildAndPushECR
		{
			name:     "BuildAndPushECR region FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.region",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REGION",
		},
		{
			name:     "BuildAndPushECR imageName FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.imageName",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushECR tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushECR dockerfile FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.dockerfile",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.DOCKERFILE",
		},
		// BuildAndPushGAR
		{
			name:     "BuildAndPushGAR imageName FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.imageName",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushGAR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushGAR tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushGAR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushGAR labels FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.labels",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushGAR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.CUSTOM_LABELS",
		},
		// BuildAndPushACR
		{
			name:     "BuildAndPushACR registry FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.registry",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REGISTRY",
		},
		{
			name:     "BuildAndPushACR imageName FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.imageName",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		{
			name:     "BuildAndPushACR tags FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.tags",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.TAGS",
		},
		{
			name:     "BuildAndPushACR subscriptionId FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.subscriptionId",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.SUBSCRIPTION_ID",
		},
		{
			name:     "BuildAndPushACR dockerfile FQN",
			path:     "pipeline.stages.build.spec.execution.steps.step1.spec.dockerfile",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushACR},
			expected: "pipeline.stages.build.steps.step1.steps.pushWithBuildx.spec.with.DOCKERFILE",
		},
		// BuildAndPushDockerRegistry relative path
		{
			name:     "BuildAndPushDockerRegistry repo relative",
			path:     "execution.steps.step1.spec.repo",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushDockerRegistry},
			expected: "steps.step1.steps.pushWithBuildx.spec.with.REPO",
		},
		// BuildAndPushECR inside stepGroup
		{
			name:     "BuildAndPushECR region in stepGroup",
			path:     "pipeline.stages.build.spec.execution.steps.group1.steps.step1.spec.region",
			context:  &ConversionContext{StepType: StepTypeBuildAndPushECR},
			expected: "pipeline.stages.build.steps.group1.steps.step1.steps.pushWithBuildx.spec.with.REGION",
		},
	}

	trie := buildPipelineTrie()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := trie.Match(tt.path, tt.context)
			if result != tt.expected {
				t.Errorf("Match() = %v, want %v", result, tt.expected)
			}
		})
	}
}
