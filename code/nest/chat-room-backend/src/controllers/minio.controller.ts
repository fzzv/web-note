import { Controller, Get, Inject, Query } from '@nestjs/common';
import * as Minio from 'minio';

@Controller('minio')
export class MinioController {

  @Inject('MINIO_CLIENT')
  private minioClient: Minio.Client;

  @Get('presignedUrl')
  presignedPutObject(@Query('name') name: string) {
    // 第一个参数是 bucketName，第二个参数是 objectName，第三个参数是 expires
    return this.minioClient.presignedPutObject('chat-room', name, 3600);
  }
}
