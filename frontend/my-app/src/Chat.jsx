import { useEffect, useRef, useContext, useState } from "react";
import connContext from "./contexts/ConnContext";
import nameContext from "./contexts/nameContext";
import "./Chat.css"; // contains styles for chat-container, chat-bubble etc.

function Chat() {
  const conn = useContext(connContext);
  const logRef = useRef(null);
  const bottomRef = useRef(null);
  const userName = useContext(nameContext);
  const origin= import.meta.env.VITE_BACKEND_ORIGIN;

  console.log("backen dis at ", origin)

  const [messages, setMessages] = useState([]);

  useEffect(() => {
    document.getElementById("fileBtn").onclick = function () {
      const file = document.getElementById("fileInput").files[0];
      if (!file) {
        alert("Please select a file first");
        return;
      }

      const formData = new FormData();
      formData.append("upload", file); // your backend should expect "upload" field
      formData.append("userName", userName.current);

      fetch(`${origin}/file`, {
        method: "POST",
        body: formData,
      })
        .then((res) => res.text())
        .then((msg) => {
          const div = document.createElement("div");
          const data = JSON.parse(msg);
          console.log("decoded file msg is ", data);
          div.className =
            msg.sender === userName ? "chat-bubble.received" : "chat-bubble";
          div.innerHTML = `
        <a href="${msg}" download style="color: #000;">
          📎 Download File
        </a>
      `;
          document.getElementById("log").appendChild(div);
        })
        .catch((err) => {
          const errDiv = document.createElement("div");
          errDiv.innerText = "Upload failed: " + err;
          console.log("error is",err);
          document.getElementById("log").appendChild(errDiv);
        });
    };
    document.getElementById("fileInput").value = null;
    if (!conn.current) return;

    conn.current.onmessage = (event) => {
      console.log(userName);
      const data = JSON.parse(event.data);
      console.log("decoded msg data is ", data, event.data);
      const isFile =
        data.message.includes("https://") || data.message.includes("s3");

      setMessages((prev) => [
        ...prev,
        {
          type: isFile ? "file" : "text",
          content: data.message,
          isMine: userName.current === data.sender, // You can enhance this later using sender info
          sender: data.sender,
        },
      ]);
    };
  }, []);

  function handleSubmit(event) {
    event.preventDefault();
    const message = event.target.message.value;
    if (
      conn.current instanceof WebSocket &&
      conn.current.readyState === WebSocket.OPEN
    ) {
      conn.current.send(message);
    }

    event.target.reset(); // Clear input
  }

  return (
    <div className="chat-container">
      <div className="chat-log" id="log" ref={logRef}>
        {messages.map((msg, index) => (
          <div
            key={index}
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: msg.isMine ? "flex-end" : "flex-start",
            }}
          >
            <div
              style={{
                fontSize: "0.75rem",
                color: "#555",
                marginBottom: "2px",
                marginRight: msg.isMine ? "4px" : "auto",
              }}
            >
              {msg.isMine ? "You" : msg.sender}
            </div>

            <div className={`chat-bubble ${msg.isMine ? "" : "received"}`}>
              {msg.type === "file" ? (
                <div className="file-message">
                  <div className="file-details">
                    <span>{msg.content.split("/").pop().split("?")[0]}</span>
                  </div>
                  <a
                    href={msg.content}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="download-btn"
                    download
                  >
                    Download
                  </a>
                </div>
              ) : (
                msg.content
              )}
            </div>
          </div>
        ))}

        <div ref={bottomRef} />
      </div>

      <form className="chat-form" onSubmit={handleSubmit}>
        <input type="text" name="message" placeholder="Type a message..." />
        <button type="submit">Send</button>
      </form>

      <div className="upload-section">
        <label className="custom-file-upload">
          📎 Choose File
          <input
            type="file"
            id="fileInput"
            onChange={(e) => {
              const fileLabel = document.getElementById("file-name");
              fileLabel.textContent =
                e.target.files[0]?.name || "No file selected";
            }}
          />
        </label>
        <span id="file-name">No file selected</span>
        <button type="button" id="fileBtn">
          Send File
        </button>
      </div>
    </div>
  );
}

export default Chat;
