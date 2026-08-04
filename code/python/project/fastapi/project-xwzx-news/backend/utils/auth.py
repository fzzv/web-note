from fastapi import Header, HTTPException, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from starlette import status

from config.db_conf import get_db
from crud import user


async def get_current_user(authorization: str = Header(..., alias="Authorization"), db: AsyncSession = Depends(get_db)):
    token = authorization.replace("Bearer ", "")
    users = await user.get_user_by_token(db, token)
    if not users:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="无效令牌或令牌失效")
    return users
