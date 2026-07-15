package auth

import (
	"iptv-spider-sh/model"
	"sort"
)

// ==================== 函数参数方式实现泛型排序 ====================

// SortChannelsByFields 使用字段提取函数实现泛型排序
// 参数说明：
//   - items: 需要排序的切片
//   - orderMap: Name 到优先级的映射
//   - getName: 提取 Name 的函数
//   - getGroup: 提取 Group 的函数
//   - getCommName: 提取 CommName 的函数
func SortChannelsByFields[T any](
	items []T,
	orderMap map[string]int,
	getName func(T) string,
	getGroup func(T) string,
	getCommName func(T) string,
) {
	if len(items) <= 1 {
		return
	}

	// 第一步：按 Group + Name(orderMap) 排序
	sort.SliceStable(items, func(i, j int) bool {
		groupI := getGroup(items[i])
		groupJ := getGroup(items[j])
		if groupI != groupJ {
			return groupI < groupJ
		}

		nameI := getName(items[i])
		nameJ := getName(items[j])
		orderI, okI := orderMap[nameI]
		orderJ, okJ := orderMap[nameJ]
		if okI != okJ {
			return okI
		}
		return orderI < orderJ
	})

	// 第二步：提取 CommName 的顺序（按第一次出现的顺序）
	commNameOrder := make(map[string]int)
	for _, item := range items {
		commName := getCommName(item)
		if _, exists := commNameOrder[commName]; !exists {
			commNameOrder[commName] = len(commNameOrder)
		}
	}

	// 第三步：按 Group + CommName 顺序排序（保持内部的 Name 顺序）
	sort.SliceStable(items, func(i, j int) bool {
		groupI := getGroup(items[i])
		groupJ := getGroup(items[j])
		if groupI != groupJ {
			return groupI < groupJ
		}
		return commNameOrder[getCommName(items[i])] <
			commNameOrder[getCommName(items[j])]
	})
}

// reorderSliceByFields 根据索引重新排列切片
func reorderSliceByFields[T any](indices []int, items []T) {
	if len(items) != len(indices) {
		return
	}
	reordered := make([]T, len(items))
	for newIdx, oldIdx := range indices {
		reordered[newIdx] = items[oldIdx]
	}
	copy(items, reordered)
}

// ==================== 便捷包装函数 ====================

// SortChannelUrlInfos 排序 ChannelUrlInfo 切片
func SortChannelUrlInfos(items []model.ChannelUrlInfo, orderMap map[string]int) {
	SortChannelsByFields(
		items,
		orderMap,
		func(item model.ChannelUrlInfo) string { return item.Name },
		func(item model.ChannelUrlInfo) string { return item.Group },
		func(item model.ChannelUrlInfo) string { return item.CommName },
	)
}

// SortChannelInfos 排序 ChannelInfo 切片
func SortChannelInfos(items []model.ChannelInfo, orderMap map[string]int) {
	SortChannelsByFields(
		items,
		orderMap,
		func(item model.ChannelInfo) string { return item.Name },
		func(item model.ChannelInfo) string { return item.Group },
		func(item model.ChannelInfo) string { return item.CommName },
	)
}
