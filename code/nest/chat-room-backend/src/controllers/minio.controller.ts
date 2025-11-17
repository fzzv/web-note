import { Controller, Get, Inject, Query } from '@nestjs/common';
import * as Minio from 'minio';
import { ApiOperation, ApiQuery, ApiResponse, ApiTags } from '@nestjs/swagger';

@ApiTags('minio')
@Controller('minio')
export class MinioController {

  @Inject('MINIO_CLIENT')
  private minioClient: Minio.Client;

  @Get('presignedUrl')
  @ApiOperation({ summary: '获取预签名上传 URL' })
  @ApiQuery({ name: 'name', description: '文件名称', required: true })
  @ApiResponse({ status: 200, description: '成功返回预签名上传 URL' })
  presignedPutObject(@Query('name') name: string) {
    // 第一个参数是 bucketName，第二个参数是 objectName，第三个参数是 expires
    return this.minioClient.presignedPutObject('chat-room', name, 3600);
  }
}
