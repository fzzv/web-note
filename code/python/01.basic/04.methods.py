import math
import keyword
import copy
from functools import reduce
from collections import namedtuple, defaultdict, Counter, OrderedDict

# int 类型的方法
# bit_length() —— 二进制位数（不含符号位和前导零）
print((0).bit_length())  # 0
print((1).bit_length())  # 1
print((255).bit_length())  # 8
print((-255).bit_length())  # 8

# bit_count() —— 二进制中 1 的个数
print((255).bit_count())
print((10).bit_count())

# to_bytes() —— 转换为字节序列
print((1024).to_bytes(2, byteorder="big"))  # b'\x04\x00'
print((1024).to_bytes(2, byteorder="little"))  # b'\x00\x04'

# from_bytes() —— 从字节序列创建整数（类方法）
print(int.from_bytes(b"\x04\x00", byteorder="big"))  # 1024

# as_integer_ratio() —— 返回分数表示
print((3).as_integer_ratio())  # (3, 1)

# float 的方法
f = 3.75

# as_integer_ratio() —— 返回精确的分数表示
print(f.as_integer_ratio())  # (15, 4)  即 15/4 = 3.75
print((0.1).as_integer_ratio())  # (3602879701896397, 36028797018963968)

# is_integer() —— 是否为整数值
print((3.0).is_integer())  # True
print((3.5).is_integer())  # False

# hex() —— 十六进制浮点表示
print((3.75).hex())  # 0x1.e000000000000p+1

# fromhex() —— 从十六进制字符串创建浮点数（类方法）
print(float.fromhex("0x1.ep+1"))  # 3.75

# abs() —— 绝对值
print(abs(-7))  # 7
print(abs(-3.14))  # 3.14
print(abs(3 + 4j))  # 5.0  （复数的模）

# round() —— 四舍五入（银行家舍入法）
print(round(3.5))  # 4
print(round(4.5))  # 4   （银行家舍入：.5 时取偶数）
print(round(3.14159, 2))  # 3.14
print(round(1234, -2))  # 1200  （负数表示小数点左边）

# pow() —— 幂运算（可选取模）
print(pow(2, 10))  # 1024
print(pow(2, 10, 100))  # 24   （等价于 2**10 % 100，但更高效）

# divmod() —— 同时返回商和余数
print(divmod(17, 5))  # (3, 2)
print(divmod(10, 3))  # (3, 1)

# max() / min()
print(max(3, 1, 4, 1, 5))  # 5
print(min(3, 1, 4, 1, 5))  # 1
print(max([3, 1, 4], key=lambda x: -x))  # 1  （按自定义规则）

# sum()
print(sum([1, 2, 3, 4, 5]))  # 15
print(sum([1, 2, 3], 10))  # 16  （初始值为 10）
print(sum(range(1, 101)))  # 5050

# 进制转换
print(bin(10))  # 0b1010
print(oct(10))  # 0o12
print(hex(255))  # 0xff
print(int("ff", 16))  # 255

# chr() / ord() —— 字符与 Unicode 码点互转
print(chr(65))  # A
print(chr(20013))  # 中
print(ord("A"))  # 65
print(ord("中"))  # 20013

# 取整
print(math.ceil(3.2))  # 4    向上取整
print(math.floor(3.8))  # 3    向下取整
print(math.trunc(3.7))  # 3    截断小数（向零取整）
print(math.trunc(-3.7))  # -3

# 常量
print(math.pi)  # 3.141592653589793
print(math.e)  # 2.718281828459045
print(math.inf)  # inf
print(math.nan)  # nan

# 常用数学函数
print(math.sqrt(16))  # 4.0       平方根
print(math.log(100, 10))  # 2.0       对数
print(math.log2(8))  # 3.0
print(math.log10(1000))  # 3.0
print(math.factorial(5))  # 120       阶乘
print(math.gcd(12, 8))  # 4         最大公约数
print(math.lcm(4, 6))  # 12        最小公倍数（Python 3.9+）

# 三角函数
print(math.sin(math.pi / 2))  # 1.0
print(math.cos(0))  # 1.0
print(math.degrees(math.pi))  # 180.0   弧度转角度
print(math.radians(180))  # 3.14... 角度转弧度

# 浮点数判断
print(math.isnan(math.nan))  # True
print(math.isinf(math.inf))  # True
print(math.isclose(0.1 + 0.2, 0.3))  # True  （解决浮点精度问题）

