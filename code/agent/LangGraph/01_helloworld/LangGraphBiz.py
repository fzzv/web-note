"""
我们先在不接入大模型的情况下构建一个加减法图工作流，
我们这里自定义两个简单函数：一个是加法函数接收当前State并将其中的x值加1，
另一个是减法函数接收当前State并将其中的x值减2，
然后添加名为addition和subtraction的节点，并关联到两个函数上，最后构建出节点之间的边。
"""

from langgraph.constants import START, END
from langgraph.graph import StateGraph


def addition(state):
    print(f"加法节点收到的初始值:{state}")
    return {"x": state["x"] + 1}


def subtraction(state):
    print(f"减法节点收到的初始值:{state}")
    return {"x": state["x"] - 2}


graph = StateGraph(dict)
# 向图构建器中添加节点
# 添加加法运算节点和减法运算节点到构建器中
graph.add_node("addition", addition)
graph.add_node("subtraction", subtraction)

# 定义节点之间的执行顺序 edges
# 设置节点间的依赖关系，形成执行流程图
graph.add_edge(START, "addition")
graph.add_edge("addition", "subtraction")
graph.add_edge("subtraction", END)
# 打印图的边和节点信息
print(graph.edges)
print()
print(graph.nodes)

# 编译图构建器生成计算图
app = graph.compile()
# invoke()方法只接收状态字典作为核心参数，定义一个初始状态字典，包含键值对"x": 5
initial_state = {"x": 5}
# 调用graph对象的invoke方法，传入初始状态，执行图计算流程
result = app.invoke(initial_state)
print(f"最后的结果是:{result}")

# # 打印图的可视化结构
print(app.get_graph().print_ascii())
print()
# 打印图的可视化结构，生成更加美观的Mermaid 代码，通过processon 编辑器查看
print(app.get_graph().draw_mermaid())
