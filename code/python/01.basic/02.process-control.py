# if / elif / else
score = 85
if score >= 90:
    grade = "A"
elif score >= 80:
    grade = "B"
elif score >= 70:
    grade = "C"
elif score >= 60:
    grade = "D"
else:
    grade = "F"
print(f"成绩等级：{grade}")

# 三元运算符
age = 20
status = "成年" if age >= 18 else "未成年"
print(status)

# 判空
data = ""
if not data:
    print("字符串为空")

# 使用 any() / all()
scores = [80, 90, 75, 60, 95]
if all(s >= 60 for s in scores):
    print("全部及格")

if any(s >= 90 for s in scores):
    print("有人得了 90 分以上")

# match-case
command = "quit"

match command:
    case "start":
        print("启动程序")
    case "stop":
        print("停止程序")
    case "quit":
        print("退出程序")
    case _:  # _ 是通配符，匹配所有情况
        print("未知命令")

point = (3, 4)

match point:
    case (x, y) if x == y:
        print(f"在对角线上：({x}, {y})")
    case (x, y) if x > 0 and y > 0:
        print(f"在第一象限：({x}, {y})")
    case (x, y):
        print(f"其他位置：({x}, {y})")

# for循环
fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    print(fruit)

for ch in "Python":
    print(ch, end=" ")

for i in range(5):
    print(i, end=" ")
print()

languages = ["Python", "Java", "Go"]
for index, lang in enumerate(languages):
    print(f"{index}: {lang}")

# while
count = 0
while count < 5:
    print(count, end=" ")
    count += 1
print()

# for-else / while-else
for n in range(2, 10):
    for x in range(2, n):
        if n % x == 0:
            break
    else:
        # 循环没有被 break，说明 n 是质数
        print(f"{n} 是质数")

print(range(1000000))

# try-except
try:
    result = 10 / 0
except ZeroDivisionError:
    print("除数不能为零")

try:
    lst = [1, 2, 3]
    print(lst[10])
except IndexError:
    print("索引越界")
except TypeError:
    print("类型错误")
except Exception as e:
    print(f"其他错误：{e}")


# try-except-else-finally
def divide(a, b):
    try:
        result = a / b
    except ZeroDivisionError:
        print("错误：除数不能为零")
        return None
    else:
        # 没有异常时执行
        print(f"{a} / {b} = {result}")
        return result
    finally:
        # 无论是否异常都会执行
        print("运算结束")


divide(10, 3)
divide(10, 0)


# raise 主动抛出异常
def set_age(age):
    if not isinstance(age, int):
        raise TypeError("年龄必须是整数")
    if age < 0 or age > 150:
        raise ValueError(f"年龄超出合理范围：{age}")
    print(f"设置年龄为：{age}")


set_age(20)


# pass、assert 和 del
def todo_function():
    pass


def calculate_average(numbers):
    assert len(numbers) > 0, "列表不能为空"
    return sum(numbers) / len(numbers)


# print(calculate_average([]))
print(calculate_average([1, 2, 3]))

lst = [1, 2, 3, 4, 5]
del lst[0]
print(lst)
