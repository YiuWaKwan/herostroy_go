package lazy_save

import (
	"hero_story.go_server/comm/log"
	"sync"
	"time"
)

// 这个就相当于 ConcurrentHashMap
var lsoMap = &sync.Map{}

//
// 初始化函数
func init() {
	startSave()
}

//
// SaveOrUpdate 保存或更新
func SaveOrUpdate(lso LazySaveObj) {
	if nil == lso {
		return
	}

	log.Info("记录延迟保存对象, lsoId = %s", lso.GetLsoId())

	nowTime := time.Now().UnixMilli()
	existRecord, _ := lsoMap.Load(lso.GetLsoId())

	if nil != existRecord {
		existRecord.(*lazySaveRecord).SetLastUpdateTime(nowTime)
		return
	}

	newRecord := &lazySaveRecord{}
	newRecord.objRef = lso
	newRecord.SetLastUpdateTime(nowTime)
	lsoMap.Store(lso.GetLsoId(), newRecord)
}

// 开始保存
func startSave() {
	go func() {
		for {
			// 先休息 1 秒
			time.Sleep(time.Second)

			nowTime := time.Now().UnixMilli()
			deleteLsoIdArray := make([]string, 64)

			lsoMap.Range(func(_, val interface{}) bool {
				if nil == val {
					return true
				}

				currRecord := val.(*lazySaveRecord)

				if nowTime-currRecord.GetLastUpdateTime() < 20000 {
					// 如果时间差小于 20 秒
					// 不进行保存
					return true
				}

				log.Info(
					"执行延迟保存, lsoId = %+v",
					currRecord.objRef.GetLsoId(),
				)

				// 执行保存逻辑
				currRecord.objRef.SaveOrUpdate()

				deleteLsoIdArray = append(deleteLsoIdArray, currRecord.objRef.GetLsoId())
				return true
			})

			for _, lsoId := range deleteLsoIdArray {
				lsoMap.Delete(lsoId)
			}
		}
	}()
}
