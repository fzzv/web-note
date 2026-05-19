from typing import TypedDict
from langgraph.constants import START, END
from langgraph.graph import StateGraph

"""图的构建流程：
1、初始化一个StateGraph实例。
2、添加节点。
3、定义边，将所有的节点连接起来。
4、设置特殊节点，入口和出口（可选）。
5、编译图。
6、执行工作流。
"""


# 定义状态
class GraphState(TypedDict):
    process_data: dict


def input_node(state: GraphState) -> GraphState:
    print(
        f"input_node节点执行state.get('process_data')方法结果:  {state.get('process_data')}"
    )
    return {"process_data": {"input": "input_value"}}


def process_node(state: dict) -> dict:
    print(
        f"process_node节点执行state.get('process_data')方法结果:  {state.get('process_data')}"
    )
    return {"process_data": {"process": "process_value9527"}}


def output_node(state: GraphState) -> GraphState:
    print(
        f"output_node节点执行state.get('process_data')方法结果:  {state.get('process_data')}"
    )
    return {"process_data": state.get("process_data")}


# 创建一个状态图StateGraph并指定状态
graph = StateGraph(GraphState)

# 添加input、process、output节点
graph.add_node("input", input_node)
graph.add_node("process", process_node)
graph.add_node("output", output_node)

# 添加固定边，执行顺序：start -> input -> process -> output -> end
graph.add_edge(START, "input")
graph.add_edge("input", "process")
graph.add_edge("process", "output")
graph.add_edge("output", END)

# 编译图，保证生成的图是正确的，如果添加了边，没添加节点，会报错
app = graph.compile()
# 执行
result = app.invoke({"process_data": {"name": "测试数据", "value": 123456}})
print(f"最后的结果是:{result}")

# 打印图的ascii可视化结构
print(app.get_graph().print_ascii())
print("=================================")
print()
# 打印图的可视化结构，生成更加美观的Mermaid 代码，通过processon 编辑器查看
print(app.get_graph().draw_mermaid())
