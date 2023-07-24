const dgram = require('dgram');

// 创建UDP服务器
const server = dgram.createSocket('udp4');

server.addListener('connect', () => {
	console.log("new client.");
})

// 监听消息接收事件
server.on('message', (msg, rinfo) => {
	const data = msg.toString();
	server.send(data, rinfo.port, rinfo.address);
	console.log(`${rinfo.address}:${rinfo.port} : ${data}`);
});

// 监听错误事件
server.on('error', (error) => {
	console.error(`UDP服务器错误：${error}`);
	server.close();
});

// 监听关闭事件
server.on('close', () => {
	console.log('UDP服务器已关闭');
});

// 绑定端口并监听
const port = 10000;
server.bind(port, () => {
	console.log(`正在监听端口 ${port}...`);
});

// 捕获 Ctrl+C 事件来关闭服务器
process.on('SIGINT', () => {
	server.close();
});