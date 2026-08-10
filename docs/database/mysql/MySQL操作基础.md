# MySQL 基础操作

## 终端中操作数据库

1. 查看数据库

```sql
show databases;
```

2. 创建数据库

```sql
create database xxx;
```

3. 使用数据库

```sql
use xxx;
```

4. 创建一张表

```sql
create table user(
	name varchar(20),
    age int,
    height double
);

# 插入数据
insert into user (name, age, height) values ('fan', 18, 1.80);
```

## SQL

SQL 语句的常用规范：

- 通常关键字是大写的，比如 CREATE、TABLE、SHOW 等等
- 一条语句结束后，需要以 `;` 结尾
- 如果遇到关键字作为表名或者字段名称，可以使用``包裹

SQL 可以分为常见的四大类：

- DDL（Data Definition Language）：数据定义语言； 

- - 可以通过DDL语句对数据库或者表进行：创建、删除、修改等操作；

- DML（Data Manipulation Language）：数据操作语言； 

- - 可以通过DML语句对表进行：添加、删除、修改等操作；

- DQL（Data Query Language）：数据查询语言； 

- - 可以通过DQL从数据库中查询记录；（重点）

- DCL（Data Control Language）：数据控制语言； 

- - 对数据库、表格的权限进行相关访问控制操作；

SQL 的数据类型

- 数字类型：INT、BIGINT、FLOAT、DOUBLE、DECIMAL
- 时间类型：DATE、YEAR、DATETIME、TIMESTAMP
- 字符串类型：CHAR、VARCHAR、BLOB、TEXT

表约束：

- 主键：PRIMARY KEY
- 唯一：UNIQUE
- 不能为空：NOT NULL
- 默认值：DEFAULT
- 自动递增：AUTO_INCREMENT

## SQL 中的函数

### 聚合函数

| 函数                  | 作用             | 示例                         |
| --------------------- | ---------------- | ---------------------------- |
| `COUNT(*)`            | 统计行数         | `SELECT COUNT(*) FROM user;` |
| `COUNT(col)`          | 统计非 NULL 数量 | `COUNT(age)`                 |
| `COUNT(DISTINCT col)` | 统计不同值数量   | `COUNT(DISTINCT age)`        |
| `SUM(col)`            | 求和             | `SUM(price)`                 |
| `AVG(col)`            | 平均值           | `AVG(score)`                 |
| `MAX(col)`            | 最大值           | `MAX(age)`                   |
| `MIN(col)`            | 最小值           | `MIN(age)`                   |

### 数学函数

| 函数         | 作用          | 示例               | 结果   |
| ------------ | ------------- | ------------------ | ------ |
| `ABS(x)`     | 绝对值        | `ABS(-10)`         | `10`   |
| `ROUND(x,n)` | 四舍五入      | `ROUND(3.14159,2)` | `3.14` |
| `ROUND(x)` | 四舍五入      | `ROUND(3.14159)` | `3` |
| `CEIL(x)`    | 向上取整      | `CEIL(3.2)`        | `4`    |
| `FLOOR(x)`   | 向下取整      | `FLOOR(3.9)`       | `3`    |
| `MOD(a,b)`   | 取余          | `MOD(10,3)`        | `1`    |
| `RAND()`     | 随机数（0~1） | `RAND()`           | 随机值 |
| `POWER(a,b)` | 幂运算        | `POWER(2,3)`       | `8`    |
| `SQRT()`     | 平方根        | `SQRT(4)`          | `2`    |

### 字符串函数

| 函数                       | 作用                     | 示例                      | 结果       |
| -------------------------- | ------------------------ | ------------------------- | ---------- |
| `CONCAT(a,b)`              | 拼接字符串               | `CONCAT('Hello','SQL')`   | `HelloSQL` |
| `LENGTH(str)`              | 返回字符串长度（字节数） | `LENGTH('abc')`           | `3`        |
| `CHAR_LENGTH(str)`         | 返回字符长度             | `CHAR_LENGTH('你好')`     | `2`        |
| `UPPER(str)`               | 转大写                   | `UPPER('abc')`            | `ABC`      |
| `LOWER(str)`               | 转小写                   | `LOWER('ABC')`            | `abc`      |
| `SUBSTRING(str,start,len)` | 截取字符串               | `SUBSTRING('abcdef',2,3)` | `bcd`      |
| `LEFT(str,n)`              | 取左边 n 个字符          | `LEFT('abcdef',2)`        | `ab`       |
| `RIGHT(str,n)`             | 取右边 n 个字符          | `RIGHT('abcdef',2)`       | `ef`       |
| `REPLACE(str,old,new)`     | 替换字符串               | `REPLACE('abc','a','x')`  | `xbc`      |
| `TRIM(str)`                | 去除首尾空格             | `TRIM(' abc ')`           | `abc`      |

### 日期时间函数

