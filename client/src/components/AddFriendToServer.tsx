import { useState } from "react";

function AddFriendToServer({ serverId, onClose }: { serverId: number; onClose: () => void }) {
  const [email, setEmail] = useState("");
  const [error, setError] = useState<string>("");
  const [successMessage, setSuccessMessage] = useState<string>("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email) {
      setError("L'email est obligatoire.");
      return;
    }

    const response = await fetch(`http://localhost:8080/servers/${serverId}/add-user`, {
      method: "POST",
      body: JSON.stringify({ email }),
      credentials: "include", // ✅ Include JWT cookie
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(text || "Erreur lors de la création du serveur");
    }

    // Simulate the logic for sending the email and inviting the friend
    console.log("Invitation envoyée à l'email :", email, "pour le serveur :", serverId);

    // Simulate success
    setSuccessMessage(`Invitation envoyée à ${email}`);

    // Reset form fields after submission
    setEmail("");
    setError(""); // Reset error messages
  };

  return (
    <div className="fixed inset-0 z-50 bg-gray-800 bg-opacity-70 flex justify-center items-center">
      <div className="bg-black text-white p-6 rounded-lg shadow-lg w-96">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">Inviter un ami au serveur</h2>
          <button
            onClick={onClose} // Close the modal when clicked
            className="text-red-400 hover:text-red-300"
          >
            ❌
          </button>
        </div>

        {error && <div className="mt-4 text-center text-red-500 text-sm">{error}</div>}
        {successMessage && <div className="mt-4 text-center text-green-500 text-sm">{successMessage}</div>}

        <form onSubmit={handleSubmit} className="mt-4">
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
            className="w-full px-3 py-2 mt-2 text-white bg-gray-800 border border-gray-600 rounded focus:outline-none focus:border-blue-400"
          />

          <div className="mt-4 text-center">
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
