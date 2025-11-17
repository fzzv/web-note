import { PrismaService } from "src/prisma/prisma.service";

export class ChatService {
  constructor(private readonly prismaService: PrismaService) { }
}
