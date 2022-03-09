package commonstruct

// 要收集日志的结构体
type LogsEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}
