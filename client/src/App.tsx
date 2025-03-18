import { BrowserRouter, Route, Routes } from "react-router-dom";
import Register from "./components/Register";
import Login from "./components/Login";
import ConnectedUserRoute from "./components/middlewares/connectedUserRoute";
import { ChatWindow } from "./components/ChatWindow";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={
            <ConnectedUserRoute>
              <ChatWindow />
            </ConnectedUserRoute>
          }
        />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route path="/chat" element={<ChatWindow />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
