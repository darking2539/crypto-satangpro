package models

type GetListRequest struct {
	Address string `json:"address"`
	Page    int64  `json:"page" binding:"required"`
	PerPage int64  `json:"perPage" binding:"required"`
}

type GetListResponse struct {
	Pagination Pagination        `json:"pagination"`
	Data       []TransactionData `json:"data"`
}

type TransactionData struct {
	BlockNo          uint64 `json:"blockNumber"`
	TransactionIndex uint64 `json:"transactionIndex"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
}
