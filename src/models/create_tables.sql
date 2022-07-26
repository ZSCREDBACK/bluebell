-- 创建一张用户表
USE `bluebell`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `username` varchar(64) COLLATE utf8_general_ci NOT NULL,
  `password` varchar(64) COLLATE utf8_general_ci NOT NULL,
  `email` varchar(64) COLLATE utf8_general_ci,
  `gender` tinyint(4) NOT NULL DEFAULT '0',
  `created_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  unique key `idx_username` (`username`) USING BTREE,
    unique key `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- bigint() 等于golang的int64

-- 为什么不使用自增id作为用户id
-- 访问者可以通过注册用户来获取应用的注册人数
-- 分库分表时,用户id可能会造成重复

-- 为什么不使用UUID作为用户id
-- UUID是随机的,会丧失按时间排序的功能(因为自增id会随着时间增加)

-- 大型企业常常使用分布式ID生成器来生成用户id