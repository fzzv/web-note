import { Controller, Get } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(private readonly userService: UserService) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }
}
