import os
import re
from langchain_openai import ChatOpenAI
from langchain.agents import create_agent
from langgraph_supervisor import create_supervisor
from dotenv import load_dotenv

load_dotenv()


# 1. 初始化大语言模型
def init_llm_model() -> ChatOpenAI:
    return ChatOpenAI(
        model="qwen-plus",
        api_key=os.getenv("QWEN_API_KEY"),
        base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
        temperature=0.1,
        max_tokens=1024,
    )


# 2. Tools（必须有 docstring）
def book_flight(from_airport: str, to_airport: str) -> str:
    """预订航班工具。根据出发机场和到达机场预订一张机票，并返回预订结果。"""
    return f"✅ 成功预订了从 {from_airport} 到 {to_airport} 的航班"


def book_hotel(hotel_name: str) -> str:
    """预订酒店工具。根据酒店名称完成酒店预订，并返回预订结果。"""
    return f"✅ 成功预订了 {hotel_name} 的住宿"


# 3. 子 Agent
flight_assistant = create_agent(
    model=init_llm_model(), tools=[book_flight], name="flight_assistant"
)

hotel_assistant = create_agent(
    model=init_llm_model(), tools=[book_hotel], name="hotel_assistant"
)

# 4. 创建 Supervisor，协调者主管
supervisor = create_supervisor(
    agents=[flight_assistant, hotel_assistant],
    model=init_llm_model(),
    prompt=(
        "你是旅行预订系统的调度主管，负责协调航班预订和酒店预订。\n\n"
        "当用户提出航班和酒店预订请求时，你的工作流程是：\n"
        "1. 首先调用flight_assistant来预订航班\n"
        "2. 然后调用hotel_assistant来预订酒店\n"
        "3. 收到两个助手的结果后，汇总并向用户报告\n"
        "4. 完成后结束对话\n\n"
        "重要规则：\n"
        "- 每个助手只能调用一次\n"
        "- 不要重复任何内容\n"
        "- 不要输出任何英文\n"
        "- 所有通信都使用中文\n"
    ),
).compile()


# 5. 消息过滤器，就是一个工具类，处理大模型返回的重复废话，直接用可以不看
def filter_messages(chunk: dict) -> str:
    """提取并过滤消息，只返回中文内容，去除重复和英文"""
    output = ""

    if isinstance(chunk, dict):
        for role, payload in chunk.items():
            if isinstance(payload, dict) and "messages" in payload:
                for msg in payload["messages"]:
                    if hasattr(msg, "content") and msg.content:
                        content = msg.content.strip()

                        # 过滤英文系统消息
                        if (
                            content
                            and not content.startswith("Successfully")
                            and not content.startswith("Transferring")
                            and "Successfully transferred" not in content
                            and "transferred back to" not in content
                            and not content.startswith("帮我预订从")
                        ):
                            # 只保留中文内容
                            chinese_content = re.sub(
                                r'[^\u4e00-\u9fff，。！？：；""、\s\d✅]', "", content
                            )
                            if chinese_content and len(chinese_content.strip()) > 5:
                                output += f"{role}: {chinese_content.strip()}\n"

    return output


