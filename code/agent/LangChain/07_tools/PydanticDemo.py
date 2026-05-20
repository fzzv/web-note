from pydantic import BaseModel, ValidationError, StrictInt


class User(BaseModel):
    # id: int
    id: StrictInt  # 改用严格整数类型，拒绝类型转换
    name: str
    age: int = 0  # 可给默认值


try:
    # 自动把字符串转成 int
    # u = User(id="41", name="z3") #  自动把字符串转成 int，可以通融。id: int
    u = User(id=42, name="z3")  # 传错类型就报错 id: StrictInt
except ValidationError as e:
    print(e)
print(u.id, type(u.id))  # 42 <class 'int'>

print()
print()


try:
    User(id="abc", name="Bob")  # 传错类型就报错 id: StrictInt
except ValidationError as e:
    print(e)
"""
1 validation error for User
id
  value is not a valid integer (type=type_error.integer)
"""
