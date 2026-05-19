"""
在LangGraph中，一个Graph除了可以单独使用，还可以作为一个Node，嵌入到一个Graph中。这种用法就称为子图。
通过子图，我们可以更好的重用Graph，构建更复杂的工作流。尤其在构建多Agent系统时非常有用。
在大型项目中，通常都是由一个团队专门开发Agent，再通过其他团队来完成Agent整合。

使用子图时，基本和使用Node没有太多的区别。唯一需要注意的是，当触发了SubGraph代表的Node后，
实际上是相当于重新调用了一次subgraph.invoke(state)方法

案例说明：
    定义一个子图节点处理函数 sub_node，它接收一个状态对象并返回包含子图响应消息的新状态。
    该函数被集成到一个使用 langgraph 构建的图结构中，最终执行图并输出结果。
"""

from operator import add
from typing import TypedDict, Annotated
from langgraph.constants import END
from langgraph.graph import StateGraph, MessagesState, START
import operator


class AtguiguState(TypedDict):
    """
    定义状态类，用于存储图节点间传递的消息状态
    messages: 使用add函数合并的字符串列表消息
    add 是 LangGraph 内置的状态合并策略，它的行为是：将新返回的列表与原状态中的列表进行拼接（而非覆盖）
    """

    messages: Annotated[list[str], add]


def sub_node(state: AtguiguState) -> AtguiguState:
    # 子图节点处理函数，接收当前状态并返回响应消息
    # @param state 当前状态对象，包含消息列表
    # @return 包含子图响应消息的新状态
    return {"messages": ["response from subgraph"]}


# 创建子图构建器并配置节点和边
subgraph_builder = StateGraph(AtguiguState)
subgraph_builder.add_node("sub_node", sub_node)

subgraph_builder.add_edge(START, "sub_node")
subgraph_builder.add_edge("sub_node", END)
subgraph = subgraph_builder.compile()

# 创建主图构建器并添加子图节点
builder = StateGraph(AtguiguState)
builder.add_node("subgraph_node", subgraph)
builder.add_edge(START, "subgraph_node")
builder.add_edge("subgraph_node", END)

# 编译主图并绘制结构图
graph = builder.compile()

# 执行图并打印结果
"""子图调用的状态传递逻辑当主图调用子图节点时，整个过程会触发两次状态合并：
第一步：主图把初始状态 {"messages": ["main-graph"]} 传递给子图

第二步：子图内部执行 sub_node，返回 {"messages": ["response from subgraph"]}，
        由于 add 策略，子图会把传入的 ["main-graph"] 和返回的 ["response from subgraph"] 拼接，
        得到 ["main-graph", "response from subgraph"]

第三步：子图执行完成后，主图会再次应用 add 策略，
    把主图原有的 ["main-graph"]
    和子图返回的 ["main-graph", "response from subgraph"] 拼接，
    最终得到 ["main-graph", "main-graph", "response from subgraph"]"""
print(graph.invoke({"messages": ["main-graph"]}))
print()  # {'messages': ['main-graph', 'main-graph', 'response from subgraph']}


# 绘制子图结构图
print(subgraph.get_graph().draw_mermaid())
print("=" * 50)
print()
