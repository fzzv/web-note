import { Injectable } from '@nestjs/common';
import { createTransport, Transporter } from 'nodemailer';
import { ConfigurationService } from 'src/service/configuration.service';

@Injectable()
export class EmailService {

  transporter: Transporter

  constructor(private readonly configurationService: ConfigurationService) {
    this.transporter = createTransport({
      host: this.configurationService.smtpHost!,
      port: this.configurationService.smtpPort!,
      secure: true, // 如果使用的是 587 端口，通常是明文传输并支持 STARTTLS。如果使用的是 465 端口，通常需要启用 secure: true 以使用 SSL/TLS 进行加密
      auth: {
        user: this.configurationService.smtpUser!,
        pass: this.configurationService.smtpPass!
      },
    });
  }

  async sendMail({ to, subject, html }) {
    await this.transporter.sendMail({
      from: {
        name: '聊天室',
        address: this.configurationService.smtpUser!
      },
      to,
      subject,
      html
    });
  }
}
