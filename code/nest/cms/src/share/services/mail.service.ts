import { Injectable } from '@nestjs/common';
import nodemailer from 'nodemailer';
import { ConfigService } from '@nestjs/config';
@Injectable()
export class MailService {
  private transporter;
  constructor(private readonly configService: ConfigService) {
    this.transporter = nodemailer.createTransport({
      host: configService.get('SMTP_HOST'),
      port: configService.get('SMTP_PORT'),
      secure: true,
      auth: {
        user: configService.get('SMTP_USER'),
        pass: configService.get('SMTP_PASS'),
      },
    });
  }

  async sendEmail(to: string, subject: string, body: string) {
    const mailOptions = {
      from: this.configService.get('SMTP_USER'), // 发件人
      to, // 收件人
      subject, // 主题
      text: body, // 邮件正文
    };
    try {
      const info = await this.transporter.sendMail(mailOptions);
      console.log(`邮件已发送: ${info.messageId}`);
    } catch (error) {
      console.error(`发送邮件失败: ${error.message}`);
    }
  }
}
