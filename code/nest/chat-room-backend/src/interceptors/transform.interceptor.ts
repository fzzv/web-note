import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ApiResponse } from '../interfaces/api-response.interface';

@Injectable()
export class TransformInterceptor<T> implements NestInterceptor<T, ApiResponse<T>> {
  intercept(context: ExecutionContext, next: CallHandler): Observable<ApiResponse<T>> {

    // 获取 HTTP 状态码
    const httpStatusCode = context.switchToHttp().getResponse().statusCode;

    return next.handle().pipe(
      map(data => ({
        success: true,
        code: httpStatusCode, // 使用 HTTP 状态码
        data: data, // data 是你 controller/service 返回的原始数据
        message: data ? 'Request successful' : 'Request failed',
      })),
    );
  }
}