# 累积运算
print(math.prod([1, 2, 3, 4, 5]))  # 120  连乘（Python 3.8+）
print(math.fsum([0.1] * 10))  # 1.0  精确浮点求和

# str 的方法
s = "hello, World! PYTHON"

print(s.upper())  # HELLO, WORLD! PYTHON
print(s.lower())  # hello, world! python
print(s.title())  # Hello, World! Python    （每个单词首字母大写）
print(s.capitalize())  # Hello, world! python    （仅首字母大写）
print(s.swapcase())  # HELLO, wORLD! python    （大小写互换）
print(s.casefold())  # hello, world! python    （更激进的小写，支持特殊字符）

# casefold vs lower 的区别
print("straße".lower())  # straße    （德语 sharp s 不变）
print("straße".casefold())  # strasse   （转为 ss）

s = "Hello, Python! Hello, World!"

# find() / rfind() —— 查找子串位置，未找到返回 -1
print(s.find("Hello"))  # 0
print(s.find("Hello", 1))  # 15   （从索引 1 开始找）
print(s.rfind("Hello"))  # 15   （从右往左找）
print(s.find("Java"))  # -1

# index() / rindex() —— 同 find，但未找到时抛出 ValueError
print(s.index("Hello"))  # 0
# s.index("Java")             # ValueError: substring not found

# count() —— 统计子串出现次数
print(s.count("Hello"))  # 2
print(s.count("l"))  # 4

# replace() —— 替换
print(s.replace("Hello", "Hi"))  # Hi, Python! Hi, World!
print(s.replace("Hello", "Hi", 1))  # Hi, Python! Hello, World!  （只替换 1 次）

# removeprefix() / removesuffix()（Python 3.9+）
filename = "test_data.csv"
print(filename.removeprefix("test_"))  # data.csv
print(filename.removesuffix(".csv"))  # test_data
print(filename.removeprefix("xyz"))  # test_data.csv  （没有前缀则不变）

# split() —— 按分隔符分割为列表
print("a,b,c,d".split(","))  # ['a', 'b', 'c', 'd']
print("a,b,c,d".split(",", 2))  # ['a', 'b', 'c,d']  （最多分割 2 次）
print(
    "  hello  world  ".split()
)  # ['hello', 'world']  （默认按空白分割，自动去除首尾空白）
print(
    "  hello  world  ".split(" ")
)  # ['', '', 'hello', '', 'world', '', '']  （严格按空格分割）

# rsplit() —— 从右侧开始分割
print("a.b.c.d".rsplit(".", 2))  # ['a.b', 'c', 'd']

# splitlines() —— 按行分割
text = "第一行\n第二行\r\n第三行"
print(text.splitlines())  # ['第一行', '第二行', '第三行']
print(text.splitlines(keepends=True))  # ['第一行\n', '第二行\r\n', '第三行']

# partition() / rpartition() —— 分成三部分（前、分隔符、后）
print("user@example.com".partition("@"))  # ('user', '@', 'example.com')
print("user@example.com".rpartition("@"))  # ('user', '@', 'example.com')
print("hello".partition("@"))  # ('hello', '', '')  （未找到分隔符）

# join() —— 拼接
print(",".join(["a", "b", "c"]))  # a,b,c
print(" -> ".join(["A", "B", "C"]))  # A -> B -> C
print("".join(["H", "e", "l", "l", "o"]))  # Hello

# 注意：join 要求元素都是字符串
nums = [1, 2, 3]
print(",".join(str(n) for n in nums))  # 1,2,3

s = "  Hello, World!  "

# strip() / lstrip() / rstrip() —— 去除首尾字符
print(s.strip())  # "Hello, World!"
print(s.lstrip())  # "Hello, World!  "
print(s.rstrip())  # "  Hello, World!"

# 可以指定要去除的字符集合
print("###Hello###".strip("#"))  # Hello
print("abcHelloabc".strip("abc"))  # Hello  （去除 a、b、c 任意组合）

# center() / ljust() / rjust() —— 填充对齐
print("hello".center(20))  # "       hello        "
print("hello".center(20, "-"))  # "-------hello--------"
print("hello".ljust(20, "."))  # "hello..............."
print("hello".rjust(20, "."))  # "...............hello"

# zfill() —— 零填充（保留正负号）
print("42".zfill(5))  # 00042
print("-42".zfill(5))  # -0042
print("3.14".zfill(7))  # 003.14

