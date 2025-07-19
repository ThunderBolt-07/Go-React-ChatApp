import { useState, useRef } from "react";
import connContext from "./contexts/ConnContext.jsx";
import nameContext from "./contexts/nameContext.jsx";
import Form from "./Form";
import Chat from "./Chat.jsx";
// import reactLogo from './assets/react.svg'
// import viteLogo from '/vite.svg'
// import './App.css'

function App() {
  const c = useRef(null);
  const n = useRef(null);

  const [isConnected, setIsConnected] = useState(false);
  const [name, setName] = useState("");

  return (
    <nameContext.Provider value={n}>
      <connContext.Provider value={c}>
        {isConnected ? (
          <Chat />
        ) : (
          <Form setIsConnected={setIsConnected} setName={setName} />
        )}
      </connContext.Provider>
    </nameContext.Provider>
  );
}

export default App;
