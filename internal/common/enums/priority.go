package enums

type ProducerOrderPriority string

const (
	PRODUCER_ORDER_PRIORITY_ORDER            ProducerOrderPriority = "ORDER"
	PRODUCER_ORDER_PRIORITY_BALANCED         ProducerOrderPriority = "BALANCED"
	PRODUCER_ORDER_PRIORITY_HIGH_PERFORMANCE ProducerOrderPriority = "HIGH_PERFORMANCE"
)

type ConsumerOrderPriority string

const (
	CONSUMER_ORDER_PRIORITY_ORDER            ConsumerOrderPriority = "ORDER"
	CONSUMER_ORDER_PRIORITY_BALANCED         ConsumerOrderPriority = "BALANCED"
	CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE ConsumerOrderPriority = "HIGH_PERFORMANCE"
	CONSUMER_ORDER_PRIORITY_RISKY            ConsumerOrderPriority = "RISKY"
)
