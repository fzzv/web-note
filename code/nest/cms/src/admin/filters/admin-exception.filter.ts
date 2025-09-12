import { ExceptionFilter, Catch, ArgumentsHost, HttpException, BadRequestException } from '@nestjs/common';
// 导入 express 的 Response 对象，用于构建 HTTP 响应
import { Response } from 'express';
// 使用 @Catch 装饰器捕获所有 HttpException 异常
@Catch(HttpException)
export class AdminExceptionFilter implements ExceptionFilter {
  // 实现 catch 方法，用于处理捕获的异常
  catch(exception: HttpException, host: ArgumentsHost) {
    // 获取当前 HTTP 请求上下文
    const ctx = host.switchToHttp();
    // 获取 HTTP 响应对象
    const response = ctx.getResponse<Response>();
    // 获取异常的 HTTP 状态码
    const status = exception.getStatus();
    // 初始化错误信息，默认为异常的消息
    let errorMessage = exception.message;
    // 如果异常是 BadRequestException 类型，进一步处理错误信息
    if (exception instanceof BadRequestException) {
      // 获取异常的响应体
      const responseBody: any = exception.getResponse();
      // 检查响应体是否是对象并且包含 message 属性
      if (typeof responseBody === 'object' && responseBody.message) {
        // 如果 message 是数组，则将其拼接成字符串，否则直接使用 message
        errorMessage = Array.isArray(responseBody.message)
          ? responseBody.message.join(', ')
          : responseBody.message;
      }
    }
    // 使用响应对象构建并发送错误页面，包含错误信息和重定向 URL
    response.status(status).render('error', {
      message: errorMessage,
      redirectUrl: ctx.getRequest().url,
    });
  }
}
