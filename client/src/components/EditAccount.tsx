import { useState } from 'react';
import '../styles/EditAccount.css';

function EditAccount() {
    const [displayName, setDisplayName] = useState('');
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState<string>('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!username || !email || !password) {
            setError('Tous les champs sont obligatoires.');
            return;
        }

        const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailPattern.test(email)) {
            setError('Veuillez entrer un email valide.');
            return;
        }

        console.log('Nom d\'affichage:', displayName);
        console.log('Nom d\'utilisateur:', username);
        console.log('Email:', email);
        console.log('Mot de passe:', password);
    };

    const togglePasswordVisibility = () => {
        setShowPassword(!showPassword);
    };

    return (
        <form onSubmit={handleSubmit} className="auth-form">
            <center>
                <div className="title-container">
                    <h2>Modifier le compte</h2>
                </div>
                <hr className="title-line" />
            </center>

            {error && <div className="error-message">{error}</div>}

            <div>
                <label htmlFor="displayName">Nom d'affichage</label>
                <input
                    id="displayName"
                    type="text"
                    placeholder="Nom d'affichage"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                />
            </div>
            <div>
                <label htmlFor="username">Nom d'utilisateur</label>
                <input
                    id="username"
                    type="text"
                    placeholder="Nom d'utilisateur"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                />
            </div>
            <div>
                <label htmlFor="email">Email</label>
                <input
                    id="email"
                    type="email"
                    placeholder="exemple@mail.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />
            </div>
            <div>
                <label htmlFor="password">Mot de passe</label>
                <div className="password-container">
                    <input
                        id="password"
                        type={showPassword ? 'text' : 'password'}
                        placeholder="Nouveau mot de passe"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                    <button type="button" onClick={togglePasswordVisibility} className="password-toggle-btn">
                        {showPassword ? 'Cacher' : 'Afficher'}
                    </button>
                </div>
            </div>
            <div>
                <input type="submit" value="Sauvegarder les modifications" />
            </div>
        </form>
    );
}

export default EditAccount;
