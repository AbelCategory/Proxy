import socks
import socket
import time
import asyncio


def send_udp_data_via_proxy(proxy_host, proxy_port, target_host, target_port, data):
    # 创建一个UDP代理socket对象
    proxy_socket = socks.socksocket(socket.AF_INET, socket.SOCK_DGRAM)
    # 设置代理服务器信息
    print(target_port)
    proxy_socket.set_proxy(socks.SOCKS5, proxy_host, proxy_port)

    # 连接到目标服务器
    proxy_socket.connect((target_host, target_port))

    print(proxy_socket.getpeername())
    # # 接收响应数据

    def recieve():
        response = proxy_socket.recv(1024).decode()
        print(response)
    # 发送数据

    def send():
        for i in range(3):
            print("???")
            proxy_socket.send((data + " " + str(i)).encode())
            # response = proxy_socket.recv(1024).decode()
            # print(response)
            recieve()
            print("!!!")
            time.sleep(1)
    send()
    # 关闭连接
    proxy_socket.close()


def main():
    # 使用代理发送UDP数据
    proxy_host = '127.0.0.1'  # 代理服务器IP地址
    proxy_port = 8080  # 代理服务器端口
    target_host = '127.0.0.1'  # 目标服务器IP地址
    target_port = 1234  # 目标服务器端口
    data = 'Hello, World!'  # 要发送的数据

    send_udp_data_via_proxy(
        proxy_host, proxy_port, target_host, target_port, data)
    # print(response)

# asyncio.get_event_loop().run_until_complete(asyncio.wait([asyncio.ensure_future(main())]))
main()