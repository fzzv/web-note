# MYSQL 高级操作

## 约束

![image-20210724104749122](MYSQL%E6%93%8D%E4%BD%9C%E9%AB%98%E7%BA%A7.assets/image-20210724104749122.png)

上面表中可以看到表中数据存在一些问题：

* id 列一般是用标示数据的唯一性的，而上述表中的id为1的有三条数据，并且 `马花疼` 没有id进行标示
* `柳白` 这条数据的age列的数据是3000，而人也不可能活到3000岁
* `马运`  这条数据的math数学成绩是-5，而数学学得再不好也不可能出现负分
* `柳青` 这条数据的english列（英文成绩）值为null，而成绩即使没考也得是0分

针对上述数据问题，我们就可以从数据库层面在添加数据的时候进行限制，这个就是约束。

### 概念

* 约束是作用于表中列上的规则，用于限制加入表的数据

  例如：我们可以给id列加约束，让其值不能重复，不能为null值。

* 约束的存在保证了数据库中数据的正确性、有效性和完整性

  添加约束可以在添加数据的时候就限制不正确的数据，年龄是3000，数学成绩是-5分这样无效的数据，继而保障数据的完整性。

### 分类

* **非空约束： 关键字是 NOT NULL**

  保证列中所有的数据不能有null值。

  例如：id列在添加 `马花疼` 这条数据时就不能添加成功。

* **唯一约束：关键字是  UNIQUE**

  保证列中所有数据各不相同。

  例如：id列中三条数据的值都是1，这样的数据在添加时是绝对不允许的。

* **主键约束： 关键字是  PRIMARY KEY**

  主键是一行数据的唯一标识，要求非空且唯一。一般我们都会给没张表添加一个主键列用来唯一标识数据。

  例如：上图表中id就可以作为主键，来标识每条数据。那么这样就要求数据中id的值不能重复，不能为null值。

* **检查约束： 关键字是  CHECK** 

  保证列中的值满足某一条件。

  例如：我们可以给age列添加一个范围，最低年龄可以设置为1，最大年龄就可以设置为300，这样的数据才更合理些。

  > 注意：MySQL不支持检查约束。
  >
  > 这样是不是就没办法保证年龄在指定的范围内了？从数据库层面不能保证，以后可以在java代码中进行限制，一样也可以实现要求。

* **默认约束： 关键字是   DEFAULT**

  保存数据时，未指定值则采用默认值。

  例如：我们在给english列添加该约束，指定默认值是0，这样在添加数据时没有指定具体值时就会采用默认给定的0。

* **外键约束： 关键字是  FOREIGN KEY**

  外键用来让两个表的数据之间建立链接，保证数据的一致性和完整性。

  外键约束现在可能还不太好理解，后面我们会重点进行讲解。

### 非空约束

非空约束用于保证列中所有数据不能有NULL值

```sql
-- 创建表时添加非空约束
CREATE TABLE 表名(
   列名 数据类型 NOT NULL,
   …
); 

-- 建完表后添加非空约束
ALTER TABLE 表名 MODIFY 字段名 数据类型 NOT NULL;

-- 删除约束
ALTER TABLE 表名 MODIFY 字段名 数据类型;
```

### 唯一约束

唯一约束用于保证列中所有数据各不相同

```sql
-- 创建表时添加唯一约束
CREATE TABLE 表名(
   列名 数据类型 UNIQUE [AUTO_INCREMENT],
   -- AUTO_INCREMENT: 当不指定值时自动增长
   …
); 
CREATE TABLE 表名(
   列名 数据类型,
   …
   [CONSTRAINT] [约束名称] UNIQUE(列名)
); 

-- 建完表后添加唯一约束
ALTER TABLE 表名 MODIFY 字段名 数据类型 UNIQUE;

-- 删除约束
ALTER TABLE 表名 DROP INDEX 字段名;
```

### 主键约束

主键是一行数据的唯一标识，要求非空且唯一。一张表只能有一个主键

```sql
-- 创建表时添加主键约束
CREATE TABLE 表名(
   列名 数据类型 PRIMARY KEY [AUTO_INCREMENT],
   …
); 
CREATE TABLE 表名(
   列名 数据类型,
   [CONSTRAINT] [约束名称] PRIMARY KEY(列名)
); 

-- 建完表后添加主键约束
ALTER TABLE 表名 ADD PRIMARY KEY(字段名);

-- 删除约束
ALTER TABLE 表名 DROP PRIMARY KEY;
```

### 默认约束

保存数据时，未指定值则采用默认值

```sql
-- 创建表时添加默认约束
CREATE TABLE 表名(
   列名 数据类型 DEFAULT 默认值,
   …
); 

-- 建完表后添加默认约束
ALTER TABLE 表名 ALTER 列名 SET DEFAULT 默认值;

-- 删除约束
ALTER TABLE 表名 ALTER 列名 DROP DEFAULT;
```

### 外键约束

外键用来让两个表的数据之间建立链接，保证数据的一致性和完整性。

```sql
-- 创建表时添加外键约束
CREATE TABLE 表名(   
    列名 数据类型,
    …   
    [CONSTRAINT] [外键名称] FOREIGN KEY(外键列名) REFERENCES 主表(主表列名)
); 

-- 建完表后添加外键约束
ALTER TABLE 表名 ADD CONSTRAINT 外键名称 FOREIGN KEY (外键字段名称) REFERENCES 主表名称(主表列名称);

-- 删除约束
ALTER TABLE 表名 DROP FOREIGN KEY 外键名称;
```

