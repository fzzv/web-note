"""
SqlitePersistence.py
在底层，检查点功能由符合BaseCheckpointSaver接口的检查点对象提供支持。
LangGraph提供了多种检查点实现，所有这些实现都通过独立的、可安装的库来完成，数据库类型的有：
	langgraph-checkpoint-sqlite：使用SQLite数据库（SqliteSaver / AsyncSqliteSaver）存储检查点。
非常适合实验和本地工作流程。需要单独安装。
	langgraph-checkpoint-postgres：使用Postgres数据库（PostgresSaver / AsyncPostgresSaver）
存储检查点，用于LangSmith。非常适合在生产环境中使用。需要单独安装。
......

本次案例，安装sqlite所需依赖
pip install langgraph-checkpoint-sqlite
"""

import sqlite3
import operator
from typing import TypedDict, Annotated
from langgraph.checkpoint.sqlite import SqliteSaver
from langgraph.graph import StateGraph, START, END


class MyState(TypedDict):
    messages: Annotated[list, operator.add]


def node_1(state: MyState):
    return {"messages": ["abc", "def"]}


def main():
    # 数据存储到D:\\workspace目录下面，需要目录存在
    conn = sqlite3.connect(
        database="D:\\workspace\\sqlite_data.db", check_same_thread=False
    )
    sqliteDB = SqliteSaver(conn=conn)

    builder = StateGraph(MyState)
    builder.add_node("node_1", node_1)

    builder.add_edge(START, "node_1")
    builder.add_edge("node_1", END)

    graph = builder.compile(checkpointer=sqliteDB)
    # 同一个用户id下，每次执行都会插入一次新数据
    config = {"configurable": {"thread_id": "user-001"}}

    initial_state = graph.get_state(config)
    print(f"Initial state: {initial_state}")

    # 执行图
    result = graph.invoke({"messages": []}, config)
    print(f"Result: {result}")

    print()
    print("====================查看执行后的状态====================")
    # 查看执行后的状态
    final_state = graph.get_state(config)
    print()
    print(f"Final state: {final_state}")

    conn.close()


if __name__ == "__main__":
    main()
