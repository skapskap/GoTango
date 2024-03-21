import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";

export default function Chat() {
  const [messages, setMessages] = useState<string[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    const newSocket = new WebSocket("ws://localhost:4869/ws");

    newSocket.onmessage = (event) => {
      const message = event.data;
      setMessages((prevMessages) => [...prevMessages, message as string]);
    };

    setSocket(newSocket);

    return () => {
      newSocket.close();
    };
  }, []);

  const handleSendMessage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    sendMessage();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const sendMessage = () => {
    if (inputMessage.trim() !== "" && socket) {
      setMessages((prevMessages) => [...prevMessages, inputMessage]);
      socket.send(inputMessage);
      setInputMessage("");
    }
  };

  return (
    <div className="flex items-center justify-center">
      <div className="flex flex-col h-[80vh] w-[60vh]">
        <header className="border-b p-4 flex items-center mt-24">
          <div className="flex-1">
            <h1 className="text-xl font-bold">DOLLARS</h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Chatting with
              <span className="font-medium"> 田中太郎</span>
            </p>
          </div>
        </header>
        <main className="flex-1 flex flex-col p-4">
          <div className="space-y-4 flex-1">
            {messages.map((message, index) => (
              <div
                key={index}
                className={`flex items-${
                  index % 2 === 0 ? "start" : "end"
                } space-x-2 justify-${index % 2 === 0 ? "start" : "end"}`}
              >
                <div
                  className={`${
                    index % 2 === 0
                      ? "bg-gray-100 dark:bg-gray-800"
                      : "bg-gray-200 dark:bg-gray-700"
                  } rounded-lg p-4 max-w-[70%] ${
                    index % 2 === 0 ? "" : "text-right"
                  }`}
                >
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    <time>{new Date().toLocaleTimeString()}</time>
                  </p>
                  <p>{message}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    Sent by: {index % 2 === 0 ? "田中太郎" : "セットン"}
                  </p>
                </div>
              </div>
            ))}
          </div>
          <div className="mt-4">
            <form
              onSubmit={handleSendMessage}
              className="flex items-center space-x-4"
            >
              <Textarea
                className="max-h-[200px] flex-1"
                placeholder="Type a message..."
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyDown={handleKeyDown}
              />
              <Button type="submit">Send</Button>
            </form>
          </div>
        </main>
      </div>
    </div>
  );
}
