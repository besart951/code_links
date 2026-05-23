package productcatalog

var IDs = []string{"infra-link", "planer-link", "loka-link"}

func IsValidID(productID string) bool {
	for _, id := range IDs {
		if productID == id {
			return true
		}
	}

	return false
}
