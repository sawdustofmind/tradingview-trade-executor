package dto

const ()

type Signal struct {
	Exchange     string `json:"exchange"`
	Symbol       string `json:"symbol"`
	StrategyName string `json:"strategy_name"`
	Action       string `json:"action"`
}
