import { useState } from "react";
import "../styles/Auth.css";
import axios from "axios";
import { useNavigate } from "react-router-dom";

function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string>("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!email || !password) {
      setError("Tous les champs sont obligatoires.");
      return;
    }

    console.log("Email:", email);
    console.log("Mot de passe:", password);

    try {
      axios.interceptors.request.use(
        (config) => {
          config.withCredentials = true;
          return config;
        },
        (error) => {
          return Promise.reject(error);
        }
      );
      const res = await axios.post("http://localhost:8080/login", {
        email,
        password,
      });
      console.log(res);
  /*     navigate("/"); */
    } catch (error) {
      console.log(error);
    }
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <form onSubmit={handleSubmit} className="auth-form">
      <center>
        <div className="title-container">
          <h2>Authentification</h2>
        </div>
        <hr className="title-line" />
      </center>

      {error && <div className="error-message">{error}</div>}

      <div>
        <label htmlFor="email">Email</label>
        <input id="email" type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required />
      </div>
      <div>
        <label htmlFor="password">Mot de passe</label>
        <div className="password-container">
          <input id="password" type={showPassword ? "text" : "password"} placeholder="********" value={password} onChange={(e) => setPassword(e.target.value)} required />
          <button type="button" onClick={togglePasswordVisibility} className="password-toggle-btn">
            {showPassword ? "Cacher" : "Afficher"}
          </button>
        </div>
      </div>
      <div>
        <input type="submit" value="Se connecter" />
      </div>
    </form>
  );
}

export default Login;
