// src/app/messages/page.js
import { useState, useEffect } from 'react';
import axios from 'axios';

export default function Messages() {
    const [message, setMessage] = useState('');
    const [messages, setMessages] = useState([]);

    const fetchMessages = async () => {
        try {
            const response = await axios.get('http://localhost:8080/messages');
            setMessages(response.data);
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        fetchMessages();
    }, []);

    const handleSendMessage = async () => {
        try {
            const token = localStorage.getItem('token');
            await axios.post('http://localhost:8080/send', { message }, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            fetchMessages();
            setMessage('');
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <div>
            <h1>Messages</h1>
            <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
            />
            <button onClick={handleSendMessage}>Send</button>
            <ul>
                {messages.map((msg) => (
                    <li key={msg.ID}>{msg.Content}</li>
                ))}
            </ul>
        </div>
    );
}
