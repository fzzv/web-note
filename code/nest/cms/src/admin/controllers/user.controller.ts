import { Controller, Get } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';

@Controller('admin/user')
export class UserController {

  constructor(private readonly userService: UserService) {}

  @Get()
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
  }
}
