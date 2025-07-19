import { useContext } from "react";
import connContext from "./contexts/ConnContext";
import "./Form.css";
import nameContext from "./contexts/nameContext";

function Form({ setIsConnected, setName }) {
  const conn = useContext(connContext);
  const name = useContext(nameContext);

  function handleSubmit(event) {
    event.preventDefault();
    const userName = event.target.elements.username.value;

    name.current = userName;
    const origin= import.meta.env.VITE_BACKEND_ORIGIN;

  console.log("backen dis at ", origin.substring(4))

    // conn.current = new WebSocket(`ws://localhost:8000/ws?username=${encodeURIComponent(userName)}`);
    // const protocol = window.location.protocol === "https:" ? "wss" : "ws";
    // conn.current = new WebSocket(
    //   `${protocol}://${
    //     window.location.hostname
    //   }:8000/ws?username=${encodeURIComponent(userName)}`
    // );
    conn.current = new WebSocket(`ws${origin.substring(4)}/ws?username=${encodeURIComponent(userName)}`);

    conn.current.onopen = () => {
      //conn.current.send("hi");
      setIsConnected(true);
    };
    conn.current.onclose = () => {
      console.log("WebSocket connection closed");
    };
  }

  return (
    <div className="form-container">
      <div className="form-box">
        <h1 className="title">🚀 FileChat</h1>
        <p className="subtitle">Secure. Simple. Seamless.</p>
        <form onSubmit={handleSubmit} className="form">
          <input
            type="text"
            name="username"
            placeholder="Enter your name"
            required
          />
          <button type="submit">Join Chat</button>
        </form>
      </div>
    </div>
  );
}

export default Form;