将两张表联系起来，我们可以将products中的brand_id关联到brand中的id。

如果是创建表添加外键约束，我们需要在创建表的()最后添加如下语句

```sql
foreing key (brand_id) peferences brand(id)
```

如果是表已经创建好，额外添加外键

```sql
alter table `products` ADD foreing key (brand_id) references brand(id);
```

将products中的brand_id关联到brand中的id的值

```sql
UPDATE `products` SET `brand_id` = 1 WHERE `brand` = '华为';
UPDATE `products` SET `brand_id` = 4 WHERE `brand` = 'OPPO';
UPDATE `products` SET `brand_id` = 3 WHERE `brand` = '苹果';
UPDATE `products` SET `brand_id` = 2 WHERE `brand` = '小米';
```

## 多表查询

查询多张表中的的数据

```sql
select * from `products`, `brand`;
```

多表的查询结果：第一张表中每一个条数据，都会和第二张表中的每一条数据结合一次，（**X  *Y*）。

使用 where 进行筛选，符合products.brand_id = brand.id条件的数据过滤出来

```sql
SELECT * FROM `products`, `brand` WHERE `products`.brand_id = `brand`.id;
```

### 连接查询

- 外连接查询：
  - 左连接：LEFT [OUTER] JOIN
  - 右连接：RIGHT [OUTER] JOIN
- 内连接：CROSS JOIN或者 JOIN都可以，相当于查询AB两个表的交集数据
- 全连接：
  - UNION，两张表上下合并，就是表头一样
  - UNION ALL，合并但保留重复
  - FULL JOIN 两边全都保留（MySQL 不支持 FULL OUTER JOIN）

#### 内连接查询

```sql
-- 隐式内连接
SELECT 字段列表 FROM 表1,表2… WHERE 条件;

-- 显式内连接
SELECT 字段列表 FROM 表1 [INNER] JOIN 表2 ON 条件;
```

> 内连接相当于查询 A B 交集数据

![image-20210724174717647](MYSQL%E6%93%8D%E4%BD%9C%E9%AB%98%E7%BA%A7.assets/image-20210724174717647.png)

```sql
-- 隐式内连接
SELECT * FROM emp, dept WHERE emp.dep_id = dept.did;

-- 显式内连接
select * from emp inner join dept on emp.dep_id = dept.did;
-- 上面语句中的inner可以省略，可以书写为如下语句
select * from emp  join dept on emp.dep_id = dept.did;
```

#### 外连接查询

```sql
-- 左外连接
SELECT 字段列表 FROM 表1 LEFT [OUTER] JOIN 表2 ON 条件;

-- 右外连接
SELECT 字段列表 FROM 表1 RIGHT [OUTER] JOIN 表2 ON 条件;
```



![image-20210724174717647](MYSQL%E6%93%8D%E4%BD%9C%E9%AB%98%E7%BA%A7.assets/image-20210724174717647-1785997870845.png)

> 左外连接：相当于查询A表所有数据和交集部分数据
>
> 右外连接：相当于查询B表所有数据和交集部分数据

查询emp表所有数据和对应的部门信息（左外连接）

```sql
select * from emp left join dept on emp.dep_id = dept.did;
```

查询dept表所有数据和对应的员工信息（右外连接）

```sql
select * from emp right join dept on emp.dep_id = dept.did;
```

### 子查询

查询中嵌套查询，称嵌套查询为子查询。

* 什么是查询中嵌套查询呢？我们通过一个例子来看：

  **需求：查询工资高于猪八戒的员工信息。**

  来实现这个需求，我们就可以通过二步实现，第一步：先查询出来 猪八戒的工资

  ```sql
  select salary from emp where name = '猪八戒'
  ```

   第二步：查询工资高于猪八戒的员工信息

  ```sql
  select * from emp where salary > 3600;
  ```

  第二步中的3600可以通过第一步的sql查询出来，所以将3600用第一步的sql语句进行替换

  ```sql
  select * from emp where salary > (select salary from emp where name = '猪八戒');
  ```

  这就是查询语句中嵌套查询语句。

* 子查询根据查询结果不同，作用不同

  * 子查询语句结果是单行单列，子查询语句作为条件值，使用 =  !=  >  <  等进行条件判断
  * 子查询语句结果是多行单列，子查询语句作为条件值，使用 in 等关键字进行条件判断
  * 子查询语句结果是多行多列，子查询语句作为虚拟表

* 案例

  * 查询 '财务部' 和 '市场部' 所有的员工信息

    ```sql
    -- 查询 '财务部' 或者 '市场部' 所有的员工的部门did
    select did from dept where dname = '财务部' or dname = '市场部';
    
    select * from emp where dep_id in (select did from dept where dname = '财务部' or dname = '市场部');
    ```

  * 查询入职日期是 '2011-11-11' 之后的员工信息和部门信息

    ```sql
    -- 查询入职日期是 '2011-11-11' 之后的员工信息
    select * from emp where join_date > '2011-11-11' ;
    -- 将上面语句的结果作为虚拟表和dept表进行内连接查询
    select * from (select * from emp where join_date > '2011-11-11' ) t1, dept where t1.dep_id = dept.did;
    ```

