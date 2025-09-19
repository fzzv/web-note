import { ApiProperty, PartialType as PartialTypeFromSwagger } from '@nestjs/swagger';
import { IsString, IsInt, MaxLength } from 'class-validator';
import { PartialType } from '@nestjs/mapped-types';
import { IdValidators, StatusValidators, SortValidators } from '../decorators/alidation-and-transform.decorators';
import { Transform } from 'class-transformer';

export class CreateArticleDto {
  @ApiProperty({ description: '标题', example: '文章标题' })
  @IsString()
  @MaxLength(50, { message: '标题不能超过50个字符' })
  title: string;

  @ApiProperty({ description: '内容', example: '文章内容' })
  content: string;

  @ApiProperty({ description: '分类ID数组', example: [1, 2] })
  @Transform(({ value }) => Array.isArray(value) ? value.map(Number) : [Number(value)])
  @IsInt({ each: true })
  categoryIds: number[];

  @ApiProperty({ description: '标签ID数组', example: [1, 2] })
  @Transform(({ value }) => Array.isArray(value) ? value.map(Number) : [Number(value)])
  @IsInt({ each: true })
  tagIds: number[];

  @StatusValidators()
  @ApiProperty({ description: '状态', example: 1 })
  status: number;

  @SortValidators()
  @ApiProperty({ description: '排序号', example: 100 })
  sort: number;
}

export class UpdateArticleDto extends PartialTypeFromSwagger(PartialType(CreateArticleDto)) {
  @IdValidators()
  id: number;
}
