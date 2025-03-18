import { useState } from "react";

function EditAccount() {
  const [displayName, setDisplayName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string>("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!username || !email || !password) {
      setError("Tous les champs sont obligatoires.");
      return;
    }

    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailPattern.test(email)) {
      setError("Veuillez entrer un email valide.");
      return;
    }

    console.log("Nom d'affichage:", displayName);
    console.log("Nom d'utilisateur:", username);
    console.log("Email:", email);
    console.log("Mot de passe:", password);
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
          <h2 className="text-xl font-bold">Modifier le compte</h2>
          <hr className="w-1/2 mx-auto mt-2 border-blue-400" />
        </div>

        {error && (
          <div className="mt-4 text-center text-red-500 text-sm">{error}</div>
        )}

        <div className="mt-4">
          <label
            htmlFor="displayName"
            className="block text-sm font-bold text-gray-300"
          >
            Nom d'affichage
          </label>
          <input
            id="displayName"
            type="text"
            placeholder="Nom d'affichage"
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

        <div className="mt-4">
          <label
            htmlFor="username"
            className="block text-sm font-bold text-gray-300"
          >
            Nom d'utilisateur
          </label>
          <input
            id="username"
            type="text"
            placeholder="Nom d'utilisateur"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />
        </div>

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
            placeholder="exemple@mail.com"
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
              placeholder="Nouveau mot de passe"
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
            value="Sauvegarder les modifications"
            className="px-6 py-2 text-white bg-blue-700 rounded cursor-pointer hover:bg-blue-400 transition"
          />
        </div>
      </form>
    </div>
  );
}

export default EditAccount;
