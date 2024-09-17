package request

type UserAddressReq struct {
	Title       string `json:"title"`
	FullAddress string `json:"full_address"`
}