# 6. 主程序
def main():
    print("=" * 60)
    print(
        "智能旅行预订系统，由于大模型每次调用，可能出现预定不成功情况，这是正常反馈,主要是2026.2.8千问赠送奶茶活动，调用失败"
    )
    print("=" * 60)
    print()

    # 收集用户信息
    print("请按顺序提供以下信息：")
    print("-" * 40)

    # 1. 询问出发机场
    from_airport = input("1. 您的出发机场是哪里？: ").strip()
    while not from_airport:
        print("请输入有效的出发机场名称")
        from_airport = input("1. 您的出发机场是哪里？: ").strip()

    # 2. 询问到达机场
    to_airport = input("\n2. 您的到达机场是哪里？: ").strip()
    while not to_airport:
        print("请输入有效的到达机场名称")
        to_airport = input("2. 您的到达机场是哪里？: ").strip()

    # 3. 询问酒店名称
    hotel_name = input("\n3. 您要预订的酒店名称是什么？: ").strip()
    while not hotel_name:
        print("请输入有效的酒店名称")
        hotel_name = input("3. 您要预订的酒店名称是什么？: ").strip()

    # 构造更明确的用户请求
    user_request = (
        f"请帮我预订以下旅行安排：\n"
        f"1. 航班：从 {from_airport} 飞往 {to_airport}\n"
        f"2. 酒店：{hotel_name}\n"
        f"请完成这两个预订。"
    )

    print("\n" + "=" * 60)
    print("正在处理您的预订请求...")
    print("=" * 60)
    print()

    # 准备输入数据
    # 创建一个字典，包含一个messages键
    # messages是一个列表，包含一个消息字典
    # 每个消息字典包含role（角色）和content（内容）字段
    input_data = {"messages": [{"role": "user", "content": user_request}]}

    # 使用流式处理
    try:
        # 创建一个空集合，用于记录已经打印过的消息内容，避免重复显示
        seen_contents = set()

        for chunk in supervisor.stream(input_data):
            # 调用filter_messages函数处理当前chunk，提取并过滤其中的消息
            filtered_output = filter_messages(chunk)
            # 如果filtered_output不为空（即有过滤后的消息内容）
            if filtered_output:
                # 将过滤后的输出按行分割成列表 strip() 去除首尾空白字符，split('\n') 按换行符分割
                lines = filtered_output.strip().split("\n")
                # 遍历每一行
                for line in lines:
                    # 检查该行是否非空且不在已见过内容的集合中
                    if line and line not in seen_contents:
                        print(line)
                        # 将该行内容添加到已见过集合中，确保不会重复打印
                        seen_contents.add(line)

        # 如果输出太少，显示总结信息
        if len(seen_contents) < 2:
            print("\n" + "=" * 60)
            print("预订已完成！")
            print(f"航班：从 {from_airport} 到 {to_airport}")
            print(f"酒店：{hotel_name}")
            print("=" * 60)
    except Exception as e:
        print(f"\n处理过程中出现错误: {e}")
        # 如果出错，直接调用工具
        print("\n正在直接执行预订...")
        flight_result = book_flight(from_airport, to_airport)
        hotel_result = book_hotel(hotel_name)
        print(flight_result)
        print(hotel_result)

    print("\n感谢使用智能旅行预订系统！")


# 7. 运行主程序
if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n程序被用户中断。")
    except Exception as e:
        print(f"\n系统出现错误: {e}")


"""
运行结果如下：

============================================================
智能旅行预订系统，由于大模型每次调用，可能出现预定不成功情况，这是正常反馈
============================================================

请按顺序提供以下信息：
----------------------------------------
1. 您的出发机场是哪里？: 北京首都

2. 您的到达机场是哪里？: 深圳宝安

3. 您要预订的酒店名称是什么？: 深圳希尔顿

============================================================
正在处理您的预订请求...
============================================================

supervisor: 请帮我预订以下旅行安排：
1 航班：从 北京首都 飞往 深圳宝安
2 酒店：深圳希尔顿
请完成这两个预订。
flight_assistant: 航班已成功预订！接下来，我将为您预订深圳希尔顿酒店。由于当前工具仅支持航班预订，我需要切换到酒店预订助手来完成此任务。
正在为您转接到酒店预订助手
supervisor: 航班已成功预订！接下来，我将为您预订深圳希尔顿酒店。由于当前工具仅支持航班预订，我需要切换到酒店预订助手来完成此任务。
hotel_assistant: ✅ 您的旅行安排已全部完成：
1 航班：已成功预订，从北京首都机场飞往深圳宝安机场由航班助手完成。
2 酒店：已成功预订 深圳希尔顿由酒店助手完成。
如有其他需求如接送、餐饮、景点门票等，欢迎随时告诉我！祝您旅途愉快！
supervisor: ✅ 您的旅行安排已全部完成：
1 航班：已成功预订，从北京首都机场飞往深圳宝安机场。
2 酒店：已成功预订深圳希尔顿酒店。
所有预订均已确认，祝您旅途愉快！

感谢使用智能旅行预订系统！
"""
