import { useState } from "react";

function AddFriendToServer({ serverId }: { serverId: number }) {
  const [email, setEmail] = useState("");
  const [error, setError] = useState<string>("");
  const [successMessage, setSuccessMessage] = useState<string>("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!email) {
      setError("L'email est obligatoire.");
      return;
    }

    // Logique pour envoyer l'email et inviter l'ami
    // Par exemple, on peut utiliser une API pour inviter l'ami.
    console.log("Invitation envoyée à l'email :", email, "pour le serveur :", serverId);

    // Simulation d'un succès
    setSuccessMessage(`Invitation envoyée à ${email}`);

    // Réinitialisation des champs après soumission
    setEmail("");
    setError(""); // Réinitialiser les erreurs
  };

  return (
    <div className="w-full h-full flex items-center justify-center">
      <div className="flex items-center justify-center min-h-screen bg-gray-900 text-white">
        <form
          onSubmit={handleSubmit}
          className="w-full max-w-lg p-6 bg-black border-4 border-blue-400 rounded-lg shadow-md"
        >
          <div className="text-center">
            <h2 className="text-xl font-bold">Inviter un ami au serveur</h2>
            <hr className="w-1/2 mx-auto mt-2 border-blue-400" />
          </div>

          {error && (
            <div className="mt-4 text-center text-red-500 text-sm">{error}</div>
          )}
          {successMessage && (
            <div className="mt-4 text-center text-green-500 text-sm">{successMessage}</div>
          )}

          <div className="mt-4">
            <label htmlFor="email" className="block text-sm font-bold text-gray-300">
              Email de l'ami
            </label>
            <input
              id="email"
              type="email"
              placeholder="Email de l'ami"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="w-full px-3 py-2 mt-1 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
            />
          </div>

          <div className="mt-6 flex justify-center">
            <input
              type="submit"
              value="Inviter"
              className="px-6 py-2 text-white bg-blue-700 rounded cursor-pointer hover:bg-blue-400 transition"
            />
          </div>
        </form>
      </div>
    </div>
  );
}

export default AddFriendToServer;
