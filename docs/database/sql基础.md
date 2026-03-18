# SQL基础

## MySQL 核心语法速查表

### 1. 基础操作 (CRUD) - 数据的增删改查

*记住：做修改和删除时，永远不要忘记 WHERE！*

| **操作** | **关键字** | **语法模板**                                     | **示例**                                             |
| -------- | ---------- | ------------------------------------------------ | ---------------------------------------------------- |
| **查询** | `SELECT`   | `SELECT 列名 FROM 表名;`                         | `SELECT name FROM students;`                         |
| **新增** | `INSERT`   | `INSERT INTO 表名 (列1, 列2) VALUES (值1, 值2);` | `INSERT INTO users (name, age) VALUES ('Jack', 20);` |
| **修改** | `UPDATE`   | `UPDATE 表名 SET 列=新值 WHERE 条件;`            | `UPDATE users SET age=21 WHERE name='Jack';`         |
| **删除** | `DELETE`   | `DELETE FROM 表名 WHERE 条件;`                   | `DELETE FROM users WHERE id=5;`                      |

### 2. 筛选与排序 - 精准定位

*就像网购时的筛选器和“按价格排序”。*

- **条件筛选 (`WHERE`)**:
  - `WHERE age > 18` (大于)
  - `WHERE name = 'Tom'` (等于文本)
  - `WHERE age > 18 AND gender = 'male'` (且 - 同时满足)
  - `WHERE color = 'red' OR color = 'blue'` (或 - 满足其一)
- **排序 (`ORDER BY`)**:
  - `ORDER BY price DESC` (降序：从大到小 ⬇️)
  - `ORDER BY price ASC` (升序：从小到大 ⬆️ - 默认)
- **限制数量 (`LIMIT`)**:
  - `LIMIT 3` (只看前 3 条)

### 3. 统计与分组 - 数据分析神器

*把数据聚在一起看宏观趋势。*

- **聚合函数**:

  - `COUNT(*)`: 数行数
  - `SUM(列)`: 求和
  - `AVG(列)`: 平均值
  - `MAX(列)` / `MIN(列)`: 最大/最小值

- **分组流程**:

  ```sql
  SELECT 类别, AVG(价格)
  FROM 商品表
  WHERE 价格 > 10       -- 1. 先筛选行 (原始数据)
  GROUP BY 类别         -- 2. 再按类别分组
  HAVING AVG(价格) < 50 -- 3. 最后筛选组 (计算后的结果)
  ```

### 4. 多表连接 (JOINS) - 打破孤岛

*把两张表拼成一张大表。*

- **`INNER JOIN` (内连接)**:
  - **口诀**: “强强联手，没对象的都不要。”
  - **效果**: 只保留两边都有数据的行（交集）。
  - **语法**: `FROM 表A JOIN 表B ON 表A.id = 表B.a_id`
- **`LEFT JOIN` (左连接)**:
  - **口诀**: “左边是老大，右边即使为空也要跟着。”
  - **效果**: 保留左表所有数据，右表没对应数据时显示 NULL。
  - **语法**: `FROM 表A LEFT JOIN 表B ON ...`

### 5. 维护与优化 - 数据库管家

- **外键 (`Foreign Key`)**:
  - **作用**: 防止脏数据录入（比如输入不存在的商品ID）。
  - **语法**: `ALTER TABLE 表名 ADD CONSTRAINT ... FOREIGN KEY ...`
- **索引 (`Index`)**:
  - **作用**: 给书加目录，让查询速度起飞 。
  - **语法**: `CREATE INDEX 索引名 ON 表名(列名);`