# expandtabs() —— 制表符替换为空格
print("a\tb\tc".expandtabs(4))  # a   b   c

# 内容判断
print("abc123".isalnum())  # True   字母或数字
print("abc".isalpha())  # True   纯字母
print("123".isdigit())  # True   纯数字（不含负号和小数点）
print("123".isnumeric())  # True   数字字符（含中文数字等）
print("三".isnumeric())  # True
print("123".isdecimal())  # True   十进制数字字符

# 大小写判断
print("hello".islower())  # True
print("HELLO".isupper())  # True
print("Hello World".istitle())  # True  （每个单词首字母大写）

# 空白判断
print("   ".isspace())  # True
print("\t\n".isspace())  # True
print("".isspace())  # False  （空字符串不是空白）

# 前缀和后缀
print("hello.py".startswith("hello"))  # True
print("hello.py".endswith(".py"))  # True
print("hello.py".startswith(("hello", "hi")))  # True  （传入元组，匹配任意一个）
print("hello.py".endswith((".py", ".js")))  # True

# 是否可作标识符
print("my_var".isidentifier())  # True
print("2name".isidentifier())  # False
print("class".isidentifier())  # True  （关键字也是合法标识符）

print(keyword.iskeyword("class"))  # True  （判断是否为关键字）

# 是否可打印
print("Hello".isprintable())  # True
print("Hello\n".isprintable())  # False  （\n 不可打印）

# isascii()（Python 3.7+）
print("Hello".isascii())  # True
print("你好".isascii())  # False

# bool() —— 转换为布尔值
print(bool(0))  # False
print(bool(1))  # True
print(bool(""))  # False
print(bool("hello"))  # True
print(bool([]))  # False
print(bool([0]))  # True  （非空列表，即使元素是 0）
print(bool(None))  # False

# 由于 bool 是 int 的子类，拥有 int 的方法
print(True.bit_length())  # 1
print(False.bit_length())  # 0
print(True.as_integer_ratio())  # (1, 1)

# 布尔值参与运算
print(True + True)  # 2
print(True * 10)  # 10
print(sum([True, False, True, True]))  # 3  （统计 True 的个数）

# 实用技巧：统计满足条件的元素数量
numbers = [1, -2, 3, -4, 5, -6]
positive_count = sum(1 for n in numbers if n > 0)
# 或者
positive_count = sum(n > 0 for n in numbers)
print(positive_count)  # 3

# any() / all() 与布尔值
print(any([False, False, True]))  # True  （有一个为 True）
print(all([True, True, True]))  # True  （全部为 True）
print(all([True, True, False]))  # False
print(any([]))  # False  （空可迭代对象）
print(all([]))  # True   （空可迭代对象，"空真"）

lst = [1, 2, 3]

# append() —— 末尾追加单个元素
lst.append(4)
print(lst)  # [1, 2, 3, 4]

lst.append([5, 6])  # type: ignore # 作为一个整体追加
print(lst)  # [1, 2, 3, 4, [5, 6]]

# insert() —— 在指定位置插入
lst = [1, 2, 3]
lst.insert(0, 0)  # 在索引 0 处插入
print(lst)  # [0, 1, 2, 3]

lst.insert(2, 1.5)  # type: ignore # 在索引 2 处插入
print(lst)  # [0, 1, 1.5, 2, 3]

lst.insert(100, 99)  # 索引超出范围，追加到末尾
print(lst)  # [0, 1, 1.5, 2, 3, 99]

# extend() —— 逐个追加可迭代对象中的元素
lst = [1, 2, 3]
lst.extend([4, 5, 6])
print(lst)  # [1, 2, 3, 4, 5, 6]

lst.extend("ab")  # type: ignore
print(lst)  # [1, 2, 3, 4, 5, 6, 'a', 'b']

# += 等价于 extend
lst = [1, 2]
lst += [3, 4]
print(lst)  # [1, 2, 3, 4]

lst = ["a", "b", "c", "b", "d"]

# remove() —— 删除第一个匹配的值
lst.remove("b")
print(lst)  # ['a', 'c', 'b', 'd']
# lst.remove("z")  # ValueError: list.remove(x): x not in list

# pop() —— 弹出指定索引的元素，默认末尾
lst = [10, 20, 30, 40, 50]
print(lst.pop())  # 50       弹出末尾
print(lst)  # [10, 20, 30, 40]
print(lst.pop(1))  # 20       弹出索引 1
print(lst)  # [10, 30, 40]

