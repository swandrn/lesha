import { BrowserRouter, Route, Routes } from "react-router-dom";
import "./App.css";
import { Chat } from "./components/Chat";
import { ServerList } from "./components/ServerList";
import Register from "./components/Register";
import Login from "./components/Login";
import ConnectedUserRoute from "./components/middlewares/connectedUserRoute";


function App() {
  return (
    <BrowserRouter>
    <Routes>
      <Route path="/" element={<ConnectedUserRoute><Chat /></ConnectedUserRoute>} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      
      <Route path="/servers" element={ <ServerList />} />
    </Routes>
  </BrowserRouter>
  
  );
}

export default App;
