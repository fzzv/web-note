// 导入 gRPC 库
const grpc = require("@grpc/grpc-js");
// 导入 proto 文件加载器，用于解析 .proto 文件
const protoLoader = require("@grpc/proto-loader");
// 同步加载 hello.proto 文件并解析为定义包
const packageDefinition = protoLoader.loadSync("hello.proto", {});
// 使用 gRPC 加载解析后的 proto 定义，生成服务对象
const helloProto = grpc.loadPackageDefinition(packageDefinition).hello;
// 定义 sayHello 方法，用于处理客户端的请求
function sayHello(call, callback) {
  // 从请求中获取客户端传递的 name 参数，并生成响应消息
  const replyMessage = `Hello, ${call.request.name}!`;
  // 通过回调返回响应消息
  callback(null, { message: replyMessage });
}

function sayHelloStream(call) {
  const names = ['Alice', 'Bob', 'Charlie'];
  names.forEach((name, index) => {
    setTimeout(() => {
      call.write({ message: `Hello, ${name}!` });
      if (index === names.length - 1) {
        call.end();
      }
    }, 1000 * (index + 1)); // 每隔 1 秒发送一个消息
  });
}
// 主函数，创建并启动 gRPC 服务器
function main() {
  // 创建一个 gRPC 服务器实例
  const server = new grpc.Server();
  // 注册 Greeter 服务，将 sayHello 方法绑定到服务中
  server.addService(helloProto.Greeter.service, {
    sayHello: sayHello,
    sayHelloStream: sayHelloStream
  });
  // 绑定服务器到指定的地址和端口，并启动服务器
  server.bindAsync(
    "0.0.0.0:50051",
    grpc.ServerCredentials.createInsecure(),
    (error, port) => {
      // 如果绑定失败，输出错误信息
      if (error) {
        console.error("绑定失败:", error);
        return;
      }
      // 成功启动服务器，打印启动信息
      console.log(`gRPC 服务器在端口 ${port} 启动`);
    }
  );
}
// 调用主函数，启动 gRPC 服务器
main();
