package readme

// OutboundIPEndpoint is the ReadMe API URL endpoint for retrieving ReadMe's outbound IP addresses.
const OutboundIPEndpoint = "/outbound-ips"

// OutboundIPService is an interface for using the outbound IPs endpoint of the ReadMe.com API.
//
// API Reference: https://docs.readme.com/main/reference/getoutboundips
type OutboundIPService interface {
	// Get all of ReadMe’s IP addresses used for outbound webhook requests and the “Try It!” button on the API Explorer.
	//
	// Although ReadMe’s outbound IP addresses may change, the IPs in this API
	// response will be valid for at least 7 days. If you configure your API or
	// webhooks to limit access based on these IPs, you should refresh the IP list
	// from this endpoint weekly.
	//
	// API Reference: https://docs.readme.com/main/reference/getoutboundips
	Get() ([]OutboundIP, *APIResponse, error)
}

// OutboundIPClient handles communication with the OutboundIP related methods of the ReadMe.com API.
type OutboundIPClient struct {
	client *Client
}

// OutboundIP represents the response from the ReadMe API when returning outbound IP addresses.
type OutboundIP struct {
	IPAddress string `json:"ipAddress"`
}

var _ OutboundIPService = &OutboundIPClient{}

// Get all of ReadMe’s IP addresses used for outbound webhook requests and the “Try It!” button on the API Explorer.
//
// API Reference: https://docs.readme.com/main/reference/getoutboundips
func (c OutboundIPClient) Get() ([]OutboundIP, *APIResponse, error) {
	ipList := []OutboundIP{}
	apiResponse, err := c.client.APIRequest(&APIRequest{
		Method:       "GET",
		Endpoint:     OutboundIPEndpoint,
		UseAuth:      false,
		OkStatusCode: []int{200},
		Response:     &ipList,
	})

	return ipList, apiResponse, err
}
