import { ApiProperty } from '@nestjs/swagger';
import { PartialType } from '@nestjs/mapped-types'
import { IdValidators, SortValidators, StatusValidators } from '../decorators/alidation-and-transform.decorators';

export class CreateRoleDto {
  @ApiProperty({ description: 'Name', example: 'role' })
  name: string;

  @StatusValidators()
  @ApiProperty({ description: '状态', example: 1 })
  status: number;

  @SortValidators()
  @ApiProperty({ description: '排序号', example: 100 })
  sort: number;

}

export class UpdateRoleDto extends PartialType(CreateRoleDto) {
  @IdValidators()
  id: number;
}

export class UpdateRoleAccessesDto {
  readonly accessIds: number[];
}
