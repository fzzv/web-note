from typing import TypedDict
from langgraph.graph import StateGraph, START, END


class BasicState(TypedDict):
    """基本的 State定义"""

    user_input: str
    response: str
    count: int
    process_data: dict


# 创建状态图，并指定状态结构
basicState = StateGraph(BasicState)
# 添加起始到结束的边（无中间节点）
basicState.add_edge(START, END)
# 编译生成计算图
app = basicState.compile()

# invoke()方法只接收状态字典作为核心参数
initial_state = {
    "user_input": "a",
    "response": "resp",
    "count": 25,
    "process_data": {"k1": "v1"},  # process_data本身是dict类型，需嵌套
}

# invoke() 仅接收 1 个核心位置参数（状态字典），可选 1 个配置参数，切勿传入多个独立参数。
result = app.invoke(initial_state)
# 打印结果验证
print("执行结果：", result)