# clear() —— 清空列表
lst = [1, 2, 3]
lst.clear()
print(lst)  # []

# del —— 删除元素或切片（语句，非方法）
lst = [0, 1, 2, 3, 4, 5]
del lst[0]
print(lst)  # [1, 2, 3, 4, 5]
del lst[1:3]
print(lst)  # [1, 4, 5]

# sort() —— 原地排序，返回 None
nums = [3, 1, 4, 1, 5, 9, 2, 6]
nums.sort()
print(nums)  # [1, 1, 2, 3, 4, 5, 6, 9]

nums.sort(reverse=True)
print(nums)  # [9, 6, 5, 4, 3, 2, 1, 1]

# 自定义排序规则
words = ["banana", "apple", "cherry", "date"]
words.sort(key=len)
print(words)  # ['date', 'apple', 'banana', 'cherry']

words.sort(key=str.lower)  # 忽略大小写排序
print(words)  # ['apple', 'banana', 'cherry', 'date']

# 多条件排序
students = [("Alice", 90), ("Bob", 85), ("Charlie", 90), ("David", 85)]
students.sort(key=lambda x: (-x[1], x[0]))  # 按成绩降序，同分按姓名升序
print(students)
# [('Alice', 90), ('Charlie', 90), ('Bob', 85), ('David', 85)]

# sorted() —— 返回新列表，不修改原列表
original = [3, 1, 2]
new_list = sorted(original)
print(original)  # [3, 1, 2]  （原列表不变）
print(new_list)  # [1, 2, 3]

# reverse() —— 原地反转
lst = [1, 2, 3, 4, 5]
lst.reverse()
print(lst)  # [5, 4, 3, 2, 1]

# reversed() —— 返回反转迭代器，不修改原列表
lst = [1, 2, 3]
print(list(reversed(lst)))  # [3, 2, 1]
print(lst)  # [1, 2, 3]  （原列表不变）

lst = [10, 20, 30, 20, 40, 20]

# index() —— 返回第一个匹配的索引
print(lst.index(20))  # 1
print(lst.index(20, 2))  # 3   （从索引 2 开始查找）
print(lst.index(20, 4))  # 5   （从索引 4 开始查找）
# lst.index(99)             # ValueError: 99 is not in list

# count() —— 统计出现次数
print(lst.count(20))  # 3
print(lst.count(99))  # 0

# in —— 成员判断（运算符，非方法）
print(20 in lst)  # True
print(99 in lst)  # False

# copy() —— 浅拷贝
lst = [1, [2, 3], 4]
copy1 = lst.copy()
copy2 = lst[:]  # 切片也是浅拷贝
copy3 = list(lst)  # 构造函数也是浅拷贝

lst[0] = 100
print(copy1)  # [1, [2, 3], 4]  （第一层不受影响）

lst[1].append(99)  # type: ignore
print(copy1)  # [1, [2, 3, 99], 4]  （嵌套对象被共享，受影响）

# 深拷贝
lst = [1, [2, 3], 4]
deep = copy.deepcopy(lst)
lst[1].append(99)  # type: ignore
print(deep)  # [1, [2, 3], 4]  （完全独立）

lst = [3, 1, 4, 1, 5, 9]

print(len(lst))  # 6
print(max(lst))  # 9
print(min(lst))  # 1
print(sum(lst))  # 23

# enumerate() —— 返回 (索引, 元素) 对
for i, v in enumerate(lst):
    print(f"[{i}] = {v}", end="  ")
print()
# [0] = 3  [1] = 1  [2] = 4  [3] = 1  [4] = 5  [5] = 9

# zip() —— 并行遍历
names = ["Alice", "Bob"]
scores = [90, 85]
print(list(zip(names, scores)))  # [('Alice', 90), ('Bob', 85)]

# map() —— 对每个元素应用函数
print(list(map(str, [1, 2, 3])))  # ['1', '2', '3']
print(list(map(lambda x: x**2, [1, 2, 3])))  # [1, 4, 9]

# filter() —— 过滤元素
print(list(filter(lambda x: x > 3, lst)))  # [4, 5, 9]

# reduce() —— 累积运算
print(reduce(lambda a, b: a + b, [1, 2, 3, 4]))  # 10  (((1+2)+3)+4)
print(reduce(lambda a, b: a * b, [1, 2, 3, 4]))  # 24

