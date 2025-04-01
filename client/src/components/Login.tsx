import { useState } from "react";
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

    try {
      axios.interceptors.request.use(
        (config) => {
          config.withCredentials = true;
          return config;
        },
        (error) => {
          return Promise.reject(error);
        },
      );
      const res = await axios.post("http://localhost:8080/login", {
        email,
        password,
      });
      console.log(res);
      navigate("/");
    } catch (error) {
      console.log(error);
    }
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-900 text-white">
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-md p-6 bg-black border-4 border-blue-400 rounded-lg shadow-md"
      >
        <div className="text-center">
          <h2 className="text-xl font-bold">Authentification</h2>
          <hr className="w-1/2 mx-auto mt-2 border-blue-400" />
        </div>

        {error && (
          <div className="mt-4 text-center text-red-500 text-sm">{error}</div>
        )}

        <div className="mt-4">
          <label
            htmlFor="email"
            className="block text-sm font-bold text-gray-300"
          >
            Email
          </label>
          <input
            id="email"
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

        <div className="mt-4">
          <label
            htmlFor="password"
            className="block text-sm font-bold text-gray-300"
          >
            Mot de passe
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              placeholder="********"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
            />
            <button
              type="button"
              onClick={togglePasswordVisibility}
              className="absolute right-3 top-1/2 transform -translate-y-1/2 text-blue-400 text-sm font-bold focus:outline-none"
            >
              {showPassword ? "Cacher" : "Afficher"}
            </button>
          </div>
        </div>

        <div className="mt-6 flex justify-center">
          <input
            type="submit"
            value="Se connecter"
            className="px-6 py-2 text-white bg-blue-700 rounded cursor-pointer hover:bg-blue-400 transition"
          />
        </div>
      </form>
    </div>
  );
}

export default Login;
