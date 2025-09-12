import { IsString, Validate } from "class-validator";
import { StartsWithConstraint, StartsWith, IsUsernameUnique, IsUsernameUniqueConstraint } from "../validators/user-validators";
import { ApiProperty, ApiPropertyOptional, PartialType } from "@nestjs/swagger"
import { IsOptionalString, IsOptionalEmail, IsOptionalNumber, IsOptionalBoolean } from "../decorators/alidation-and-transform.decorators";

export class CreateUserDto {
  @ApiProperty({ description: '用户名，必须唯一且以指定前缀开头', example: 'user_john_doe' })
  @IsString()
  @Validate(StartsWithConstraint, ['user_'], {
    message: `用户名必须以 "user_" 开头`,
  })
  @Validate(IsUsernameUniqueConstraint, { message: '用户名已存在' })
  // @StartsWith('user_', { message: '用户名必须以 "user_" 开头' })
  // @IsUsernameUnique({ message: '用户名已存在' })
  readonly username: string;

  @ApiProperty({ description: '密码', example: 'securePassword123' })
  @IsString()
  readonly password: string;

  @ApiPropertyOptional({ description: '手机号', example: '13124567890' })
  @IsOptionalString()
  readonly mobile?: string;

  @ApiPropertyOptional({ description: '邮箱地址', example: 'john.doe@example.com' })
  @IsOptionalEmail()
  readonly email?: string;

  @ApiPropertyOptional({ description: '用户状态', example: 1 })
  @IsOptionalNumber()
  readonly status?: number;

  @ApiPropertyOptional({ description: '是否为超级管理员', example: true })
  @IsOptionalBoolean()
  readonly is_super?: boolean;
}

export class UpdateUserDto extends PartialType(CreateUserDto) {
  @ApiProperty({ description: '用户ID', example: 1 })
  @IsOptionalNumber()
  readonly id: number;
}
