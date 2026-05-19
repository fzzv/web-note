import time
from typing_extensions import TypedDict
from langgraph.graph import StateGraph
from langgraph.cache.memory import InMemoryCache
from langgraph.types import CachePolicy


# 定义状态类，也就是你的业务实体entity
class State(TypedDict):
    x: int
    result: int


# 创建图
builder = StateGraph(State)


# 定义节点：模拟耗时计算（sleep3秒）
def expensive_node(state: State) -> dict[str, int]:
    time.sleep(3)
    return {"result": state["x"] * 2}


#     builder.add_node("node1", node_default_1)

# 添加节点
builder.add_node(
    node="expensive_node",
    action=expensive_node,
    # 不用传key_fn，底层自动用默认逻辑
    cache_policy=CachePolicy(ttl=8),
)

# 设置入口和出口
builder.set_entry_point("expensive_node")
builder.set_finish_point("expensive_node")

# 编译图，指定内存缓存
app = builder.compile(cache=InMemoryCache())

# 第一次执行：耗时3秒（无缓存）
print("第一次执行（无缓存，耗时3秒）：")
print(app.invoke({"x": 5}))
# 第二次执行：瞬间返回（利用缓存，8秒内有效）
print("\n1111111111111111111111111111")
print("第二次运行利用缓存并快速返回：")
print(app.invoke({"x": 5}))

# 可选：测试8秒后缓存过期（取消注释查看）
print("\n等待8秒，缓存过期...")
time.sleep(8)
print("8秒后第三次执行（重新计算，耗时3秒）：")
print(app.invoke({"x": 5}))
