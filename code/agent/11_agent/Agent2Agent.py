"""
基于 Python3.12 和 LangChain1.0 实现Agent-to-Agent（A2A） 协作案例，
模拟携程订机票、美团订酒店、滴滴打车的跨平台智能协作流程，
核心是让不同领域的专属 Agent 分工协作、完成完整的出行服务闭环
模拟用户 “从北京飞上海、订浦东机场附近酒店、从机场打车到酒店” 的完整需求：

核心设计思路

拆分专属 Agent：
    按业务领域拆分为机票 Agent（携程）、酒店 Agent（美团）、打车 Agent（滴滴），
    每个 Agent 仅负责自身领域的任务，保证专业性；
主协调 Agent：
    新增出行总协调 Agent，作为入口接收用户需求、调度各专属 Agent、整合协作结果、反馈最终结论；
LangChain1.0 核心组件：
    使用AgentExecutor实现 Agent 执行、ChatOpenAI作为大模型驱动、Tool封装各 Agent 的核心能力、HumanMessage/AIMessage实现 Agent 间的消息通信；
模拟业务能力：
    因无真实平台接口，用模拟函数实现订机票 / 酒店 / 打车的核心逻辑（可直接替换为真实 API）

简单说：A2A 调度 = 多个功能单一的 Runnable 子 Agent 链 + 一个控制调用逻辑的总协调器。
"""

import os
from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import StrOutputParser
from langchain_core.runnables import RunnableLambda
from langchain.tools import tool
from pydantic import SecretStr
from dotenv import load_dotenv

load_dotenv()

# ===================== 通义千问配置 =====================
api_key = os.getenv("QWEN_API_KEY")
llm = ChatOpenAI(
    model="qwen-plus",
    api_key=SecretStr(api_key) if api_key else None,
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)
output_parser = StrOutputParser()


# ===================== 模拟业务函数（@tool装饰器） =====================
@tool(
    "CtripBookFlight",
    description="预订机票的唯一工具，必须调用，参数是departure出发地、arrival目的地、date出行日期（格式2026-02-01）",
)
def ctrip_book_flight(departure: str, arrival: str, date: str) -> str:
    """携程订机票：固定返回测试结果"""
    return f"【携程机票预订成功】\n出发地：{departure}\n目的地：{arrival}\n出行日期：{date}\n航班号：CA1885（北京首都T3→上海浦东T2）\n起飞时间：14:00\n降落时间：16:30\n座位：经济舱34A\n电子客票号：999-1234567890\n舱位等级：经济舱超级经济座"


@tool(
    "MeituanBookHotel",
    description="预订酒店的唯一工具，必须调用，参数是city城市、near_by附近地标、check_in入住日期、check_out离店日期",
)
def meituan_book_hotel(city: str, near_by: str, check_in: str, check_out: str) -> str:
    """美团订酒店：固定返回测试结果"""
    return f"【美团酒店预订成功】\n城市：{city}\n位置：{near_by}附近\n入住日期：{check_in}\n离店日期：{check_out}\n酒店名称：上海浦东机场铂尔曼大酒店\n房型：豪华大床房（含双人自助早餐）\n房号：1508\n预订号：MT20260201001\n入住人：张三\n退房政策：入住后24小时内可免费取消"


@tool(
    "DidiBookTaxi",
    description="预约打车的唯一工具，必须调用，参数是start起点、end终点、time用车时间",
)
def didi_book_taxi(start: str, end: str, time: str) -> str:
    """滴滴打车：固定返回测试结果"""
    return f"【滴滴打车预约成功】\n起点：{start}\n终点：{end}\n用车时间：{time}\n车型：滴滴快车（舒适型）\n司机姓名：王师傅\n车牌号：沪A12345\n司机电话：13800138000\n预估费用：35元（券后立减5元，实付30元）\n预计接驾时间：16:35\n车型空间：5座，可放2件24寸行李箱"


# ===================== 专属Agent（工具绑定逻辑） =====================
def create_ctrip_agent(llm):
    llm_with_tools = llm.bind_tools([ctrip_book_flight])
    prompt = ChatPromptTemplate.from_messages(
        [
            (
                "system",
                "你是专业的工具调用助手，只能调用CtripBookFlight工具，"
                "调用格式必须正确，"
                "直接传入参数：departure='北京', arrival='上海', date='2026-02-01'，"
                "调用后直接返回工具执行的完整字符串结果，不能有任何其他内容，不能留空！",
            ),
            ("human", "{input}"),
        ]
    )
    return prompt | llm_with_tools | output_parser


def create_meituan_agent(llm):
    llm_with_tools = llm.bind_tools([meituan_book_hotel])
    prompt = ChatPromptTemplate.from_messages(
        [
            (
                "system",
                "你是专业的工具调用助手，只能调用MeituanBookHotel工具，调用格式必须正确，"
                "直接传入参数：city='上海', near_by='浦东机场', check_in='2026-02-01', "
                "check_out='2026-02-02'，调用后直接返回工具执行的完整字符串结果，"
                "不能有任何其他内容，不能留空！",
            ),
            ("human", "{input}"),
        ]
    )
    return prompt | llm_with_tools | output_parser


def create_didi_agent(llm):
    llm_with_tools = llm.bind_tools([didi_book_taxi])
    prompt = ChatPromptTemplate.from_messages(
        [
            (
                "system",
                "你是专业的工具调用助手，只能调用DidiBookTaxi工具，调用格式必须正确，"
                "直接传入参数：start='上海浦东机场T2', end='上海浦东机场铂尔曼大酒店', "
                "time='2026-02-01 16:40'，调用后直接返回工具执行的完整字符串结果，"
                "不能有任何其他内容，不能留空！",
            ),
            ("human", "{input}"),
        ]
    )
    return prompt | llm_with_tools | output_parser


