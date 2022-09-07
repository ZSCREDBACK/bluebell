package redis

// 用于存放 redis key

// 实战经验-1:
// 为了避免与其他服务共用redis时,发生的key名重复冲突的问题
// 可以使用符号分割key的名称,实现类似于命名空间的隔离效果,也方便了查询和拆分

// 实战经验-2:
// 都以Key开头,方便写代码调用时进行检索(与状态码类似,虽然状态码我没有设置成这样)
// 尽量名称设置得见名知意

const (
	KeyPrefix      = "bluebell:"   // 一般使用:分割比较多,约定俗成,非必须
	KeyPostTime    = "post:time"   // ZSet;帖子及发帖时间
	KeyPostScore   = "post:score"  // ZSet;帖子及投票分数
	KeyVotedPrefix = "post:voted:" // ZSet;记录用户及投票类型(赞成or反对) // post:voted:post_id 最后一项是变化的,所以在设置Key名时加上了prefix来明晰其含义
)

// 为key值加上前缀
func getKey(key string) string {
	return KeyPrefix + key
}
