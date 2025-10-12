package yahoofinance

type SubscribeMessage struct {
	Subscribe []string `json:"subscribe,omitempty"`
}

type ResponseMessage struct {
	Symbol           string  `json:"symbol,omitempty" protoIndex:"1"`
	Price            float32 `json:"price,omitempty" protoIndex:"2"`
	FiatSymbol       string  `json:"fiat_symbol,omitempty" protoIndex:"4"`
	Change24hPercent float32 `json:"change_24h_percent,omitempty" protoIndex:"8"`
	Change24hFiat    float32 `json:"change_24h_fiat,omitempty" protoIndex:"12"`
}

type rawResponse struct {
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
}