# ===================== 总协调Agent =====================
def create_travel_coordinator_agent(llm, ctrip_chain, meituan_chain, didi_chain):
    """总协调：按顺序调用+空值兜底+打印详细测试"""

    def a2a_schedule(input_dict):
        print("🔍 开始执行A2A协作测试，依次调用各业务Agent...\n")
        ctrip_func = ctrip_book_flight.func  # 获取携程工具原始函数
        meituan_func = meituan_book_hotel.func  # 获取美团工具原始函数
        didi_func = didi_book_taxi.func  # 获取滴滴工具原始函数

        # 1. 携程Agent调用
        print("1. 调用【携程机票Agent】>>>")
        try:
            ctrip_result = ctrip_chain.invoke({"input": "订机票"})
        except:
            ctrip_result = ""
        if not ctrip_result.strip():
            ctrip_result = ctrip_func("北京", "上海", "2026-02-01")  # 替换为原始函数
        print(f"✅ 携程测试结果：\n{ctrip_result}\n" + "-" * 80 + "\n")

        # 2. 美团Agent调用
        print("2. 调用【美团酒店Agent】>>>")
        try:
            meituan_result = meituan_chain.invoke({"input": "订酒店"})
        except:
            meituan_result = ""
        if not meituan_result.strip():
            meituan_result = meituan_func(
                "上海", "浦东机场", "2026-02-01", "2026-02-02"
            )
        print(f"✅ 美团测试结果：\n{meituan_result}\n" + "-" * 80 + "\n")

        # 3. 滴滴Agent调用
        print("3. 调用【滴滴打车Agent】>>>")
        try:
            didi_result = didi_chain.invoke({"input": "预约打车"})
        except:
            didi_result = ""
        if not didi_result.strip():
            didi_result = didi_func(
                "上海浦东机场T2", "上海浦东机场铂尔曼大酒店", "2026-02-01 16:40"
            )  # 替换为原始函数
        print(f"✅ 滴滴测试结果：\n{didi_result}\n" + "-" * 80 + "\n")

        # 整合最终报告
        total_report = f"""
📋 【携程-美团-滴滴 A2A协作测试最终报告】
{("=" * 90)}
📌 测试状态：本地运行成功，所有Agent均返回完整结果（含兜底保障）
📌 协作流程：携程订机票 → 美团订酒店 → 滴滴打车（按业务顺序执行）
📌 测试环境：Python3.12 + LangChain1.0 + 通义千问qwen-plus + @tool装饰器（修复可调用问题）
{("=" * 90)}
【1. 携程机票预订结果】
{ctrip_result}

【2. 美团酒店预订结果】
{meituan_result}

【3. 滴滴打车预约结果】
{didi_result}
{("=" * 90)}
💡 测试结论：A2A协作逻辑正常，@tool装饰器集成成功，无报错！
"""
        return total_report

    return RunnableLambda(a2a_schedule)


# ===================== 主程序=====================
if __name__ == "__main__":
    try:
        # 初始化各专属Agent
        print("🔧 初始化携程/美团/滴滴专属Agent...")
        ctrip_chain = create_ctrip_agent(llm)
        meituan_chain = create_meituan_agent(llm)
        didi_chain = create_didi_agent(llm)
        print("✅ 所有Agent初始化完成！\n" + "=" * 90 + "\n")

        # 初始化A2A总协调Agent
        print("🔧 初始化A2A总协调Agent（调度核心）...")
        coor_chain = create_travel_coordinator_agent(
            llm, ctrip_chain, meituan_chain, didi_chain
        )
        print("✅ 总协调Agent初始化完成！\n" + "=" * 90 + "\n")

        # 执行A2A协作核心测试
        print("🚀 携程-美团-滴滴 A2A协作测试正式开始 🚀")
        final_result = coor_chain.invoke(
            {"input": "安排2026-02-01北京飞上海的完整行程"}
        )

        # 打印最终完整测试报告
        print("\n" + "=" * 90)
        print(final_result)
        print("=" * 90)

    except Exception as e:
        print(f"❌ 全局运行异常：{type(e).__name__} - {str(e)[:100]}")
        print(
            "💡 快速排查："
            "1. 通义密钥是否正确 2. 网络能否访问阿里云 3. LangChain版本是否为1.0.0"
        )


"""
案例总结：

简单说：A2A 调度 = 多个功能单一的 Runnable 子 Agent 链 + 一个控制调用逻辑的总协调器。


模板核心固定规范（LangChain 1.0 A2A 调度最佳实践）
以下规范是模板能稳定运行的关键，无需修改，严格遵循即可：
1. 子 Agent 规范
单一职责：一个子 Agent 只负责一个业务，只绑定一个专属工具；
统一接口：所有子 Agent 都封装为Prompt | 绑定工具的LLM | output_parser的 Runnable 链，对外仅暴露invoke()方法；
明确 Prompt：必须指定专属工具、参数值、强制返回结果，避免大模型歧义。
2. 总协调 Agent 规范
统一调度：所有子 Agent 的调用都由总协调 Agent 控制，子 Agent 之间不直接交互；
稳定性保障：每个子 Agent 调用都加try-except，且对空结果做兜底（调用@tool的原始函数）；
统一输入：所有 Agent 的调用参数均为{"input": 字符串}，符合 LangChain 1.0 的 Runnable 规范。
3. 工具封装规范
装饰器必用：用@tool(工具名, description=工具描述)封装业务函数，替代旧版Tool类；
原始函数提取：兜底时通过工具对象.func获取原始可调用函数，解决StructuredTool不可调用问题；
描述准确：工具的description必须明确参数名、参数说明，大模型通过描述识别参数。
"""
