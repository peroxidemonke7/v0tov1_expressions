package convertexpressions

var StageConversionRules = []ConversionRule{
	{"identifier", "id"},
}

var DeploymentStageSpecConversionRules = []ConversionRule{
	{"env{alias: env}.identifier", "env.id"},
	{"env{alias: env}.envGroupName", "env.group.name"},
	{"env{alias: env}.envGroupRef", "env.group.id"},
	{"infra{alias: infra}.connectorRef", "infra.connector"},
	{"service{alias: service}.identifier", "service.id"},
}
