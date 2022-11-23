using Grpc.Core;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace grpctest_csharp
{
    internal class Program
    {
        static void Main(string[] args)
        {
            Console.Title="C# gRPC客户端";
            var channel = new Channel("127.0.0.1:9091",ChannelCredentials.Insecure);
            var client=new TestService.TestServiceClient(channel);
            var callEachotherContext = client.CallEachOther();
            Task.Factory.StartNew(async () =>
            {
                while (await callEachotherContext.ResponseStream.MoveNext())
                {
                    Console.WriteLine("接收：{0}", callEachotherContext.ResponseStream.Current.Data);
                }
            });
            while (true)
            {
                Console.WriteLine("请输入要发送的内容");
                var str = Console.ReadLine();
                callEachotherContext.RequestStream.WriteAsync(new CallRequest() { Data = str }).Wait();
            }
        }
    }
}
