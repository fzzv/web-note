// 导入 gRPC 库
const grpc = require("@grpc/grpc-js");
// 导入 proto 文件加载器，用于解析 .proto 文件
const protoLoader = require("@grpc/proto-loader");
// 同步加载 hello.proto 文件并解析为定义包
const packageDefinition = protoLoader.loadSync("hello.proto", {});
// 使用 gRPC 加载解析后的 proto 定义，生成服务对象
const helloProto = grpc.loadPackageDefinition(packageDefinition).hello;
// 主函数，创建并调用 gRPC 客户端
function main() {
  // 创建一个 gRPC 客户端，连接到本地的 gRPC 服务器 'localhost:50051'
  // 使用不加密的凭据进行通信（createInsecure）
  const client = new helloProto.Greeter(
    "localhost:50051",
    grpc.credentials.createInsecure()
  );
  // 调用 sayHello 方法，发送一个包含 name 字段的请求对象
  client.sayHello({ name: "World" }, (error, response) => {
    // 如果调用过程中发生错误，打印错误信息
    if (error) {
      console.error("调用服务出错:", error);
    } else {
      // 如果没有错误，打印服务端返回的响应消息
      console.log("服务响应:", response.message);
    }
  });
  client.sayHelloStream({ name: 'World' }).on('data', (response) => {
    console.log('服务流响应:', response.message);
  });
}
// 调用主函数，启动客户端并调用服务
main();
