package terraform

type State struct {
	FormatVersion    string       `json:"format_version,omitempty"`
	TerraformVersion string       `json:"terraform_version,omitempty"`
	Values           *StateValues `json:"values,omitempty"`
	Checks           any          `json:"checks,omitempty"`
}

type StateValues struct {
	Outputs    map[string]Output `json:"outputs,omitempty"`
	RootModule Module            `json:"root_module,omitempty"`
}

type Plan struct {
	FormatVersion      string            `json:"format_version,omitempty"`
	TerraformVersion   string            `json:"terraform_version,omitempty"`
	Variables          Variables         `json:"variables,omitempty"`
	PlannedValues      StateValues       `json:"planned_values,omitempty"`
	ResourceDrift      []ResourceChange  `json:"resource_drift,omitempty"`
	ResourceChanges    []ResourceChange  `json:"resource_changes,omitempty"`
	OutputChanges      map[string]Change `json:"output_changes,omitempty"`
	PriorState         any               `json:"prior_state,omitempty"`
	Config             any               `json:"configuration,omitempty"`
	RelevantAttributes []ResourceAttr    `json:"relevant_attributes,omitempty"`
	Checks             any               `json:"checks,omitempty"`
	Timestamp          string            `json:"timestamp,omitempty"`
	Errored            bool              `json:"errored"`
}

type Variables any

type Output struct {
	Sensitive bool `json:"sensitive"`
	Type      any  `json:"type,omitempty"`
	Value     any  `json:"value,omitempty"`
}

type Module struct {
	Resources    []Resource `json:"resources,omitempty"`
	Address      string     `json:"address,omitempty"`
	ChildModules []Module   `json:"child_modules,omitempty"`
}

type ResourceChange any

type Change struct {
	Actions         []string `json:"actions,omitempty"`
	Before          any      `json:"before,omitempty"`
	After           any      `json:"after,omitempty"`
	AfterUnknown    any      `json:"after_unknown,omitempty"`
	BeforeSensitive any      `json:"before_sensitive,omitempty"`
	AfterSensitive  any      `json:"after_sensitive,omitempty"`
	ReplacePaths    any      `json:"replace_paths,omitempty"`
	Importing       any      `json:"importing,omitempty"`
	GeneratedConfig string   `json:"generated_config,omitempty"`
}

type ResourceAttr any

type Resource struct {
	Addr            string `json:"addr"`
	Module          string `json:"module"`
	Resource        string `json:"resource"`
	ImpliedProvider string `json:"implied_provider"`
	ResourceType    string `json:"resource_type"`
	ResourceName    string `json:"resource_name"`
	ResourceKey     string `json:"resource_key"`
}
