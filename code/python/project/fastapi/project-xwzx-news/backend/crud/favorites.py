from sqlalchemy import select, func, delete
from sqlalchemy.ext.asyncio import AsyncSession

from models.favorites import Favorite
from models.news import News
from models.user import User


async def get_favorites_list(db: AsyncSession, user_data: User, offset: int = 0, page_size: int = 10):
    # 总量 + 收藏的新闻列表
    count_query = select(func.count()).where(Favorite.user_id == user_data.id)
    count_result = await db.execute(count_query)
    total = count_result.scalar_one()

    # 获取收藏列表 - 联表查询 join() + 收藏时间排序 + 分页
    # select(查询主体模型类, 字段别名).join(联合查询的模型类, 联合查询的条件).where().order_by().offset().limit()
    # 别名： Favorite.created_at.label("favorite_time")
    # [
    #   (新闻对象, 收藏时间, 收藏id)
    # ]
    query = (select(News, Favorite.created_at.label("favorite_time"), Favorite.id.label("favorite_id"))
             .join(Favorite, Favorite.news_id == News.id)
             .where(Favorite.user_id == user_data.id)
             .order_by(Favorite.created_at.desc())
             .offset(offset).limit(page_size)
             )
    result = await db.execute(query)
    rows = result.all()
    return rows, total


async def check_status(db: AsyncSession, news_id: int, user: User):
    query = select(Favorite).where(Favorite.user_id == user.id, Favorite.news_id == news_id)
    result = await db.execute(query)
    # 是否有收藏记录
    return result.scalar_one_or_none() is not None

async def add_news_to_favorite(db: AsyncSession, news_id: int, user: User):
    favorite = Favorite(user_id=user.id, news_id=news_id)
    db.add(favorite)
    await db.commit()
    await db.refresh(favorite)
    return favorite

async def cancel_favorite(db: AsyncSession, news_id: int, user: User):
    query = delete(Favorite).where(Favorite.user_id == user.id, Favorite.news_id == news_id)
    result = await db.execute(query)
    await db.commit()
    return result.rowcount > 0

async def clear_favorite(db: AsyncSession, user_id: int):
    query = delete(Favorite).where(Favorite.user_id == user_id)
    result = await db.execute(query)
    await db.commit()
    # 返回一个删除的数量
    return result.rowcount or 0