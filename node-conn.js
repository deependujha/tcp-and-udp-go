const net = require("net");

function nodeConnector() {
    let client = new net.Socket(); // Create a new socket for each connection

    client.connect(6969, "localhost", () => {
        console.log("Connected");
        let rawHex = "hello world from node";
        client.write(rawHex); // This will send the byte buffer over TCP
    });

    client.on('error', (err) => {
        console.error(`Connection error: ${err}`);
    });

    client.on('close', () => {
        console.log('Connection closed');
    });
}

// Loop to create 1000 connections
for (let i = 0; i < 1000; i++) {
    nodeConnector();
}

