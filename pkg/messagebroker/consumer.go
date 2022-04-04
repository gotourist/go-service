package messagebroker

var ListenTopics []string = []string{"sale_order"}

// Consumer ...
type Consumer interface {
	Start()
}
