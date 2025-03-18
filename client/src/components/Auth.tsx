import { useState } from 'react';
import '../styles/Auth.css';

function Connect() {
    const [name, setName] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const [error, setError] = useState<string>('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        
        if (!name || !password) {
            setError('Tous les champs sont obligatoires.');
            return;
        }

        setError('');
        console.log('Nom:', name);
        console.log('Mot de passe:', password);
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
                <label htmlFor="username">Nom d'utilisateur ou email</label>
                <input
                    id="username"
                    type="text"
                    placeholder="Nom / Email"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                />
            </div>
            <div>
                <label htmlFor="password">Mot de passe</label>
                <div className="password-container">
                    <input
                        id="password"
                        type={showPassword ? 'text' : 'password'}
                        placeholder="********"
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
                <input type="submit" value="Se connecter" />
            </div>
        </form>
    );
}

export default Connect;
