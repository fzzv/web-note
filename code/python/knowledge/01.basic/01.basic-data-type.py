import copy

# number
a = 10
b = -1
c = 0

print(type(a))

binary = 0b1010  # 二进制
octal = 0o12  # 八进制
hexadecimal = 0xA  # 十六进制

print(binary, octal, hexadecimal)

f1 = 3.14
f2 = -0.5
f3 = 1.0
f4 = 1.5e2
f5 = 3e-3

print(f1, f2, f3, f4, f5)

print(0.1 + 0.2)
print(0.1 + 0.2 == 0.3)

print(round(0.1 + 0.2))

c1 = 3 + 4j
c2 = complex(1, 2)
print(c1 + c2)

print(abs(-7))
print(max(1, 5, 3))
print(min(1, 5, 3))
print(divmod(10, 3))

# str
# 单引号、双引号、三引号
s1 = "hello"
s2 = "world"
s3 = """多行
字符串"""
s4 = """也是
多行字符串"""
print(type(s1))
print(s3)

print(r"C:\new\test")  # C:\new\test

s = "Hello, Python!"

print(s[0])
print(s[-1])
print(s[0:2])
print(s[:5])
print(s[5:])
print(s[::2])
print(s[::-1])  # 反转字符串

name = "Alice"
age = 25
score = 95.55

print(f"我叫{name}，今年{age}岁")
print(f"成绩：{score:.1f}")
print(f"十六进制：{255:#x}")
print(f"{'居中':=^20}")

# list
list1 = [1, 2, 3, 4, 5]
list2 = list()
list3 = list("hello")
list4 = list(range(5))
list5 = list(range(1, 6))
print(list1, list2, list3, list4, list5)

fruits = ["apple", "banana"]
fruits.append("cherry")
print(fruits)
fruits.insert(1, "blueberry")
print(fruits)
fruits.extend(["date", "fig"])
print(fruits)
fruits.remove("blueberry")
print(fruits)
del fruits[0]
print(fruits)

squares = [x**2 for x in range(1, 6)]
print(squares)

# tuple
t1 = (1, 2, 3)
t2 = ("a", "b", "c")
t3 = tuple([1, 2, 3])
t4 = (42,)  # 单元素元组必须加逗号
print(t1, t2, t3, t4)
print(t1.count(1))
print(t1.index(1))

first, *rest = (1, 2, 3, 4, 5)
print(first, rest)

# dict
d1 = {"name": "Alice", "age": 25, "city": "Beijing"}
d2 = dict(name="Bob", age=30)
d3 = dict([("a", 1), ("b", 2)])

d4 = dict.fromkeys(["a", "b", "c"], 0)
print(d4)
print(type(d1))
print(len(d1))

scores = {"Alice": 90, "Bob": 85, "Charlie": 92}
for name in scores:
    print(name, end=" ")

for score in scores.values():
    print(score, end=" ")

for name, score in scores.items():
    print(f"{name}: {score}")

print("Alice" in scores)
print("Alice" not in scores)

# set 无序、不重复的元素集 元素必须是不可变类型
s1 = {1, 2, 3, 4, 5}
s2 = set([1, 2, 2, 3, 3])
s3 = set("hello")
print(type(s1), s1, s2, s3)

# 拷贝
a = [1, 2, 3]
b = a
a.append(4)
print(a)
print(b)

a = [1, 2, 3]
b = a.copy()
a.append(4)
print(a)
print(b)

a = [[1, 2], [3, 4]]
b = copy.deepcopy(a)
a[0].append(5)
print(a)
print(b)
