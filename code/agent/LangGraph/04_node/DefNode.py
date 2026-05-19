from functools import partial
from typing import TypedDict

from langgraph.graph import StateGraph, START, END
from langgraph.types import RetryPolicy
from requests import RequestException, Timeout


# 定义状态
class GraphState(TypedDict):
    process_data: dict  # 默认更新策略


# 定义一个节点，入参为state
def input_node(state: GraphState) -> GraphState:
    print(f"input_node收到的初始值:{state}")
    return {"process_data": {"input": "input_value"}}


# 定义带参数的node节点
def process_node(state: dict, param1: int, param2: str) -> dict:
    print(state, param1, param2)
    return {"process_data": {"process": "process_value"}}


# 重试策略,add_node方法时可选
retry_policy = RetryPolicy(
    max_attempts=3,  # 最大重试次数
    initial_interval=1,  # 初始间隔
    jitter=True,  # 抖动（添加随机性避免重试风暴）
    backoff_factor=2,  # 退避乘数（每次重试间隔时间的增长倍数）
    retry_on=[RequestException, Timeout],  # 只重试这些异常
)


stateGraph = StateGraph(GraphState)
# 添加inpu节点
stateGraph.add_node("input", input_node)
# 给process_node节点绑定参数
process_with_params = partial(process_node, param1=100, param2="test")
# 添加带参数的node节点
stateGraph.add_node("process", process_with_params, retry=retry_policy)

# 定义节点之间的执行顺序 edges
# 设置节点间的依赖关系，形成执行流程图
stateGraph.add_edge(START, "input")
stateGraph.add_edge("input", "process")
stateGraph.add_edge("process", END)

# 编译图构建器生成计算图
graph = stateGraph.compile()


# # 打印图的边和节点信息
print(stateGraph.edges)
print(stateGraph.nodes)
# 打印图的可视化结构
print(graph.get_graph().print_ascii())

print()

# 定义一个初始状态字典，包含键值对"x": 5
initial_state = {"process_data": 5}
# 调用graph对象的invoke方法，传入初始状态，执行图计算流程
result = graph.invoke(initial_state)
print(f"最后的结果是:{result}")
