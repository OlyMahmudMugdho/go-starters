package types

type StringOpsRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StringOpsGetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StringOpsSetResponse struct {
	Status string `json:"status"`
}

type HashRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
