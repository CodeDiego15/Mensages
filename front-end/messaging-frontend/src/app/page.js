// src/app/page.js
"use client"; // Marca este archivo como un Client Component

import { useState } from 'react';
import axios from 'axios';
import './login.css'; // Asegúrate de crear este archivo para los estilos

export default function Home() {
    const [view, setView] = useState('login'); // 'login' or 'verify'
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [password, setPassword] = useState('');
    const [verificationCode, setVerificationCode] = useState('');
    const [isVerifying, setIsVerifying] = useState(false);

    const handleLogin = async () => {
        try {
            if (view === 'login') {
                const response = await axios.post('http://localhost:8080/login', { email, phone, password });
                localStorage.setItem('token', response.data.token);
                setView('verify');
            } else if (view === 'verify') {
                await axios.post('http://localhost:8080/verify', { verificationCode });
                window.location.href = '/messages'; // Redirige a la página de mensajes
            }
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <div className="login-container">
            {view === 'login' && (
                <div className="login-form">
                    <h1>Login</h1>
                    <div className="input-group">
                        <input
                            type="email"
                            placeholder="Email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        <button onClick={handleLogin}>Login with Email</button>
                    </div>
                    <p>or</p>
                    <div className="input-group">
                        <input
                            type="text"
                            placeholder="Phone Number"
                            value={phone}
                            onChange={(e) => setPhone(e.target.value)}
                        />
                        <button onClick={handleLogin}>Login with Phone</button>
                    </div>
                </div>
            )}

            {view === 'verify' && (
                <div className="verify-form">
                    <h1>Verify Your Account</h1>
                    <input
                        type="text"
                        placeholder="Verification Code"
                        value={verificationCode}
                        onChange={(e) => setVerificationCode(e.target.value)}
                    />
                    <button onClick={handleLogin}>Verify</button>
                </div>
            )}
        </div>
    );
}


