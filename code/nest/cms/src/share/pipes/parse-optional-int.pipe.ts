import { Injectable, PipeTransform, ArgumentMetadata, BadRequestException } from '@nestjs/common';

/**
 * 解析可选的整数参数
 * 如果参数为空（undefined、null 或 ''），返回默认值
 * 如果参数不是有效整数，则抛出 400 错误
 * 否则返回解析后的整数
 */
@Injectable()
export class ParseOptionalIntPipe implements PipeTransform<string, number> {
  constructor(private readonly defaultValue: number) { }

  transform(value: string, metadata: ArgumentMetadata): number {
    // 1. 如果参数为空（undefined、null 或 ''），返回默认值
    if (!value) {
      return this.defaultValue;
    }

    // 2. 尝试解析为整数
    const parsedValue = parseInt(value, 10);

    // 3. 如果不是有效整数，则抛出 400 错误
    if (isNaN(parsedValue)) {
      throw new BadRequestException(`Validation failed. "${value}" is not an integer.`);
    }

    // 4. 否则返回解析后的整数
    return parsedValue;
  }
}
