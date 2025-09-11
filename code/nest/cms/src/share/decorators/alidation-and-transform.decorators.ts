import { applyDecorators } from "@nestjs/common";
import { Type } from "class-transformer";
import { IsBoolean, IsEmail, IsNumber, IsOptional, IsString } from "class-validator";

// 可选字符串
export function IsOptionalString() {
  return applyDecorators(IsOptional(), IsString())
}

// 可选邮箱
export function IsOptionalEmail() {
  return applyDecorators(IsOptional(), IsEmail())
}

// 可选数字 并转换为数字
export function IsOptionalNumber() {
  return applyDecorators(IsOptional(), IsNumber(), Type(() => Number))
}

// 可选布尔值 并转换为布尔值
export function IsOptionalBoolean() {
  return applyDecorators(IsOptional(), IsBoolean(), Type(() => Boolean))
}