# 元组 是不可变序列，只有两个方法。
t = (1, 2, 3, 2, 4, 2, 5)

# count() —— 统计出现次数
print(t.count(2))  # 3
print(t.count(99))  # 0

# index() —— 返回第一个匹配的索引
print(t.index(2))  # 1
print(t.index(2, 2))  # 3   （从索引 2 开始查找）
print(t.index(2, 4))  # 5   （从索引 4 开始查找）
# t.index(99)            # ValueError: tuple.index(x): x not in tuple

# 虽然方法只有两个，但元组支持很多通用操作
t = (1, 2, 3, 4, 5)

# 索引和切片
print(t[0])  # 1
print(t[-1])  # 5
print(t[1:3])  # (2, 3)
print(t[::-1])  # (5, 4, 3, 2, 1)

# 拼接和重复
print((1, 2) + (3, 4))  # (1, 2, 3, 4)
print((1, 2) * 3)  # (1, 2, 1, 2, 1, 2)

# 拆包
a, b, c = (10, 20, 30)
print(a, b, c)  # 10 20 30

first, *middle, last = (1, 2, 3, 4, 5)
print(first)  # 1
print(middle)  # [2, 3, 4]
print(last)  # 5

# 成员判断
print(3 in t)  # True

# 内置函数
print(len(t))  # 5
print(max(t))  # 5
print(min(t))  # 1
print(sum(t))  # 15
print(sorted(t, reverse=True))  # [5, 4, 3, 2, 1]  （返回列表）

# 命名元组 —— 给元组的字段命名
Point = namedtuple("Point", ["x", "y"])
p = Point(3, 4)
print(p)  # Point(x=3, y=4)
print(p.x)  # 3
print(p.y)  # 4
print(p[0])  # 3  （仍然可以用索引访问）

# 转换
print(p._asdict())  # {'x': 3, 'y': 4}
p2 = p._replace(x=10)
print(p2)  # Point(x=10, y=4)
print(Point._make([5, 6]))  # Point(x=5, y=6)

# 字典
d = {"name": "Alice", "age": 25, "city": "Beijing"}

# get() —— 安全获取，不存在时返回默认值
print(d.get("name"))  # Alice
print(d.get("email"))  # None
print(d.get("email", "未设置"))  # 未设置

# [] —— 直接访问，不存在时抛出 KeyError
print(d["name"])  # Alice
# print(d["email"])  # KeyError: 'email'

# keys() / values() / items() —— 返回视图对象
print(list(d.keys()))  # ['name', 'age', 'city']
print(list(d.values()))  # ['Alice', 25, 'Beijing']
print(list(d.items()))  # [('name', 'Alice'), ('age', 25), ('city', 'Beijing')]

# 视图对象是动态的，反映字典的变化
keys = d.keys()
d["email"] = "a@test.com"
print(list(keys))  # ['name', 'age', 'city', 'email']  （自动更新）

d = {"a": 1, "b": 2}

# 直接赋值 —— 存在则修改，不存在则添加
d["a"] = 10
d["c"] = 3
print(d)  # {'a': 10, 'b': 2, 'c': 3}

# update() —— 批量更新
d.update({"b": 20, "d": 4})
print(d)  # {'a': 10, 'b': 20, 'c': 3, 'd': 4}

d.update(e=5, f=6)  # 也可以用关键字参数
print(d)  # {'a': 10, 'b': 20, 'c': 3, 'd': 4, 'e': 5, 'f': 6}

# |= 运算符（Python 3.9+）
d |= {"g": 7}
print(d)  # {'a': 10, 'b': 20, 'c': 3, 'd': 4, 'e': 5, 'f': 6, 'g': 7}

# setdefault() —— 键存在则返回对应值，不存在则设置默认值并返回
d = {"a": 1, "b": 2}
print(d.setdefault("a", 99))  # 1   （已存在，不修改，返回已有值）
print(d.setdefault("c", 3))  # 3   （不存在，设置并返回）
print(d)  # {'a': 1, 'b': 2, 'c': 3}

# setdefault 的常见用途：给列表型的值追加元素
groups = {}
for name, dept in [("Alice", "IT"), ("Bob", "HR"), ("Charlie", "IT")]:
    groups.setdefault(dept, []).append(name)
print(groups)  # {'IT': ['Alice', 'Charlie'], 'HR': ['Bob']}

d = {"a": 1, "b": 2, "c": 3, "d": 4}

