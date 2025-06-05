import asyncio
import websockets

async def connect():
    uri = "ws://localhost:8080/ws?token=session_testuser" 
    async with websockets.connect(uri) as websocket:
        print("Connected to websocket server.")
        await websocket.send("Hello, server!")
        response = await websocket.recv()
        print(f"Received: {response}")

if __name__ == "__main__":
    asyncio.run(connect())