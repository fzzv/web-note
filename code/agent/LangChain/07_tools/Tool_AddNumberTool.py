from langchain.tools import tool


@tool
def add_number(a: int, b: int) -> int:
    """两个整数相加"""
    return a + b


result = add_number.invoke({"a": 1, "b": 12})
print(result)

print()

print(f"{add_number.name=}\n{add_number.description=}\n{add_number.args=}")
# add_number.name='add_number'
# add_number.description='两个整数相加'
# add_number.args={'a': {'title': 'A', 'type': 'integer'}, 'b': {'title': 'B', 'type': 'integer'}}
