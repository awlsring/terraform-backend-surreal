package state

type TfResourceInstance struct {
	IndexKey int `json:"index_key,"`
	SchemaVersion int `json:"schema_version"`
	Attributes map[string]interface{} `json:"attributes"`
}

type TfResource struct {
	Mode string `json:"mode"`
	Type string `json:"type"`
	Name string `json:"name"`
	ProviderName string `json:"provider"`
	Instances []TfResourceInstance `json:"instances"`
}

type TfState struct {
	Version int `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial int `json:"serial"`
	Lineage string `json:"lineage"`
	Resources []TfResource `json:"resources"`
	CheckResult *TfCheckResult `json:"check_result"`
	Outputs map[string]TfOutput `json:"outputs"`
}

type TfOutput struct {
	Sensitive bool `json:"sensitive,omitempty"`
	Type string `json:"type"`
	Value interface{} `json:"value"`
}

type TfCheckResult struct {
	Checked bool `json:"checked"`
	Valid bool `json:"valid"`
	Errors []string `json:"errors"`
}