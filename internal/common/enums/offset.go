package enums

type AutoOffsetReset string

const (
	OFFSET_RESET_ERROR     AutoOffsetReset = "error"
	OFFSET_RESET_SMALLEST  AutoOffsetReset = "smallest"
	OFFSET_RESET_EARLIEST  AutoOffsetReset = "earliest"
	OFFSET_RESET_BEGINNING AutoOffsetReset = "beginning"
	OFFSET_RESET_LARGEST   AutoOffsetReset = "largest"
	OFFSET_RESET_LATEST    AutoOffsetReset = "latest"
	OFFSET_RESET_END       AutoOffsetReset = "end"
)
