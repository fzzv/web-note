"""
从节点调用图，本案例是 LangGraph 中跨图状态交互的标准做法。

核心逻辑解释
1. 状态结构差异设计
父图状态（ParentState）：仅包含 user_query（用户输入）和 final_answer（最终结果），聚焦业务层；
子图状态（SubgraphState）：
包含 analysis_input（分析输入）、analysis_result（分析结果）、intermediate_steps（中间步骤），
聚焦分析层；两者无重叠字段，完全独立，必须通过代理节点手动转换。

2. 代理节点核心作用（call_subgraph_proxy）
  2.1 步骤1：父→子状态转换（按子图要求构造输入）
  2.2 步骤2：手动调用子图（而非直接将子图作为父图节点）
  2.3 步骤3：子→父状态映射（提取子图结果，赋值给父图字段）

核心方案：
    父子图状态不同时，通过父图的代理节点而非直接添加子图节点，
    手动完成「父状态→子输入」转换、调用子图、「子输出→父状态」映射；
关键要点：
    代理节点必须接收父图状态、返回父图状态，子图调用通过 compiled_subgraph.invoke() 手动触发；
灵活性：
    该模式可适配任意结构的父子图状态，只需在代理节点中自定义转换逻辑
"""

from langgraph.graph import StateGraph, START, END
from typing import TypedDict


# ====================== 1. 定义不同结构的父子图状态 ======================
# 父图状态：仅包含用户查询和最终答案（与子图状态完全不同）
class ParentState(TypedDict):
    user_query: str  # 父图独有：用户输入的查询
    final_answer: str | None  # 父图独有：子图处理后的最终结果


# 子图状态：专注于分析逻辑（与父图状态无重叠字段）
class SubgraphState(TypedDict):
    analysis_input: str  # 子图独有：分析输入
    analysis_result: str  # 子图独有：分析结果
    intermediate_steps: list  # 子图独有：中间步骤（私有数据）


# ====================== 2. 定义子图核心逻辑 ======================
def subgraph_analysis_node(state: SubgraphState) -> SubgraphState:
    """子图核心节点：处理分析逻辑，生成结果"""
    # 模拟子图的分析过程
    query = state["analysis_input"]
    state["intermediate_steps"] = [f"解析查询：{query}", "执行分析逻辑", "生成结果"]
    state["analysis_result"] = f"针对「{query}」的分析结果：这是子图处理后的内容"
    return state


def build_subgraph() -> StateGraph:
    """构建并编译子图"""
    sub_builder = StateGraph(SubgraphState)
    sub_builder.add_node("subgraph_analysis_node", subgraph_analysis_node)

    sub_builder.add_edge(START, "subgraph_analysis_node")
    sub_builder.add_edge("subgraph_analysis_node", END)
    return sub_builder.compile()


# 提前编译子图（供父图代理节点调用）
compiled_subgraph = build_subgraph()


# ============ 3. 定义父图代理节点（核心：状态转换+调用子图）从节点调用图=======
def call_subgraph_proxy(state: ParentState) -> ParentState:
    """
    父图的代理节点：
    1. 将父图状态转换为子图所需的输入格式
    2. 手动调用子图
    3. 将子图输出映射回父图状态
    """

    # 步骤1：父图状态 → 子图输入（状态转换）,提取父图的user_query，转换为子图需要的analysis_input
    subgraph_input = {
        "analysis_input": state["user_query"],
        "intermediate_steps": [],  # 初始化子图的私有字段
        "analysis_result": "",  # 初始化子图结果字段
    }

    # 步骤2：手动调用编译后的子图，手动调用子图（而非直接将子图作为父图节点）
    subgraph_response = compiled_subgraph.invoke(subgraph_input)

    # 步骤3：子图输出 → 父图状态（结果映射）
    # 提取子图的analysis_result，赋值给父图的final_answer
    return {
        "user_query": state["user_query"],  # 保留父图原有字段
        "final_answer": subgraph_response["analysis_result"],
    }


def build_parent_graph() -> StateGraph:
    """构建并编译父图（添加代理节点，而非直接添加子图）"""
    parent_builder = StateGraph(ParentState)
    # 添加代理节点（核心：手动处理状态转换+调用子图）
    parent_builder.add_node("call_subgraph_proxy", call_subgraph_proxy)
    # 父图执行链路：START → 代理节点 → END
    parent_builder.add_edge(START, "call_subgraph_proxy")
    parent_builder.add_edge("call_subgraph_proxy", END)
    return parent_builder.compile()


# ====================== 4. 主方法 ======================
def main():
    """主函数：执行父图，验证跨图状态转换逻辑"""
    # 1. 构建父图
    parent_graph = build_parent_graph()

    # 2. 定义父图初始状态（仅包含user_query，符合父图状态结构）
    initial_state = {
        "user_query": "请分析Python中StateGraph的使用场景",
        "final_answer": None,
    }
    print("父图初始状态：", initial_state)

    # 3. 执行父图，实际而言父图调用了call_subgraph_proxy
    final_state = parent_graph.invoke(initial_state)

    # 4. 输出结果
    print("\n父图最终状态：", final_state)
    print("\n子图处理后的最终答案：", final_state["final_answer"])


if __name__ == "__main__":
    main()
