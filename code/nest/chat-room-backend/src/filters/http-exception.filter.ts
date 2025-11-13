import { ExceptionFilter, Catch, ArgumentsHost, HttpException, HttpStatus } from '@nestjs/common';
import { Request, Response } from 'express';
import { PrismaClientKnownRequestError } from '@prisma/client/runtime/library';

@Catch() // 捕获所有类型的异常
export class AllExceptionsFilter implements ExceptionFilter {
  catch(exception: unknown, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();

    let status = HttpStatus.INTERNAL_SERVER_ERROR;
    let message = 'Internal server error';
    let code = 500;

    // 1. 处理 Nest.js 的 HttpException
    if (exception instanceof HttpException) {
      status = exception.getStatus();
      message = exception.message;
      code = status;
    }
    // 2. 处理 Prisma 特定的已知请求错误
    else if (exception instanceof PrismaClientKnownRequestError) {
      message = `[Prisma Error ${exception.code}] ${exception.message.split('\n').pop()}`; // 简化 Prisma 错误信息

      switch (exception.code) {
        case 'P2002': // 唯一约束失败 (e.g., email already exists)
          status = HttpStatus.CONFLICT; // 409
          message = 'A record with this value already exists.';
          break;
        case 'P2025': // 记录未找到 (e.g., for update/delete)
          status = HttpStatus.NOT_FOUND; // 404
          message = 'Record to update/delete not found.';
          break;
        // ... 在这里添加其他 Prisma 错误码
        default:
          status = HttpStatus.BAD_REQUEST; // 400
          break;
      }
      code = status; // 也可以使用 Prisma code 作为业务码
    }
    // 3. 处理其他未知错误
    else if (exception instanceof Error) {
      message = exception.message;
    }

    response.status(status).json({
      success: false,
      code: code,
      data: null,
      message: message,
      timestamp: new Date().toISOString(), // 额外信息
      path: request.url, // 额外信息
    });
  }
}
