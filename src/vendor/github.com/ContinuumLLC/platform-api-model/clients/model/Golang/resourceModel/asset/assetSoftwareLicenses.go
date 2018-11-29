package asset

//SoftwareLicenses keeps the product license info
type SoftwareLicenses struct {
	ProductName string `json:"productName,omitempty" cql:"product_name"`
	ProductID   string `json:"productID,omitempty" cql:"product_id"`
	ProductKey  string `json:"productKey,omitempty" cql:"product_key"`
}
