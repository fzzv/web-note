def my_decorator(func):
    def wrapper():
        print("--- 函数执行前 ---")
        func()
        print("--- 函数执行后 ---")

    return wrapper


# 手动装饰
def say_hello():
    print("Hello!")


# 装饰器装饰
@my_decorator
def say_no():
    print("No!")


say_hello = my_decorator(say_hello)  # 用新函数替换原函数
say_hello()
say_no()


def custom_decorator(func):
    def wrapper(*args, **kwargs):
        print(f"调用 {func.__name__}，参数:args={args}, kwargs={kwargs}")
        result = func(*args, **kwargs)
        print(f"{func.__name__} 的返回值: {result}")
        return result

    return wrapper


@custom_decorator
def add(a, b):
    return a + b


@custom_decorator
def greet(name, greeting="Hello"):
    return f"{greeting}, {name}!"


print(add(1, 2))
print(greet("Alice", greeting="Hi"))
