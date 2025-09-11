import { Injectable } from "@nestjs/common";
import { registerDecorator, ValidationArguments, ValidationOptions, ValidatorConstraint, ValidatorConstraintInterface, } from "class-validator";

// 定义一个自定义验证器，名为 'startsWith'，不需要异步验证
@ValidatorConstraint({ name: 'startsWith', async: false })
// 使用 Injectable 装饰器使这个类可被依赖注入
@Injectable()
// 定义 StartsWithConstraint 类并实现 ValidatorConstraintInterface 接口
export class StartsWithConstraint implements ValidatorConstraintInterface {
  // 定义验证逻辑，检查值是否以指定的前缀开头
  validate(value: any, args: ValidationArguments) {
    const [prefix] = args.constraints;
    return typeof value === 'string' && value.startsWith(prefix);
  }
  // 定义默认消息，当验证失败时返回的错误信息
  defaultMessage(args: ValidationArguments) {
    const [prefix] = args.constraints;
    return `${args.property} must start with ${prefix}`;
  }
}

// 定义一个自定义验证器，名为 'isUsernameUnique'，需要异步验证
@ValidatorConstraint({ name: 'isUsernameUnique', async: true })
// 使用 Injectable 装饰器使这个类可被依赖注入
@Injectable()
// 定义 IsUsernameUniqueConstraint 类并实现 ValidatorConstraintInterface 接口
export class IsUsernameUniqueConstraint implements ValidatorConstraintInterface {
  // 定义验证逻辑，检查用户名是否唯一
  async validate(value: any, args: ValidationArguments) {
    const existingUsernames = ['ADMIN', 'USER', 'GUEST']; // 模拟已存在的用户名列表
    return !existingUsernames.includes(value);
  }
  // 定义默认消息，当验证失败时返回的错误信息
  defaultMessage(args: ValidationArguments) {
    return `${args.property} must be unique`;
  }
}

// 创建 StartsWith 装饰器工厂函数，用于给属性添加 'startsWith' 验证逻辑
export function StartsWith(prefix: string, validationOptions?: ValidationOptions) {
  return function (object: Object, propertyName: string) {
    registerDecorator({
      target: object.constructor, // 目标类
      propertyName: propertyName, // 目标属性名
      options: validationOptions, // 验证选项
      constraints: [prefix], // 传递给验证器的参数，如前缀
      validator: StartsWithConstraint, // 指定使用的验证器类
    });
  };
}

// 创建 IsUsernameUnique 装饰器工厂函数，用于给属性添加 'isUsernameUnique' 验证逻辑
export function IsUsernameUnique(validationOptions?: ValidationOptions) {
  return function (object: Object, propertyName: string) {
    registerDecorator({
      target: object.constructor, // 目标类
      propertyName: propertyName, // 目标属性名
      options: validationOptions, // 验证选项
      constraints: [], // 传递给验证器的参数，这里不需要
      validator: IsUsernameUniqueConstraint, // 指定使用的验证器类
    });
  };
}
