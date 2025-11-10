import { Controller } from '@nestjs/common';
import { EmailService } from 'src/service/email.service';

@Controller('email')
export class EmailController {
  constructor(private readonly emailService: EmailService) {}
}