| 函数                            | 作用           | 示例                                  | 结果                  |
| ------------------------------- | -------------- | ------------------------------------- | --------------------- |
| `NOW()`                         | 当前日期和时间 | `NOW()`                               | `2026-08-06 16:30:00` |
| `CURDATE()`                     | 当前日期       | `CURDATE()`                           | `2026-08-06`          |
| `CURTIME()`                     | 当前时间       | `CURTIME()`                           | `16:30:00`            |
| `YEAR(date)`                    | 获取年份       | `YEAR(NOW())`                         | `2026`                |
| `MONTH(date)`                   | 获取月份       | `MONTH(NOW())`                        | `8`                   |
| `DAY(date)`                     | 获取日期       | `DAY(NOW())`                          | `6`                   |
| `DATEDIFF(d1,d2)`               | 日期差（天）   | `DATEDIFF('2026-08-10','2026-08-06')` | `4`                   |
| `DATE_ADD(date,INTERVAL n DAY)` | 日期加减       | `DATE_ADD(NOW(),INTERVAL 7 DAY)`      | 7 天后                |
| `DATE_SUB(date,INTERVAL n DAY)` | 日期减         | `DATE_SUB(NOW(),INTERVAL 7 DAY)`      | 7 天前                |

### 空值处理函数

| 函数              | 作用                           | 示例                      | 结果   |
| ----------------- | ------------------------------ | ------------------------- | ------ |
| `COALESCE(a,b,c)` | 返回第一个非 NULL 值           | `COALESCE(NULL,NULL,100)` | `100`  |
| `IFNULL(a,b)`     | 如果 a 为 NULL 返回 b（MySQL） | `IFNULL(NULL,0)`          | `0`    |
| `NULLIF(a,b)`     | 如果 a=b 返回 NULL，否则返回 a | `NULLIF(10,10)`           | `NULL` |

### 条件函数

| 函数                                  | 作用                   | 示例                            |
| ------------------------------------- | ---------------------- | ------------------------------- |
| `IF(condition,a,b)`                   | 条件判断（MySQL）      | `IF(score>=60,'及格','不及格')` |
| `CASE WHEN ... THEN ... ELSE ... END` | 多条件判断（标准 SQL） | 见下例                          |

## SQL操作：DDL(操作表)

### 查询表

```sql
# 查询当前数据库下所有表名称
SHOW TABLES;
# 查询表结构
DESC 表名称;
```

### 创建表

创建数据表

```sql
CREATE TABLE 表名 (
	字段名1  数据类型1,
	字段名2  数据类型2,
	…
	字段名n  数据类型n
);
```

```sql
create table if not exists `users` (
	id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(20) NOT NULL,
    age INT DEFAULT 0,
    telPhone VARCHAR(20) DEFAULT '' UNIQUE NOT NULL
);
```

### 删除表

```sql
DROP TABLE 表名;
# 删除时判断表是否存在
DROP TABLE IF EXISTS 表名;
```

```sql
drop table users;
drop table if exists users;
```

### 修改表

```sql
# 修改表名
ALTER TABLE 表名 RENAME TO 新的表名;
# 添加一个新的列
ALTER TABLE 表名 ADD 列名 数据类型;
# 删除一列数据
ALTER TABLE 表名 DROP 列名;
# 修改列的名称
ALTER TABLE 表名 CHANGE 列名 新列名 新数据类型;
# 修改列的数据类型
ALTER TABLE 表名 MODIFY 列名 新数据类型;
```

```sql
# 修改表名
alter table `moments` rename to `moment`;
# 添加一个新的列
alter table `moment` ADD `publishTime` DATETIME;
# 删除一列数据
alter table `moment` drop `publishTime`;
# 修改列的名称
alter table `moment` change `publishTime` `publishDate` DATE;
# 修改列的数据类型
alter table `moment` modify `id` INT;
```

## SQL操作：DML

### 表中插入数据

```sql
# 给指定列添加数据
INSERT INTO 表名(列名1,列名2,…) VALUES(值1,值2,…);
# 给全部列添加数据
INSERT INTO 表名 VALUES(值1,值2,…);
# 批量添加数据
INSERT INTO 表名(列名1,列名2,…) VALUES(值1,值2,…),(值1,值2,…),(值1,值2,…)…;
INSERT INTO 表名 VALUES(值1,值2,…),(值1,值2,…),(值1,值2,…)…;
```

```sql
insert into `products` (`title`, `description`, `price`, `publishTime`) VALUES ('iPhone', 'iPhone14主要798', 798, '2020-11-11')
```

### 表中更新操作，修改和删除

```sql
# 修改表数据
UPDATE 表名 SET 列名1=值1,列名2=值2,… [WHERE 条件] ;
# 删除数据
DELETE FROM 表名 [WHERE 条件] ;
```

```sql
# 删除数据
# 删除表中所有的数据
delete from `products`;
# 删除符合条件的数据
delete from `products` where `title` = 'iPhone';
# # 修改数据
# 会修改表中所有的数据
UPDATE `products` SET `title` = 'iPhone12', `price` = 1299.88;
# 会修改符合条件的数据
UPDATE `products` SET `title` = 'iPhone12', `price` = 1299.88 WHERE `title` = 'iPhone';
```

