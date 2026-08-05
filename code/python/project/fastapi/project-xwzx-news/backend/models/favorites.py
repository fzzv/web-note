from datetime import datetime

from sqlalchemy import Integer, ForeignKey, DateTime, UniqueConstraint, Index
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column

from models.news import News
from models.user import User


class Base(DeclarativeBase):
    pass

class Favorite(Base):
    __tablename__ = "favorite"

    # 创建索引
    # UniqueConstraint: 唯一约束, 当前用户，当前新闻，只能收藏一次
    __table_args__ = (
        UniqueConstraint('user_id', 'news_id', name='user_news_unique'),
        Index('fk_favorite_user_idx', 'user_id'),
        Index('fk_favorite_news_idx', 'news_id'),
    )

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True, comment="收藏id")
    user_id: Mapped[int] = mapped_column(Integer, ForeignKey(User.id), nullable=False, comment="用户id")
    news_id: Mapped[int] = mapped_column(Integer, ForeignKey(News.id), nullable=False, comment="新闻id")
    created_at: Mapped[datetime] = mapped_column(DateTime, default=datetime.now(), nullable=False, comment="收藏时间")

    def __repr__(self):
        return f"<Favorites(id={self.id}, user_id={self.user_id}, news_id={self.news_id}, created_at={self.created_at})>"