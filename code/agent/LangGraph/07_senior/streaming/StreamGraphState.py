"""
StreamGraphState.py

流图状态
使用流模式，并在图执行时流式传输其状态。updatesvalues
updates在图的每一步后，将更新流向状态。
values在图的每一步后，流出状态的---->全部值。
"""

from typing import TypedDict
from langgraph.graph import StateGraph, START, END


class AtguiguState(TypedDict):
    topic: str
    joke: str


def refine_topic(state: AtguiguState):
    return {"topic": state["topic"] + " and cats"}


def generate_joke(state: AtguiguState):
    return {"joke": f"This is a joke about {state['topic']}"}


def main():
    graph = (
        StateGraph(AtguiguState)
        .add_node(refine_topic)
        .add_node(generate_joke)
        .add_edge(START, "refine_topic")
        .add_edge("refine_topic", "generate_joke")
        .add_edge("generate_joke", END)
        .compile()
    )

    # updates在图的每一步后，将更新流向状态。
    for chunk in graph.stream({"topic": "ice cream"}, stream_mode="updates"):
        print(chunk)

    print()

    # values在图的每一步后，流出状态的全部值。
    for chunk in graph.stream({"topic": "ice cream"}, stream_mode="values"):
        print(chunk)


if __name__ == "__main__":
    main()