如果我们希望修改完数据后，直接可以显示最新的更新时间

```sql
alter table `products` ADD `updateTime` TIMESTAMP default current_timestamp on update current_timestamp;
```

## SQL操作：DQL

### 查询表中的数据

```sql
SELECT select_expr [, select_expr]...
    [FROM table_references]
    [WHERE where_condition]
    [ORDER BY expr [ASC | DESC]]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
    [GROUP BY expr]
    [HAVING where_condition]
```

查询所有的数据并且显示所有的字段

```sql
SELECT * FROM `products`;
```

查询title、brand、price

```sql
SELECT title, brand, price FROM `products`;
```

### 查询结果去重

1.使用DISTINCT关键字

```sql
SELECT DISTINCT brand FROM `products`;
```
2.使用GROUP BY关键字

```sql
SELECT brand FROM `products` GROUP BY brand;
```

### 条件查询

条件查询会使用 WEHRE查询子句

```sql
# 查询价格小于1000的手机
SELECT * FROM `products` WHERE price < 1000;
# 查询品牌是华为，并且小于2000元的手机
SELECT * FROM `products` WHERE `brand` = '华为' and `price` < 2000;
SELECT * FROM `products` WHERE `brand` = '华为' && `price` < 2000;
# 查询1000到2000的手机（不包含1000和2000）
SELECT * FROM `products` WHERE price > 1000 and price < 2000;
# 包含1000和2000
SELECT * FROM `products` WHERE price BETWEEN 1000 and 2000;
# OR: 符合一个条件即可
# 查询所有的华为手机或者价格小于1000的手机
SELECT * FROM `products` WHERE brand = '华为' or price < 1000;
# 查询1000到2000的手机（包含1000和2000）
SELECT * FROM `products` WHERE price BETWEEN 1000 and 2000;
# 查看多个结果中的一个
SELECT * FROM `products` WHERE brand in ('华为', '小米');
```

模糊查询使用LIKE关键字，结合两个特殊的符号：

- % 表示匹配任意个的任意字符
- _ 表示匹配一个的任意字符

```sql
# 查询所有以v开头的title
SELECT * FROM `products` WHERE title LIKE 'v%';
# 查询带M的title
SELECT * FROM `products` WHERE title LIKE '%M%';
# 查询带M的title必须是第三个字符
SELECT * FROM `products` WHERE title LIKE '__M%';
```

### 查询结果排序

当我们查询到结果的时候，我们希望讲结果按照某种方式进行排序，这个时候使用的是ORDER BY；

ORDER BY有两个常用的值：

- ASC：升序排列；
- DESC：降序排列；

```sql
SELECT * FROM `products` WHERE brand = '华为' or price < 1000 ORDER BY price ASC;
```

### 分页查询

当数据库中的数据非常多时，一次性查询到所有的结果进行显示是不太现实的，我们可以在数据库中进行分页查询

用法：[LIMIT {[offset,] row_count | row_count OFFSET offset}]

```sql
SELECT * FROM `products` LIMIT 30 OFFSET 0;
SELECT * FROM `products` LIMIT 30 OFFSET 30;
SELECT * FROM `products` LIMIT 30 OFFSET 60;
# 另外一种写法：offset, row_count
SELECT * FROM `products` LIMIT 90, 30;
```

### 聚合函数

聚合函数表示对值集合进行操作的组（集合）函数

```sql
# 华为手机价格的平均值
SELECT AVG(price) FROM `products` WHERE brand = '华为';
# 计算所有手机的平均分
SELECT AVG(score) FROM `products`;
# 手机中最低和最高分数
SELECT MAX(score) FROM `products`;
SELECT MIN(score) FROM `products`;
# 计算总投票人数
SELECT SUM(voteCnt) FROM `products`;
# 计算所有条目的数量
SELECT COUNT(*) FROM `products`;
# 华为手机的个数
SELECT COUNT(*) FROM `products` WHERE brand = '华为';
```

### Group By 分组查询

```sql
# 根据品牌进行分组；
# 计算各个品牌中：商品的个数、平均价格
# 也包括：最高价格、最低价格、平均评分
SELECT brand, 
    COUNT(*) as count, 
    ROUND(AVG(price),2) as avgPrice,
    MAX(price) as maxPrice,
    MIN(price) as minPrice,
    AVG(score) as avgScore
FROM `products` GROUP BY brand;
```

给Group By查询到的结果添加一些约束，那么我们可以使用：HAVING。

```sql
# 分组并筛选出平均价格在4000以下，并且平均分在7以上的品牌
SELECT brand, 
    COUNT(*) as count, 
    ROUND(AVG(price),2) as avgPrice,
    MAX(price) as maxPrice,
    MIN(price) as minPrice,
    AVG(score) as avgScore
FROM `products` GROUP BY brand 
HAVING avgPrice < 4000 and avgScore > 7;
```

