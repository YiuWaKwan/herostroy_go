package lazy_save

// 延迟保存记录
type lazySaveRecord struct {
	// 延迟保存对象
	objRef LazySaveObj
	// 最后修改时间
	lastUpdateTime int64
}

//
// GetLastUpdateTime 获取最后更新时间
func (record *lazySaveRecord) GetLastUpdateTime() int64 {
	return record.lastUpdateTime
}

//
// SetLastUpdateTime 设置最后更新时间
func (record *lazySaveRecord) SetLastUpdateTime(val int64) {
	record.lastUpdateTime = val
}