# pop() —— 删除指定键并返回值
print(d.pop("a"))  # 1
print(d)  # {'b': 2, 'c': 3, 'd': 4}
print(d.pop("z", "默认"))  # 默认  （键不存在，返回默认值）
# d.pop("z")                # KeyError: 'z'  （无默认值时报错）

# popitem() —— 删除并返回最后一个键值对（Python 3.7+ 按插入顺序）
print(d.popitem())  # ('d', 4)
print(d)  # {'b': 2, 'c': 3}

# del —— 删除指定键
del d["b"]
print(d)  # {'c': 3}

# clear() —— 清空字典
d.clear()
print(d)  # {}

# copy() —— 浅拷贝
d = {"a": 1, "b": [2, 3]}
d2 = d.copy()
d["a"] = 100
print(d2)  # {'a': 1, 'b': [2, 3]}  （第一层不受影响）

d["b"].append(4)
print(d2)  # {'a': 1, 'b': [2, 3, 4]}  （嵌套对象被共享）

# fromkeys() —— 用键序列创建字典（类方法）
keys = ["x", "y", "z"]
d = dict.fromkeys(keys)
print(d)  # {'x': None, 'y': None, 'z': None}

d = dict.fromkeys(keys, 0)
print(d)  # {'x': 0, 'y': 0, 'z': 0}

# 注意：默认值是共享的，可变对象要小心
d = dict.fromkeys(["a", "b"], [])
d["a"].append(1)
print(d)  # {'a': [1], 'b': [1]}  （所有值指向同一个列表）

# | 运算符合并字典（Python 3.9+）
d1 = {"a": 1, "b": 2}
d2 = {"b": 3, "c": 4}
print(d1 | d2)  # {'a': 1, 'b': 3, 'c': 4}  （d2 的值覆盖 d1）
print(d2 | d1)  # {'b': 2, 'c': 4, 'a': 1}  （d1 的值覆盖 d2）

d = {"Alice": 90, "Bob": 85, "Charlie": 92}

# 遍历键
for key in d:
    print(key, end=" ")
print()  # Alice Bob Charlie

# 遍历值
for val in d.values():
    print(val, end=" ")
print()  # 90 85 92

# 遍历键值对
for key, val in d.items():
    print(f"{key}={val}", end=" ")
print()  # Alice=90 Bob=85 Charlie=92

# 成员判断（默认判断键）
print("Alice" in d)  # True
print(90 in d)  # False  （判断的是键，不是值）
print(90 in d.values())  # True   （判断值要显式调用 values()）

# defaultdict —— 带默认值的字典
# 默认值为 0
dd = defaultdict(int)
for ch in "hello world":
    dd[ch] += 1
print(dict(dd))  # {'h': 1, 'e': 1, 'l': 3, 'o': 2, ' ': 1, 'w': 1, 'r': 1, 'd': 1}

# 默认值为 []
dd = defaultdict(list)
pairs = [("fruit", "apple"), ("fruit", "banana"), ("veggie", "carrot")]
for category, item in pairs:
    dd[category].append(item)
print(dict(dd))  # {'fruit': ['apple', 'banana'], 'veggie': ['carrot']}

# Counter —— 计数器

counter = Counter("abracadabra")
print(counter)  # Counter({'a': 5, 'b': 2, 'r': 2, 'c': 1, 'd': 1})
print(counter.most_common(3))  # [('a', 5), ('b', 2), ('r', 2)]
print(counter["a"])  # 5
print(counter["z"])  # 0  （不存在的键返回 0）

# Counter 的运算
c1 = Counter(a=3, b=1)
c2 = Counter(a=1, b=2)
print(c1 + c2)  # Counter({'a': 4, 'b': 3})
print(c1 - c2)  # Counter({'a': 2})  （只保留正数）

# Python 3.7+ 普通 dict 也保持插入顺序
# OrderedDict 额外支持 move_to_end 和比较时考虑顺序
od = OrderedDict([("a", 1), ("b", 2), ("c", 3)])

# move_to_end() —— 移到末尾或开头
od.move_to_end("a")
print(list(od.keys()))  # ['b', 'c', 'a']

od.move_to_end("a", last=False)  # 移到开头
print(list(od.keys()))  # ['a', 'b', 'c']

# 顺序不同的 OrderedDict 不相等
od1 = OrderedDict([("a", 1), ("b", 2)])
od2 = OrderedDict([("b", 2), ("a", 1)])
print(od1 == od2)  # False

