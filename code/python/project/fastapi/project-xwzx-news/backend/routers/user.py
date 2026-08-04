from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from starlette import status

from config.db_conf import get_db
from crud import user
from schemas.user import UserRequest, UserAuthResponse, UserInfoResponse
from utils.response import success_response

router = APIRouter(prefix="/api/user", tags=["user"])

@router.post("/register")
async def register_user(user_data: UserRequest, db: AsyncSession = Depends(get_db)):
    """
    用户注册
    :param user_data:
    :param db:
    :return:
    """
    # 检查用户名是否存在
    existing_user = await user.get_user_by_username(db, user_data.username)
    if existing_user:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="用户已存在")
    # 创建新用户
    users = await user.create_user(db, user_data)
    token = await user.create_token(db, users.id)
    response_data = UserAuthResponse(token=token, user_info=UserInfoResponse.model_validate(users))
    return success_response(message="注册成功", data=response_data)
