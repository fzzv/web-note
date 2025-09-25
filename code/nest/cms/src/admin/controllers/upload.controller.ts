import { Controller, Get, Post, Query, UploadedFile, UseInterceptors } from '@nestjs/common';
// 导入文件上传拦截器
import { FileInterceptor } from '@nestjs/platform-express';
// 导入multer的磁盘存储配置
import { diskStorage } from 'multer';
// 使用Node内置的randomUUID生成唯一文件名，避免ESM/CJS兼容问题
import { randomUUID } from 'crypto';
// 导入Node.js路径处理模块
import path from 'path';
import sharp from 'sharp';
import fs from 'fs';

/**
 * 文件上传控制器
 * 负责处理管理后台的文件上传功能
 * 支持图片文件上传，包括jpg、jpeg、png、gif格式
 */
@Controller('admin')
export class UploadController {

  /**
   * 文件上传接口
   * POST /admin/upload
   * 
   * 功能说明：
   * 1. 接收客户端上传的文件
   * 2. 验证文件类型（仅支持图片格式）
   * 3. 生成唯一文件名避免冲突
   * 4. 将文件保存到服务器磁盘
   * 5. 返回文件访问URL
   * 
   * @param file 上传的文件对象，包含文件信息和元数据
   * @returns 返回包含文件访问URL的响应对象
   */
  @Post('upload')
  @UseInterceptors(FileInterceptor('upload', {
    // 配置文件存储方式为磁盘存储
    storage: diskStorage({
      // 设置文件保存目录为项目根目录下的uploads文件夹
      destination: './uploads',
      // 自定义文件名生成规则
      filename: (_req, file, callback) => {
        // 使用Node内置的randomUUID生成唯一标识符，保留原文件扩展名
        // 这样可以避免文件名冲突，同时保持文件类型信息
        const filename: string = randomUUID() + path.extname(file.originalname);
        callback(null, filename);
      }
    }),
    // 文件类型过滤器，只允许特定格式的图片文件
    fileFilter: (req, file, callback) => {
      // 使用正则表达式验证MIME类型
      // 只允许jpg、jpeg、png、gif格式的图片文件
      if (!file.mimetype.match(/\/(jpg|jpeg|png|gif)$/)) {
        // 如果文件类型不支持，返回错误信息
        return callback(new Error('不支持的文件类型'), false);
      }
      // 文件类型验证通过，允许上传
      callback(null, true);
    }
  }))
  async uploadFile(@UploadedFile() file: Express.Multer.File) {
    // 生成压缩后的文件名，扩展名为 .min.jpeg
    const filename = `${path.basename(file.filename, path.extname(file.filename))}.min.jpeg`;
    // 压缩后的文件路径
    const outputFilePath = path.resolve('./uploads', filename);
    // 先读入 buffer，避免 sharp 占用源文件句柄
    const buffer = await fs.promises.readFile(file.path);
    // 使用 sharp 压缩
    await sharp(buffer)
      .resize(800, 600, {
        fit: sharp.fit.inside,
        withoutEnlargement: true,
      })
      .toFormat('jpeg')
      .jpeg({ quality: 80 })
      .toFile(outputFilePath);
    // safe unlink（删除原始上传文件）
    try {
      await fs.promises.unlink(file.path);
    } catch (err) {
      console.warn(`⚠️ 删除原文件失败: ${file.path}`, err);
    }
    // 返回压缩后的 URL
    return { url: `/uploads/${filename}` };
  }
}