# 普通 dict 顺序不同也相等
print({"a": 1, "b": 2} == {"b": 2, "a": 1})  # True

s = {1, 2, 3}

# add() —— 添加单个元素
s.add(4)
print(s)  # {1, 2, 3, 4}
s.add(3)  # 已存在，不变
print(s)  # {1, 2, 3, 4}

# update() —— 添加多个元素（接受可迭代对象）
s.update([5, 6])
print(s)  # {1, 2, 3, 4, 5, 6}

s.update([7, 8], {9, 10})  # 可以传多个参数
print(s)  # {1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

# remove() —— 删除元素，不存在则报错
s = {1, 2, 3, 4, 5}
s.remove(5)
print(s)  # {1, 2, 3, 4}
# s.remove(99)  # KeyError: 99

# discard() —— 删除元素，不存在不报错
s.discard(4)
print(s)  # {1, 2, 3}
s.discard(99)  # 不报错
print(s)  # {1, 2, 3}

# pop() —— 随机弹出一个元素
elem = s.pop()
print(f"弹出：{elem}，剩余：{s}")

# clear() —— 清空集合
s.clear()
print(s)  # set()

a = {1, 2, 3, 4, 5}
b = {4, 5, 6, 7, 8}

# ---- 并集 ----
print(a | b)  # {1, 2, 3, 4, 5, 6, 7, 8}
print(a.union(b))  # {1, 2, 3, 4, 5, 6, 7, 8}

# ---- 交集 ----
print(a & b)  # {4, 5}
print(a.intersection(b))  # {4, 5}

# ---- 差集 ----
print(a - b)  # {1, 2, 3}  （在 a 中但不在 b 中）
print(a.difference(b))  # {1, 2, 3}

# ---- 对称差集 ----
print(a ^ b)  # {1, 2, 3, 6, 7, 8}
print(a.symmetric_difference(b))  # {1, 2, 3, 6, 7, 8}

# union / intersection / difference 可以接受任意可迭代对象
print(a.union([10, 11]))  # {1, 2, 3, 4, 5, 10, 11}
print(a.intersection(range(3, 8)))  # {3, 4, 5}

a = {1, 2, 3, 4, 5}
b = {4, 5, 6, 7, 8}

# intersection_update() —— 原地取交集
s = a.copy()
s.intersection_update(b)
print(s)  # {4, 5}
# 等价于 s &= b

# difference_update() —— 原地取差集
s = a.copy()
s.difference_update(b)
print(s)  # {1, 2, 3}
# 等价于 s -= b

# symmetric_difference_update() —— 原地取对称差集
s = a.copy()
s.symmetric_difference_update(b)
print(s)  # {1, 2, 3, 6, 7, 8}
# 等价于 s ^= b

a = {1, 2, 3, 4, 5}
b = {1, 2, 3}
c = {6, 7, 8}

# issubset() —— 是否为子集
print(b.issubset(a))  # True   b 是 a 的子集
print(b <= a)  # True   运算符写法
print(b < a)  # True   真子集（子集且不相等）

# issuperset() —— 是否为超集
print(a.issuperset(b))  # True   a 是 b 的超集
print(a >= b)  # True
print(a > b)  # True   真超集

# isdisjoint() —— 是否没有交集
print(a.isdisjoint(c))  # True   a 和 c 无公共元素
print(a.isdisjoint(b))  # False  a 和 b 有公共元素

# 相等判断
print({1, 2, 3} == {3, 2, 1})  # True  （集合无序，元素相同即相等）

# frozenset —— 不可变集合
fs = frozenset([1, 2, 3, 2, 1])
print(fs)  # frozenset({1, 2, 3})
print(type(fs))  # <class 'frozenset'>

# 支持所有不修改集合的方法
a = frozenset([1, 2, 3, 4])
b = frozenset([3, 4, 5, 6])
print(a | b)  # frozenset({1, 2, 3, 4, 5, 6})
print(a & b)  # frozenset({3, 4})
print(a - b)  # frozenset({1, 2})

# 不支持修改操作
# fs.add(4)       # AttributeError
# fs.remove(1)    # AttributeError

# 可以作为字典的键或集合的元素
d = {frozenset([1, 2]): "pair"}
print(d[frozenset([1, 2])])  # pair

s = {frozenset([1, 2]), frozenset([3, 4])}
print(s)  # {frozenset({1, 2}), frozenset({3, 4})}
